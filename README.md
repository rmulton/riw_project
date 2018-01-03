# Golang information retrieval engine
Project done during Céline Hudelot class on Information Retrieval.

# Installation
```git clone https://github.com/rmulton/riw_project```

```cd riw_project```

```go build```

```./riw_project -collection=cacm```

# Architecture
NB : if the collection (not the index) cannot be held in memory, check if it would be more efficient to use BSBI.
## Use cases
The goal of this project is to parse a collection of documents of any kind (for now only files), then build a reversed index to handle search request on the collection. It is designed to be used in three cases :

1. The index can be held in memory
2. The index cannot be held in memory
3. The index is filled and requested on a distributed network

### 1. The index can be held in memory
This is the most simple case. There is no need to persist intermediary results. The program fills the index while reading the collection, then the index is available for requests. It is the case for stanford or CACM collection on a machine that has 4 or 8GB of memory.

### 2. The index cannot be held in memory
Since the index cannot be held in memory, we need to store it on the disk. In this case, we have chosen to store the postings list for each term in a separate file. This version is much slower since it needs to open, read and write a lot of files.

The objective here, in order to reduce the time to have the index ready, is to reduce the number of time postings files are read, filled with new postings, then written. For this reason, we only move from the memory to the disk the current longest postings list when a postings lists buffer is full. 

The limit would be a collection for which one term would have a posting list that would be impossible to be held in memory. In that case, we would have to break this file down into several files.

NB : 
- Might be better to use the swap partition of the computer when possible.
- It is almost a single-pass in-memory indexing. Compare the results with real single-pass in-memory indexing.

### 3. The index is filled and requested on a distributed network (not implemented yet)
In that case we need to have a machine that distributes documents to parse from a document queue. Then we also need a machine that merges postings lists sent from the machines that are parsing the documents.

The previous system can be augmented to handle this case. However, it might be more clever to use a distributed database system to avoid reinventing the wheel.

Since we need all parsers to be able to write to the database at the same time, availability is required. Hence a A-P database system is required. A document-oriented database would suit this use case (SimpleDB for instance).

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