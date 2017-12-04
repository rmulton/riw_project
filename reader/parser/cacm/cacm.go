package cacm

import (
	"io/ioutil"
	"regexp"
	"strings"
	"strconv"
)

// Change here the path to CACM folder
var folderPath = "./consignes/Data/CACM/"

type part struct {
	category string
	tokens []string
}

// Document : Structure to store a document
type Document struct {
	id int
	partList []part
}

// Read stop words file
var stopListFile = fileToString(folderPath + "common_words")
var stopList = strings.Split(stopListFile, "\n")
// Output variable
var reversedIndex = make(map[string][]int)

// ParseDocuments : Parse CACM documents
func ParseDocuments(folderPath string) map[string][]int {
	// Read documents file
	var dataFile = fileToString(folderPath + "cacm.all") // Change this to handle bigger files
	// Create output variable
	var documents []Document

	// Split the files in documents
	regexDoc := regexp.MustCompile(".I ([0-9]*)\n")
	docs := regexDoc.Split(dataFile, -1)
	docsNum := regexDoc.FindAllStringSubmatch(dataFile, -1)
	
	// Iterate over the documents and parse them
	for i, doc := range docs {
		if doc != "" { // TODO: Check how to avoid having an empty document
			// Create the document data structure
			docID, err := strconv.Atoi(docsNum[i-1][1])
			partList := parseDocument(doc, docID)	
			if err == nil {
				parsedDoc := Document{docID, partList}
				documents = append(documents, parsedDoc)
			}
		}
	}

	return reversedIndex
}

func parsePart(content string, docID int) []string {
	// Split content into tokens
	tokens := strings.FieldsFunc(content, func(r rune) bool {
		return r == ' ' || r == '.' || r == '\n' || r == ',' || r == '?' || r == '!' || r == '(' || r == ')' || r == '*' || r == ';' || r == '"' || r == '\'' || r == ':' || r == '{' || r == '}'
	})

	addSignificantTokens(tokens, docID)

	return tokens
}

func fileToString(filePath string) string {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fileString := string(dat)
	return fileString
}

func addSignificantTokens(tokens []string, docID int) {
	// Get significant words
	for _, token := range tokens {
		token = strings.ToLower(token)
		if isSignificant(token) {
			_, exists := reversedIndex[token]
			if exists {
				reversedIndex[token] = append(reversedIndex[token], docID)
			} else {
				reversedIndex[token] = []int{docID}
			}
		}
	}

}

func isSignificant(token string) bool {
	// TODO : inefficient
	for _, unsignificantWord := range stopList {
		if token == unsignificantWord {
			return false
		}
	}
	return true
}

func parseDocument(doc string, docID int) []part {
	// Create output variable
	var partList []part

	// Split the doc in parts
	regexDocPart := regexp.MustCompile("\\.([A-Z])\n")
	partsContent := regexDocPart.Split(doc, -1)
	partsName := regexDocPart.FindAllStringSubmatch(doc, -1)
	
	// Add part to doc data structure
	for j, partName := range partsName {
		partName := partName[1]
		if partName == "T" || partName == "W" || partName == "K" {
			partContent := partsContent[j+1]
			partTokens := parsePart(partContent, docID)
			part := part{partName, partTokens}
			partList = append(partList, part)
		}
	}
	
	return partList
}