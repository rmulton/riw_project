package requests

import (
	"fmt"
	// "sort"

)

type UserOutput struct {
	// collection Collection //TODO : add the collection to get a printable paragraph
	sortedDocIDs []int
	docsScore DocsScore
}

func NewUserOutput(sortedDocIDs []int, docsScore DocsScore) UserOutput {
	// Sort doc ids
	userOutput := UserOutput{sortedDocIDs, docsScore}
	return userOutput
}
	
func (userOutput *UserOutput) Print() {
	fmt.Println("Results:")
	for i, docID := range userOutput.sortedDocIDs {
		userOutput.printDoc(i, docID)
	}
}

func (userOutput *UserOutput) printDoc(ranking int, docID int) {
	score := userOutput.docsScore[docID]
	fmt.Printf("%d: Document %d with score %f\n", ranking+1, docID, score)
}