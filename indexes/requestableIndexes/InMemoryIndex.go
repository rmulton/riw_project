package requestableIndexes

import (
	"github.com/rmulton/riw_project/utils"
	"github.com/rmulton/riw_project/indexes"
	"log"
)
type InMemoryIndex struct {
	index *indexes.Index
}
func InMemoryIndexFromIndex(index *indexes.Index) *InMemoryIndex {
	return &InMemoryIndex{
		index: index,
	}
}
func InMemoryIndexFromFile(filePath string) *InMemoryIndex {
	index := indexes.NewEmptyIndex()
	err := utils.ReadGob("./saved/index.gob", &index)
	if err != nil {
		log.Println(err)
	}
	return &InMemoryIndex{
		index: index,
	}
}

func (index *InMemoryIndex) GetPostingListsForTerms(terms []string) map[string]indexes.PostingList {
	postingListsForTerms := make(map[string]indexes.PostingList)
	for _, term := range terms {
		postingList, _ := index.index.GetPostingListForTerm(term)
		postingListsForTerms[term] = postingList
	}
	return postingListsForTerms
}

func (index *InMemoryIndex) GetDocIDToPath() map[int]string {
	return index.index.GetDocIDToFilePath()
}