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
	flag.Parse()

	// Run the engine
	var index *indexes.ReversedIndex
	switch *collectionType {
	case "stanford":
		engine := requests.StanfordEngine{}
		index = engine.LoadEngine()
		for {
			input := readInput()
			output := engine.Request(input, index)
			output.Print()
		}
	default:
		engine := requests.CacmEngine{}
		index = engine.LoadEngine()
		for {
			input := readInput()
			output := engine.Request(input, index)
			output.Print()
		}
	}

}