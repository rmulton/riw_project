package buildingIndexes

import (
	"github.com/rmulton/riw_project/utils"
	"github.com/rmulton/riw_project/indexes"
	"sync"
	"log"
	"io/ioutil"
)

type BufferIndex struct {
	Mux *sync.Mutex
	bufferSize int
	currentSize int
	docCounter int
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
func (buffer *BufferIndex) AddDocToTerm(docID int, term string) {
	buffer.index.AddDocToTerm(docID, term)
	buffer.currentSize++
}

// Add a new document in the index so that index keep trace of docID -> doc
func (buffer *BufferIndex) AddDocToIndex(docID int, docPath string) {
	buffer.docCounter++
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

	buffer.appendToTermFile(longestPostingList, termWithLongestPostingList, false)
}

// TODO : avoid code repition by buildingIndexes buffer.appendPostingListOnDisk(term)

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
		postingList, exists := buffer.index.GetPostingListForTerm(term)
		if !exists {
			log.Printf("Trying to writing the posting list of %s that is not in the index", term)
		}
		buffer.appendToTermFile(postingList, term, false)
	}
}

// When no more documents are to be read
// Used only for InMemoryBuilder
func (buffer *BufferIndex) writeAllPostingLists() {
	defer close(buffer.writingChannel)
	for term, postingList := range buffer.index.GetPostingLists() {
		buffer.appendToTermFile(postingList, term, true)
	}
}

func (buffer *BufferIndex) toTfIdf(corpusSize int) {
	buffer.index.ToTfIdf(corpusSize)
}

func (buffer *BufferIndex) toTfIdfTerms(terms map[string]bool) {
	buffer.index.ToTfIdfTerms(buffer.docCounter, terms)
}

func (buffer *BufferIndex) writeDocIDToFilePath(path string) {
	utils.WriteGob(path, buffer.index.GetDocIDToFilePath())
}

func (buffer *BufferIndex) getPostingListForTerm(term string) (indexes.PostingList, bool) {
	return buffer.index.GetPostingListForTerm(term)
}

/* Find out which terms are in memory, on disk or both */

func (buffer *BufferIndex) categorizeTerms() (map[string]bool, map[string]bool, map[string]bool) {
	onDiskTerms := getOnDiskTerms()
	inMemoryTerms := buffer.getInMemoryTerms()
	onDiskOnly, inMemoryOnly, onDiskAndInMemory := separate(onDiskTerms, inMemoryTerms)
	return onDiskOnly, inMemoryOnly, onDiskAndInMemory
}

func getOnDiskTerms() map[string]bool {
	onDiskTerms := make(map[string]bool)
	files, err := ioutil.ReadDir("./saved/postings/")
    if err != nil {
        log.Println(err)
    }
    for _, f := range files {
			onDiskTerms[f.Name()] = true
    }
	return onDiskTerms
}

func (buffer *BufferIndex) getInMemoryTerms() map[string]bool {
	inMemoryTerms := make(map[string]bool)
	for term, _ := range buffer.index.GetPostingLists() {
		inMemoryTerms[term] = true
	}
	return inMemoryTerms
}

func separate(first map[string]bool, second map[string]bool) (map[string]bool, map[string]bool, map[string]bool) {
	onlyFirst := make(map[string]bool)
	onlySecond := make(map[string]bool)
	both := make(map[string]bool)
	for firstKey, _ := range first {
		_, exists := second[firstKey]
		if exists {
			both[firstKey] = true
		} else {
			onlyFirst[firstKey] = true
		}
	}
	for secondKey, _ := range second {
		_, exists := both[secondKey]
		if !exists {
			onlySecond[secondKey] = true
		}
	}
	return onlyFirst, onlySecond, both
}

func (buffer *BufferIndex) GetDocCounter() int {
	return buffer.docCounter
}