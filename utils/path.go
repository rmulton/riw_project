package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// 0777 is used to avoid access denied errors on linux when execution permission is not given
var PERMISSION = os.FileMode(0777)

// FileToString gets the content of a file as a string
func FileToString(filePath string) string {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error \"%v\" while trying to read \"%s\"\n", err, filePath)
	}
	fileString := string(dat)
	return fileString
}

// ClearFolder clears a folder's content
func ClearFolder(folderPath string) {
	err := os.RemoveAll(folderPath)
	if err != nil {
		fmt.Println(err)
	}
	err = os.MkdirAll(folderPath, PERMISSION)
	if err != nil {
		fmt.Println(err)
	}
}

// CheckPathExists returns true if a path exists, false otherwise
func CheckPathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println(err)
		return false
	}
	return true
}

// ClearOrCreatePersistedIndex clears a persisted index as structured in this project:
// ./saved/postings/ and ./saved/meta/
func ClearOrCreatePersistedIndex(indexPath string) {
	fmt.Println("Clearing the previous index saved on the disk")
	ClearFolder(indexPath)
	err := os.MkdirAll(fmt.Sprintf("%s/postings", indexPath), PERMISSION)
	if err != nil {
		fmt.Println(err)
	}
	err = os.MkdirAll(fmt.Sprintf("%s/meta", indexPath), PERMISSION)
	if err != nil {
		fmt.Println(err)
	}
}
