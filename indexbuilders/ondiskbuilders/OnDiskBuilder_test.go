package ondiskbuilders

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/utils"
)

var onDiskWaitGroup sync.WaitGroup
var onDiskReadingChannel = make(indexes.ReadingChannel)

var someDocuments = []indexes.Document{
	indexes.Document{
		ID:   0,
		Path: "./mock_path/test.test",
		NormalizedTokens: []string{
			"blabla",
			"bleble",
			"blublu",
		},
	},
	indexes.Document{
		ID:   12,
		Path: "ssljfsd",
		NormalizedTokens: []string{
			"aaa",
			"bbb",
			"ccc",
			"blabla",
			"blabla",
		},
	},
	indexes.Document{
		ID:   64,
		Path: "sdlkajfsdflkdfjd",
		NormalizedTokens: []string{
			"aaa",
			"bleble",
			"lskdjfldsjf",
			"dkjfdsljf",
			"aaa",
			"aaa",
			"dkjfdsljf",
			"dkjfdsljf",
			"dkjfdsljf",
		},
	},
	indexes.Document{
		ID:   27,
		Path: "dlfkdsl",
		NormalizedTokens: []string{
			"llll",
			"blabla",
			"blibli",
		},
	},
	indexes.Document{
		ID:   324,
		Path: "dflkjdsl",
		NormalizedTokens: []string{
			"lalala",
			"bebebe",
			"toto",
			"toto",
		},
	},
	indexes.Document{
		ID:   3674,
		Path: "djdsl",
		NormalizedTokens: []string{

			"aaa",
			"aaa",
			"aaa",
			"aaa",
			"aaa",

			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",
			"blabla",

			"bleble",
			"bleble",

			"bebebe",
			"bebebe",
			"bebebe",
			"bebebe",
			"bebebe",
			"bebebe",

			"toto",
			"toto",
			"toto",
			"toto",
			"toto",
			"toto",
			"toto",
			"toto",
			"toto",
		},
	},
}

var expectedPostingLists = map[string]indexes.PostingList{
	"blabla": indexes.PostingList{
		0:    1,
		12:   2,
		27:   1,
		3674: 36,
	},
	"bleble": indexes.PostingList{
		0:    1,
		64:   1,
		3674: 2,
	},
	"blublu": indexes.PostingList{
		0: 1,
	},
	"aaa": indexes.PostingList{
		12:   1,
		64:   3,
		3674: 5,
	},
	"bbb": indexes.PostingList{
		12: 1,
	},
	"ccc": indexes.PostingList{
		12: 1,
	},
	"lskdjfldsjf": indexes.PostingList{
		64: 1,
	},
	"dkjfdsljf": indexes.PostingList{
		64: 4,
	},
	"llll": indexes.PostingList{
		27: 1,
	},
	"blibli": indexes.PostingList{
		27: 1,
	},
	"lalala": indexes.PostingList{
		324: 1,
	},
	"bebebe": indexes.PostingList{
		324:  1,
		3674: 6,
	},
	"toto": indexes.PostingList{
		324:  2,
		3674: 9,
	},
}

var expectedDocIDToFilePath = map[int]string{
	0:    "./mock_path/test.test",
	12:   "ssljfsd",
	64:   "sdlkajfsdflkdfjd",
	27:   "dlfkdsl",
	324:  "dflkjdsl",
	3674: "djdsl",
}

func TestBuildOnDisk(t *testing.T) {
	utils.ClearOrCreatePersistedIndex("./saved")
	var builder = NewOnDiskBuilder(3, onDiskReadingChannel, 2, &onDiskWaitGroup)
	onDiskWaitGroup.Add(1)
	go builder.Build()
	for _, doc := range someDocuments {
		builder.readingChannel <- doc
	}
	close(builder.readingChannel)
	onDiskWaitGroup.Wait()
	// NB: the tf-idf functionality is tested in ./indexes. Here we rely on it and keep
	// scores as integers for clarity
	for _, postingList := range expectedPostingLists {
		postingList.TfIdf(len(someDocuments))
	}
	index := builder.GetIndex()
	// For terms that are in both indexes
	for term, expectedPostingList := range expectedPostingLists {
		postingList := index.GetPostingListsForTerms([]string{term})[term]
		if !reflect.DeepEqual(postingList, expectedPostingList) {
			// t.Errorf("\nFor %s:\n   - Should be %v\n   - Not %v\n", term, expectedPostingList, postingList)
			for docID, expectedScore := range expectedPostingList {
				score := postingList[docID]
				if score != expectedScore {
					t.Errorf("Score for doc %d for term %s should be %f, not %f", docID, term, expectedScore, score)
				}
			}
		}
	}
	// Check the docID to path map
	if !reflect.DeepEqual(index.GetDocIDToPath(), expectedDocIDToFilePath) {
		t.Errorf("DocID to path is not correct. It should be %#v, not %#v", expectedDocIDToFilePath, index.GetDocIDToPath())
	}
	// Clear the folder after
	os.RemoveAll("./saved")
}
