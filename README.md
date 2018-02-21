# Golang information retrieval engine
Project done during Céline Hudelot class on Information Retrieval.

# Installation
```
// Clone this repository
git clone https://github.com/rmulton/riw_project
// Move to the folder
cd riw_project
// Build the application
go build
// Run the program
./riw_project -build_for <path_to_collection> -collection <collection_type>
```





# Design
## Priority
The main design choice of this program is that **the priority is to have the shortest response to request time**. A possible drawback could be a slower index building time.
## Use cases
The goal of this project is to parse a collection of documents of any kind, then build a reversed index to handle search request on the collection. Three use cases are considered:
1. The index can be held in memory
2. The index cannot be held in memory
3. The index cannot be held on one machine (not implemented)

### 1. The index can be held in memory
This is the most simple use case. There is no need to persist intermediate results. The program fills the index while reading the collection, then the index is available to handle queries. For example, stanford or CACM collections can be indexed in this way on a machine that has 4 or 8GB of memory.

### 2. The index cannot be held in memory
Since the index cannot be completly held in memory, it needs to be stored on the disk. This version is slower since it needs to open, read and write a lot of files.

In order to reduce the index building time, it is necessary to reduce the number of time data is read or written.

<!-- For this reason, we only move from the memory to the disk the current longest postings list when a postings lists buffer is full. 
The limit would be a collection for which one term would have a posting list that would be impossible to be held in memory. In that case, we would have to break this file down into several files.
NB : if the collection (not the index) cannot be held in memory, check if it would be more efficient to use BSBI.
NB : 
- Might be better to use the swap partition of the computer when possible.
- It is almost a single-pass in-memory indexing. Compare the results with real single-pass in-memory indexing.
- B+ tree ?
- However, it might be more clever to use a distributed database system to avoid reinventing the wheel. Since we need all parsers to be able to write to the database at the same time, availability is required. Hence a A-P database system is required. A document-oriented database would suit this use case (SimpleDB for instance).
 -->
 ### 3. The index cannot be held on one machine (not implemented)
Since the index cannot be held on one machine, it cannot be held on one machine, it needs to be stored on a distributed network. This version is slower than the previous because of the networking layer.

In order to reduce the index building time, it is necessary to reduce the quantity of data that needs to be sent through the network
<!-- In this case, in addition to the components that have been previously discussed, we need:
- a component that distributes documents to parse from a document queue
- a component that merges intermediate results from the parsing machines -->

# Architecture
This project is organized in layers in order to keep different parts independant from each other. The data structures shared by all the program can be found in ./indexes.
Data goes through a pipeline :
Reader -[BufferPostingList]-> IndexBuilder -[Index]-> RequestHandler -[PostingList]-> OutputFormater
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