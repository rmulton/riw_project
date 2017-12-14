package cacm

import (
	"strings"
	"../../utils"
)

func GetStopListFromFolder(path string) []string {
	var stopListFile = utils.FileToString(path + "common_words")
	var stopList = strings.Split(stopListFile, "\n")
	return stopList
}