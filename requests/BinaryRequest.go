package requests

import (
	"regexp"
	"../indexes"
)

type BinaryRequest struct {
	input string
	index indexes.ReversedIndex
	ands [][]string
}

func NewBinaryRequest(input string, index indexes.ReversedIndex) BinaryRequest {
	return BinaryRequest{input, index, [][]string{}}
}

func (request BinaryRequest) parse() BinaryRequest {
	andRegex := regexp.MustCompile("[a-z]+[ AND [a-z]+]*")	
	ands := andRegex.FindAllString(request.input, -1)
	for _, and := range ands {
		wordRegex := regexp.MustCompile("[a-z]+")
		terms := wordRegex.FindAllString(and, -1)
		request.ands = append(request.ands, terms)
	}
	return request
}

func (request BinaryRequest) Compute() map[int]int {
	// Parse the request
	request = request.parse()
	// Find the output
	output := request.computeRequest()
	return output
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
