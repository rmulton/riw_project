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
	fmt.Printf("%v", res)
}