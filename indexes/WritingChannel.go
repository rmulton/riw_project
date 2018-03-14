package indexes

// BufferPostingList is used to move posting list from the index in memory to the disk
type BufferPostingList struct {
	Term               string
	PostingList        PostingList
	ReplaceCurrentFile bool
}

// WritingChannel is a channel for BufferPostingList
type WritingChannel chan BufferPostingList

// NewBufferPostingList creates a BufferPostingList that is going to be appended to the current
// posting list file
func NewBufferPostingList(term string, postingList PostingList) BufferPostingList {
	return BufferPostingList{
		Term:               term,
		PostingList:        postingList,
		ReplaceCurrentFile: false,
	}
}

// NewReplacingBufferPostingList creates a BufferPostingList that is going to erase and replace
// the current posting list file
func NewReplacingBufferPostingList(term string, postingList PostingList) BufferPostingList {
	return BufferPostingList{
		Term:               term,
		PostingList:        postingList,
		ReplaceCurrentFile: true,
	}
}
