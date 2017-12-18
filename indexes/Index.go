package indexes

import (

)

// Index represents the interface you have to implement to have an index
type Index interface {
	AddParagraphForDoc(paragraph string, docID int) // Count words in paragraph for docID
	Finish() // Finishes the index. (Applies transformation like getting tf-idf from frequency)
}