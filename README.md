# Golang information retrieval engine
Project done during Céline Hudelot class on Information Retrieval.

# Installation
```git clone https://github.com/rmulton/riw_project```

```cd riw_project```

```go build```

```./riw_project -collection=cacm```

# Design
NB : if the collection (not the index) cannot be held in memory, check if it would be more efficient to use BSBI.
## Assumptions
The priority is to have the fastest response to request time. (Which can harm the index building time).
Any document of the collection can be held in memory.
Compute data for the user asap.
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

# Architecture

## Choices
- Specific reading procedures for a collection is implemented in /parsers
- Common calculation procedures to index the documents are handled in /indexes
- Specific request parsing procedures are handled in /requests
- Common calculation procedures to answer the request are handled in /engines

## /parsers
A Parser handles reading and parsing procedure for a collection. It needs to send the parsed documents to the index builder through a channel. Channels are used and not Mutex so that a parser routine wouldn't be blocked waiting for the lock. See the parser interface.

Implemented :
- CACM
- Stanford

## /indexes
An Index stores the parsed documents. When the index is finished, it applies a procedure to transform the postings lists to scores.

Implemented :
- TfIdf index : term frequency - inverse document frequency index

## /requests
A Request handles a certain kind of request. It takes a string input from the user and compute the output according to an Index.

Implemented :
- Binary requests : requests using "and" and "or" conditions
- Vectorized requests : requests using the angle between the request and the documents

## /engines
An Engine stores an index and responds to requests with a sorted list of documents that correspond to the request.

## /normalizers
A Normalizer maps words with terms using NLP procedures.

## /utils
Utils contains functions to
- Transform files to string
- Read and write Gob files to persist golang structures

# Further work
- Compare persisting maps with persisting tuples or binaries that represent the posting lists and use an index
- Write the distributed version (separate the writer from the index builder; add the network layer)