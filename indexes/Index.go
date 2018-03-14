package indexes

import (
	"fmt"
)

// Index stores an index for a collection
type Index struct {
	currentSize     int
	postingLists    map[string]PostingList
	docIDToFilePath map[int]string
}

// NewEmptyIndex generates an new and empty index ready to be filled
func NewEmptyIndex() *Index {
	postingLists := make(map[string]PostingList)
	docIDToFilePath := make(map[int]string)
	return &Index{
		postingLists:    postingLists,
		docIDToFilePath: docIDToFilePath,
	}
}

// NewIndexWithDocIDToPath allows recovering an index from the hard drive
func NewIndexWithDocIDToPath(docIDToPath map[int]string) *Index {
	postingLists := make(map[string]PostingList)
	return &Index{
		postingLists:    postingLists,
		docIDToFilePath: docIDToPath,
	}
}

// GetPostingLists returns all the posting lists stored in the index
func (index *Index) GetPostingLists() map[string]PostingList {
	return index.postingLists
}

// GetPostingListForTerm returns the posting list of a specific term
func (index *Index) GetPostingListForTerm(term string) (PostingList, bool) {
	postingList, exists := index.postingLists[term]
	return postingList, exists
}

// SetPostingListForTerm allows modifying a posting list
func (index *Index) SetPostingListForTerm(postingList PostingList, term string) {
	index.postingLists[term] = postingList
}

// GetDocIDToFilePath returns the docID to filepath map
func (index *Index) GetDocIDToFilePath() map[int]string {
	return index.docIDToFilePath
}

// PrintPostings prints the posting lists
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

// ToTfIdfTerms transforms the frequencies in the posting lists into tf-idf scores
func (index *Index) ToTfIdfTerms(corpusSize int, terms map[string]bool) {
	for term, _ := range terms {
		index.postingLists[term].TfIdf(corpusSize)
	}
}

// NB: In some cases, it might be usefull to merge AddDocToTerm and AddDocToIndex.
// However, we chose not to because of the on disk index use case. Indeed, the document counter
// shouldn't be held by the index but by the buffer index.

// AddDocToTerm is used to fill in the posting lists
func (index *Index) AddDocToTerm(docID int, term string) {
	_, exists := index.postingLists[term]
	if !exists {
		index.postingLists[term] = make(PostingList)
	}
	index.postingLists[term][docID]++
	index.currentSize++
}

// AddDocToIndex adds a new document in the index so that index keep trace of docID -> doc
func (index *Index) AddDocToIndex(docID int, docPath string) {
	index.docIDToFilePath[docID] = docPath
}

// ClearPostingListFor removes a posting list from the index
func (index *Index) ClearPostingListFor(term string) {
	delete(index.postingLists, term)
}

// GetDocCounter returns the document counter of the index
func (index *Index) GetDocCounter() int {
	return len(index.docIDToFilePath)
}

// IsInTheIndex returns whether the term has a posting list in the index
func (index *Index) IsInTheIndex(term string) bool {
	_, exists := index.postingLists[term]
	return exists
}
