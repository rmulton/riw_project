package inversers

import (
)

type PostingList map[int]int

type toWrite struct {
	term string
	postingList PostingList
}

func (postingList PostingList) appendToTermFile(term string, writingChannel writingChannel) {
	writingChannel <- &toWrite{term, postingList}
}

func (postingList PostingList) mergeWith(otherPostingList PostingList) {
	for docID, frqc := range otherPostingList {
		_, exists := postingList[docID]
		if exists {
			postingList[docID] += frqc
		} else {
			postingList[docID] = frqc
		}
	}
}