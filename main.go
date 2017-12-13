package main

import (
	"fmt"
	"time"
	"./parsers/cacm"
	"./requests"
)

func main() {
	start := time.Now()

	// Get a reversed dictionnary of relevant terms
	collection := cacm.NewCollection("./consignes/Data/CACM/")

	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("Index computed in %f seconds\n", elapsed.Seconds())

	request := requests.NewBinaryRequest("computer AND science", collection.Index)
	fmt.Print("Found: ", request.Output)

}