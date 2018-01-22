package readers

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"log"
	"strings"
	"../normalizers"
	"../utils"
	"../indexes"
)

type CACMReader struct {
	collectionPath string
	// WaitGroup for main program
	parentWaitGroup *sync.WaitGroup
	Docs chan *indexes.Document
	ReadCounter int
	Mux sync.Mutex
	sem chan bool
	stopList []string
}

func NewCACMReader(collectionPath string, routines int, parentWaitGroup *sync.WaitGroup) *CACMReader {
	var mux sync.Mutex
	docs := make(chan *indexes.Document)
	sem := make(chan bool, routines)
	// Stop list
	stopListFile := utils.FileToString(collectionPath + "/common_words")
	stopList := strings.Split(stopListFile, "\n")
	reader := CACMReader{
		collectionPath: collectionPath,
		parentWaitGroup: parentWaitGroup,
		Docs: docs,
		ReadCounter: 0,
		Mux: mux,
		sem: sem,
		stopList: stopList,
	}
	return &reader
}

// TODO: Change it to Read(path)
func (reader *CACMReader) Read() {
	log.Print("Reading CACM collection")
	// Close output channel when done, tell the main program that the thread is done when returns
	defer reader.parentWaitGroup.Done()
	defer close(reader.Docs)
	stringFile := utils.FileToString(reader.collectionPath + "/cacm.all")
	regexDoc := regexp.MustCompile("\\.I ([0-9]+)\n")
	// Documents ID. Important since there might be missing ids
	docsNum := regexDoc.FindAllStringSubmatch(stringFile, -1)
	documents := regexDoc.Split(stringFile, -1)
	for i, strID := range docsNum {
		reader.sem <- true
		go reader.read(i, strID[1], documents[i])
	}
	// Wait that all files have been read
	for i := 0; i < cap(reader.sem); i++ {
		reader.sem <- true
	}
	log.Printf("Done reading the %d documents", len(docsNum))
}

func (reader *CACMReader) read(counter int, ID string, document string) {
	// Tell to the reader that the thread is done when read() returns
	defer func() {<-reader.sem}()
	// docID might be a source of error, for instance if there are two docs with the same docID
	// this is why we are using the counter as a unique identifier
	docID, err := strconv.Atoi(ID)
	if err != nil {
		panic(err)
	}
	// Get the usefull parts of the doc
	var docContent string
	regexDocPart := regexp.MustCompile("\\.([A-Z])\n")
	partsContent := regexDocPart.Split(document, -1)
	partsName := regexDocPart.FindAllStringSubmatch(document, -1)
	for j, partName := range partsName {
		partName := partName[1]
		// Only use text from T, W and K parts
		if partName == "T" || partName == "W" || partName == "K" {
			partContent := partsContent[j+1]
			docContent += partContent
		}
	}
	// Transform it to a list of normalized tokens
	normalizedTokens := normalizers.Normalize(docContent, reader.stopList)
	documentPath := fmt.Sprintf("%s/cacm.all#%d", reader.collectionPath, docID)
	readDoc := &indexes.Document{
		ID: counter,
		Path: documentPath,
		NormalizedTokens: normalizedTokens,
	}
	reader.Docs <- readDoc
}