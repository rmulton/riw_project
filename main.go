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

func parseDocument (filePath string) []document {
	dat, err := ioutil.ReadFile(filePath)

	var output []document

	if err != nil {
		panic(err)
	}

	fileString := string(dat) // Change this to handle bigger files

	// Split the files in documents
	regexDoc := regexp.MustCompile(".I ([0-9]*)\n")
	docs := regexDoc.Split(fileString, -1)
	docsNum := regexDoc.FindAllStringSubmatch(fileString, -1)
	
	var partList []part
	for i, doc := range docs {
		if i>0 {
			// Split the doc in parts
			regexDocPart := regexp.MustCompile("\\.([A-Z])\n")
			partsContent := regexDocPart.Split(doc, -1)
			partsName := regexDocPart.FindAllStringSubmatch(doc, -1)
	
	
			// Add part to doc data structure
			fmt.Println("Document:", docsNum[i-1][1])
			for j, partName := range partsName {
				partName := partName[1]
				partContent := partsContent[j+1]
				part := part{partName, partContent}
				partList = append(partList, part)
			}
	
			// Create doc data structure
			parsedDoc := document{docsNum[i-1][1], partList}
			output = append(output, parsedDoc)
		}
	}
	return output
}

func main() {
	filePath := "./consignes/Data/CACM/cacm.all"
	parsedDoc := parseDocument(filePath)
	fmt.Println("Done", parsedDoc[0])
}