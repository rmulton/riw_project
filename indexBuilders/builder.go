package indexBuilders

import (
	"github.com/rmulton/riw_project/indexes/requestableIndexes"
)

// IndexBuilder builds an index. It gets data from a reader and output a requestableIndexes index.
type IndexBuilder interface {
	Build()
	GetIndex() requestableIndexes.RequestableIndex
}