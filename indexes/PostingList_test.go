package indexes

import (
	"reflect"
	"testing"
	"math"
)

var postingLists = []PostingList {
	PostingList {
		1: 342,
		543: 452,
		3: 3,
		76: 45,
		654: 3,
	},
	PostingList {
		23: 34,
		1: 544,
		0: 324,
		92: 2,
	},
}

var expectedTfIdfPostingLists = []PostingList {
	PostingList {
		1: (1 + math.Log(342.)) * math.Log(8./5.),
		543: (1 + math.Log(452.)) * math.Log(8./5.),
		3: (1 + math.Log(3.)) * math.Log(8./5.),
		76: (1 + math.Log(45.)) * math.Log(8./5.),
		654: (1 + math.Log(3.)) * math.Log(8./5.),
	},
	PostingList {
		23: (1 + math.Log(34.)) * math.Log(8./4.),
		1: (1 + math.Log(544.)) * math.Log(8./4.),
		0: (1 + math.Log(324.)) * math.Log(8./4.),
		92: (1 + math.Log(2.)) * math.Log(8./4.),
	},
}

func TestTfIdf(t *testing.T) {
	for i, postingList := range postingLists {
		postingList.TfIdf(8)
		expectedTfIdfPostingList := expectedTfIdfPostingLists[i]
		if !reflect.DeepEqual(expectedTfIdfPostingList, postingList) {
			t.Errorf("Posting list %v should be %v after caculating tf-idf score", postingList, expectedTfIdfPostingList)
		}
	}
}

var someFrqcPostingLists = []PostingList {
	PostingList {
		1: 342,
		543: 452,
		3: 3,
		76: 45,
		654: 3,
	},
	PostingList {
		23: 34,
		1: 544,
		0: 324,
		92: 2,
	},
}

var expectedMergedPostingList = PostingList {
	0: 324,
	1: 544 + 342,
	3: 3,
	23: 34,
	76: 45,
	92: 2,
	543: 452,
	654: 3,
}

func TestMergeWith(t *testing.T) {
	postingList := someFrqcPostingLists[0]
	otherPostingList := someFrqcPostingLists[1]
	postingList.MergeWith(otherPostingList)
	if !reflect.DeepEqual(postingList, expectedMergedPostingList) {
		t.Errorf("Merged posting list should be %v, not %v", expectedMergedPostingList, postingList)
	}
}