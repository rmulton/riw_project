package indexBuilders

import (
	"sync"
	"fmt"
	"os"
	"log"
	"path/filepath"
	"../utils"
	"../indexes"
)

type OnDiskBuilder struct {
	index *BufferIndex
	parentWaitGroup *sync.WaitGroup
	waitGroup *sync.WaitGroup
	folder string
	bufferSize int
	readingChannel indexes.ReadingChannel
	writingChannel indexes.WritingChannel
	docCounter int
}


func NewOnDiskBuilder(bufferSize int, folder string, readingChannel chan *indexes.Document, routines int, parentWaitGroup *sync.WaitGroup) *OnDiskBuilder {
	var waitGroup sync.WaitGroup
	writingChannel := make(chan *indexes.BufferPostingList)
	index := NewBufferIndex(1000000000, writingChannel)
	return &OnDiskBuilder{
		index: index,
		bufferSize: bufferSize,
		parentWaitGroup: parentWaitGroup,
		waitGroup: &waitGroup,
		folder: folder,
		readingChannel: readingChannel,
		writingChannel: writingChannel,
	}
}

func (builder *OnDiskBuilder) Build() {
	defer builder.parentWaitGroup.Done()
	builder.waitGroup.Add(2)
	go builder.readDocs()
	go builder.writePostingLists() // the block is copied to allow continuing operations on the builder
	builder.waitGroup.Wait()
	log.Printf("Done filling with %d documents", builder.docCounter)
}
 
func (builder *OnDiskBuilder) GetIndex() *indexes.OnDiskIndex {
	return indexes.OnDiskIndexFromFolder("./saved/")
}


func (builder *OnDiskBuilder) readDocs() {
	defer builder.waitGroup.Done()
	for doc := range builder.readingChannel {
		builder.docCounter++
		builder.addDoc(doc)
	}
	builder.finish()
	log.Printf("Done getting %d documents", builder.docCounter)
}

func (builder *OnDiskBuilder) finish() {
	builder.index.writeAllPostingLists()
	// Warning: Do not use go routine otherwise it will conflict with the next line
	builder.toTfIdf()
	builder.writeDocIDToFilePath("./saved/idToPath.meta")
}

func (builder *OnDiskBuilder) writeDocIDToFilePath(path string) {
	builder.index.writeDocIDToFilePath(path)
}

func (builder *OnDiskBuilder) toTfIdf() {
	log.Println("Computing tf-idf scores from frequencies")
	filepath.Walk("./saved/", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			err, postingList := indexes.PostingListFromFile(path)
			if err != nil {
				return err
			}
			// TODO: Parallelize
			postingList.TfIdf(builder.docCounter) // TODO ? check that docCounter = len(os.listdir)
			utils.WriteGob(path, postingList)
		}
		return nil
	})
}

// Add a document to the current block
// NB: Might evolve to allow filling several blocks at the same time
func (builder *OnDiskBuilder) addDoc(doc *indexes.Document) {

	// TODO: use routines to parallelize addDocToIndex and addDocToTerm

	// Add the path to the doc to the map
	builder.index.addDocToIndex(doc.ID, doc.Path)

	for _, term := range doc.NormalizedTokens {
		//	Add this document to the posting list of this term
		builder.index.addDocToTerm(doc.ID, term)
	}
}

func (builder *OnDiskBuilder) writePostingLists() {
	defer builder.waitGroup.Done()
	for toWrite := range builder.writingChannel {
		// log.Printf("Getting posting list for %s", toWrite.term)
		// Write it to the disk
		termFile := fmt.Sprintf("./saved/%s.postings", toWrite.Term)
		// If the file exists append it
		if _, err := os.Stat(termFile); err == nil {
			err, postingListSoFar := indexes.PostingListFromFile(termFile)
			if err != nil {
				panic(err)
			}
			// Merge the current posting list
			postingListSoFar.MergeWith(toWrite.PostingList)
			// Write it to file
			utils.WriteGob(termFile, postingListSoFar)
		// Otherwise create it
		} else {
			utils.WriteGob(termFile, toWrite.PostingList)
		}
	}
}