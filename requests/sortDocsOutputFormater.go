package requests

import (
	"math"
	"fmt"
	"sort"
	"../inversers"
)

type sortDocsOutputFormater struct {

}

type scoresToDocs map[float64][]int

func NewSortDocsOutputFormater() *sortDocsOutputFormater {
	return &sortDocsOutputFormater{}
}

func (fmter *sortDocsOutputFormater) output(res *inversers.PostingList) {
	var rank int
	scores, scoresToDocs := fmter.sort(res)
	if scoresToDocs != nil && scores != nil {
		for _, score := range scores {
			for _, docID := range scoresToDocs[score] {
				rank++
				normalizedScore := score/math.Acos(0)*100
				if rank <= 20 {
					fmt.Printf("%d: Doc %d with score %f %%\n", rank, docID, normalizedScore)
				}
			}
		}
	} else {
		fmt.Println("The result is empty")
	}
}

func (fmter *sortDocsOutputFormater) sort(res *inversers.PostingList) ([]float64, scoresToDocs) {
	scores := make([]float64, 0)
	scoresToDocs := make(scoresToDocs)
	var docCounter int

	for docID, score := range *res {
		docCounter ++
		_, exists := scoresToDocs[score]
		if !exists {
			scoresToDocs[score] = make([]int, 0)
			scores = append(scores, score)
		}
		scoresToDocs[score] = append(scoresToDocs[score], docID)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(scores)))
	fmt.Printf("Found %d documents:\n", docCounter)
	return scores, scoresToDocs
	
}