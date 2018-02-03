package indexBuilders

import (
	"../indexes"
	"../utils"
	"sync"
	// "log"
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
	// log.Printf("Appending to term file from writeBiggestPostingList term %s", termWithLongestPostingList)
	buffer.appendToTermFile(longestPostingList, termWithLongestPostingList, false)
}

// TODO : avoid code repition by building buffer.appendPostingListOnDisk(term)

// Should be done by the buffer index instead
func (buffer *BufferIndex) appendToTermFile(postingList indexes.PostingList, term string, replace bool) {
	// Here is the problem: the score is added to the file instead of replacing it
	// TODO: Clean the mechanics that's below
	var bufferPostingList indexes.BufferPostingList
	if replace {
		bufferPostingList = indexes.NewReplacingBufferPostingList(term, postingList)
	} else {
		bufferPostingList = indexes.NewBufferPostingList(term, postingList)
	}
	buffer.writingChannel <- bufferPostingList
}

func (buffer *BufferIndex) writePostingListForTerms(terms map[string]bool) {
	for term, _ := range terms {
		postingList := buffer.index.GetPostingListForTerm(term)
		buffer.appendToTermFile(postingList, term, false)
	}
}

// When no more documents are to be read
// Used only for InMemoryBuilder
func (buffer *BufferIndex) writeAllPostingLists() {
	defer close(buffer.writingChannel)
	// log.Printf("Writing remaining posting lists")
	for term, postingList := range buffer.index.GetPostingLists() {
		// fmt.Printf("Writing posting list for %s", term)
		// log.Printf("Appending to term file from writeAllPostingList term %s", term)
		buffer.appendToTermFile(postingList, term, true)
		// go postingList.appendToTermFile(term, index.writingChannel)
	}
}

func (buffer *BufferIndex) toTfIdf(corpusSize int) {
	buffer.index.ToTfIdf(corpusSize)
}

func (buffer *BufferIndex) toTfIdfTerms(corpusSize int, terms map[string]bool) {
	buffer.index.ToTfIdfTerms(corpusSize, terms)
}

func (buffer *BufferIndex) writeDocIDToFilePath(path string) {
	utils.WriteGob(path, buffer.index.GetDocIDToFilePath())
}

func (buffer *BufferIndex) getPostingListForTerm(term string) indexes.PostingList {
	return buffer.index.GetPostingListForTerm(term)
}