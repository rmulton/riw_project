package requests

import (
	"fmt"
	"../inversers"
)

type Engine struct {
	index *Index
}

func NewEngine(folder string) *Engine {
	return &Engine{
		index: NewIndex(folder),
	}
}

func (engine *Engine) RequestAnd(request []string) {
	err := engine.index.LoadTerms(request)
	if err != nil {
		fmt.Print("One of the terms is not in the collection")
		return
	}

	var min int
	var bestTerm string
	for i, term := range request {
		postingListSize := len(engine.index.postingLists[term])
		if postingListSize <= min || i == 0 {
			min = postingListSize
			bestTerm = term
		}
	}

	// Copy the shortest posting list
	accurateDocs := make(inversers.PostingList)
	for k, v := range engine.index.postingLists[bestTerm] {
		accurateDocs[k] = v
	}

	fmt.Printf("%v", accurateDocs)

	// Add frequencies and remove inaccurate terms
	for docID, _ := range accurateDocs {
		for _, term := range request {
			if term != bestTerm {
				newFrqc, exists := engine.index.postingLists[term][docID]
				if !exists {
					delete(accurateDocs, docID)
				} else {
					accurateDocs[docID] += newFrqc
				}
			}
		}
	}

	fmt.Printf("%v", accurateDocs)

}