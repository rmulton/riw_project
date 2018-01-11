package inversers

import (
	"sync"
	"log"
)

// BufferIndex is used both as a reversed index for when index can be held in memory and as a buffer when it can't
type BufferIndex struct {
	Mux *sync.Mutex
	bufferSize int
	corpusSize int
	bufferCounter int
	postingLists map[string]PostingList
	docIDToFilePath map[int]string
	writingChannel writingChannel
}

func NewBufferIndex(bufferSize int, writingChannel writingChannel) *BufferIndex {
	var mux sync.Mutex
	postingLists := make(map[string]PostingList)
	docIDToFilePath := make(map[int]string)
	return &BufferIndex{
		writingChannel: writingChannel,
		Mux: &mux,
		bufferSize: bufferSize,
		postingLists: postingLists,
		docIDToFilePath: docIDToFilePath,
	}
}

// Used to fill the posting lists
func (index *BufferIndex) addDocToTerm(docID int, term string) {
	_, exists := index.postingLists[term]
	if !exists {
		index.postingLists[term] = make(PostingList)
	}
	index.postingLists[term][docID]++
	index.bufferCounter++
}

// Add a new document in the index so that index keep trace of docID -> doc
func (index *BufferIndex) addDocToIndex(docID int, docPath string) {
	if index.bufferCounter >= index.bufferSize && index.bufferSize != -1 { // NB: This is a very important decision for the system. Using the biggest posting list might not be the best one.
		index.writeBiggestPostingList()
	}
	index.docIDToFilePath[docID] = docPath
	index.corpusSize++
}

func (index *BufferIndex) writeBiggestPostingList() {
	index.Mux.Lock()
	// Find the longest posting list
	var termWithLongestPostingList string
	max := -1
	for term, postingList := range index.postingLists {
		if len(postingList) > max {
			termWithLongestPostingList = term
			max = len(postingList)
		}
	}

	// Copy it to let other routines get access to the index, and remove it from the index
	longestPostingList := index.postingLists[termWithLongestPostingList]
	emptyPostingList := make(PostingList)
	index.postingLists[termWithLongestPostingList] = emptyPostingList

	index.bufferCounter -= max
	index.Mux.Unlock()

	// log.Printf("Writing posting list for %s", termWithLongestPostingList)
	// go longestPostingList.appendToTermFile(termWithLongestPostingList, index.writingChannel)
	longestPostingList.appendToTermFile(termWithLongestPostingList, index.writingChannel)
}

// When no more documents are to be read
func (index *BufferIndex) writeAllPostingLists() {
	defer close(index.writingChannel)
	log.Printf("Writing remaining posting lists")
	for term, postingList := range index.postingLists {
		// fmt.Printf("Writing posting list for %s", term)
		postingList.appendToTermFile(term, index.writingChannel)
		// go postingList.appendToTermFile(term, index.writingChannel)
	}
}

func (index *BufferIndex) toTfIdf() {
	for _, postingList := range index.postingLists {
		postingList.tfIdf(index.corpusSize)
	}
}