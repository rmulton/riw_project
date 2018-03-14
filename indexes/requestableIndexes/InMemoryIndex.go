package requestableIndexes

import (
	"github.com/rmulton/riw_project/indexes"
)

// InMemoryIndex is a requestable index used for indexes built in memory
type InMemoryIndex struct {
	index *indexes.Index
}

// InMemoryIndexFromIndex returns an InMemoryIndex that contains information from the input index
func InMemoryIndexFromIndex(index *indexes.Index) *InMemoryIndex {
	return &InMemoryIndex{
		index: index,
	}
}

// GetPostingListsForTerms returns the posting lists of some terms
func (index *InMemoryIndex) GetPostingListsForTerms(terms []string) map[string]indexes.PostingList {
	postingListsForTerms := make(map[string]indexes.PostingList)
	for _, term := range terms {
		postingList, _ := index.index.GetPostingListForTerm(term)
		postingListsForTerms[term] = postingList
	}
	return postingListsForTerms
}

// GetDocIDToPath returns the docID to filepath map
func (index *InMemoryIndex) GetDocIDToPath() map[int]string {
	return index.index.GetDocIDToFilePath()
}

// GetDocCounter returns the number of documents used to build the index
func (index *InMemoryIndex) GetDocCounter() int {
	return index.index.GetDocCounter()
}

// IsInTheIndex returns whether the term has a posting list in the index
func (index *InMemoryIndex) IsInTheIndex(term string) bool {
	return index.index.IsInTheIndex(term)
}
