package requests

import (
	"testing"
	"reflect"
	"../parsers/cacm"
)

type testRequest struct {
	input string
	docIDs []int
}

type parsedRequest struct {
	input string
	ands [][]string
}

var testRequests = []testRequest {
	{"extraction & roots & subtractions & digital & computers", []int{2}},
	{"chicken", []int{}},
	{"chicken | bacon", []int{}},
	{"36-bit", []int{1295, 1691, 3026}},
}

var equalTestRequests = [][2]string {
	[2]string{"computer & science", "science & computer"},
	[2]string{"computer & science | science | computer", "science & computer"},
	[2]string{"computer & science | chicken", "science & computer"},
	[2]string{"Computer & scIEnce", "SciencE & comPUter"},
	// [2]string{"36-bit", "36-BiT"}, // TODO: minimize
	[2]string{"chicken", ""},
}

var unequalTestRequests = [][2]string {
	[2]string{"toto", "computer & science"},
	[2]string{"computer & science", "computer | science"},
	[2]string{"computer & science", "computer & science & chicken"},
	[2]string{"computer", "science"},
}

var parseTestRequests = []parsedRequest {
	parsedRequest{
		"chicken & bacon & eggs",
		[][]string{
			{"chicken", "bacon", "eggs"},
		},
	},
	parsedRequest{
		"computer & science & chicken | ham",
		[][]string{
			{"computer", "science", "chicken"},
			{"ham"},
		},
	},
	parsedRequest{
		"36-BiT",
		[][]string{
			{"36-BiT"},
		},
	},
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
		}
		for key, _ := range foundDocIDs {
			if len(foundDocIDs) != 0 {
				t.Errorf("Found unexpected document with id %d for request %s", key, binaryRequest.input)
			}
		}
	}
}
func TestEqualBinaryRequests(t *testing.T) {
	index := cacm.NewCollection("../consignes/Data/CACM/").Index
	for _, equalTestRequest := range equalTestRequests {
		binaryRequest1 := NewBinaryRequest(equalTestRequest[0], index)
		binaryRequest2 := NewBinaryRequest(equalTestRequest[1], index)
		for docID, _ := range binaryRequest1.Output {  // TODO : replace with reflect.DeepEqual to take frqc into account, and change to avoid adding the same frqc several times
			_, exists := binaryRequest2.Output[docID]
			if !exists {
				t.Errorf("'%s' don't output doc %d as '%s' does", binaryRequest1.input, docID, binaryRequest2.input)
			}
		}
	}
}

func TestUnequalBinaryRequests(t *testing.T) {
	index := cacm.NewCollection("../consignes/Data/CACM/").Index
	for _, unequalTestRequest := range unequalTestRequests {
		binaryRequest1 := NewBinaryRequest(unequalTestRequest[0], index)
		binaryRequest2 := NewBinaryRequest(unequalTestRequest[1], index)
		if equal := reflect.DeepEqual(binaryRequest1.Output, binaryRequest2.Output); equal {
			t.Errorf("'%s' outputs the same documents as '%s'", binaryRequest1.input, binaryRequest2.input)
		}
	}
}

func TestParseRequest(t *testing.T) {
	index := cacm.NewCollection("../consignes/Data/CACM/").Index //Move to avoid creating the index several times
	
	for _, parseTestRequest := range parseTestRequests {
		input := parseTestRequest.input
		binaryRequest := NewBinaryRequest(input, index)
		if equal := reflect.DeepEqual(binaryRequest.ands, parseTestRequest.ands); !equal {
			t.Errorf("Parser output for '%s' is %v, expected %v", input, binaryRequest.ands, parseTestRequest.ands)
		}
	}
}