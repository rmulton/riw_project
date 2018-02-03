package indexes

type BufferPostingList struct {
	Term string
	PostingList PostingList
	ReplaceCurrentFile bool
}
type WritingChannel chan BufferPostingList

// NewBufferPostingList creates a BufferPostingList that is going to be appended to the current
// posting list file
func NewBufferPostingList(term string, postingList PostingList) BufferPostingList {
	return BufferPostingList {
		Term: term,
		PostingList: postingList,
		ReplaceCurrentFile: false,
	}
}

// NewReplacingBufferPostingList creates a BufferPostingList that is going to erase and replace
// the current posting list file
func NewReplacingBufferPostingList(term string, postingList PostingList) BufferPostingList {
	return BufferPostingList {
		Term: term,
		PostingList: postingList,
		ReplaceCurrentFile: true,
	}
}
