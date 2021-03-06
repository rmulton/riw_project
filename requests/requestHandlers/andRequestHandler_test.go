package requestHandlers

import (
	"math"
	"reflect"
	"sync"
	"testing"

	"github.com/rmulton/riw_project/indexbuilders/inmemorybuilders"
	"github.com/rmulton/riw_project/indexes"
)

// Warning : only use words normalized with snowball algorithm
var testInMemSomeDocuments = []indexes.Document{
	indexes.Document{
		ID:   0,
		Path: "./mock_path/test.test",
		NormalizedTokens: []string{
			"blabla",
			"blebl",
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
			"blebl",
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
}

var someRequests = map[string]indexes.PostingList{
	"aaa bbb": indexes.PostingList{
		12: (1+math.Log(1.))*math.Log(5./2.) + (1+math.Log(1.))*math.Log(5./1.),
	},
}

func TestAndRequestHandler(t *testing.T) {
	// Get the index
	// TODO: check how to avoid duplicate code with ./indexbuilders
	readingChan := make(indexes.ReadingChannel)
	var wg sync.WaitGroup
	var builder = inmemorybuilders.NewInMemoryBuilder(readingChan, 2, &wg)
	wg.Add(1)
	go builder.Build()
	for _, doc := range testInMemSomeDocuments {
		readingChan <- doc
	}
	close(readingChan)
	// NB: the tf-idf functionality is tested in ./indexes. Here we rely on it and keep
	// scores as integers for clarity
	wg.Wait()
	index := builder.GetIndex()
	// Test the request handler
	andRequestHandler := NewAndRequestHandler(index)
	for request, expectedResponse := range someRequests {
		res := andRequestHandler.Request(request, []string{})
		if !reflect.DeepEqual(*res, expectedResponse) {
			t.Errorf("Response to %s should be %v, not %v", request, expectedResponse, *res)
		}
	}
}
