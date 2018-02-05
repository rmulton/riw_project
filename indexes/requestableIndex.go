package indexes

type RequestableIndex interface {
	GetPostingListsForTerms([]string) map[string]PostingList
	GetDocIDToPath() map[int]string
}