package indexBuilders

import (
	"sync"
	"fmt"
	"os"
	"log"
	"path/filepath"
	"io/ioutil"
	"../utils"
	"../indexes"
)

type OnDiskBuilder struct {
	index *BufferIndex
	parentWaitGroup *sync.WaitGroup
	waitGroup *sync.WaitGroup
	folder string
	bufferSize int
	readingChannel indexes.ReadingChannel
	writingChannel indexes.WritingChannel
	docCounter int
}


func NewOnDiskBuilder(bufferSize int, folder string, readingChannel indexes.ReadingChannel, routines int, parentWaitGroup *sync.WaitGroup) *OnDiskBuilder {
	var waitGroup sync.WaitGroup
	writingChannel := make(indexes.WritingChannel)
	index := NewBufferIndex(bufferSize, writingChannel)
	return &OnDiskBuilder{
		index: index,
		bufferSize: bufferSize,
		parentWaitGroup: parentWaitGroup,
		waitGroup: &waitGroup,
		folder: folder,
		readingChannel: readingChannel,
		writingChannel: writingChannel,
	}
}

// Build method starts one that handle reading the documents from the collection and one
// routine that handles writing data to the disk
func (builder *OnDiskBuilder) Build() {
	defer builder.parentWaitGroup.Done()
	builder.waitGroup.Add(2)
	go builder.readDocs()
	go builder.writePostingLists() // the block is copied to allow continuing operations on the builder
	builder.waitGroup.Wait()
	log.Printf("Done filling with %d documents", builder.docCounter)
}
 
func (builder *OnDiskBuilder) GetIndex() *indexes.OnDiskIndex {
	return indexes.OnDiskIndexFromFolder("./saved/")
}


func (builder *OnDiskBuilder) readDocs() {
	defer builder.waitGroup.Done()
	for doc := range builder.readingChannel {
		builder.docCounter++
		builder.addDoc(doc)
	}
	builder.finish()
	log.Printf("Done getting %d documents", builder.docCounter)
}

func (builder *OnDiskBuilder) writeTfIdfInMemoryTerms(inMemoryTerms map[string]bool) {
	defer builder.waitGroup.Done()
	// TODO: worker pool
	builder.index.toTfIdf(builder.docCounter)
	builder.index.writeAllPostingLists()

}

func (builder *OnDiskBuilder) tfIdfOnDiskTerms(onDiskTerms map[string]bool) {
	defer builder.waitGroup.Done()
	builder.fileToTfIdfForTerms(onDiskTerms)
}

func (builder *OnDiskBuilder) mergeDiskMemoryThenTfIdfTerms(toMergeThenTfIdf map[string]bool) {
	defer builder.waitGroup.Done()
	for term, _ := range toMergeThenTfIdf {
		// worker pool
		builder.mergeDiskMemoryThenTfIdfTerm(term)
	}
}

func (builder *OnDiskBuilder) mergeDiskMemoryThenTfIdfTerm(term string) {
	postingList := builder.index.getPostingListForTerm(term)
	postingListSoFar := builder.currentPostingListOnDisk(term)
	
	postingList.MergeWith(postingListSoFar)
	postingList.TfIdf(builder.docCounter)
	builder.index.appendToTermFile(postingList, term)
}


func (builder *OnDiskBuilder) getOnDiskTerms() map[string]bool {
	onDiskTerms := make(map[string]bool)
	files, err := ioutil.ReadDir("./saved/postings/")
    if err != nil {
        log.Fatal(err)
    }

    for _, f := range files {
			onDiskTerms[f.Name()] = true
    }
	return onDiskTerms
}

func (builder *OnDiskBuilder) getInMemoryTerms() map[string]bool {
	inMemoryTerms := make(map[string]bool)
	for term, _ := range builder.index.index.GetPostingLists() {
		inMemoryTerms[term] = true
	}
	return inMemoryTerms
}

