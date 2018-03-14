package requestableIndexes

import (
	"github.com/rmulton/riw_project/indexes"
)

// RequestableIndex needs to be implemented to get an index that can be used to answer to the user's requests
type RequestableIndex interface {
	GetPostingListsForTerms([]string) map[string]indexes.PostingList
	GetDocIDToPath() map[int]string
	GetDocCounter() int
	IsInTheIndex(string) bool
}
