package outputFormaters

import (
	"fmt"
	"github.com/rmulton/riw_project/indexes"
)

type dumbOutputFormater struct {

}

func NewDumbOutputFormater() dumbOutputFormater {
	return dumbOutputFormater{}
}

func (fmter *dumbOutputFormater) output(res *indexes.PostingList) {
	if res != nil {
		fmt.Printf("%v\n", res)
		fmt.Printf("Length: %v\n", len(*res))
	} else {
		fmt.Println("The result is empty")
	}
}