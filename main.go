package main

import (
	"fmt"
	"./parsers/cacm"
	"./indexes"
)

func main() {
	// Get a reversed dictionnary of relevant terms
	// reversedIndex := cacm.ParseDocuments()
	collection := cacm.Collection{"./consignes/Data/CACM/", make(indexes.ReversedIndex), []string{}}
	collection.ComputeIndex()
	fmt.Print(collection.Index)

}