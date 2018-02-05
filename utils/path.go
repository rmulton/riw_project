package utils

import (
	"io/ioutil"
	"log"
	"os"
	"encoding/gob"
)

// WriteGob writes the content of an interface variable to a file
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

// ReadGob reads an interface variable from a file
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