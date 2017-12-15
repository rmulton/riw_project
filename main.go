package main

import (
	"bufio"
	"fmt"
	"os"
	"./requests"
	"flag"
	"./indexes"
)

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n>> Please enter request: ")
	text, _ := reader.ReadString('\n')
	return text
}

func main() {
	// Input engine type
	collectionType := flag.String("collection", "cacm", "Choose which collection to use. \"cacm\" and \"stanford\" are implemented")
	refresh := flag.Bool("refresh", false, "Choose whether to computer the index or load it from file if it already exists")
	flag.Parse()

	// Run the engine
	var index *indexes.ReversedIndex
	switch *collectionType {
	case "stanford":
		engine := requests.StanfordEngine{}
		index = engine.LoadEngine(*refresh)
		fmt.Printf("Found %v documents", index.CorpusSize)
		for {
			input := readInput()
			output := engine.Request(input, index)
			output.Print()
		}
	default:
		engine := requests.CacmEngine{}
		index = engine.LoadEngine(*refresh)
		fmt.Printf("Found %v documents", index.CorpusSize)
		for {
			input := readInput()
			output := engine.Request(input, index)
			output.Print()
		}
	}

}