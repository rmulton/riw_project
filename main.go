package main

import (
	// "sync"
	// "./readers"
	// "./inversers"
	// "log"
	// "time"
	"./requests"
)

func main() {
	// start := time.Now()
	// var waitGroup sync.WaitGroup
	// reader := readers.NewStanfordReader(
	// 	"../consignes/Data/CS276/pa1-data/",
	// 	10,
	// 	&waitGroup,
	// )
	// filler := inversers.NewFiller(
	// 	35000,
	// 	"./saved/",
	// 	reader.Docs,
	// 	10,
	// 	&waitGroup,
	// )
	// log.Println("Starting")
	// waitGroup.Add(2)
	// go reader.Read(
	// go filler.Fill()
	// waitGroup.Wait()
	// done := time.Now()
	// elapsed := done.Sub(start)
	// log.Printf("Done in %v", elapsed)
	engine := requests.NewEngine("./saved/")
	engine.RequestAnd([]string{"stanford", "doctor", "medic", "mathemat"})
	
}