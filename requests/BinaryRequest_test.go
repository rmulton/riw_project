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
	[2]string{"36-bit", "36-BiT"},
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
			{"chicken", "bacon", "egg"},
		},
	},
	parsedRequest{
		"computer & science & chicken | ham",
		[][]string{
			{"comput", "scienc", "chicken"},
			{"ham"},
		},
	},
	parsedRequest{
		"36-BiT",
		[][]string{
			{"36-bit"},
		},
	},
}

func TestExampleBinaryRequests(t *testing.T) {
	// Get the index
	index := cacm.NewCollection("../consignes/Data/CACM/").Index

	for _, testRequest := range testRequests {
		input := testRequest.input
		binaryRequest := NewBinaryRequest(input, index)
		if !reflect.DeepEqual(binaryRequest.Output, testRequest.docIDs) {
			if len(binaryRequest.Output)>0 || len(testRequest.docIDs)>0 {
				t.Errorf("Found documents %v for '%s', should be %v", binaryRequest.Output, input, testRequest.docIDs)
			}
		}
	}
}
func TestEqualBinaryRequests(t *testing.T) {
	index := cacm.NewCollection("../consignes/Data/CACM/").Index
	for _, equalTestRequest := range equalTestRequests {
		binaryRequest1 := NewBinaryRequest(equalTestRequest[0], index)
		binaryRequest2 := NewBinaryRequest(equalTestRequest[1], index)
		for _, docID := range binaryRequest1.Output {  // TODO : replace with reflect.DeepEqual to take frqc into account, and change to avoid adding the same frqc several times
			exists := false
			for _, foundDocID := range binaryRequest2.Output {
				if docID == foundDocID {
					exists = true
				}
			}
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

type testSortedDocID struct {
	docIDs map[int]float64
	sortedDocIDList []int // Can it has a problem since there are several solutions ?
}

var testSortedDocIDs = []testSortedDocID{
	testSortedDocID{
		map[int]float64{
			23: 1,
			24: 2,
			834823: 56,
			1: 9,
		},
		[]int{834823, 1, 24, 23},
	},
}

func TestNewUserOutput(t *testing.T) {
	for _, testSortedDocID := range testSortedDocIDs {
		sortedDocIDs := getSortedDocList(testSortedDocID.docIDs)
		if !reflect.DeepEqual(sortedDocIDs, testSortedDocID.sortedDocIDList) {
			t.Errorf("%v is sorted %v, should be %v", testSortedDocID.docIDs, sortedDocIDs, testSortedDocID.sortedDocIDList)
		}
	}

}