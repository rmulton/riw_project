package indexBuilders

import (
	"github.com/rmulton/riw_project/indexes/requestable"
)

// IndexBuilder builds an index. It gets data from a reader and output a requestable index.
type IndexBuilder interface {
	Build()
	GetIndex() requestable.RequestableIndex
}