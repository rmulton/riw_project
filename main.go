package main

import (
	"fmt"
	"time"
	"./parsers/cacm"
	// "./indexes"
)

func main() {
	start := time.Now()

	// Get a reversed dictionnary of relevant terms
	collection := cacm.NewCollection("./consignes/Data/CACM/")
	collection.ComputeIndex()
	fmt.Print(collection.Index["written"])

	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("\n>> Done in %f seconds", elapsed.Seconds())

}