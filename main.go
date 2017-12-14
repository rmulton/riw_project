package main

import (
	"bufio"
	"fmt"
	"time"
	"os"
	"./parsers/cacm"
	"./requests"
	"./utils"
	"./parsers/stanford"
)

func analyseCacmCollection(path string) *cacm.Collection {
	start := time.Now()
	// Get a reversed dictionnary of relevant terms
	collection := cacm.NewCollection(path)
	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("[Index computed in %f seconds]\nReady to get queries !\n", elapsed.Seconds())
	return collection
}

func analyseStanfordCollection(path string) *stanford.Collection {
	start := time.Now()
	// Get a reversed dictionnary of relevant terms
	collection := stanford.NewCollection(path)
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

func handleInput(collection *stanford.Collection, query string) {
	start := time.Now()
	request := requests.NewBinaryRequest(query, collection.Index)
	userOutput := requests.NewUserOutput(request.Output, request.DocsScore)
	userOutput.Print()
	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("[Result computed in %f seconds]\n", elapsed.Seconds())
}
func main() {

	// collection := analyseStanfordCollection("./consignes/Data/CS276/pa1-data/")
	var collection = new(stanford.Collection)
	err := utils.ReadGob("./stanford_index.gob", collection)
	if err != nil {
		panic(err)
	}
	for {
		input := readInput()
		handleInput(collection, input)
	}
}