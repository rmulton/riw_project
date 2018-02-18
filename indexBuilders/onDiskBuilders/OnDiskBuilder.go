package onDiskBuilders

// TODO: Add folders handler
import (
	"github.com/rmulton/riw_project/indexes/buildingIndexes"
	"github.com/rmulton/riw_project/indexes/requestableIndexes"
	"github.com/rmulton/riw_project/indexBuilders"
	"sync"
	"os"
	"log"
	"path/filepath"
	"github.com/rmulton/riw_project/indexes"
)

// OnDiskBuilder handles
//    - getting documents parsed by the readers
//    - feeding a buffer index with it
//    - getting the index from frequency scores to tf-idf scores
type OnDiskBuilder struct {
	index *buildingIndexes.BufferIndex
	parentWaitGroup *sync.WaitGroup
	bufferSize int
	readingChannel indexes.ReadingChannel
	writingChannel indexes.WritingChannel
}

// NewOnDiskBuilder creates a OnDiskBuilder
func NewOnDiskBuilder(bufferSize int, readingChannel indexes.ReadingChannel, routines int, parentWaitGroup *sync.WaitGroup) *OnDiskBuilder {
	writingChannel := make(indexes.WritingChannel)
	index := buildingIndexes.NewBufferIndex(bufferSize, writingChannel)
	return &OnDiskBuilder{
		index: index,
		bufferSize: bufferSize,
		parentWaitGroup: parentWaitGroup,
		readingChannel: readingChannel,
		writingChannel: writingChannel,
	}
}

// Build method starts one routine that handle reading the documents from the collection and one
// routine that handles writing data to the disk
func (builder *OnDiskBuilder) Build() {
	defer builder.parentWaitGroup.Done()
	// Fill the index with the documents the reader sends
	var wg sync.WaitGroup
	wg.Add(1)
	go indexBuilders.FillIndex(builder.index, builder.readingChannel, &wg)
	// Start the disk writer
	go indexBuilders.WritePostingLists(builder.writingChannel, &wg) // the block is copied to allow continuing operations on the builder
	// Wait for the index filling to be done before finishing the index
	wg.Wait()
	wg.Add(1)
	// Finish the index
	builder.finish()
	// Wait for the disk writing to be done
	wg.Wait()
	log.Printf("Done filling with %d documents", builder.index.GetDocCounter())
}
 
// GetIndex returns a OnDiskIndexFromFolder that uses the posting lists from ./saved to respond queries
func (builder *OnDiskBuilder) GetIndex() requestableIndexes.RequestableIndex {
	return requestableIndexes.OnDiskIndexFromFolder("./saved/")
}

// finish handles:
//    - getting the index from frequency scores to tf-idf scores
//    - writing the docID to path map
func (builder *OnDiskBuilder) finish() {
	// Write the docID to path map
	go builder.index.WriteDocIDToFilePath("./saved/meta/idToPath")
	// Find out which terms are in memory, on disk or both
	onDiskOnly, inMemoryOnly, onDiskAndInMemory := builder.index.CategorizeTerms()
	// Get from frequency scores to tf-idf
	var wg sync.WaitGroup
	wg.Add(3)
	go builder.writeTfIdfInMemoryTerms(inMemoryOnly, &wg)
	go builder.tfIdfOnDiskTerms(onDiskOnly, &wg)
	go builder.mergeDiskMemoryThenTfIdfTerms(onDiskAndInMemory, &wg)
	wg.Wait()
	close(builder.writingChannel)
}

// writeTfIdfInMemoryTerms handles:
//    - getting from frequency scores to tf-idf scores for terms that are only in memory
//    - writing their posting lists
func (builder *OnDiskBuilder) writeTfIdfInMemoryTerms(inMemoryTerms map[string]bool, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	builder.index.ToTfIdfTerms(inMemoryTerms)
	builder.index.WritePostingListForTerms(inMemoryTerms)
}

// tfIdfOnDiskTerms handles getting from frequency scores to tf-idf scores for terms that are only on the disk
func (builder *OnDiskBuilder) tfIdfOnDiskTerms(onDiskTerms map[string]bool, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	builder.fileToTfIdfForTerms(onDiskTerms)
}

func (builder *OnDiskBuilder) fileToTfIdfForTerms(terms map[string]bool) {
	filepath.Walk("./saved/postings/", func(path string, info os.FileInfo, err error) error { // NB: according to the doc, Walk might not be the most efficient option
		term := info.Name()
		if !info.IsDir() && terms[term] == true {
			go builder.fileToTfIdfForTerm(term, path)
		}
		return nil
	})
}

func (builder *OnDiskBuilder) fileToTfIdfForTerm(term string, path string) {
	err, postingList := indexes.PostingListFromFile(path)
	if err != nil {
		log.Println(err)
	}
	postingList.TfIdf(builder.index.GetDocCounter()) // TODO ? check that docCounter = len(os.listdir)
	builder.writingChannel <- indexes.NewReplacingBufferPostingList(term, postingList)
}

// mergeDiskMemoryThenTfIdfTerms handles:
//    - merging scores from disk and memory 
//    - getting from frequency scores to tf-idf
//    - writing the resulting posting list to the disk
func (builder *OnDiskBuilder) mergeDiskMemoryThenTfIdfTerms(toMergeThenTfIdf map[string]bool, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for term := range toMergeThenTfIdf {
		builder.mergeDiskMemoryThenTfIdfTerm(term)
	}
}

func (builder *OnDiskBuilder) mergeDiskMemoryThenTfIdfTerm(term string) {
	postingList, _ := builder.index.GetPostingListForTerm(term)
	postingListSoFar := indexBuilders.CurrentPostingListOnDisk(term)
	
	// Here it is faster to load the persisted scores then get to tf-idf score rather than
	// appending the score in memory to the file then use the functionnality to tf-idf a file
	// even though it would use existing functionnalities. Maybe it would be wise to have a
	// Finisher struct that handle finishing indexes and a Writer that handles writing for the
	// next version.
	postingList.MergeWith(postingListSoFar)
	postingList.TfIdf(builder.index.GetDocCounter())
	builder.index.AppendToTermFile(postingList, term, true)
}