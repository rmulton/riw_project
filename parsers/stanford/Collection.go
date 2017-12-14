package stanford

import (
	"io/ioutil"
	"../../indexes"
	"../../utils"
	"log"
)

type Collection struct {
	path string
	Index *indexes.ReversedIndex
}

func NewCollection(dataFolderPath string) *Collection {
	collection := Collection{dataFolderPath, indexes.NewReversedIndex()}
	collection.computeIndex()
	collection.Index.Finish()
	utils.WriteGob("./stanford_index.gob", collection)

	return &collection
}

func (collection *Collection) computeIndex() {
	// List folders											// TODO: use functionnal style
	elements, err := ioutil.ReadDir(collection.path)
	if err != nil {
		panic(err)
	}

	// Loop over folders
	for j , element := range elements {
		if element.IsDir() {
			folderPath := collection.path + element.Name()
			collection.computeIndexForFolder(folderPath, j)
		}
	}
}

func (collection *Collection) computeIndexForFolder(folderpath string, folderID int) {
	log.Printf("Reading folder %v", folderpath)
	// List files
	elements, err := ioutil.ReadDir(folderpath)
	if err != nil {
		panic(err)
	}
	// Loop over files
	for i, element := range elements {
		if !element.IsDir() {
			filePath := folderpath + "/" + element.Name()
			fileID := i*10 + folderID
			collection.computeIndexForFile(filePath, fileID)
		}
	}
}

func (collection *Collection) computeIndexForFile(path string, fileID int) {
	// Read doc
	stringFile := utils.FileToString(path)
	// Add to index
	collection.Index.AddParagraphForDoc(stringFile, fileID) // Warning: id must use both folder and files, otherwise there are some conflicts
}