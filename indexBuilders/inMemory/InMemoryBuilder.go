package inMemory

import (
	"sync"
	"log"
	"github.com/rmulton/riw_project/indexBuilders/inMemory"
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
	var wg sync.WaitGroup
	wg.Add(1)
	fillIndex(builder.index, builder.readingChannel, &wg)
	wg.Wait()
	builder.finish()
	log.Printf("Done filling with %d documents", builder.docCounter)
}

func (builder *InMemoryBuilder) GetIndex() indexes.RequestableIndex {
	// TODO : Add cached version
	return indexes.InMemoryIndexFromIndex(builder.index)
}

func (builder *InMemoryBuilder) finish() {
	builder.index.ToTfIdf(builder.docCounter)
}