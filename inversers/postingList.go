package inversers

import (
)

type postingList map[int]int

type toWrite struct {
	term string
	postingList postingList
}

func (postingList postingList) appendToTermFile(term string, writingChannel writingChannel) {
	writingChannel <- &toWrite{term, postingList}
}

func (postingList postingList) mergeWith(otherPostingList postingList) {
	for docID, frqc := range otherPostingList {
		_, exists := postingList[docID]
		if exists {
			postingList[docID] += frqc
		} else {
			postingList[docID] = frqc
		}
	}
}