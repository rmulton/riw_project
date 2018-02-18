package inMemory

import (
	"testing"
	"sync"
	"reflect"
	"github.com/rmulton/riw_project/indexes"
)

var testInMemWaitGroup sync.WaitGroup
var testInMemReadingChannel = make(chan indexes.Document)

var someDocuments2 = []indexes.Document {
	indexes.Document {
		ID: 0,
		Path: "./mock_path/test.test",
		NormalizedTokens: []string {
			"blabla",
			"bleble",
			"blublu",
		},
	},
	indexes.Document {
		ID: 12,
		Path: "ssljfsd",
		NormalizedTokens: []string {
			"aaa",
			"bbb",
			"ccc",
			"blabla",
			"blabla",
		},
	},
	indexes.Document {
		ID: 64,
		Path: "sdlkajfsdflkdfjd",
		NormalizedTokens: []string {
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
	indexes.Document {
		ID: 27,
		Path: "dlfkdsl",
		NormalizedTokens: []string {
			"llll",
			"blabla",
			"blibli",
		},
	},
	indexes.Document {
		ID: 324,
		Path: "dflkjdsl",
		NormalizedTokens: []string {
			"lalala",
			"bebebe",
			"toto",
			"toto",
		},
	},
	indexes.Document {
		ID: 3674,
		Path: "djdsl",
		NormalizedTokens: []string {

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

var expectedPostingLists2 = map[string]indexes.PostingList {
	"blabla": indexes.PostingList {
		0: 1,
		12: 2,
		27: 1,
		3674: 36,
	},
	"bleble": indexes.PostingList {
		0: 1,
		64: 1,
		3674: 2,
	},
	"blublu": indexes.PostingList {
		0: 1,
	},
	"aaa": indexes.PostingList {
		12: 1,
		64: 3,
		3674: 5,
	},
	"bbb": indexes.PostingList {
		12: 1,
	},
	"ccc": indexes.PostingList {
		12: 1,
	},
	"lskdjfldsjf": indexes.PostingList {
		64: 1,
	},
	"dkjfdsljf": indexes.PostingList {
		64: 4,
	},
	"llll": indexes.PostingList {
		27: 1,
	},
	"blibli": indexes.PostingList {
		27: 1,
	},
	"lalala": indexes.PostingList {
		324: 1,
	},
	"bebebe": indexes.PostingList {
		324: 1,
		3674: 6,
	},
	"toto": indexes.PostingList {
		324: 2,
		3674: 9,
	},
}

var expectedDocIDToFilePath2 = map[int]string {
	0: "./mock_path/test.test",
	12: "ssljfsd",
	64: "sdlkajfsdflkdfjd",
	27: "dlfkdsl",
	324: "dflkjdsl",
	3674: "djdsl",
}

func TestBuildInMemory(t *testing.T) {
	var builder = NewInMemoryBuilder(testInMemReadingChannel, 2, &testInMemWaitGroup)
	builder.parentWaitGroup.Add(1)
	go builder.Build()
	for _, doc := range someDocuments2 {
		builder.readingChannel <- doc
	}
	close(builder.readingChannel)
	// NB: the tf-idf functionality is tested in ./indexes. Here we rely on it and keep
	// scores as integers for clarity
	builder.parentWaitGroup.Wait()
	for _, postingList := range expectedPostingLists2 {
		postingList.TfIdf(len(someDocuments2))
	}
	index := builder.GetIndex()
	// For terms that are in both indexes
	for term, expectedPostingList := range expectedPostingLists2 {
		postingList := index.GetPostingListsForTerms([]string{term})[term]
		if !reflect.DeepEqual(postingList, expectedPostingList) {
			for docID, expectedScore := range expectedPostingList {
				score := postingList[docID]
				if score != expectedScore {
					t.Errorf("Score for doc %d for term %s should be %f, not %f", docID, term, expectedScore, score)
				}
			}
		}
	}
	// Check the docID to path map
	if !reflect.DeepEqual(index.GetDocIDToPath(), expectedDocIDToFilePath2) {
		t.Errorf("DocID to path is not correct")
	}
}