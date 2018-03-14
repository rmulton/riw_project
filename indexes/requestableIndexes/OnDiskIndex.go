package requestableIndexes

import (
	"fmt"
	"log"

	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/utils"
)

// OnDiskIndex is a requestable index that is used when the index has been built on the disk
type OnDiskIndex struct {
	folderPath string
	index      *indexes.Index
}

// OnDiskIndexFromFolder returns an OnDiskIndex built from the data saved on the disk
func OnDiskIndexFromFolder(folderPath string) *OnDiskIndex {
	docIDToFilePath := make(map[int]string)
	err := utils.ReadGob("./saved/meta/idToPath", &docIDToFilePath)
	if err != nil {
		log.Printf("Error while building the index from on disk files: %v", err)
	}
	index := indexes.NewIndexWithDocIDToPath(docIDToFilePath)
	return &OnDiskIndex{
		folderPath: folderPath,
		index:      index,
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

// GetPostingListsForTerms returns the posting lists for the input terms
func (odi *OnDiskIndex) GetPostingListsForTerms(terms []string) map[string]indexes.PostingList {
	postingListsForTerms := make(map[string]indexes.PostingList)
	err := odi.LoadTerms(terms)
	if err != nil {
		log.Println(err)
	}
	for _, term := range terms {
		postingList, _ := odi.index.GetPostingListForTerm(term)
		postingListsForTerms[term] = postingList
	}
	return postingListsForTerms
}

// LoadTerms loads the posting lists for the input term from the disk
func (odi *OnDiskIndex) LoadTerms(terms []string) error {
	for _, term := range terms {
		err := odi.loadTerm(term)
		if err != nil {
			fmt.Printf("%s is not in the index, it won't be taken into account\n", term)
		}
	}
	return nil
}

// GetDocIDToPath returns the docID to filepath map
func (odi *OnDiskIndex) GetDocIDToPath() map[int]string {
	return odi.index.GetDocIDToFilePath()
}

// GetDocCounter returns the number of documents used to build the index
func (odi *OnDiskIndex) GetDocCounter() int {
	return odi.index.GetDocCounter()
}

// IsInTheIndex returns whether the term has a posting list in the index
func (odi *OnDiskIndex) IsInTheIndex(term string) bool {
	odi.LoadTerms([]string{term})
	return odi.index.IsInTheIndex(term)
}
