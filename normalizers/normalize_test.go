package normalizers

import (
	"testing"
	"../parsers/cacm"
)

type word struct {
	word string
	normalized string
}

var stopList = cacm.GetStopListFromFolder("../consignes/Data/CACM/")
var testWords = []word{
	word{"horse", "hors"}, // test word
	word{" horse ", "hors"}, // test word
	word{"accordingly", ""}, // stop word
}

func TestNormalizeWord(t *testing.T) {
	for _, word := range testWords {
		normalizedWord := NormalizeWord(word.word, stopList)
		if normalizedWord != word.normalized {
			t.Errorf("%s normalized in %s, expected %s", word.word, word.normalized, normalizedWord)
		}
	}
}