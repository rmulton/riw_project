package requests

import (
	"../inversers"
)

type outputFormater interface {
	output(*inversers.PostingList)
}