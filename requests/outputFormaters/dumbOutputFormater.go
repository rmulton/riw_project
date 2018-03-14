package outputFormaters

import (
	"fmt"

	"github.com/rmulton/riw_project/indexes"
)

type dumbOutputFormater struct {
}

// NewDumbOutputFormater returns a new dumb output formater
func NewDumbOutputFormater() *dumbOutputFormater {
	return &dumbOutputFormater{}
}

// Output prints the results for the user
func (fmter *dumbOutputFormater) Output(res *indexes.PostingList) {
	if res != nil {
		fmt.Printf("%v\n", res)
		fmt.Printf("Length: %v\n", len(*res))
	} else {
		fmt.Println("The result is empty")
	}
}
