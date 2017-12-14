package requests

import (
	"sort"
	"regexp"
	"math"
	"../indexes"
	"../normalizers"
)

type BinaryRequest struct {
	input string
	index *indexes.ReversedIndex
	ands [][]string
	DocsScore map[int]float64
	Output []int
}

func NewBinaryRequest(input string, index *indexes.ReversedIndex) BinaryRequest {
	request := BinaryRequest{input, index, [][]string{}, map[int]float64{}, []int{}}
	request.parse()
	request.computeOutput()
	return request
}

func (request *BinaryRequest) parse() { // TODO : retirer stopwords, mots en double, 36-Bit = 36-bit
	andRegex := regexp.MustCompile("[0-9A-z\\-]+[ & [0-9A-z\\-]+]*")	
	ands := andRegex.FindAllString(request.input, -1)
	for _, and := range ands {
		wordRegex := regexp.MustCompile("[0-9A-z\\-]+")
		terms := wordRegex.FindAllString(and, -1)
		// Normalize words
		terms = normalizers.NormalizeWords(terms, request.index.StopList)
		request.ands = append(request.ands, terms)
	}
}

func (request *BinaryRequest) computeOutput() {
	// Find the output
	docsFrqc := request.computeRequest()
	request.DocsScore = docsFrqc
	request.Output = getSortedDocList(docsFrqc)
}

// Transform a dictionnary of docIDS:frqc to a docIDs list sorted on frqc
func getSortedDocList(docsFrqc map[int]float64) []int {
	// Get the reverse map
	reversedDocsFrqc := make(map[float64][]int)
	for k, v := range docsFrqc {
		reversedDocsFrqc[v] = append(reversedDocsFrqc[v], k)
	}
	// Sort the values in decreasing order
	var frqcs []float64
	for k, _ := range reversedDocsFrqc {
		frqcs = append(frqcs, k)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(frqcs)))
	// Get the output in a list
	var sortedDocIDs []int
	for _, frqc := range frqcs {
		// Append to output list
		for _, docID := range reversedDocsFrqc[frqc] {
			sortedDocIDs = append(sortedDocIDs, docID)
		}
	}
	return sortedDocIDs
}

func (request *BinaryRequest) computeRequest() map[int]float64 {
	// Output variable
	var output = make(map[int]float64)

	// For each and condition
	for _, andCondition:= range request.ands {
		// Use the first word as a reference
		referenceWord := andCondition[0]
		docsForWord := request.index.DocsForWords[referenceWord]
		// Remove documents that don't contain the other words from the and condition
		// Add the frequency from other words
		for _, word := range andCondition[1:] {
			for docID, frqc := range docsForWord {
				_, exists := request.index.DocsForWords[word][docID]
				if !exists {
					delete(docsForWord, docID)
				} else {
					docsForWord[docID] += frqc
				}
			}
		}
		// Output score is the max of all the docForWords[docID] score
		for docID, frqc := range docsForWord {
			output[docID] = math.Max(frqc, output[docID])
		}
	}


	return output
}
