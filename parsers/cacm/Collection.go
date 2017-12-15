package cacm

import (
	"fmt"
	"regexp"
	"strconv"
	"../../indexes"
	"../../utils"
)

// Collection stores the path to the data folder, the index and the stop list.
// Initiate with empty Index and stopList, they are computed in ComputeIndex()
type Collection struct {
	path string
	Index *indexes.ReversedIndex
}

// NewCollection create a collection structure and fill it with data from the dataFolderPath
func NewCollection(dataFolderPath string) *Collection {
	// Create an emtpy collection
	collection := Collection{dataFolderPath, indexes.NewReversedIndex()}
	// Fille the collection Index with the reversed index computed on the collection
	collection.computeIndex() // NB: the index is stored in a collection object to avoid multiple function arguments
	// When the index is done, get from linear frequency to log frequency
	collection.Index.Finish() // See if it is possible to move it to a ReversedIndex method
	fmt.Printf("Found %v documents", collection.Index.CorpusSize)

	return &collection
}

func (collection *Collection) computeIndex() {
	// Read the documents from the folder
	var dataFile = collection.getData() // TODO: Change this to handle bigger files
	// Split the files in documents
	regexDoc := regexp.MustCompile("\\.I ([0-9]+)\n")
	// Documents content
	docs := regexDoc.Split(dataFile, -1)
	// Documents ID. Important since there might be missing ids
	docsNum := regexDoc.FindAllStringSubmatch(dataFile, -1)
	// Fill the index with the data from the documents	
	collection.computeIndexForDocs(docs, docsNum)
}

func (collection *Collection) computeIndexForDocs(docs []string, docsNum [][]string) {
	// Iterate over the documents and parse them
	for i, doc := range docs {
		if doc != "" { // TODO: Check how to avoid having an empty document
			docID, err := strconv.Atoi(docsNum[i-1][1])
			if err != nil {
				panic(err)
			}
			// Fill the index with the data from this document
			collection.computeIndexForDoc(doc, docID)	
		}
	}
		
}

func (collection *Collection) computeIndexForDoc(doc string, docID int) {
	// Split the doc in parts
	regexDocPart := regexp.MustCompile("\\.([A-Z])\n")
	partsContent := regexDocPart.Split(doc, -1)
	partsName := regexDocPart.FindAllStringSubmatch(doc, -1)
	
	// Fill the index with the data from these parts
	for j, partName := range partsName {
		partName := partName[1]
		// Only use text from T, W and K parts
		if partName == "T" || partName == "W" || partName == "K" {
			partContent := partsContent[j+1]
			collection.computeIndexForPart(partContent, docID)
		}
	}

	collection.Index.CorpusSize++
}

func (collection *Collection) computeIndexForPart(partContent string, docID int) {
	// Split content into tokens
	collection.Index.AddParagraphForDoc(partContent, docID)
}

func (collection *Collection) getData() string {
	return utils.FileToString(collection.path + "cacm.all")
}


