package requestHandlers

import (
	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/indexes/requestableIndexes"
	"github.com/rmulton/riw_project/normalizers"
)

type vectorizedRequestHandler struct {
	index requestableIndexes.RequestableIndex
}

// NewVectorizedRequestHandler returns a new vectorized request handler
func NewVectorizedRequestHandler(index requestableIndexes.RequestableIndex) *vectorizedRequestHandler {
	return &vectorizedRequestHandler{index}
}

func (reqHandler *vectorizedRequestHandler) removeNotInTheIndexTerms(terms []string) []string {
	var inTheIndexTerms []string
	for _, term := range terms {
		if reqHandler.index.IsInTheIndex(term) {
			inTheIndexTerms = append(inTheIndexTerms, term)
		}
	}
	return inTheIndexTerms
}

// Request returns the response to a request
func (reqHandler *vectorizedRequestHandler) Request(request string, stopList []string) *indexes.PostingList {
	// Tokenize and normalize the request
	terms := normalizers.Normalize(request, stopList)
	terms = reqHandler.removeNotInTheIndexTerms(terms)
	postingListsForTerms := reqHandler.index.GetPostingListsForTerms(terms)

	// Create the vector representing the request
	reqVector := reqVectorFromTerms(terms)

	// Compute the angles between the request and the docs
	vectorizedPostingList := indexes.MergeToVector(postingListsForTerms)
	docScores := vectorizedPostingList.ToAnglesTo(reqVector)
	return &docScores
}

func reqVectorFromTerms(terms []string) map[string]float64 {
	reqVector := make(map[string]float64)
	for _, term := range terms {
		reqVector[term] = 1.
	}
	return reqVector
}

func values(postingListsForTerms map[string]indexes.PostingList) []indexes.PostingList {
	var postingLists []indexes.PostingList
	for _, postingList := range postingListsForTerms {
		postingLists = append(postingLists, postingList)
	}
	return postingLists
}
