package requestHandlers

import (
	"strings"

	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/indexes/requestableIndexes"
	"github.com/rmulton/riw_project/normalizers"
)

type binaryRequestHandler struct {
	index requestableIndexes.RequestableIndex
}

// NewBinaryRequestHandler returns a new binary request handler
func NewBinaryRequestHandler(index requestableIndexes.RequestableIndex) *binaryRequestHandler {
	return &binaryRequestHandler{index}
}

// Request returns the response to a request
func (reqHandler *binaryRequestHandler) Request(request string, stopList []string) *indexes.PostingList {
	conjClauses := strings.Split(request, "|")
	res := make(indexes.PostingList)
	for _, clause := range conjClauses {
		accurateDocs := reqHandler.requestAnd(clause, stopList)
		if accurateDocs != nil {
			res.MergeWith(*accurateDocs) // TODO: Check that there is no memory wasted
		}
	}
	return &res
}

func (reqHandler *binaryRequestHandler) requestAnd(request string, stopList []string) *indexes.PostingList {
	terms := normalizers.Normalize(request, stopList)
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
