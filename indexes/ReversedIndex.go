package indexes

import (
	"fmt"
	"sort"
	"math"
	"../normalizers"
)



// DocsForWord stores the accurate documents for a word, together with their score.
type DocsForWord map[int]float64

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
//
type ReversedIndex struct {
	DocsForWords map[string]DocsForWord
	StopList []string
}

// Create an empty ReversedIndex to fill
func NewReversedIndex() *ReversedIndex { // change it to Collection, an interface
	docsForWords := make(map[string]DocsForWord)
	stopList := []string{}
	return &ReversedIndex{docsForWords, stopList}
}

func (index *ReversedIndex) Finish() {
	// Get score to logscore, then normalize it
	index.frqcToLogFrqc()
	index.normalizeDocsScore()
}

func (index *ReversedIndex) normalizeDocsScore() {
	// Get the sum of the scores for a word
	for word, docsScore := range index.DocsForWords {
		// Get the sum of the scores
		var wordSum float64 = 0
		for _, score := range docsScore {
			wordSum += score
		}
		// Normalize the scores
		for docID, score := range docsScore {
			index.DocsForWords[word][docID] = score/wordSum
		}
	}
}

// FrqcToLogFrqc transforms the linear frequency score to a log frequency score
// The linear frequency score for a document and a word is the numer of occurence of the word in the document
func (index *ReversedIndex) frqcToLogFrqc() {
	// Iterate over the index
	for word, docFrqcs := range index.DocsForWords {
		// Iterate over the documents/frqc map of the word
		for docID, frqc := range docFrqcs {
			index.DocsForWords[word][docID] = 1 + math.Log10(frqc)
		}

	}
}

func (index *ReversedIndex) AddTokensForDoc(tokens []string, docID int) {
	normalizedTokens := normalizers.NormalizeWords(tokens, index.StopList) // Check ""
	index.addFrequenciesForDoc(normalizedTokens, docID)
}

func (index *ReversedIndex) addFrequenciesForDoc(wordRoots []string, docID int) {
	for _, wordRoot := range wordRoots {
		if _, exists := index.DocsForWords[wordRoot]; !exists {
			index.DocsForWords[wordRoot] = make(DocsForWord)
		}
		index.DocsForWords[wordRoot][docID]++
	}
}

func (index *ReversedIndex) String() string {
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