package indexes

import (
	"fmt"
	"log"
	"github.com/rmulton/riw_project/utils"
	"github.com/rmulton/riw_project/indexes"
)
type OnDiskIndex struct {
	folderPath string
	index *indexes.Index
}

// Only for indexes persisted to disk

func OnDiskIndexFromFolder(folderPath string) *OnDiskIndex {
	docIDToFilePath := make(map[int]string)
	err := utils.ReadGob("./saved/meta/idToPath", &docIDToFilePath)
	if err != nil {
		log.Println("Error while building the index from on disk files: %v", err)
	}
	index := indexes.NewEmptyIndex()
	return &OnDiskIndex{
		folderPath: folderPath,
		index: index,
	}
}

func (odi *OnDiskIndex) loadTerm(term string) error {
	termFile := fmt.Sprintf("./saved/postings/%s", term)
	_, exists := odi.index.GetPostingListForTerm(term)
	if !exists {
		err, postingList := indexes.PostingListFromFile(termFile)	
		if err != nil {
			return err
		}
		odi.index.SetPostingListForTerm(postingList, term)
	}
	return nil
}

// TODO: make it safe
func (odi *OnDiskIndex) unloadTerm(term string) {
	odi.index.ClearPostingListFor(term)
}

func (odi *OnDiskIndex) GetPostingListsForTerms(terms []string) map[string]PostingList {
	postingListsForTerms := make(map[string]PostingList)
	err := odi.LoadTerms(terms)
	if err != nil {
		panic(err)
	}
	for _, term := range terms {
		postingListsForTerms[term] = odi.index.postingLists[term]
	}
	return postingListsForTerms
}

func (odi *OnDiskIndex) LoadTerms(terms []string) error {
	for _, term := range terms {
		err := odi.loadTerm(term)
		if err != nil {
			fmt.Printf("%s is not in the index, it won't be taken into account\n", term)
		}
	}
	return nil
}

func (index *OnDiskIndex) GetDocIDToPath() map[int]string {
	return index.index.docIDToFilePath
}