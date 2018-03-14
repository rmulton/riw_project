package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rmulton/riw_project/utils"

	"github.com/rmulton/riw_project/indexes"
)

func writeFrqcsToCSV(filename string, frqcs [][]string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, frqc := range frqcs {
		err := writer.Write(frqc)
		if err != nil {
			log.Println(err)
		}
	}
}

// getFrqcs returns the frequencies for the terms in the index as a slice of string slices
func getFrqcs(postingsFolderPath string) [][]string {
	var frequencies [][]string
	err := filepath.Walk(postingsFolderPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			var postingList indexes.PostingList
			err = utils.ReadGob(path, &postingList)
			if err != nil {
				log.Println(err)
			}
			frequence := len(postingList)
			frequenceStr := strconv.Itoa(frequence)
			// log.Printf("%s: %s\n", info.Name(), frequenceStr)
			frequencies = append(frequencies, []string{frequenceStr})
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return frequencies
}

func main() {
	// Get the frequencies
	frqcs := getFrqcs("../saved/postings")
	// Write to csv
	writeFrqcsToCSV("./frequencies.csv", frqcs)

}
