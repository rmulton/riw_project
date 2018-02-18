package indexes

import (
	"math"
	"github.com/rmulton/riw_project/utils"
)

// PostingList is a docID -> score map
type PostingList map[int]float64

func PostingListFromFile(path string) (error, PostingList) {
	postingList := make(PostingList)
	err := utils.ReadGob(path, &postingList)
	if err != nil {
		return err, nil
	}
	return nil, postingList
}

func (postingList PostingList) TfIdf(corpusSize int) {
	idf := float64(corpusSize)/float64(len(postingList)) // Inverse of the proportion of documents that contain the term
	for docID, frqc := range postingList {
		tf := frqc // Frequency of the term in the document
		postingList[docID] = (1 + math.Log(tf)) * math.Log(idf)
	}
}

// MergeWith merges two posting list to a posting list. The score of the output is
// the addition of the scores from the two posting lists - if a document is not in
// one of the posting lists, we do as if the score was zero.
func (postingList PostingList) MergeWith(otherPostingList PostingList) {
	for docID, frqc := range otherPostingList {
		_, exists := postingList[docID]
		if exists {
			postingList[docID] += frqc
		} else {
			postingList[docID] = frqc
		}
	}
}