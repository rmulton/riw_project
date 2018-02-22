# Information retrieval engine
Project done during Céline Hudelot class on Information Retrieval.

# Installation
```sh
# clone this repository
git clone https://github.com/rmulton/riw_project
# move to the folder
cd riw_project
# build the application
go build
# run the program
./riw_project -build_for <path_to_collection> -collection <collection_type>
```

# Design
## Priorities
The main priority of this program is to have **the shortest response to request time**. A possible drawback could be a slower index building time.
The second priority is to **allow the user to easily extend the program**.

## Target use cases
The goal of this project is to parse a collection of documents of any kind, then build a reversed index to handle search request on the collection. Three use cases are considered:
1. The index can be held in memory
2. The index cannot be held in memory
3. The index cannot be held on one machine (not implemented)

### 1. The index can be held in memory
This is the most simple use case. There is no need to persist intermediate results. The program fills the index while reading the collection, then the index is available to handle queries. For example, stanford or CACM collections can be indexed in this way on a machine that has 4 or 8GB of memory.

### 2. The index cannot be held in memory
Since the index cannot be completly held in memory, it needs to be stored on the disk. This version is slower since it needs to open, read and write a lot of files.

In order to reduce the index building time, it is necessary to reduce the number of time data is read or written.

### 3. The index cannot be held on one machine (not implemented)
Since the index cannot be held on one machine, it cannot be held on one machine, it needs to be stored on a distributed network. This version is slower than the previous because of the networking layer.

In order to reduce the index building time, it is necessary to reduce the quantity of data that needs to be sent through the network

# Implementation
This project is organized in layers in order to keep different parts independant from each other. The data structures shared by all the program can be found in ./indexes.
Data goes through a pipeline :
Reader -[BufferPostingList]-> IndexBuilder -[Index]-> RequestHandler -[PostingList]-> OutputFormater

## Structure
## ./indexBuilders
A Builder interface must be implemented in order to give the building procedure for an index. The builders currently implemented are:
- in memory builders: build the index in memory
- on disk builders: use the hard disk to extend the maximum index size.

Note Bene:
- The on disk builder implemented doesn't use Block Search Based Indexing (BSBI). Instead, in order to have a shorter querying time, the index is persisted in a document-based manner. Each term of the collection has a file that contains its posting list. This way, when querying the index, the program can load the smallest amount of data possible. When using BSBI, all the blocks that contain postings for the queried terms are loaded, resulting in a lot of unnecessary reading.
- The current rule to use the hard disk is to give a maximum size to the index, then write the biggest posting list to the disk when the index's size exceed the maximum size.

## ./readers
A Reader interface must be implemented in order to give the reading procedure of a collection. It sends parsed documents to the index builder through a channel.
Channels are used here, not Mutexes, so that a parser routine wouldn't be blocked waiting for the lock. The readers implemented are :
- cacm: specific to CACM documents collection
- stanford: every file contained in the input folder is considered as a document

## ./indexes
This folder contains all the **data structures shared throughout this program**.
- PostingList is the core data structure of this program. It stores a document id to document score map
- A document parsed by a reader is sent as a Document through a ReadingChannel
- It is then used to fill a group of PostingList of a BuildingIndex
- When the Index needs to write a PostingList because it is full, it sends it through the WritingChannel
- When the user enter a request, accurate documents are found using a RequestableIndex. Then angle between the request and a given document is output in a VectorPostingList

## ./requests
- A RequestHandler interface implementation gives a procedure to query a requestable index
- An OutputFormater interface implementation gives a procedure to display a request's response to the user
- The Engine structure is the interface between the command line and the querying procedure

Implemented :
- RequestHandler: binary, vectorized, and
- OutputFormater: dumb, sort documents on score

## ./normalizers
- The Normalize(paragraph, stopList) function gives the procedure to normalize a paragraph

## ./utils
Contains helper functions
- Read and write .gob files
- Check whether a path correspond to a file or a folder
- Clear a folder
- Clear the current index saved in ./saved
- Get the content of a file as a string

# Further work
- Compare maps with tuples or binaries to represent the posting lists
- Write the distributed version (separate the writer from the index builder, add the networking layer)
- Compare performance using B+ trees instead of Document-oriented data structure
- Compare writing the biggest posting list when a BufferIndex has exceeded its maximum size with other heuristics to choose what to write
- Extend PostingList to allow posting lists that cannot be held in memory