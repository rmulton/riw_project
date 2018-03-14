package outputFormaters

import (
	"fmt"
	"sort"

	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/utils"
)

type sortDocsOutputFormater struct {
	docIDToPath map[int]string
}

type scoresToDocs map[float64][]int

// NewSortDocsOutputFormater returns a new SortDocsOutputFormater that format output
// by sorting the documents according to their scores
func NewSortDocsOutputFormater(docIDToPath map[int]string) *sortDocsOutputFormater {
	return &sortDocsOutputFormater{
		docIDToPath: docIDToPath,
	}
}

// Output print the output for the user
func (fmter *sortDocsOutputFormater) Output(res *indexes.PostingList) {
	var rank int
	scores, scoresToDocs := Sort(*res)
	fmt.Printf("\n> Found %d documents:\n\n", len(*res))
	if scoresToDocs != nil && scores != nil {
		for _, score := range scores {
			for _, docID := range scoresToDocs[score] {
				rank++
				// TODO: move to vectorized requests
				normalizedScore := score * 100
				if rank <= 3 {
					docPath := fmter.docIDToPath[docID]
					content := utils.FileToString(docPath)
					if len(content) > 400 {
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

// Sort sorts the documents according to their scores
func Sort(res indexes.PostingList) ([]float64, scoresToDocs) {
	scores := make([]float64, 0)
	scoresToDocs := make(scoresToDocs)

	for docID, score := range res {
		_, exists := scoresToDocs[score]
		if !exists {
			scoresToDocs[score] = make([]int, 0)
			scores = append(scores, score)
		}
		scoresToDocs[score] = append(scoresToDocs[score], docID)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(scores)))
	for _, scores := range scoresToDocs {
		sort.Ints(scores)
	}
	return scores, scoresToDocs

}
