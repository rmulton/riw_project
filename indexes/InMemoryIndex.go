package indexes

import (
	"../utils"
)
type InMemoryIndex struct {
	index *Index
}
func InMemoryIndexFromIndex(index *Index) *InMemoryIndex {
	return &InMemoryIndex{
		index: index,
	}
}
func InMemoryIndexFromFile(filePath string) *InMemoryIndex {
	index := NewEmptyIndex()
	err := utils.ReadGob("./saved/index.gob", &index)
	if err != nil {
		panic(err)
	}
	return &InMemoryIndex{
		index: index,
	}
}

func (index *InMemoryIndex) GetPostingListsForTerms(terms []string) map[string]PostingList {
	postingListsForTerms := make(map[string]PostingList)
	for _, term := range terms {
		postingListsForTerms[term] = index.index.postingLists[term]
	}
	return postingListsForTerms
}
