package requests

import (
	"../indexes"
)

type outputFormater interface {
	output(*indexes.PostingList)
}