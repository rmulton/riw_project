package stanford

import (
	"io/ioutil"
	"sync"
	"log"
	"../../indexes"
	"../../utils"
)

type ThreadedCollection struct {
	path string
	SafeIndex *indexes.SafeReversedIndex
}

func NewThreadedCollection(dataFolderPath string) *ThreadedCollection{
	threadedCollection := ThreadedCollection{
		path: dataFolderPath, 
		SafeIndex: indexes.NewSafeReversedIndex(),
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	threadedCollection.computeIndex(&waitGroup)
	// Wait that everything is done
	waitGroup.Wait()
	log.Println("Done parsing")
	threadedCollection.SafeIndex.Finish()
	return &threadedCollection
}

func (threadedCollection *ThreadedCollection) computeIndex(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	// List folders
	elements, err := ioutil.ReadDir(threadedCollection.path)
	if err != nil {
		panic(err)
	}
	
	var wg sync.WaitGroup
	// Loop over folders
	for j , element := range elements {
		if element.IsDir() { // TODO: use a path list instead of checking the conditions the same way
			folderPath := threadedCollection.path + element.Name()
			wg.Add(1)
			go threadedCollection.computeIndexForFolder(folderPath, j, &wg)
		}
	}
	wg.Wait()
}

func (threadedCollection *ThreadedCollection) computeIndexForFolder(folderpath string, folderID int, wg *sync.WaitGroup) {
	// Update waiting group
	defer wg.Done()
	// Log reading process
	log.Printf("Started reading folder %v", folderID)
	// List files
	elements, err := ioutil.ReadDir(folderpath)
	if err != nil {
		panic(err)
	}
	// Loop over files
	for i, element := range elements {
		if !element.IsDir() {
			wg.Add(1)
			filePath := folderpath + "/" + element.Name()
			fileID := i*10 + folderID
			threadedCollection.computeIndexForFile(filePath, fileID)
			wg.Done()
		}
	}
	log.Printf("Done parsing %d", folderID)
}

func (threadedCollection *ThreadedCollection) computeIndexForFile(path string, fileID int) {
	// Read doc
	stringFile := utils.FileToString(path)
	// Add to index
	threadedCollection.SafeIndex.AddParagraphForDoc(stringFile, fileID) // Warning: id must use both folder and files, otherwise there are some conflicts
	threadedCollection.SafeIndex.Index.CorpusSize++
}