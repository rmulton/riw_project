package indexes

import (
	"math"
)
// TODO : Should we move elsewhere the following functions ?
// TODO : Should we juste have a function ToAngleScoresTo([]postingLists, request) ??

// MergeToVector merges two posting lists to a vectorized posting list.
// For an input of this kind :
// [
// 	"blabla" : [1: 3.],
// 	"bleble" : [1: 4.],
// ]
// The output looks like this :
// {
// 	1: ["blabla": 3., "bleble": 4.]
// }
func MergeToVector(postingLists map[string]PostingList) VectorPostingList{
	vectorPostingList := make(VectorPostingList)
	var i int
	for term, postingList := range postingLists {
		for docID, frqc := range postingList {
			_, exists := vectorPostingList[docID]
			if !exists {
				vectorPostingList[docID] = make(map[string]float64, len(postingLists))
			}
			vectorPostingList[docID][term] = frqc
		}
		i++
	}
	return vectorPostingList
}

// ToAnglesTo gets the angles between the documents contained in a VectorPostingList and a document.
// This is used to compute angles between a request and the documents containing at leat on of the
// term in the request, for now.
func (vecPostingList VectorPostingList) ToAnglesTo(vector map[string]float64) PostingList {
	output := make(PostingList)
	for docID, docVector := range vecPostingList {
		output[docID] = angle(docVector, vector)
	}
	return output
}

func angle(v1 map[string]float64, v2 map[string]float64) float64 {
	n1 := norm(v1)
	n2 := norm(v2)
	s := scalar(v1, v2)
	cos := s/(n1*n2)
	angle := math.Acos(cos)
	return angle
}

func norm(v map[string]float64) float64 {
	var norm float64
	for _, coord := range v {
		norm += coord*coord
	}
	floatNorm := math.Sqrt(norm)
	return floatNorm
}

func scalar(v1 map[string]float64, v2 map[string]float64) float64 {
	var output float64
	for term, score := range v1 {
		output += score*v2[term]
	}
	return output
}