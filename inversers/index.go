package inversers

import (
	"fmt"
	"../utils"
)

// Index and BufferIndex are different structure so that to be sure to avoid conflicts between the index used by the user and the one used to build it.
type Index struct {
	folderPath string
	corpusSize int
	postingLists map[string]postingList
	docIDToFilePath map[int]string
}

func newIndex (folderPath string) *Index {
	// corpusSize := 1
	docIDToFilePath := make(map[int]string)
	err := utils.ReadGob("./saved/IDToPath.meta", &docIDToFilePath)
	if err != nil {
		panic(err)
	}
	postingLists := make(map[string]postingList)
	return &Index{
		folderPath: folderPath,
		docIDToFilePath: docIDToFilePath,
		postingLists: postingLists,
	}
}

func (index *Index) loadTerm(term string) {
	termFile := fmt.Sprintf("./saved/%s.postings", term)
	postingList := make(postingList)
	err := utils.ReadGob(termFile, &postingList)
	if err != nil {
		panic(err)
	}
	index.postingLists[term] = postingList
}

// TODO: make it safe
func (index *Index) unloadTerm(term string) {
	index.postingLists[term] = nil
}