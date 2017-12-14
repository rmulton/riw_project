package normalizers

import (
	"github.com/kljensen/snowball"
)

// NormalizeWords normalizes a list of words
func NormalizeWords(words []string, stopWords []string) []string {
	normalizedWords := []string{}
	for _, word := range words {
		normalizedWord := NormalizeWord(word, stopWords)
		if normalizedWord != "" { // TODO: Check what happends for stopwords
			normalizedWords = append(normalizedWords, normalizedWord)
		}
	}
	return normalizedWords
}

// NormalizeWord normalizes a single word
func NormalizeWord(word string, stopWords []string) string {
	if contains(stopWords, word) {
		return ""
	}
	stemed, err := snowball.Stem(word, "english", true)
	if err != nil {
		panic(err)
	}
	return stemed
}

func contains(list []string, element string) bool {
	for _, el := range list {
		if el == element {
			return true
		}
	}
	return false
}