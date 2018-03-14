package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rmulton/riw_project/requests"
	"github.com/rmulton/riw_project/requests/outputFormaters"
	"github.com/rmulton/riw_project/utils"
)

// NB: some of the expectedDocIDs are empty because of qrels.txt that define some requests that should return no document

// evalRequest contains a request and the IDs of the document expected in the response
type evalRequest struct {
	req            string
	expectedDocIDs map[int]bool
}

var stopList = []string{
	"interested",
	"considerations",
	"topic",
	"topics",
	"subtopic",
	"subtopics",
	"resources",
	"addressing",
	"find",
	"discussions",
	"descriptions",
	"especially",
	"applicable",
	"specific",
	"interest",
	"article",
	"articles",
	"describing",
	"finding",
	"relating",
	"results",
	"list",
	"special",
	"dealing",
}

// getRecall calculates the precision score for a response given the search engine response and the expected response
func getPrecision(res map[int]bool, expected map[int]bool) float64 {
	truePos := 0
	for docFound := range res {
		_, exists := expected[docFound]
		if exists {
			truePos++
		}
	}
	truePosFalsePos := len(res)
	precision := 0.
	if truePosFalsePos > 0 {
		precision = float64(truePos) / float64(truePosFalsePos)
	} else {
		log.Printf("Error, truePosFalsePos is negative or null: %d", truePosFalsePos)
	}
	return precision
}

// getRecall calculates the recall score for a response given the search engine response and the expected response
func getRecall(res map[int]bool, expected map[int]bool) float64 {
	truePos := 0
	for docFound := range res {
		_, exists := expected[docFound]
		if exists {
			truePos++
		}
	}
	truePosFalseNeg := len(expected)
	recall := 0.
	if truePosFalseNeg > 0 {
		recall = float64(truePos) / float64(truePosFalseNeg)
	}
	return recall
}

// getQrels parses the qrels.text file that contains the labels for the requests
func getQrels(collectionPath string) map[int]map[int]bool {
	output := make(map[int]map[int]bool)
	stringFileRels := utils.FileToString(collectionPath + "qrels.text")
	lines := strings.Split(stringFileRels, "\n")
	for _, line := range lines {
		if line != "" {
			// Get the ID of the request and the doc refered in this line
			els := strings.Split(line, " ")
			reqIDstr := els[0]
			reqID, err := strconv.Atoi(reqIDstr)
			if err != nil {
				log.Println(err)
			}
			docIDstr := els[1]
			docID, err := strconv.Atoi(docIDstr)
			if err != nil {
				log.Println(err)
			}
			// Append the result to the output
			_, exists := output[reqID]
			if !exists {
				output[reqID] = make(map[int]bool)
			}
			output[reqID][docID] = true
		}
	}
	return output
}

// getEvalRequests parses the test set
func getEvalRequests(collectionPath string) []evalRequest {
	// Output variable
	var evalRequests []evalRequest
	regexDoc := regexp.MustCompile("\\.I ([0-9]+)\n")
	regexDocPart := regexp.MustCompile("\\.([A-Z])\n")

	// Get Qrels
	qrels := getQrels(collectionPath)

	// Split requests
	stringFileQueries := utils.FileToString(collectionPath + "query.text")
	documents := regexDoc.Split(stringFileQueries, -1)
	for i, doc := range documents {
		var req string
		partsContent := regexDocPart.Split(doc, -1)
		partsName := regexDocPart.FindAllStringSubmatch(doc, -1)
		for j, partName := range partsName {
			partName := partName[1]
			// Only use text from the W part
			if partName == "W" {
				req = partsContent[j+1]
			}
		}
		expectedDocIDs := qrels[i]

		// Append to the output variable
		evalRequest := evalRequest{
			req:            req,
			expectedDocIDs: expectedDocIDs,
		}
		evalRequests = append(evalRequests, evalRequest)
	}
	return evalRequests
}

func writeScoresToCSV(filename string, scores [][]string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, score := range scores {
		err := writer.Write(score)
		if err != nil {
			log.Println(err)
		}
	}
}

func getNbDocsReturned(scores []float64, scoreToDocs map[float64][]int) int {
	nDocs := 0
	for _, score := range scores {
		nDocs += len(scoreToDocs[score])
	}
	return nDocs
}

func getIBestDocs(scores []float64, scoreToDocs map[float64][]int, i int) map[int]bool {
	iBestDocs := make(map[int]bool)

	// Try to use j scores
	for _, score := range scores {
		// Last score that has documents for us
		if len(iBestDocs)+len(scoreToDocs[score]) >= i {
			docsLeftToAppend := i - len(iBestDocs)
			for _, doc := range scoreToDocs[score][:docsLeftToAppend] {
				iBestDocs[doc] = true
			}
			// log.Printf("ibestdocs has size %d, i is %d", len(iBestDocs), i)
			return iBestDocs
		}
		// Add the documents to the list
		for _, doc := range scoreToDocs[score] {
			iBestDocs[doc] = true
		}
	}

	// TODO: return Error instead
	fmt.Printf("Error: i should be < nb of docs returned, not %d", i)
	return nil
}

// getRecallPrecisionCurve returns the recall-precision points as a slice of string slices
func getRecallPrecisionCurve(scores []float64, scoreToDocs map[float64][]int, expectedDocs map[int]bool, numberOfDocs int) [][]string {
	// points are (recall, precision)
	precisionAllDocs := float64(len(expectedDocs)) / float64(numberOfDocs)
	if precisionAllDocs <= 0 {
		fmt.Printf("Precision for all docs: %f\n", precisionAllDocs)
		fmt.Printf("Expected %d docs\n", len(expectedDocs))
	}
	output := [][]string{
		[]string{
			"recall",
			"precision",
		},
	}

	// Score if you return every document of the collection
	output = append(output, []string{
		"1.",
		strconv.FormatFloat(precisionAllDocs, 'f', -1, 64),
	})

	currentPrec := precisionAllDocs
	nbDocsReturned := getNbDocsReturned(scores, scoreToDocs)
	i := nbDocsReturned
	for i > 0 {
		iBestDocs := getIBestDocs(scores, scoreToDocs, i)
		recall := getRecall(iBestDocs, expectedDocs)
		precision := getPrecision(iBestDocs, expectedDocs)
		if precision < currentPrec {
			precision = currentPrec
		} else {
			currentPrec = precision
		}
		point := []string{
			strconv.FormatFloat(recall, 'f', -1, 64),
			strconv.FormatFloat(precision, 'f', -1, 64),
		}
		output = append(output, point)
		i--
	}
	return output
}

// evaluateCacm is a procedure to calculate the recall and precision for an index on cacm using the test set
func evaluateCacm(collectionPath string) {
	evalRequests := getEvalRequests(collectionPath)
	// Build the index
	index := buildIndex(collectionPath, "cacm", true)
	// Get the engine
	engine := requests.NewEngine(index, "vectorial", "dumb") // TODO: here we shouldn't have to enter an output formater type
	// Read the requests
	for i, evalRequest := range evalRequests {
		// Get the response
		res := engine.GetRes(evalRequest.req, stopList)
		if len(res) > 0 {
			// Sort it
			scores, scoreToDocs := outputFormaters.Sort(res)
			recallPrecisionCurve := getRecallPrecisionCurve(scores, scoreToDocs, evalRequest.expectedDocIDs, index.GetDocCounter())

			writeScoresToCSV(fmt.Sprintf("./evaluation/data/recall_precision%d.csv", i), recallPrecisionCurve)
		}
	}
}
