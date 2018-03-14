package inmemorybuilders

import (
	"log"
	"sync"

	"github.com/rmulton/riw_project/indexbuilders"
	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/indexes/requestableIndexes"
)

// InMemoryBuilder is an index builder that only uses the memory, nothing is stored on the disk
type InMemoryBuilder struct {
	index           *indexes.Index
	parentWaitGroup *sync.WaitGroup
	readingChannel  indexes.ReadingChannel
	docCounter      int
}

// NewInMemoryBuilder returns a InMemoryBuilder with an empty index
func NewInMemoryBuilder(readingChannel indexes.ReadingChannel, routines int, parentWaitGroup *sync.WaitGroup) *InMemoryBuilder {
	index := indexes.NewEmptyIndex()
	return &InMemoryBuilder{
		index:           index,
		parentWaitGroup: parentWaitGroup,
		readingChannel:  readingChannel,
	}
}

// Build fills the index in with the reading channel
func (builder *InMemoryBuilder) Build() {
	defer builder.parentWaitGroup.Done()
	var wg sync.WaitGroup
	wg.Add(1)
	indexbuilders.FillIndex(builder.index, builder.readingChannel, &wg)
	wg.Wait()
	builder.finish()
	log.Printf("Done filling with %d documents", builder.docCounter)
}

// GetIndex returns a requestable index based on the built index
func (builder *InMemoryBuilder) GetIndex() requestableIndexes.RequestableIndex {
	return requestableIndexes.InMemoryIndexFromIndex(builder.index)
}

// finish finishes the index by computing tf-idf scores based on the frequency scores
func (builder *InMemoryBuilder) finish() {
	builder.index.ToTfIdf(builder.index.GetDocCounter())
}
