package requests

import (
	"../indexes"
)

type requestHandler interface {
	request(string) *indexes.PostingList
}