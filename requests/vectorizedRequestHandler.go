package requests

import (
	"../inversers"
	"../normalizers"
	"fmt"
)

type vectorizedRequestHandler struct {

}


func NewVectorizedRequestHandler() *vectorizedRequestHandler {
	return &vectorizedRequestHandler{}
}

func (reqHandler *vectorizedRequestHandler) request(request string, index *Index) *inversers.PostingList {
	// Tokenize and normalize the request
	terms := *normalizers.Normalize(&request, &[]string{})
	err, postingLists := index.GetTerms(terms)
	if err != nil {
		fmt.Println("One of the terms is not in the collection")
		return nil
	}

	// Create the vector representing the request
	nTerms := len(terms)
	reqVector := make([]float64, nTerms)
	for i, _ := range reqVector {
		reqVector[i] = 1
	}

	// Compute the angles between the request and the docs
	vectorizedPostingList := inversers.MergeToVector(postingLists)
	docScores := vectorizedPostingList.ToAngleTo(reqVector)
	return &docScores
}