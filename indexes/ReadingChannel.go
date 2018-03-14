package indexes

// Document contains a document parsed by a reader
type Document struct {
	Path             string
	ID               int
	NormalizedTokens []string
}

// ReadingChannel is a channel for Document
type ReadingChannel chan Document
