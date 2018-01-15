package utils

import (
	"io/ioutil"
	"os"
	"encoding/gob"
)

// NB: it is often faster to pass by value rather than difference in golang
func FileToString(filePath string) string {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	fileString := string(dat)
	return fileString
}

func WriteGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func ReadGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}