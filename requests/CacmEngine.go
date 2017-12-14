package requests

import (
	"../parsers/cacm"
	"time"
	"fmt"
	"os"
	"../utils"
	"../indexes"
)

type CacmEngine struct {
}

func (engine *CacmEngine) LoadEngine() *indexes.ReversedIndex {
	// Start timer
	start := time.Now()
	indexFile := "./saved/cacm_index.gob"
	var collection = new(cacm.Collection)
	// Load the engine
	if _, err := os.Stat(indexFile); err == nil {
		// If the collection is saved, load it
		err := utils.ReadGob(indexFile, collection)
		if err != nil {
			panic(err)
		}
	} else {
		// Otherwise compute it
		collection = cacm.NewCollection("./consignes/Data/CACM/")
		utils.WriteGob(indexFile, collection)
	}
	// Display loading time
	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("[Index computed for CACM in %f seconds]\nReady to get queries !\n", elapsed.Seconds())
	return collection.Index
}

func (engine *CacmEngine) Request(query string, index *indexes.ReversedIndex) UserOutput {
	// Compute the request
	binaryRequest := NewBinaryRequest(query, index)
	// Return the output
	return NewUserOutput(binaryRequest.Output, binaryRequest.DocsScore) // Simplify to binaryRequest
}