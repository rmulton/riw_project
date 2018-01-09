package requests

import (
	"fmt"
	"../inversers"
)

type dumbOutputFormater struct {

}

func NewDumbOutputFormater() *dumbOutputFormater {
	return &dumbOutputFormater{}
}

func (fmter *dumbOutputFormater) output(res *inversers.PostingList) {
	if res != nil {
		fmt.Println("%v", res)
	} else {
		fmt.Println("The result is empty")
	}
}