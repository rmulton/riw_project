package requests

import (
	"fmt"
	"strings"
	"../inversers"
)

type andRequestHandler struct {

}

func NewAndRequestHandler() *andRequestHandler {
	return &andRequestHandler{}
}

func (reqHandler *andRequestHandler) request(request string, index *Index) *inversers.PostingList {
	terms := strings.Split(request, " ")
	err := index.LoadTerms(terms)
	if err != nil {
		fmt.Println("One of the terms is not in the collection")
		return nil
	}

	var min int
	var bestTerm string
	for i, term := range terms {
		postingListSize := len(index.postingLists[term])
		if postingListSize <= min || i == 0 {
			min = postingListSize
			bestTerm = term
		}
	}

	// Copy the shortest posting list
	accurateDocs := make(inversers.PostingList)
	for k, v := range index.postingLists[bestTerm] {
		accurateDocs[k] = v
	}

	fmt.Printf("%v", accurateDocs)

	// Add frequencies and remove inaccurate terms
	for docID, _ := range accurateDocs {
		for _, term := range terms {
			if term != bestTerm {
				newFrqc, exists := index.postingLists[term][docID]
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