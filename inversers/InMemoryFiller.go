package inversers

import (
	"sync"
	"log"
	"../readers"
)

type InMemoryFiller struct {
	index *BufferIndex
	parentWaitGroup *sync.WaitGroup
	readingChannel readingChannel
	docCounter int
}

func NewInMemoryFiller(readingChannel chan *readers.Document, routines int, parentWaitGroup *sync.WaitGroup) *InMemoryFiller {
	index := NewBufferIndex(-1, nil) // simple hack to get the maximum int value
	return &InMemoryFiller{
		index: index,
		parentWaitGroup: parentWaitGroup,
		readingChannel: readingChannel,
	}
	// return an index
}

func (filler *InMemoryFiller) Fill() {
	defer filler.parentWaitGroup.Done()
	filler.readDocs()
	log.Printf("Done filling with %d documents", filler.docCounter)
}

func (filler *InMemoryFiller) readDocs() {
	for doc := range filler.readingChannel {
		filler.docCounter++
		filler.addDoc(doc)
	}
	filler.finish()
	log.Printf("Done getting %d documents", filler.docCounter)
}

func (filler *InMemoryFiller) finish() {
	filler.index.toTfIdf()
}

// Add a document to the current block
// NB: Might evolve to allow filling several blocks at the same time
func (filler *InMemoryFiller) addDoc(doc *readers.Document) {

	// TODO: use routines to parallelize addDocToIndex and addDocToTerm

	// Add the path to the doc to the map
	filler.index.addDocToIndex(doc.Id, doc.Path)

	for _, term := range *doc.NormalizedTokens {
		//	Add this document to the posting list of this term
		filler.index.addDocToTerm(doc.Id, term)
	}
}
