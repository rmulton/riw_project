package outputFormaters

import (
	"math"
	"fmt"
	"sort"
	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/utils"
)

type sortDocsOutputFormater struct {
	docIDToPath map[int]string
}

type scoresToDocs map[float64][]int

func NewSortDocsOutputFormater(docIDToPath map[int]string) sortDocsOutputFormater {
	return sortDocsOutputFormater{
		docIDToPath: docIDToPath,
	}
}

func (fmter *sortDocsOutputFormater) output(res *indexes.PostingList) {
	var rank int
	scores, scoresToDocs := fmter.sort(res)
	if scoresToDocs != nil && scores != nil {
		for _, score := range scores {
			for _, docID := range scoresToDocs[score] {
				rank++
				// TODO: move to vectorized requests
				normalizedScore := score/math.Acos(0)*100
				if rank <= 3 {
					docPath := fmter.docIDToPath[docID]
					content := utils.FileToString(docPath)
					if len(content) > 400{
						content = content[:400]
					} else if len(content) > 0 {
						content = content[:len(content)-1]
					} else {
						content = "" // TODO: implement for collections for wich documents are mixed up in a file (like cacm)
					}
					fmt.Printf("> %d | Doc %d | Score %f%%\n%s ...\n%s\n\n", rank, docID, normalizedScore, content, docPath)
				}
			}
		}
	} else {
		fmt.Println("The result is empty")
	}
}

func (fmter *sortDocsOutputFormater) sort(res *indexes.PostingList) ([]float64, scoresToDocs) {
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
	fmt.Printf("\nFound %d documents:\n\n", docCounter)
	return scores, scoresToDocs
	
}