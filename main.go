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
		"../riw_project/consignes/Data/CS276/pa1-data/",
		&waitGroup,
	)
	blockFiller := blocks.NewFiller(
		35000,
		"./saved/",
		reader.Docs,
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