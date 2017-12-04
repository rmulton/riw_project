package main

import (
	"fmt"
	"./reader/parser/cacm"
)

func main() {
	filePath := "./consignes/Data/CACM/"
	parsedDocs := cacm.ParseDocuments(filePath)
	fmt.Println("Done", parsedDocs[62])
}