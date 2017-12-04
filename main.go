package main

import (
	"fmt"
	"./reader/parser/cacm"
)

func main() {
	filePath := "./consignes/Data/CACM/"
	reversedIndex := cacm.ParseDocuments(filePath)
	for k, v := range reversedIndex {
		fmt.Println(k, ":", v, "\n\n")
	}
}