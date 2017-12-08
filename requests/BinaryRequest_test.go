package requests

import (
	"testing"
	"../parsers/cacm"
)

type testRequest struct {
	input string
	docIDs []int
}

var testRequests = []testRequest {
	{"extraction AND roots AND subtractions AND digital AND computers", []int{2}},
	{"chicken", []int{}}, // TODO : add caps
}

var equalTestRequests = [][2]string {
	[2]string{"computer AND science", "science AND computer"},
	[2]string{"Computer AND scIEnce", "SciencE AND comPUter"},
}

var unequalTestRequests = [][2]string {
	[2]string{"toto", "computer AND science"},
	[2]string{"computer AND science", "computer OR science"},
	[2]string{"computer", "science"},
}

func TestExampleBinaryRequests(t *testing.T) {
	// Get the index
	index := cacm.NewCollection("../consignes/Data/CACM/").Index

	for _, testRequest := range testRequests {
		input := testRequest.input
		binaryRequest := NewBinaryRequest(input, index)
		foundDocIDs := binaryRequest.Output
		for key, _ := range foundDocIDs {
			for _, expectedDocId := range testRequest.docIDs {
				if key == expectedDocId {
					delete(foundDocIDs, key)
				}
			}
			if len(foundDocIDs) != 0 {
				t.Errorf("Found unexpected document with id %d", key)
			}
		}
	}
}
func TestEqualBinaryRequests(t *testing.T) {
	index := cacm.NewCollection("../consignes/Data/CACM/").Index
	for _, equalTestRequest := range equalTestRequests {
		binaryRequest1 := NewBinaryRequest(equalTestRequest[0], index)
		binaryRequest2 := NewBinaryRequest(equalTestRequest[1], index)
		for docID, _ := range binaryRequest1.Output {
			_, exists := binaryRequest2.Output[docID]
			if !exists {
				t.Errorf("'%s' don't output doc %d as '%s' does", binaryRequest1.input, docID, binaryRequest2.input)
			}
		}
	}
}

func TestUnequalBinaryRequests(t *testing.T) {
	index := cacm.NewCollection("../consignes/Data/CACM/").Index
	for _, unequalTestRequest := range equalTestRequests {
		binaryRequest1 := NewBinaryRequest(unequalTestRequest[0], index)
		binaryRequest2 := NewBinaryRequest(unequalTestRequest[1], index)
		similarDocIDs := binaryRequest1.Output
		for docID, _ := range binaryRequest1.Output {
			_, exists := binaryRequest2.Output[docID]
			if exists {
				delete(similarDocIDs,docID)
			}
		}
		if len(similarDocIDs)!=0 {
			t.Errorf("'%s' outputs the same documents as '%s'", binaryRequest1.input, binaryRequest2.input)
		}
	}
}