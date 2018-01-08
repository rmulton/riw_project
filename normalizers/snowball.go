package normalizers

import (
	"strings"
	"github.com/kljensen/snowball"
)

// Get a normalized token list from a string
func Normalize(paragraph *string, stopwords *[]string) *[]string {
	tokens := strings.FieldsFunc(*paragraph, func(r rune) bool {
		return r == ' ' || r == '.' || r == '\n' || r == ',' || r == '?' || r == '!' || r == '(' || r == ')' || r == '*' || r == ';' || r == '"' || r == '\'' || r == ':' || r == '{' || r == '}' || r == '/' || r == '|'
	})
	return normalizeWords(&tokens, stopwords)
}

// normalizeWords normalizes a list of words
func normalizeWords(words *[]string, stopWords *[]string) *[]string {
	normalizedWords := []string{}
	for _, word := range *words {
		normalizedWord := normalizeWord(word, stopWords)
		if normalizedWord != "" { // TODO: Check what happends for stopwords
			normalizedWords = append(normalizedWords, normalizedWord)
		}
	}
	return &normalizedWords
}

// normalizeWord normalizes a single word
func normalizeWord(word string, stopWords *[]string) string {
	if !contains(stopWords, word) {
		stemed, err := snowball.Stem(word, "english", true)
		if err != nil {
			panic(err)
		}
		return stemed
	}
	empty := ""
	return empty
}

func contains(list *[]string, element string) bool {
	for _, el := range *list {
		if el == element {
			return true
		}
	}
	return false
}