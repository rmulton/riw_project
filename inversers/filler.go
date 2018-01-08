package inversers

import (
	"sync"
	"fmt"
	"log"
	"os"
	"../utils"
	"../readers"
)

type Filler struct {
	index *Index
	parentWaitGroup *sync.WaitGroup
	waitGroup *sync.WaitGroup
	sem chan bool
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
	sem := make(chan bool, routines)
	index := NewIndex(50000, writingChannel)
	return &Filler{
		index: index,
		bufferSize: bufferSize,
		parentWaitGroup: parentWaitGroup,
		waitGroup: &waitGroup,
		folder: folder,
		readingChannel: readingChannel,
		writingChannel: writingChannel,
		sem: sem,
	}
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
		// filler.sem <- true
		// go filler.addDoc(doc)
		filler.addDoc(doc)
		filler.docCounter++
	}
	// for i := 0; i < cap(filler.sem); i++ {
	// 	<- filler.sem
	// }
	log.Printf("Done getting %s documents", filler.docCounter)
}

// Add a document to the current block
// NB: Might evolve to allow filling several blocks at the same time
func (filler *Filler) addDoc(doc *readers.Document) {
	// defer func(){<-filler.sem}()

	// Add the path to the doc to the map
	filler.index.addDocToIndex(doc.Id, doc.Path)

	for _, term := range *doc.NormalizedTokens {
		//	Add this document to the posting list of this term
		filler.index.addDocToTerm(doc.Id, term)
	}
}

func (filler *Filler) writePostingLists() {
	for toWrite := range filler.writingChannel {
		// Write it to the disk
		termFile := fmt.Sprintf("./saved/%s", toWrite.term)
		// If the file exists append it
		if _, err := os.Stat(termFile); err == nil {
			postingListSoFar := make(postingList)
			err := utils.ReadGob(termFile, postingListSoFar)
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