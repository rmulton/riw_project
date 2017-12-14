package indexes

import (
	"fmt"
	"sort"
	"math"
)

// ReversedIndex stores the index built by a parser.
// 
//		{
//				"toto": {
//					1234: 4,
//					23: 1,				// ReversedIndex structure
//					34: 10,
//					...
//				},
//		
//				...
//			}
type DocsScore map[int]float64

type ReversedIndex struct {
	DocsForWords map[string]DocsScore
	StopList []string
}

func NewReversedIndex() *ReversedIndex { // change it to Collection, an interface
	docsForWords := make(map[string]DocsScore)
	stopList := []string{}
	return &ReversedIndex{docsForWords, stopList}
}

func (index ReversedIndex) String() string {
	// Output variable
	var output string

	// Get all the keys and order them
	var keys []string
	for k := range index.DocsForWords {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Append key-value pairs to the output
	for _, key := range keys {
		termDict := index.DocsForWords[key]
		output += fmt.Sprintf("%s : %s\n", key, fmt.Sprint(termDict))
	}
		
	return output
}

func (index ReversedIndex) FrqcToLogFrqc() {
	// Iterate over the index
	for word, docFrqcs := range index.DocsForWords {
		// Iterate over the documents/frqc map of the word
		for docID, frqc := range docFrqcs {
			index.DocsForWords[word][docID] = 1 + math.Log10(frqc)
		}

	}

}