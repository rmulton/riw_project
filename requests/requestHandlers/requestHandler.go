package requestHandlers

import (
	"github.com/rmulton/riw_project/indexes"
)

// RequestHandler needs to be implemented in order to implement a new type of request
type RequestHandler interface {
	Request(string, []string) *indexes.PostingList
}
