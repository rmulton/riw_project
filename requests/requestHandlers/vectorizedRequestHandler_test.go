package requestHandlers

import (
	"math"

	"github.com/rmulton/riw_project/indexes"
)

// Warning: the request must be normalized,
var someVecRequests = map[string]indexes.PostingList{
	"aaa bbb": indexes.VectorPostingList{
		12: map[string]float64{
			"aaa": (1 + math.Log(1.)) * math.Log(5./2.),
			"bbb": (1 + math.Log(1.)) * math.Log(5./1.),
		},
		64: map[string]float64{
			"aaa": (1 + math.Log(3.)) * math.Log(5./2.),
		},
	}.ToAnglesTo(map[string]float64{"aaa": 1., "bbb": 1.}),
}

// func TestVecRequestHandler(t *testing.T) {
// 	// Get the index
// 	// TODO: check how to avoid duplicate code with ./indexbuilders
// 	readingChan := make(indexes.ReadingChannel)
// 	var wg sync.WaitGroup
// 	var builder = inmemorybuilders.NewInMemoryBuilder(readingChan, 2, &wg)
// 	wg.Add(1)
// 	go builder.Build()
// 	for _, doc := range testInMemSomeDocuments {
// 		readingChan <- doc
// 	}
// 	close(readingChan)
// 	// NB: the tf-idf functionality is tested in ./indexes. Here we rely on it and keep
// 	// scores as integers for clarity
// 	wg.Wait()
// 	index := builder.GetIndex()
// 	// Test the request handler
// 	vecRequestHandler := NewVectorizedRequestHandler(index)
// 	for request, expectedResponse := range someVecRequests {
// 		res := vecRequestHandler.Request(request, []string{})
// 		if !reflect.DeepEqual(*res, expectedResponse) {
// 			t.Errorf("Response to %s should be %v, not %v", request, expectedResponse, *res)
// 		}
// 	}
// }
