package requests

import (
	"../parsers/stanford"
	"time"
	"fmt"
	"os"
	"../utils"
	"../indexes"
)

type StanfordEngine struct {
	collection *stanford.Collection
}
func (engine *StanfordEngine) LoadEngine(refresh bool) *indexes.ReversedIndex{
	// Start timer
	start := time.Now()
	indexFile := "./saved/stanford_index.gob"
	var collection = new(stanford.Collection)
	// Load the engine
	if _, err := os.Stat(indexFile); err == nil && !refresh {
		// If the collection is saved, load it
		err := utils.ReadGob(indexFile, collection)
		if err != nil {
			panic(err)
		}
	} else {
		// Otherwise compute it
		collection = stanford.NewCollection("./consignes/Data/CS276/pa1-data/")
		utils.WriteGob(indexFile, collection)
	}
	// Display loading time
	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("[Index computed for Stanford in %f seconds]\nReady to get queries !\n", elapsed.Seconds())
	return collection.Index
}

func (engine *StanfordEngine) Request(query string, index *indexes.ReversedIndex) UserOutput {
	binaryRequest := NewBinaryRequest(query, index)
	return NewUserOutput(binaryRequest.Output, binaryRequest.DocsScore) // Simplify to binaryRequest
}