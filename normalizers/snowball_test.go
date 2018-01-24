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
		"horse  tree accordingly",
		[]string{"accordingly"},
		[]string{"hors", "tree"},
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