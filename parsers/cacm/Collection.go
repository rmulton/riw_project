package cacm

import (
	"strings"
	"io/ioutil"
	"regexp"
	"strconv"
	"../../indexes"
	"github.com/kljensen/snowball"
)

// Collection
// Stores the path to the data folder, the index and the stop list.
// Initiate with empty Index and stopList, they are computed in ComputeIndex()
type Collection struct {
	path string
	Index *indexes.ReversedIndex
	stopList []string // Need to be stored to avoid computation
}

func NewCollection(dataFolderPath string) Collection {
	collection := Collection{dataFolderPath, indexes.NewReversedIndex(), []string{}}
	collection.computeIndex() // the index is stored in a collection object to avoid multiple function arguments
	return collection
}

func (collection *Collection) computeIndex() {
	// Read documents file
	var dataFile = collection.getData() // Change this to handle bigger files
	
	// Split the files in documents
	regexDoc := regexp.MustCompile("\\.I ([0-9]+)\n")
	docs := regexDoc.Split(dataFile, -1)
	docsNum := regexDoc.FindAllStringSubmatch(dataFile, -1)
	
	collection.computeIndexForDocs(docs, docsNum)
}

func (collection *Collection) computeIndexForDocs(docs []string, docsNum [][]string) {
	collection.setStopList()
	// Iterate over the documents and parse them
	for i, doc := range docs {
		if doc != "" { // TODO: Check how to avoid having an empty document
			// Create the document data structure
			docID, err := strconv.Atoi(docsNum[i-1][1])
			if err != nil {
				panic(err)
			}
			collection.computeIndexForDoc(doc, docID)	
		}
	}
	// When the index is done, get from linear frequency to log frequency
	collection.Index.FrqcToLogFrqc()
	
}

func (collection *Collection) computeIndexForDoc(doc string, docID int) {
	// Split the doc in parts
	regexDocPart := regexp.MustCompile("\\.([A-Z])\n")
	partsContent := regexDocPart.Split(doc, -1)
	partsName := regexDocPart.FindAllStringSubmatch(doc, -1)
	
	// Add part to doc data structure
	for j, partName := range partsName {
		partName := partName[1]
		if partName == "T" || partName == "W" || partName == "K" {
			partContent := partsContent[j+1]
			collection.computeIndexForPart(partContent, docID)
		}
	}
}

func (collection *Collection) computeIndexForPart(partContent string, docID int) {
	// Split content into tokens
	tokens := strings.FieldsFunc(partContent, func(r rune) bool {
		return r == ' ' || r == '.' || r == '\n' || r == ',' || r == '?' || r == '!' || r == '(' || r == ')' || r == '*' || r == ';' || r == '"' || r == '\'' || r == ':' || r == '{' || r == '}' || r == '/' || r == '|'
	})
	collection.addSignificantTokensToIndex(tokens, docID)
}

func (collection *Collection) addSignificantTokensToIndex(tokens []string, docID int) {
	// Copy the index
	index := collection.Index
	// Get significant words
	for _, token := range tokens {
		// done by the stemmatizer
		// token = strings.ToLower(token)
		stemmedToken, err := snowball.Stem(token, "english", true)
		if err != nil {
			panic(err)
		}
		if collection.isSignificant(stemmedToken) {
			// Add token to the list of keys if necessary
			_, exists := index.DocsForWords[stemmedToken]
			if !exists {
				tokenDict := make(map[int]float64)
				index.DocsForWords[stemmedToken] = tokenDict
				index.DocsForWords[stemmedToken][docID] = 0
			}
			index.DocsForWords[stemmedToken][docID]++
		}
	}
}

func (collection *Collection) isSignificant(token string) bool {
	// TODO : inefficient
	for _, unsignificantWord := range collection.stopList {
		if token == unsignificantWord {
			return false
		}
	}
	return true
}

func (collection *Collection) setStopList() {
	// Read stop words file
	collection.stopList = GetStopListFromFolder(collection.path)
}

func (collection *Collection) getData() string {
	return fileToString(collection.path + "cacm.all")
}

func fileToString(filePath string) string {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fileString := string(dat)
	return fileString
}
