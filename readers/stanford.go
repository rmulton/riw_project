package readers

import (
	"os"
	"path/filepath"
	"sync"
	"log"
	"../normalizers"
	"../utils"
)

type StanfordReader struct {
	collectionPath string
	// WaitGroup for main program
	parentWaitGroup *sync.WaitGroup
	Docs chan *Document
	ReadCounter int
	Mux sync.Mutex
	sem chan bool
}

type Document struct {
	Path string
	Id int
	NormalizedTokens *[]string
}

func NewStanfordReader(collectionPath string, routines int, parentWaitGroup *sync.WaitGroup) *StanfordReader {
	var mux sync.Mutex
	docs := make(chan *Document)
	sem := make(chan bool, routines)
	reader := StanfordReader{
		collectionPath: collectionPath,
		parentWaitGroup: parentWaitGroup,
		Docs: docs,
		ReadCounter: 0,
		Mux: mux,
		sem: sem,
	}
	return &reader
}

// TODO: Change it to Read(path)
func (reader *StanfordReader) Read() {
	// Close output channel when done, tell the main program that the thread is done when returns
	defer reader.parentWaitGroup.Done()
	defer close(reader.Docs)
	// Walk over collection files and read them
	// TODO: Compare walk with other methods. For instance, using folders as a tree could speed up the process
	filepath.Walk(reader.collectionPath, func(path string, info os.FileInfo, err error) error {
		reader.sem <- true
		go reader.read(info, path) // use goroutines to maximize disk usage. TODO: check that too many routines won't slow the program
		return nil
	})
	// Wait that all files have been read
	for i := 0; i < cap(reader.sem); i++ {
		reader.sem <- true
	}
	log.Printf("Done reading %d files", reader.ReadCounter)
}

func (reader *StanfordReader) read(info os.FileInfo, path string) {
	// Tell to the reader that the thread is done when read() returns
	defer func() {<-reader.sem}()
	if !info.IsDir() {
		// Read and update counter used to get the document ID
		reader.Mux.Lock()
		counter := reader.ReadCounter
		reader.ReadCounter++
		reader.Mux.Unlock()
		// Get file content as a string
		stringFile := utils.FileToString(path)
		// Transform it to a list of normalized tokens
		normalizedTokens := normalizers.Normalize(stringFile, &[]string{})
		readDoc := &Document{
			Id: counter,
			Path: path,
			NormalizedTokens: normalizedTokens,
		}
		reader.Docs <- readDoc
	}
}