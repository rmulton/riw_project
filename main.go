package main

import (
	"sync"
	"./readers"
	"./indexBuilders"
	"log"
	"time"
	"bufio"
	"os"
	"fmt"
	"./requests"
)

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n>> Please enter request: ")
	text, _ := reader.ReadString('\n')
	return text
}

func main() {
	start := time.Now()
	var waitGroup sync.WaitGroup
	reader := readers.NewStanfordReader(
		"../consignes/Data/CS276/pa1-data/",
		10,
		&waitGroup,
	)
	builder := indexBuilders.NewInMemoryBuilder(
		reader.Docs,
		10,
		&waitGroup,
	)
	log.Println("Starting")
	waitGroup.Add(2)
	go reader.Read()
	go builder.Build()
	waitGroup.Wait()
	done := time.Now()
	elapsed := done.Sub(start)
	log.Printf("Done in %v", elapsed)
	index := builder.GetIndex()
	engine := requests.NewEngine(index)
	var req string
	for {
		req = readInput()
		engine.Request(req)
	}
	
}
