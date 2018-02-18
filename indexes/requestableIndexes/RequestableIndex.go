package requestable

import (
	"github.com/rmulton/riw_project/indexes"
)


type RequestableIndex interface {
	GetPostingListsForTerms([]string) map[string]indexes.PostingList
	GetDocIDToPath() map[int]string
}