package requests

import (
	"github.com/rmulton/riw_project/indexes"
)

type requestHandler interface {
	request(string) *indexes.PostingList
}