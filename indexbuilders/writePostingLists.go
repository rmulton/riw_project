package indexbuilders

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/utils"
)

// CurrentPostingListOnDisk loads the posting list for a term currently saved on the disk
func CurrentPostingListOnDisk(term string) indexes.PostingList {
	termFile := fmt.Sprintf("./saved/postings/%s", term)
	err, postingListSoFar := indexes.PostingListFromFile(termFile)
	if err != nil {
		log.Println(err)
	}
	return postingListSoFar
}

// WritePostingLists writes the posting lists that arrive on the writing channel
func WritePostingLists(writingChannel indexes.WritingChannel, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	for toWrite := range writingChannel {
		// Write it to the disk
		// TODO : duplicate code with CurrentPostingListOnDisk()
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
				postingListSoFar = CurrentPostingListOnDisk(toWrite.Term)
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
