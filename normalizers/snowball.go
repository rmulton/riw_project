package normalizers

import (
	"regexp"
	"log"
	"github.com/kljensen/snowball"
)

// Get a normalized token list from a string
func Normalize(paragraph string, stopwords []string) []string {
	wordRegex := regexp.MustCompile("[A-z0-9][A-z0-9\\.\\_\\/]+[A-z0-9]") // avoids having dots or slashes at the begining or the end of the name of the file
	tokens := wordRegex.FindAllString(paragraph, -1)

	return normalizeWords(tokens, stopwords)
}

// normalizeWords normalizes a list of words
func normalizeWords(words []string, stopWords []string) []string {
	normalizedWords := []string{}
	for _, word := range words {
		normalizedWord := normalizeWord(word, stopWords)
		if normalizedWord != "" { // TODO: Check what happends for stopwords
			normalizedWords = append(normalizedWords, normalizedWord)
		}
	}
	return normalizedWords
}

// normalizeWord normalizes a single word
func normalizeWord(word string, stopWords []string) string {
	if !contains(stopWords, word) {
		stemed, err := snowball.Stem(word, "english", true)
		if err != nil {
			log.Println(err)
		}
		return stemed
	}
	empty := ""
	return empty
}

func contains(list []string, element string) bool {
	for _, el := range list {
		if el == element {
			return true
		}
	}
	return false
}