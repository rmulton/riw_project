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
var reversedIndex map[string]map[string]int

// ParseDocuments : Parse CACM documents
func ParseDocuments(folderPath string) []Document {
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
			partList := parseDocument(doc)	
			docID, err := strconv.Atoi(docsNum[i-1][1])
			if err == nil {
				parsedDoc := Document{docID, partList}
				documents = append(documents, parsedDoc)
			}
		}
	}

	return documents
}

func parsePart(content string) []string {
	// Create output variable
	var loweredTokens []string

	// Split content into tokens
	tokens := strings.FieldsFunc(content, func(r rune) bool {
		return r == ' ' || r == '.' || r == '\n' || r == ','
	})

	// Tokens to lower case
	for _, token := range tokens {
		lowered := strings.ToLower(token)
		loweredTokens = append(loweredTokens, lowered)
	}
	
	loweredTokens = removeUnsignificantTokens(loweredTokens)
	return loweredTokens
}

func fileToString(filePath string) string {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fileString := string(dat)
	return fileString
}

func removeUnsignificantTokens(tokens []string) []string {
	// Create output variable
	var significantTokens []string

	// Get significant words
	for _, token := range tokens {
		if isSignificant(token) {
			significantTokens = append(significantTokens, token)
		}
	}

	return significantTokens
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

func parseDocument(doc string) []part {
	// Create output variable
	var partList []part

	// Split the doc in parts
	regexDocPart := regexp.MustCompile("\\.([A-Z])\n")
	partsContent := regexDocPart.Split(doc, -1)
	partsName := regexDocPart.FindAllStringSubmatch(doc, -1)
	
	// Add part to doc data structure
	for j, partName := range partsName {
		partName := partName[1]
		partContent := partsContent[j+1]
		partTokens := parsePart(partContent)
		part := part{partName, partTokens}
		partList = append(partList, part)
	}
	
	return partList
}