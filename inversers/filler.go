package inversers

import (
	"sync"
	"fmt"
	"os"
	"log"
	"../utils"
	"../readers"
)

type Filler struct {
	index *BufferIndex
	parentWaitGroup *sync.WaitGroup
	waitGroup *sync.WaitGroup
	folder string
	bufferSize int
	readingChannel readingChannel
	writingChannel writingChannel
	docCounter int
}

type readingChannel chan *readers.Document
type writingChannel chan *toWrite

func NewFiller(bufferSize int, folder string, readingChannel chan *readers.Document, routines int, parentWaitGroup *sync.WaitGroup) *Filler {
	var waitGroup sync.WaitGroup
	writingChannel := make(chan *toWrite)
	index := NewBufferIndex(500000, writingChannel)
	return &Filler{
		index: index,
		bufferSize: bufferSize,
		parentWaitGroup: parentWaitGroup,
		waitGroup: &waitGroup,
		folder: folder,
		readingChannel: readingChannel,
		writingChannel: writingChannel,
	}
	// return an index
}

func (filler *Filler) Fill() {
	defer filler.parentWaitGroup.Done()
	filler.waitGroup.Add(2)
	go filler.readDocs()
	go filler.writePostingLists() // the block is copied to allow continuing operations on the filler
	filler.waitGroup.Wait()
	log.Printf("Done filling with %d documents", filler.docCounter)
}

func (filler *Filler) readDocs() {
	defer filler.waitGroup.Done()
	for doc := range filler.readingChannel {
		filler.docCounter++
		filler.addDoc(doc)
	}
	filler.finish()
	log.Printf("Done getting %d documents", filler.docCounter)
}

func (filler *Filler) finish() {
	filler.index.writeRemainingPostingLists()
	utils.WriteGob("./saved/idToPath.meta", filler.index.docIDToFilePath)
}

// Add a document to the current block
// NB: Might evolve to allow filling several blocks at the same time
func (filler *Filler) addDoc(doc *readers.Document) {

	// TODO: use routines to parallelize addDocToIndex and addDocToTerm

	// Add the path to the doc to the map
	filler.index.addDocToIndex(doc.Id, doc.Path)

	for _, term := range *doc.NormalizedTokens {
		//	Add this document to the posting list of this term
		filler.index.addDocToTerm(doc.Id, term)
	}
}

func (filler *Filler) writePostingLists() {
	defer filler.waitGroup.Done()
	for toWrite := range filler.writingChannel {
		// log.Printf("Getting posting list for %s", toWrite.term)
		// Write it to the disk
		termFile := fmt.Sprintf("./saved/%s.postings", toWrite.term)
		// If the file exists append it
		if _, err := os.Stat(termFile); err == nil {
			postingListSoFar := make(postingList)
			err := utils.ReadGob(termFile, &postingListSoFar)
			if err != nil {
				panic(err)
			}
			// Merge the current posting list
			postingListSoFar.mergeWith(toWrite.postingList)
			// Write it to file
			utils.WriteGob(termFile, postingListSoFar)
		// Otherwise create it
		} else {
			utils.WriteGob(termFile, toWrite.postingList)
		}
	}
}