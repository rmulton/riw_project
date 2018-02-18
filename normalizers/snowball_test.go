package normalizers

import (
	"reflect"
	"testing"
)

type testParagraph struct {
	content string
	stopList []string
	normalized []string
}


var testParagraphs = []testParagraph{
	testParagraph{
		"horse  tree accordingly HoRse 36-BiT   64_for robin.com er#3",
		[]string{"accordingly"},
		[]string{"hors", "tree", "hors", "36-bit", "64_for", "robin.com", "er", "3"},
	},
}

func TestNormalizeWord(t *testing.T) {
	for _, testParagraph := range testParagraphs {
		normalizedContent := Normalize(testParagraph.content, testParagraph.stopList)
		if !reflect.DeepEqual(normalizedContent, testParagraph.normalized) {
			t.Errorf("%s normalized in %v, expected %v", testParagraph.content, normalizedContent, testParagraph.normalized)
		}
	}
}