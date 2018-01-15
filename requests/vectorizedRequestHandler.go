package requests

import (
	"../indexes"
	"../normalizers"
	// "fmt"
)

type vectorizedRequestHandler struct {
	index indexes.RequestableIndex
}


func NewVectorizedRequestHandler(index indexes.RequestableIndex) *vectorizedRequestHandler {
	return &vectorizedRequestHandler{index}
}

func (reqHandler *vectorizedRequestHandler) request(request string) *indexes.PostingList {
	// Tokenize and normalize the request
	terms := normalizers.Normalize(request, []string{})
	postingListsForTerms := reqHandler.index.GetPostingListsForTerms(terms)
	// err, postingLists := reqHandler.index.GetTerms(terms) // TODO: Move to the index
	// if err != nil {
		// fmt.Println("One of the terms is not in the collection")
		// return nil
	// }

	// Create the vector representing the request
	nTerms := len(terms)
	reqVector := make([]float64, nTerms)
	for i, _ := range reqVector {
		reqVector[i] = 1
	}

	// Compute the angles between the request and the docs
	vectorizedPostingList := indexes.MergeToVector(postingListsForTerms)
	docScores := vectorizedPostingList.ToAngleTo(reqVector)
	return &docScores
}