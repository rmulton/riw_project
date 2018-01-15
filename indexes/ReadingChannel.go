package indexes

type Document struct {
	Path string
	ID int
	NormalizedTokens []string
}
type ReadingChannel chan *Document
