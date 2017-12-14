package cacm

import (
	"strings"
)

func GetStopListFromFolder(path string) []string {
	var stopListFile = fileToString(path + "common_words")
	var stopList = strings.Split(stopListFile, "\n")
	return stopList
}