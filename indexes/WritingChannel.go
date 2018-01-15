package indexes

type BufferPostingList struct {
	Term string
	PostingList PostingList
}
type WritingChannel chan *BufferPostingList

func NewBufferPostingList(term string, postingList PostingList) *BufferPostingList {
	return &BufferPostingList{
		Term: term,
		PostingList: postingList,
	}
}
