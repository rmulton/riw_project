package parsers

import (
	"../indexes"
)

type Parser interface {
	Parse(collectionName string) *indexes.ReversedIndex
}

// No common functions to implement for the parsers. The common computations are all done by the index.