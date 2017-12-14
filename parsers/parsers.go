package parsers

import (
	"../indexes"
)

type Parser interface {
	Parse(collectionName string) *indexes.ReversedIndex
}
