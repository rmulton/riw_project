package requests

import (
	"../indexes"
)

type Engine interface {
	LoadEngine() *indexes.ReversedIndex
	Request(query string) UserOutput
}