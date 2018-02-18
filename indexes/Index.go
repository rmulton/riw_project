package indexes

import (
	"fmt"
)

// Index and BufferIndex are different structure so that to be sure to avoid conflicts between the index used by the user and the one used to build it.
type Index struct {
	currentSize int
	postingLists map[string]PostingList
	docIDToFilePath map[int]string
}

func NewEmptyIndex() *Index {
	postingLists := make(map[string]PostingList)
	docIDToFilePath := make(map[int]string)
	return &Index{
		postingLists: postingLists,
		docIDToFilePath: docIDToFilePath,
	}
}

func (index *Index) GetPostingLists() map[string]PostingList {
	return index.postingLists
}

func (index *Index) GetPostingListForTerm(term string) (PostingList, bool) {
	postingList, exists := index.postingLists[term]
	return postingList, exists
}

func (index *Index) SetPostingListForTerm(postingList PostingList, term string) {
	index.postingLists[term] = postingList
}

func (index *Index) GetDocIDToFilePath() map[int]string {
	return index.docIDToFilePath
}

func (index *Index) PrintPostings() {
	for _, postingList := range index.postingLists {
		fmt.Printf("%v\n", postingList)
	}
}

// NB: In some cases, it would be smarter to compute the corpus' size instead of getting
// it as an argument. We chose not to because of the on disk index use case. Indeed, since
// there are two corpus sizes (one that might be calculated by Index and the one of the
// collection), there might be some confusion.
// Our choice was to write Index as a dumb data store with little functionality and to 
// delegate it to the BufferIndex and the builders. This way, Index can be used in a variety
// of cases
func (index *Index) ToTfIdf(corpusSize int) {
	for _, postingList := range index.postingLists {
		postingList.TfIdf(corpusSize)
	}
}

func (index *Index) ToTfIdfTerms(corpusSize int, terms map[string]bool) {
	for term, _ := range terms {
		index.postingLists[term].TfIdf(corpusSize)
	}
}

// NB: In some cases, it might be usefull to merge AddDocToTerm and AddDocToIndex.
// However, we chose not to because of the on disk index use case. Indeed, the document counter
// shouldn't be held by the index but by the buffer index.

// Used to fill the posting lists
func (index *Index) AddDocToTerm(docID int, term string) {
	_, exists := index.postingLists[term]
	if !exists {
		index.postingLists[term] = make(PostingList)
	}
	index.postingLists[term][docID]++
	index.currentSize++
}

// Add a new document in the index so that index keep trace of docID -> doc
func (index *Index) AddDocToIndex(docID int, docPath string) {
	index.docIDToFilePath[docID] = docPath
}

func (index *Index) ClearPostingListFor(term string) {
	delete(index.postingLists, term)
}

func (index *Index) GetDocCounter() int {
	fmt.Println(len(index.docIDToFilePath))
	return len(index.docIDToFilePath)
}