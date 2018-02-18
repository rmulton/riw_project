package indexBuilders

import (
	"github.com/rmulton/riw_project/indexes"
)

type IndexBuilder interface {
	Build()
	GetIndex() indexes.RequestableIndex
}