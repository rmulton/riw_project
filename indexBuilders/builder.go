package indexBuilders

import (
	"../indexes"
)

type IndexBuilder interface {
	Build()
	GetIndex() indexes.RequestableIndex
}