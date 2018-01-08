package inversers

import (
	"sync"
	"fmt"
)

// Index is used both as a reversed index for when index can be held in memory and as a buffer when it can't
type Index struct {
	Mux *sync.Mutex
	bufferSize int
	DocCounter int
	TermCounter int
	postingLists map[string]postingList
	DocIDToFilePath map[int]string
	writingChannel writingChannel
}

func NewIndex(bufferSize int, writingChannel chan *toWrite) *Index {
	var mux sync.Mutex
	postingLists := make(map[string]postingList)
	docIDToFilePath := make(map[int]string)
	return &Index{
		writingChannel: writingChannel,
		Mux: &mux,
		bufferSize: bufferSize,
		postingLists: postingLists,
		DocIDToFilePath: docIDToFilePath,
	}
}

// Used to fill the posting lists
func (index *Index) addDocToTerm(docID int, term string) {
	_, exists := index.postingLists[term]
	if !exists {
		index.postingLists[term] = make(postingList)
	}
	index.postingLists[term][docID]++
}

// Add a new document in the index so that index keep trace of docID -> doc
func (index *Index) addDocToIndex(docID int, docPath string) {
	index.DocIDToFilePath[docID] = docPath
	index.DocCounter++
}

func (index *Index) writeBiggestPostingList() {
	index.Mux.Lock()
	// Find the longest posting list
	termWithLongestPostingList := ""
	max := -1
	for term, postingList := range index.postingLists {
		if len(postingList) > max {
			termWithLongestPostingList = term
			max = len(postingList)
		}
	}

	// Copy it to let other routines get access to the index, and remove it from the index
	longestPostingList := index.postingLists[termWithLongestPostingList]
	emptyPostingList := make(postingList)
	index.postingLists[termWithLongestPostingList] = emptyPostingList

	index.Mux.Unlock()

	fmt.Printf("Writing posting list for %s", termWithLongestPostingList)
	go longestPostingList.appendToTermFile(termWithLongestPostingList, index.writingChannel)
}

// When no more documents are to be read
func (index *Index) writeRemainingPostingLists() {
	fmt.Printf("Writing remaining posting lists")
	for term, postingList := range index.postingLists {
		fmt.Printf("Writing posting list for %s", term)
		go postingList.appendToTermFile(term, index.writingChannel)
	}
	close(index.writingChannel)
}

