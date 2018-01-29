package indexes

import(
	"testing"
	"reflect"
	"math"
)

type testDoc struct {
	id int
	path string
	terms []string
}

var someDocs = []testDoc{
	testDoc{
		0,
		"./mock_doc.test",
		[]string{
			"lala",
			"lala",
			"lili",
			"lulu",
		},
	},
	testDoc{
		1,
		"./mock_doc_1.test",
		[]string{
			"toto",
			"toto",
			"toto",
			"toto",
			"toto",
			"lala",
			"lili",
			"1",
		},
	},
	testDoc{
		234,
		"./mock_doc_234.test",
		[]string {
			"lala",
			"1",
			"1",
			"lala",
		},
	},
}
var expectedFrqc = map[string]PostingList {
	"lala": PostingList {
		0: 2.,
		1: 1.,
		234: 2.,
	},
	"lili": PostingList {
		0: 1.,
		1: 1.,
	},
	"toto": PostingList {
		1: 5.,
	},
	"1": PostingList {
		1: 1.,
		234: 2.,
	},
	"lulu": PostingList {
		0: 1.,
	},
}

var expectedTfIdf = map[string]PostingList {
	"lala": PostingList {
		0: (1 + math.Log(2.)) * math.Log(3./3.),
		1: (1 + math.Log(1.)) * math.Log(3./3.),
		234: (1 + math.Log(2.)) * math.Log(3./3.),
	},
	"lili": PostingList {
		0: (1 + math.Log(1.)) * math.Log(3./2.),
		1: (1 + math.Log(1.)) * math.Log(3./2.),
	},
	"toto": PostingList {
		1: (1 + math.Log(5.)) * math.Log(3./1.),
	},
	"1": PostingList {
		1: (1 + math.Log(1.)) * math.Log(3./2.),
		234: (1 + math.Log(2.)) * math.Log(3./2.),
	},
	"lulu": PostingList {
		0: (1 + math.Log(1.)) * math.Log(3./1.),
	},
}

var expectedClearLulu = map[string]PostingList {
	"lala": PostingList {
		0: 2.,
		1: 1.,
		234: 2.,
	},
	"lili": PostingList {
		0: 1.,
		1: 1.,
	},
	"toto": PostingList {
		1: 5.,
	},
	"1": PostingList {
		1: 1.,
		234: 2.,
	},
}

func TestIndexFrequencyAndTfIdf(t *testing.T) {
	// Simple test
	var testIndex = NewEmptyIndex()
	for _, doc := range someDocs {
		testIndex.AddDocToIndex(doc.id, doc.path)
		for _, term := range doc.terms {
			testIndex.AddDocToTerm(doc.id, term)
		}
	}
	if !reflect.DeepEqual(expectedFrqc, testIndex.postingLists) {
		t.Errorf("%v is different from %v", expectedFrqc, testIndex.postingLists)
	}
	// Tf idf
	testIndex.ToTfIdf(len(someDocs))
	if !reflect.DeepEqual(expectedTfIdf, testIndex.postingLists) {
		t.Errorf("%v is different from %v", expectedTfIdf, testIndex.postingLists)
	}
}

func TestClearPostingList(t *testing.T) {
	// Simple test
	var testIndex = NewEmptyIndex()
	for _, doc := range someDocs {
		testIndex.AddDocToIndex(doc.id, doc.path)
		for _, term := range doc.terms {
			testIndex.AddDocToTerm(doc.id, term)
		}
	}
	// Clear "lulu" posting list
	testIndex.ClearPostingListFor("lulu")
	if !reflect.DeepEqual(testIndex.postingLists, expectedClearLulu) {
		t.Errorf("Index should contain %v, not %v", expectedClearLulu, testIndex.postingLists)
	}
	// Clear all posting lists
	for term, _ := range testIndex.postingLists {
		testIndex.ClearPostingListFor(term)
	}
	if !reflect.DeepEqual(testIndex.postingLists, make(map[string]PostingList)) {
		t.Errorf("Index %v should be empty,", testIndex.postingLists)
	}
}