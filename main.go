package main

import (
	"bufio"
	"fmt"
	"time"
	"os"
	"./parsers/cacm"
	"./requests"
)

func analyseCollection(path string) *cacm.Collection {
	start := time.Now()
	// Get a reversed dictionnary of relevant terms
	collection := cacm.NewCollection(path)
	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("[Index computed in %f seconds]\nReady to get queries !\n", elapsed.Seconds())
	return collection
}
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n>> Please enter request: ")
	text, _ := reader.ReadString('\n')
	return text
}

func handleInput(collection *cacm.Collection, query string) {
	start := time.Now()
	request := requests.NewBinaryRequest(query, collection.Index)
	userOutput := requests.NewUserOutput(request.Output, request.DocsScore)
	userOutput.Print()
	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("[Result computed in %f seconds]\n", elapsed.Seconds())
}
func main() {

	collection := analyseCollection("./consignes/Data/CACM/")
	for {
		input := readInput()
		handleInput(collection, input)
	}

}