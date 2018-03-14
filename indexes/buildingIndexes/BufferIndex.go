package buildingIndexes

import (
	"io/ioutil"
	"log"
	"sync"

	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/utils"
)

// BufferIndex is used to build an index using the hard disk to extend the possible size of the index
type BufferIndex struct {
	Mux            *sync.Mutex
	bufferSize     int
	currentSize    int
	docCounter     int
	writingChannel indexes.WritingChannel
	index          *indexes.Index
}

// NewBufferIndex returns a new BufferIndex
func NewBufferIndex(bufferSize int, writingChannel indexes.WritingChannel) *BufferIndex {
	var mux sync.Mutex
	index := indexes.NewEmptyIndex()
	return &BufferIndex{
		writingChannel: writingChannel,
		Mux:            &mux,
		bufferSize:     bufferSize,
		index:          index,
	}
}

// AddDocToTerm is used to fill the posting lists
func (buffer *BufferIndex) AddDocToTerm(docID int, term string) {
	buffer.index.AddDocToTerm(docID, term)
	buffer.currentSize++
}

// AddDocToIndex adds a new document in the index so that index keep trace of docID -> doc
func (buffer *BufferIndex) AddDocToIndex(docID int, docPath string) {
	buffer.docCounter++
	if buffer.currentSize >= buffer.bufferSize && buffer.bufferSize != -1 { // NB: This is a very important decision for the system. Using the biggest posting list might not be the best one.
		buffer.WriteBiggestPostingList()
	}
	buffer.index.AddDocToIndex(docID, docPath)
}

// WriteBiggestPostingList moves the biggest posting list in memory to the disk to free some memory
// /!\ MAYBE IT WOULD BE BETTER TO HAVE A ROUTINE WORKING ON WRITING THE POSTING LISTS
// AND USE RWLock instead of Lock
func (buffer *BufferIndex) WriteBiggestPostingList() {
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

	buffer.AppendToTermFile(longestPostingList, termWithLongestPostingList, false)
}

// TODO : avoid code repition by buildingIndexes buffer.appendPostingListOnDisk(term)

// AppendToTermFile handles sending the BufferPostingList to be written on the writing channel
func (buffer *BufferIndex) AppendToTermFile(postingList indexes.PostingList, term string, replace bool) {
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

// WritePostingListForTerms writes the posting lists of the input terms on the disk
func (buffer *BufferIndex) WritePostingListForTerms(terms map[string]bool) {
	for term, _ := range terms {
		postingList, exists := buffer.index.GetPostingListForTerm(term)
		if !exists {
			log.Printf("Trying to writing the posting list of %s that is not in the index", term)
		}
		buffer.AppendToTermFile(postingList, term, false)
	}
}

// WriteAllPostingLists writes all the posting lists in the index on the disk
func (buffer *BufferIndex) WriteAllPostingLists() {
	defer close(buffer.writingChannel)
	for term, postingList := range buffer.index.GetPostingLists() {
		buffer.AppendToTermFile(postingList, term, true)
	}
}

func (buffer *BufferIndex) toTfIdf(corpusSize int) {
	buffer.index.ToTfIdf(corpusSize)
}

// ToTfIdfTerms computes tf-idf scores from frequency scores
func (buffer *BufferIndex) ToTfIdfTerms(terms map[string]bool) {
	buffer.index.ToTfIdfTerms(buffer.docCounter, terms)
}

// WriteDocIDToFilePath writes the docID to filepath map on the disk
func (buffer *BufferIndex) WriteDocIDToFilePath(path string) {
	utils.WriteGob(path, buffer.index.GetDocIDToFilePath())
}

// GetPostingListForTerm returns the posting list of a term
func (buffer *BufferIndex) GetPostingListForTerm(term string) (indexes.PostingList, bool) {
	return buffer.index.GetPostingListForTerm(term)
}

/* Find out which terms are in memory, on disk or both */

// CategorizeTerms find out which terms have their posting lists on the disk, in memory or both
func (buffer *BufferIndex) CategorizeTerms() (map[string]bool, map[string]bool, map[string]bool) {
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

// GetDocCounter returns the number of documents that have been sent to the BufferIndex
func (buffer *BufferIndex) GetDocCounter() int {
	return buffer.docCounter
}
