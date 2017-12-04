package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

type part struct {
	category string
	content string
}
type document struct {
	id string
	partList []part
}

func parseDocuments (filePath string) []document {
	dat, err := ioutil.ReadFile(filePath)
	
	if err != nil {
		panic(err)
	}
	
	var output []document
	
	fileString := string(dat) // Change this to handle bigger files

	// Split the files in documents
	regexDoc := regexp.MustCompile(".I ([0-9]*)\n")
	docs := regexDoc.Split(fileString, -1)
	docsNum := regexDoc.FindAllStringSubmatch(fileString, -1)
	
	for i, doc := range docs {
		if doc != "" {
			// Create doc data structure
			partList := parseDocument(doc)	
			parsedDoc := document{docsNum[i-1][1], partList}
			output = append(output, parsedDoc)
		}
	}
	return output
}

func parseDocument(doc string) []part {
	var partList []part
	// Split the doc in parts
	regexDocPart := regexp.MustCompile("\\.([A-Z])\n")
	partsContent := regexDocPart.Split(doc, -1)
	partsName := regexDocPart.FindAllStringSubmatch(doc, -1)
	
	// Add part to doc data structure
	for j, partName := range partsName {
		partName := partName[1]
		partContent := partsContent[j+1]
		part := part{partName, partContent}
		partList = append(partList, part)
	}
	
	return partList
}

func main() {
	filePath := "./consignes/Data/CACM/cacm.all"
	parsedDoc := parseDocuments(filePath)
	fmt.Println("Done", parsedDoc[0])
}