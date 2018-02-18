package main

import (
	"sync"
	"github.com/rmulton/riw_project/readers"
	"github.com/rmulton/riw_project/indexBuilders"
	"github.com/rmulton/riw_project/indexBuilders/onDiskBuilders"
	"github.com/rmulton/riw_project/indexBuilders/inMemoryBuilders"
	"log"
	"time"
	"bufio"
	"os"
	"fmt"
	"flag"
	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/indexes/requestableIndexes"
	"github.com/rmulton/riw_project/requests"
	"github.com/rmulton/riw_project/utils"
)

func readInput(helpMessage string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(helpMessage)
	text, _ := reader.ReadString('\n')
	return text
}

func checkFolderExistsOrAskForAnother(folder string, helpMessage string) string {
	exists := utils.CheckPathExists(folder)
	for !exists {
		folder = readInput(helpMessage)
		exists = utils.CheckPathExists(folder)
	}
	return folder
}

func buildIndex(dataFolder string, collection string, inMemoryIndex bool) requestableIndexes.RequestableIndex {
	start := time.Now()
	// Check that the given dataFolder exists or ask for another one
	dataFolderDoesNotExistsMsg := fmt.Sprintf("The data folder %s does not exist.\n>> Please input an existing folder: ", dataFolder)
	dataFolder = checkFolderExistsOrAskForAnother(dataFolder, dataFolderDoesNotExistsMsg)
	// Erase the folder "./saved" or create it
	utils.ClearOrCreatePersistedIndex("./saved")
	// Create the reader
	var waitGroup sync.WaitGroup
	var reader readers.Reader
	readingChannel := make(indexes.ReadingChannel)
	if collection=="stanford" {
		reader = readers.NewStanfordReader(
			readingChannel,
			dataFolder,
			10,
			&waitGroup,
		)
		
	} else {
		reader = readers.NewCACMReader(
			readingChannel,
			dataFolder,
			5,
			&waitGroup,
			)
	}
	// Create the builder
	// Clean "./saved" folder
	var builder indexBuilders.IndexBuilder
	if inMemoryIndex {
		builder = inMemoryBuilders.NewInMemoryBuilder(
			readingChannel,
			10,
			&waitGroup,
		)
	} else {
		builder = onDiskBuilders.NewOnDiskBuilder(
			1000000000,
			readingChannel,
			10,
			&waitGroup,
			)
	}
	// Build the index
	fmt.Printf("Building an index for the %s collection using data from %s\n", collection, dataFolder)
	waitGroup.Add(2)
	go reader.Read()
	go builder.Build()
	waitGroup.Wait()
	// Display building time
	done := time.Now()
	elapsed := done.Sub(start)
	log.Printf("Done building the index in %v", elapsed)
	// Get the index
	return builder.GetIndex()
}

func loadIndexFromDisk() requestableIndexes.RequestableIndex {
	if utils.CheckPathExists("./saved/postings/") && utils.CheckPathExists("./saved/meta/idToPath") {
		index := requestableIndexes.OnDiskIndexFromFolder("./saved/")
		fmt.Println("Loaded existing index")
		return index
	} else {
		return nil	
	}
}

func getFlags() (string, bool, string, string, bool) {
	// Get the command line arguments
	fromScratchFlag := flag.Bool("from_scratch", false, "Use -from_scratch flag if you want the builder to erase the persisted index.")
	dataFolderFlag := flag.String("build_for", "", "Use -build flag to build or rebuild the index.")
	inMemoryIndexFlag := flag.Bool("index_in_memory", false, "Use -index_in_memory flag if you want to keep the index builder from using the hard disk. Don't use this option if your collection is too big")
	collectionFlag := flag.String("collection", "cacm", "Use -collection to choose which document collection to work on. For now, only \"stanford\" and \"cacm\" are working.")
	requestTypeFlag := flag.String("request", "vectorial", "Use -requestType to choose which kind of request you want to use on the index. For now, \"binary\" and \"vectorial\" are implemented")

	// Get the input from command line arguments
	flag.Parse()
	dataFolder := *dataFolderFlag
	inMemoryIndex := *inMemoryIndexFlag
	collection := *collectionFlag
	requestType := *requestTypeFlag
	fromScratch := *fromScratchFlag

	// Check the command line arguments
	// Check that the chosen collection's reader is implemented
	for !(collection=="stanford" || collection=="cacm") {
		collectionDoesNotExistMsg := fmt.Sprintf("The reader for the collection %s is not implemented. Implement it or choose one of the following collection:\n - stanford\n - cacm\n>> Please enter a collection: ", collection)
		collection = readInput(collectionDoesNotExistMsg)
	}

	// Check that the chosen requestType is implemented
	for !(requestType=="binary" || requestType=="vectorial") {
		requestTypeDoesNotExistMsg := fmt.Sprintf("The request type %s is not implemented. Implement it or choose one of the following request type:\n - binary\n - vectorial\n>> Please enter a request type: ", collection)
		requestType = readInput(requestTypeDoesNotExistMsg)
	}

	return dataFolder, inMemoryIndex, collection, requestType, fromScratch
}

func main() {

	// Get the command line interface arguments
	dataFolder, inMemoryIndex, collection, requestType, fromScratch := getFlags()

	// Creating the index
	var index requestableIndexes.RequestableIndex
	// If the from scratch option is set to true, erase everything in "./saved/meta" and "./saved/postings"
	if fromScratch && dataFolder==""{
		utils.ClearOrCreatePersistedIndex("./saved")
	}
	// If no dataFolder is given as a flag, try loading the index from "./saved"
	if dataFolder=="" {
		// Load the index
		index = loadIndexFromDisk()
		// If empty, a folder to build the index from is required
		if index==nil {
			for dataFolder=="" {
				dataFolder = readInput("No index can be found in memory. We need a folder to build an index from.\n>> Please input a path to the folder you want to build the index from: ")
			}
		}
	}

	// If a dataFolder is given as a flag, build the index for it
	// TODO: for cacm, return an error if the folder doesn't have the accurate format
	if dataFolder!="" {
		index = buildIndex(dataFolder, collection, inMemoryIndex)
	}
			
	// Load the request engine
	engine := requests.NewEngine(index, requestType, "sorted")

	// Respond requests
	var req string
	for {
		req = readInput("\n>> Please enter a request: ")
		engine.Request(req)
	}
}
