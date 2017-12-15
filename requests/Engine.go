package requests

import (
	"../indexes"
)

type Engine interface {
	LoadEngine(refresh bool) *indexes.ReversedIndex
	Request(query string) UserOutput
}