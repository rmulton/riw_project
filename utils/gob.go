package utils

import (
	"io/ioutil"
	"log"
	"os"
	"encoding/gob"
)

// NB: it is often faster to pass by value rather than difference in golang
func FileToString(filePath string) string {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error \"%v\" while trying to read \"%s\"", err, filePath)
	}
	fileString := string(dat)
	return fileString
}

func WriteGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err!=nil {
		log.Printf("Error trying to create \"%s\": %v", filePath, err)
		return err
	}
	encoder := gob.NewEncoder(file)
	encErr := encoder.Encode(object)
	if encErr!=nil {
		log.Printf("Error trying to encode \"%s\": %v", filePath, encErr)
		return encErr
	}
	file.Close()
	return nil
}

func ReadGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error trying to read %s: %v", filePath, err)
		return err
	}
	decoder := gob.NewDecoder(file)
	decErr := decoder.Decode(object)
	if decErr != nil {
		log.Printf("Error trying to decode %s, %v", filePath, decErr)
		return decErr
	}
	file.Close()
	return nil
}