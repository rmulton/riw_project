package inversers

import (
	"math"
	"../utils"
)

type PostingList map[int]float64
type VectorPostingList map[int][]float64

type toWrite struct {
	term string
	postingList PostingList
}

func (postingList PostingList) appendToTermFile(term string, writingChannel writingChannel) {
	writingChannel <- &toWrite{term, postingList}
}

func PostingListFromFile(path string) (error, PostingList) {
	postingList := make(PostingList)
	err := utils.ReadGob(path, &postingList)
	if err != nil {
		return err, nil
	}
	return nil, postingList
}

func (postingList PostingList) tfIdf(corpusSize int) {
	idf := float64(corpusSize)/float64(len(postingList)) // Inverse of the proportion of documents that contain the term
	for docID, frqc := range postingList {
		tf := frqc // Frequency of the term in the document
		postingList[docID] = (1 + math.Log(tf)) * math.Log(idf)
	}
}

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

func MergeToVector(postingLists []PostingList) VectorPostingList{
	vectorPostingList := make(VectorPostingList)
	for i, postingList := range postingLists {
		for docID, frqc := range postingList {
			_, exists := vectorPostingList[docID]
			if !exists {
				vectorPostingList[docID] = make([]float64, len(postingLists))
			}
			vectorPostingList[docID][i] = frqc
		}
	}
	return vectorPostingList
}

func (vecPostingList VectorPostingList) ToAngleTo(vector []float64) PostingList {
	output := make(PostingList)
	for docID, docVector := range vecPostingList {
		output[docID] = angle(docVector, vector)
	}
	return output
}

func angle(v1 []float64, v2 []float64) float64 {
	n1 := norm(v1)
	n2 := norm(v2)
	s := scalar(v1, v2)
	cos := s/(n1*n2)
	angle := math.Acos(cos)
	return angle
}

func norm(v []float64) float64 {
	var norm float64
	for _, coord := range v {
		norm += coord*coord
	}
	floatNorm := math.Sqrt(norm)
	return floatNorm
}

func scalar(v1 []float64, v2 []float64) float64 {
	if len(v1) != len(v2) {
		panic("Trying to compute scalar product of vectors of different sizes")
	}
	var output float64
	for i, score := range v1 {
		output += score*v2[i]
	}
	return output
}