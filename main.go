package main

import (
	"sync"
	"./parsers"
	"./blocks"
	"log"
	"time"
)

func main() {
	start := time.Now()
	var waitGroup sync.WaitGroup
	reader := parsers.NewStanfordReader(
		"../consignes/Data/CS276/pa1-data/",
		10,
		&waitGroup,
	)
	blockFiller := blocks.NewFiller(
		35000,
		"./saved/",
		reader.Docs,
		10,
		&waitGroup,
	)
	log.Println("Starting")
	waitGroup.Add(2)
	go reader.Read()
	go blockFiller.Fill()
	waitGroup.Wait()
	done := time.Now()
	elapsed := done.Sub(start)
	log.Printf("Done in %v", elapsed)
}