package requests

import (
	"fmt"
	"regexp"
)

type BinaryRequest struct {
	input string
	ands [][]string
}

func NewBinaryRequest(input string) BinaryRequest {
	return BinaryRequest{input, [][]string{}}
}

func (request BinaryRequest) Parse() {
	andRegex := regexp.MustCompile("[a-z]+[ AND [a-z]+]*")	
	ands := andRegex.FindAllString(request.input, -1)
	for _, and := range ands {
		wordRegex := regexp.MustCompile("[a-z]+")
		terms := wordRegex.FindAllString(and, -1)
		request.ands = append(request.ands, terms)
	}
}