func separate(first map[string]bool, second map[string]bool) (map[string]bool, map[string]bool, map[string]bool) {
	onlyFirst := make(map[string]bool)
	onlySecond := make(map[string]bool)
	both := make(map[string]bool)
	for firstKey, _ := range first {
		_, exists := second[firstKey]
		if exists {
			both[firstKey] = true
		} else {
			onlyFirst[firstKey] = true
		}
	}
	for secondKey, _ := range second {
		_, exists := both[secondKey]
		if !exists {
			onlySecond[secondKey] = true
		}
	}
	return onlyFirst, onlySecond, both
}
func (builder *OnDiskBuilder) categorizeTerms() (map[string]bool, map[string]bool, map[string]bool) {
	onDiskTerms := builder.getOnDiskTerms()
	inMemoryTerms := builder.getInMemoryTerms()
	onDiskOnly, inMemoryOnly, onDiskAndInMemory := separate(onDiskTerms, inMemoryTerms)
	return onDiskOnly, inMemoryOnly, onDiskAndInMemory
}
func (builder *OnDiskBuilder) finish() {
	// NB : use several folders
	// On disk terms and in memory terms
	// Compute the intersection
	// Split in folders ? No too long
	// Keep trace of the list
	onDiskOnly, inMemoryOnly, onDiskAndInMemory := builder.categorizeTerms()
	builder.waitGroup.Add(4)
	// Here we try to compute tf-idf scores in memory as much as possible
	// NB : use workers pool
	go builder.writeTfIdfInMemoryTerms(inMemoryOnly)
	// NB : use workers pool and walk for this one
	go builder.tfIdfOnDiskTerms(onDiskOnly)
	go builder.mergeDiskMemoryThenTfIdfTerms(onDiskAndInMemory)
	go builder.writeDocIDToFilePath("./saved/meta/idToPath")
}

func (builder *OnDiskBuilder) writeDocIDToFilePath(path string) {
	defer builder.waitGroup.Done()
	builder.index.writeDocIDToFilePath(path)
}

func (builder *OnDiskBuilder) fileToTfIdfForTerms(terms map[string]bool) {
	filepath.Walk("./saved/postings/", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && terms[info.Name()] == true {
			// log.Println(path)
			err, postingList := indexes.PostingListFromFile(path)
			if err != nil {
				return err
			}
			postingList.TfIdf(builder.docCounter) // TODO ? check that docCounter = len(os.listdir)
			utils.WriteGob(path, postingList)
		}
		return nil
	})
}
// Use only if you want to wait for everything to be written before computing tf-idf score
// func (builder *OnDiskBuilder) allFilesToTfIdf() {
// 	log.Println("Computing tf-idf scores from frequencies")
// 	filepath.Walk("./saved/postings/", func(path string, info os.FileInfo, err error) error {
// 		if !info.IsDir()  {
// 			err, postingList := indexes.PostingListFromFile(path)
// 			if err != nil {
// 				return err
// 			}
// 			// TODO: Parallelize
// 			postingList.TfIdf(builder.docCounter) // TODO ? check that docCounter = len(os.listdir)
// 			utils.WriteGob(path, postingList)
// 		}
// 		return nil
// 	})
// }

// Add a document to the current block
// NB: Might evolve to allow filling several blocks at the same time
func (builder *OnDiskBuilder) addDoc(doc indexes.Document) {

	// TODO: use routines to parallelize addDocToIndex and addDocToTerm

	// Add the path to the doc to the map
	builder.index.addDocToIndex(doc.ID, doc.Path)

	for _, term := range doc.NormalizedTokens {
		//	Add this document to the posting list of this term
		builder.index.addDocToTerm(doc.ID, term)
	}
}

func (builder *OnDiskBuilder) writePostingLists() {
	defer builder.waitGroup.Done()
	for toWrite := range builder.writingChannel {
		// Write it to the disk
		// TODO : duplicate code with currentPostingListOnDisk()
		termFile := fmt.Sprintf("./saved/postings/%s", toWrite.Term)
		// If the file exists append it
		if _, err := os.Stat(termFile); err == nil {
			postingListSoFar := builder.currentPostingListOnDisk(toWrite.Term)
			
			// Merge the current posting list
			postingListSoFar.MergeWith(toWrite.PostingList)
			// Write it to file
			utils.WriteGob(termFile, postingListSoFar)
		// Otherwise create it
		} else {
			utils.WriteGob(termFile, toWrite.PostingList)
		}
	}
}

func (builder *OnDiskBuilder) currentPostingListOnDisk(term string) indexes.PostingList {
	termFile := fmt.Sprintf("./saved/postings/%s", term)
	err, postingListSoFar := indexes.PostingListFromFile(termFile)
	if err != nil {
		panic(err)
	}
	return postingListSoFar
}