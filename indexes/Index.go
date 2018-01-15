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

func (index *Index) GetDocIDToFilePath() map[int]string {
	return index.docIDToFilePath
}

func (index *Index) PrintPostings() {
	for _, postingList := range index.postingLists {
		fmt.Printf("%v\n", postingList)
	}
}

func (index *Index) ToTfIdf(corpusSize int) {
	// TODO: Parallelize
	for _, postingList := range index.postingLists {
		postingList.TfIdf(corpusSize)
	}
}

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
	emptyPostingList := make(PostingList)
	index.postingLists[term] = emptyPostingList
}

// TODO: computations on the posting lists should be done by the index
