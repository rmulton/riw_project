package requests

import (
	"fmt"
	"strings"
	"../inversers"
	"../normalizers"
)

type binaryRequestHandler struct {

}

func NewBinaryRequestHandler() *binaryRequestHandler {
	return &binaryRequestHandler{}
}

func (reqHandler *binaryRequestHandler) request(request string, index *Index) *inversers.PostingList {
	conjClauses := strings.Split(request, "|")
	res := make(inversers.PostingList)
	for _, clause := range conjClauses {
		accurateDocs := requestAnd(clause, index)
		if accurateDocs != nil {
			res.MergeWith(*accurateDocs) // TODO: Check that there is no memory wasted
		}
	}
	return &res
}

func requestAnd(request string, index *Index) *inversers.PostingList{
	terms := *normalizers.Normalize(&request, &[]string{})
	err := index.LoadTerms(terms)
	if err != nil {
		fmt.Println(err)
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