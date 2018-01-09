package requests

import (
	"../inversers"
)

type requestHandler interface {
	request(string, *Index) *inversers.PostingList
}