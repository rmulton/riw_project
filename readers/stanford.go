package readers

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/normalizers"
	"github.com/rmulton/riw_project/utils"
)

// StanfordReader is a reader for any collection contained in a single documents, in which a document is a file
type StanfordReader struct {
	collectionPath string
	// WaitGroup for main program
	parentWaitGroup *sync.WaitGroup
	Docs            indexes.ReadingChannel
	ReadCounter     int
	Mux             sync.Mutex
	sem             chan bool
	stopwords       []string
}

// NewStanfordReader returns a new StanfordReader
func NewStanfordReader(docs indexes.ReadingChannel, collectionPath string, routines int, parentWaitGroup *sync.WaitGroup) *StanfordReader {
	var mux sync.Mutex
	sem := make(chan bool, routines)
	stopwords := []string{
		"of",
		"the",
		"and",
		"to",
		"in",
		"for",
		"a",
		"on",
		"by",
		"this",
		"at",
		"is",
		"with",
		"s",
		"about",
		"from",
		"are",
		"us",
		"all",
		"be",
		"that",
		"it",
		"as",
		"or",
		"an",
		"you",
		"i",
		"your",
		"can",
		"will",
		"we",
		"how",
		"what",
		"where",
		"has",
		"have",
		"which",
		"if",
		"not",
		"e",
		"may",
		"our",
		"no",
		"here",
		"their",
		"do",
		"who",
		"it",
		"been",
		"but",
		"when",
		"some",
		"they",
		"there",
		"through",
		"take",
		"into",
		"well",
		"he",
		"she",
		"him",
		"her",
		"my",
		"such",
		"off",
		"then",
		"his",
		"via",
		"so",
		"am",
		"would",
		"without",
		"everything",
		"them",
		"were",
		"per",
		"be",
		"her",
		"which",
		"me",
		"much",
	}
	reader := StanfordReader{
		collectionPath:  collectionPath,
		parentWaitGroup: parentWaitGroup,
		Docs:            docs,
		ReadCounter:     0,
		Mux:             mux,
		sem:             sem,
		stopwords:       stopwords,
	}
	return &reader
}

// Read handles the procedure to read the collection
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
	defer func() { <-reader.sem }()
	if !info.IsDir() {
		// Read and update counter used to get the document ID
		reader.Mux.Lock()
		counter := reader.ReadCounter
		reader.ReadCounter++
		reader.Mux.Unlock()
		// Get file content as a string
		stringFile := utils.FileToString(path)
		// Transform it to a list of normalized tokens
		normalizedTokens := normalizers.Normalize(stringFile, reader.stopwords)
		readDoc := indexes.Document{
			ID:               counter,
			Path:             path,
			NormalizedTokens: normalizedTokens,
		}
		reader.Docs <- readDoc
	}
}
