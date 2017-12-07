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
	collection.ComputeIndex()
	// fmt.Print(collection.Index["written"])

	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("Index computed in %f seconds\n", elapsed.Seconds())

	request := requests.NewBinaryRequest("logic AND written AND ad", collection.Index)
	output := request.Compute()
	fmt.Print("Found: ", output)

}