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

	start = time.Now()
	request := requests.NewBinaryRequest("36-bit & memory & calculation | computer", collection.Index)
	fmt.Println(request.DocsScore)
	userOutput := requests.NewUserOutput(request.Output)
	userOutput.Print()
	done = time.Now()
	elapsed = done.Sub(start)
	fmt.Printf("Result computed in %f seconds\n", elapsed.Seconds())
}