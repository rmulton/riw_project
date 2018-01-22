package requests

import (
	"../indexes"
	"../normalizers"
)

type andRequestHandler struct {
	index indexes.RequestableIndex
}

func NewAndRequestHandler(index indexes.RequestableIndex) *andRequestHandler {
	return &andRequestHandler{index}
}

func (reqHandler *andRequestHandler) request(request string) *indexes.PostingList {
	terms := normalizers.Normalize(request, []string{})
	postingListsForTerms := reqHandler.index.GetPostingListsForTerms(terms)

	var min int
	var bestTerm string
	for i, term := range terms {
		postingListSize := len(postingListsForTerms[term])
		if postingListSize <= min || i == 0 {
			min = postingListSize
			bestTerm = term
		}
	}

	// Copy the shortest posting list
	accurateDocs := make(indexes.PostingList)
	for k, v := range postingListsForTerms[bestTerm] {
		accurateDocs[k] = v
	}

	// Add frequencies and remove inaccurate terms
	for docID, _ := range accurateDocs {
		for _, term := range terms {
			if term != bestTerm {
				newFrqc, exists := postingListsForTerms[term][docID]
				if !exists {
					delete(accurateDocs, docID)
				} else {
					accurateDocs[docID] += newFrqc
				}
			}
		}
	}

	return &accurateDocs
}