# Golang information retrieval engine
Project done during Céline Hudelot class on Information Retrieval.

# Installation
```git clone https://github.com/rmulton/riw_project```

```cd riw_project```

```go build```

```./riw_project```

# Architecture
## /parsers
A Parser handles everything that is specific to a document collection. It needs to parse the collection and fill an Index doing so.
## /indexes
An Index is filled by a Parser. Then it applies some rules to transform the score the Parser has given to every document for every word.
## /requests
A Request handles a certain kind of request. It takes a string input from the user and compute the output according to an Index.
## /normalizers
A Normalizer maps words with terms using NLP procedures.
## /utils
Utils contains functions to
- Transform files to string
- Read and write Gob files to persist golang structures

# Bugs
- Binary request with spaces between words work, which is weird