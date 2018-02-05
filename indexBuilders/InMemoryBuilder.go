package indexBuilders

import (
	"sync"
	"log"
	"../indexes"
)

type InMemoryBuilder struct {
	index *indexes.Index
	parentWaitGroup *sync.WaitGroup
	readingChannel indexes.ReadingChannel
	docCounter int
}

func NewInMemoryBuilder(readingChannel indexes.ReadingChannel, routines int, parentWaitGroup *sync.WaitGroup) *InMemoryBuilder {
	index := indexes.NewEmptyIndex()
	return &InMemoryBuilder{
		index: index,
		parentWaitGroup: parentWaitGroup,
		readingChannel: readingChannel,
	}
}

func (builder *InMemoryBuilder) Build() {
	defer builder.parentWaitGroup.Done()
	builder.readDocs()
	log.Printf("Done filling with %d documents", builder.docCounter)
}

func (builder *InMemoryBuilder) GetIndex() indexes.RequestableIndex {
	// TODO : Add cached version
	return indexes.InMemoryIndexFromIndex(builder.index)
}

func (builder *InMemoryBuilder) readDocs() {
	for doc := range builder.readingChannel {
		builder.addDoc(doc)
		builder.docCounter++
	}
	builder.finish()
	log.Printf("Done getting %d documents", builder.docCounter)
}

func (builder *InMemoryBuilder) finish() {
	builder.index.ToTfIdf(builder.docCounter)
}

// Add a document to the current block
// NB: Might evolve to allow filling several blocks at the same time
func (builder *InMemoryBuilder) addDoc(doc indexes.Document) {

	// TODO: use routines to parallelize addDocToIndex and addDocToTerm

	// Add the path to the doc to the map
	builder.index.AddDocToIndex(doc.ID, doc.Path)

	for _, term := range doc.NormalizedTokens {
		//	Add this document to the posting list of this term
		builder.index.AddDocToTerm(doc.ID, term)
	}
}
