package indexBuilders

import (
	"testing"
	"sync"
	"reflect"
	"../indexes"
)

var testInMemWaitGroup sync.WaitGroup
var testInMemReadingChannel = make(chan indexes.Document)

var testInMemSomeDocuments = []indexes.Document {
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
}

var expectedInMemPostingLists = map[string]indexes.PostingList {
	"blabla": indexes.PostingList {
		0: 1,
		12: 2,
		27: 1,
	},
	"bleble": indexes.PostingList {
		0: 1,
		64: 1,
	},
	"blublu": indexes.PostingList {
		0: 1,
	},
	"aaa": indexes.PostingList {
		12: 1,
		64: 3,
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
	},
	"toto": indexes.PostingList {
		324: 2,
	},
}

var expectedInMemDocIDToFilePath = map[int]string {
	0: "./mock_path/test.test",
	12: "ssljfsd",
	64: "sdlkajfsdflkdfjd",
	27: "dlfkdsl",
	324: "dflkjdsl",
}

func TestBuildInMemory(t *testing.T) {
	var builder = NewInMemoryBuilder(testInMemReadingChannel, 2, &testInMemWaitGroup)
	builder.parentWaitGroup.Add(1)
	go builder.Build()
	for _, doc := range testInMemSomeDocuments {
		builder.readingChannel <- doc
	}
	close(builder.readingChannel)
	// NB: the tf-idf functionality is tested in ./indexes. Here we rely on it and keep
	// scores as integers for clarity
	builder.parentWaitGroup.Wait()
	for _, postingList := range expectedInMemPostingLists {
		postingList.TfIdf(5)
	}
	index := builder.GetIndex()
	// For terms that are in both indexes
	for term, expectedPostingList := range expectedInMemPostingLists {
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
}