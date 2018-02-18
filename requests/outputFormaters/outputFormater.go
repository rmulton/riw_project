package outputFormaters

import (
	"github.com/rmulton/riw_project/indexes"
)

type OutputFormater interface {
	Output(*indexes.PostingList)
}