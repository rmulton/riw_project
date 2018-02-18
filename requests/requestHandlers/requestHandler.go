package requests

import (
	"github.com/rmulton/riw_project/indexes"
)

type RequestHandler interface {
	Request(string) *indexes.PostingList
}