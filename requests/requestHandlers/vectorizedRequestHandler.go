package requestHandlers

import (
	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/indexes/requestableIndexes"
	"github.com/rmulton/riw_project/normalizers"
)

type vectorizedRequestHandler struct {
	index requestableIndexes.RequestableIndex
}


func NewVectorizedRequestHandler(index requestableIndexes.RequestableIndex) *vectorizedRequestHandler {
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
	reqVector := reqVectorFromTerms(terms)

	// Compute the angles between the request and the docs
	// postingListsForTermsSlice := values(postingListsForTerms)
	vectorizedPostingList := indexes.MergeToVector(postingListsForTerms)
	docScores := vectorizedPostingList.ToAnglesTo(reqVector)
	return &docScores
}

func reqVectorFromTerms(terms []string) map[string]float64{
	reqVector := make(map[string]float64)
	for _, term := range terms {
		reqVector[term] = 1.
	}
	return reqVector
}

func values(postingListsForTerms map[string]indexes.PostingList) []indexes.PostingList{
	var postingLists []indexes.PostingList
	for _, postingList := range postingListsForTerms {
		postingLists = append(postingLists, postingList)
	}
	return postingLists
}