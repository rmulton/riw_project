package requests

import (
	"github.com/rmulton/riw_project/indexes"
)

type outputFormater interface {
	output(*indexes.PostingList)
}