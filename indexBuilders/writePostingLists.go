package indexBuilders

import (
	"sync"
	"github.com/rmulton/riw_project/utils"
	"fmt"
	"os"
	"log"
	"github.com/rmulton/riw_project/indexes"
)

func currentPostingListOnDisk(term string) indexes.PostingList {
	termFile := fmt.Sprintf("./saved/postings/%s", term)
	err, postingListSoFar := indexes.PostingListFromFile(termFile)
	if err != nil {
		log.Println(err)
	}
	return postingListSoFar
}

func writePostingLists(writingChannel indexes.WritingChannel, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for toWrite := range writingChannel {
		// Write it to the disk
		// TODO : duplicate code with currentPostingListOnDisk()
		termFile := fmt.Sprintf("./saved/postings/%s", toWrite.Term)
		// If the file exists, either append it or replace it according to ReplaceCurrentFile field
		// from toWrite
		if _, err := os.Stat(termFile); err == nil {
			postingListSoFar := make(indexes.PostingList)
			// If the field is false, it means that the current posting list persisted should be merged
			// with the new posting list sent to the writingChannel
			if toWrite.ReplaceCurrentFile {
				postingListSoFar = toWrite.PostingList
			} else {
				postingListSoFar = currentPostingListOnDisk(toWrite.Term)
				// Merge the current posting list
				postingListSoFar.MergeWith(toWrite.PostingList)
			}
			
			// Write it to file
			utils.WriteGob(termFile, postingListSoFar)
		// Otherwise create it
		} else {
			utils.WriteGob(termFile, toWrite.PostingList)
		}
	}
}