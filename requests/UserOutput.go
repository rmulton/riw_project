package requests

import (
	"fmt"
	// "sort"

)

type UserOutput struct {
	// collection Collection //TODO : add the collection to get a printable paragraph
	sortedDocIDs []int
}

func NewUserOutput(sortedDocIDs []int) UserOutput {
	// Sort doc ids
	userOutput := UserOutput{sortedDocIDs}
	return userOutput
}
	
func (userOutput *UserOutput) Print() {
	fmt.Println("Result:")
	for _, docID := range userOutput.sortedDocIDs {
		printDoc(docID)
	}
}

func printDoc(docID int) {
	fmt.Println(docID)
}