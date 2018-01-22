package indexBuilders

import (
	"../indexes"
	"../utils"
	"sync"
	"log"
)

type BufferIndex struct {
	Mux *sync.Mutex
	bufferSize int
	currentSize int
	writingChannel indexes.WritingChannel
	index *indexes.Index

}

func NewBufferIndex(bufferSize int, writingChannel indexes.WritingChannel) *BufferIndex {
	var mux sync.Mutex
	index := indexes.NewEmptyIndex()
	return &BufferIndex{
		writingChannel: writingChannel,
		Mux: &mux,
		bufferSize: bufferSize,
		index : index,
	}
}

// Used to fill the posting lists
func (buffer *BufferIndex) addDocToTerm(docID int, term string) {
	buffer.index.AddDocToTerm(docID, term)
	buffer.currentSize++
}

// Add a new document in the index so that index keep trace of docID -> doc
func (buffer *BufferIndex) addDocToIndex(docID int, docPath string) {
	if buffer.currentSize >= buffer.bufferSize && buffer.bufferSize != -1 { // NB: This is a very important decision for the system. Using the biggest posting list might not be the best one.
		buffer.writeBiggestPostingList()
	}
	buffer.index.AddDocToIndex(docID, docPath)
}

// /!\ MAYBE IT WOULD BE BETTER TO HAVE A ROUTINE WORKING ON WRITING THE POSTING LISTS
// AND USE RWLock instead of Lock
func (buffer *BufferIndex) writeBiggestPostingList() {
	buffer.Mux.Lock()
	// Find the longest posting list
	var termWithLongestPostingList string
	max := -1
	postingLists := buffer.index.GetPostingLists()
	for term, postingList := range postingLists {
		if len(postingList) > max {
			termWithLongestPostingList = term
			max = len(postingList)
		}
	}

	// Copy it to let other routines get access to the index, and remove it from the index
	longestPostingList := postingLists[termWithLongestPostingList]
	buffer.index.ClearPostingListFor(termWithLongestPostingList)
	
	buffer.currentSize -= max
	buffer.Mux.Unlock()

	// log.Printf("Writing posting list for %s", termWithLongestPostingList)
	// go longestPostingList.appendToTermFile(termWithLongestPostingList, index.writingChannel)
	buffer.appendToTermFile(longestPostingList, termWithLongestPostingList)
}

// TODO : avoid code repition by building buffer.appendPostingListOnDisk(term)

// Should be done by the buffer index instead
func (buffer *BufferIndex) appendToTermFile(postingList indexes.PostingList, term string) {
	buffer.writingChannel <- indexes.NewBufferPostingList(term, postingList)
}

// When no more documents are to be read
func (buffer *BufferIndex) writeAllPostingLists() {
	defer close(buffer.writingChannel)
	log.Printf("Writing remaining posting lists")
	for term, postingList := range buffer.index.GetPostingLists() {
		// fmt.Printf("Writing posting list for %s", term)
		buffer.appendToTermFile(postingList, term)
		// go postingList.appendToTermFile(term, index.writingChannel)
	}
}

func (buffer *BufferIndex) toTfIdf(corpusSize int) {
	buffer.index.ToTfIdf(corpusSize)
}

func (buffer *BufferIndex) writeDocIDToFilePath(path string) {
	utils.WriteGob(path, buffer.index.GetDocIDToFilePath())
}