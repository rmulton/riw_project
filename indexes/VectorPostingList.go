package indexes

import (
	"math"
)

// VectorPostingList is a docID -> (term -> score) map
type VectorPostingList map[int]map[string]float64

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
func MergeToVector(postingLists map[string]PostingList) VectorPostingList {
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

// ToScore gets the scores from the angles between the documents contained in a VectorPostingList and a document.
func (vecPostingList VectorPostingList) ToScore(vector map[string]float64) PostingList {
	output := make(PostingList)
	for docID, docVector := range vecPostingList {
		angle := angle(docVector, vector)
		percentage := (math.Acos(0) - angle) / math.Acos(0)
		output[docID] = percentage
	}
	return output
}

func angle(v1 map[string]float64, v2 map[string]float64) float64 {
	n1 := norm(v1)
	n2 := norm(v2)
	s := scalar(v1, v2)
	cos := s / (n1 * n2)
	angle := math.Acos(cos)
	return angle
}

func norm(v map[string]float64) float64 {
	var norm float64
	for _, coord := range v {
		norm += coord * coord
	}
	floatNorm := math.Sqrt(norm)
	return floatNorm
}

func scalar(v1 map[string]float64, v2 map[string]float64) float64 {
	var output float64
	for term, score := range v1 {
		output += score * v2[term]
	}
	return output
}
