package requests

import (
	"regexp"
	"../indexes"
)

type BinaryRequest struct {
	input string
	index indexes.ReversedIndex
	ands [][]string
	Output map[int]int 
}

func NewBinaryRequest(input string, index indexes.ReversedIndex) BinaryRequest {
	request := BinaryRequest{input, index, [][]string{}, make(map[int]int)}
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
		request.ands = append(request.ands, terms)
	}
}

func (request *BinaryRequest) computeOutput() {
	// Find the output
	output := request.computeRequest()
	request.Output = output
}

func (request BinaryRequest) computeRequest() map[int]int {
	// Output variable
	var output = make(map[int]int)

	// For each and condition
	for _, andCondition:= range request.ands {
		// Use the first word as a reference
		referenceWord := andCondition[0]
		docsForWord := request.index[referenceWord]
		// Remove documents that don't contain the other words from the and condition
		// Add the frequency from other words
		for _, word := range andCondition[1:] {
			for docID, frqc := range docsForWord {
				_, exists := request.index[word][docID]
				if !exists {
					delete(docsForWord, docID)
				} else {
					docsForWord[docID] += frqc
				}
			}
		}

		// Add remaining words to output dict
		for docID, frqc := range docsForWord {
			output[docID] += frqc
		}
	}
	return output
}
