package indexes

import (
	"sync"
)

// SafeReversedIndex is used for threading
type SafeReversedIndex struct {
	Index *ReversedIndex // if it was a pointer, some other part of the code might get access to it
	mux sync.Mutex // Mutex was choosed instead of channels to avoid computation time (only one reading step and one writing step)
}

// Create an empty ReversedIndex to fill
func NewSafeReversedIndex() *SafeReversedIndex { // change it to Collection, an interface
	reversedIndex := NewReversedIndex()
	var mux sync.Mutex
	return &SafeReversedIndex{Index: reversedIndex, mux: mux}
}

// Finish applies transformations to get from docs-words mapping to docs[word] = score
func (safeIndex *SafeReversedIndex) Finish() {
	safeIndex.mux.Lock()
	defer safeIndex.mux.Unlock()
	safeIndex.Index.Finish()
}

func (safeIndex *SafeReversedIndex) AddParagraphForDoc(paragraph string, docID int) {
	safeIndex.mux.Lock()
	defer safeIndex.mux.Unlock()
	safeIndex.Index.AddParagraphForDoc(paragraph, docID)
}
