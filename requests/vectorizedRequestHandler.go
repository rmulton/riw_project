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
	reqVector := reqVectorFromTerms(term)	

	// Compute the angles between the request and the docs
	postingListForTermsSlice := values(postingListsForTerms)
	vectorizedPostingList := indexes.MergeToVector(postingListsForTermsSlice)
	docScores := vectorizedPostingList.ToAngleTo(reqVector)
	return &docScores
}

func reqVectorFromTerms(terms []string) map[string]float64{
	reqVector := make(map[string]float64)
	for _, term := range terms {
		reqVector[term] = 1.
	}
	return reqVector
}

func values(postingListsForTerms map[string]PostingList) {
	var postingLists []PostingList
	for _, postingList := range postingListsForTerms {
		postingLists = append(postingList, postingLists)
	}
	return postingLists
}