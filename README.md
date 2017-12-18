# Golang information retrieval engine
Project done during Céline Hudelot class on Information Retrieval.

# Installation
```git clone https://github.com/rmulton/riw_project```

```cd riw_project```

```go build```

```./riw_project -collection=cacm```

# Architecture
## Choices
- Specific reading procedures for a collection is implemented in /parser
- Common calculation procedures to index the documents are handled in /indexes
- Specific request parsing procedures are handled in /requests
- Common calculation procedures to answer the request are handled in /engines
## /parsers
A Parser handles everything that is specific to a document collection. It needs to parse the collection and fill an Index doing so.
## /indexes
An Index is filled by a Parser. Then it applies some rules to transform the score the Parser has given to every document for every word.
- ReversedIndex implements a simple TF-IdF index procedure.
- SafeReversedIndex implements a threaded version of ReversedIndex
## /requests
A Request handles a certain kind of request. It takes a string input from the user and compute the output according to an Index.
## /engines
An engine is supposed to take a request and an index and return an output (not already done). It should be request-type and collection agnostic.
## /normalizers
A Normalizer maps words with terms using NLP procedures.
## /utils
Utils contains functions to
- Transform files to string
- Read and write Gob files to persist golang structures

# Bugs
- Binary request with spaces between words work, which is weird