package requests

import (
	"fmt"
	"time"

	"github.com/rmulton/riw_project/indexes"
	"github.com/rmulton/riw_project/indexes/requestableIndexes"
	"github.com/rmulton/riw_project/requests/outputFormaters"
	"github.com/rmulton/riw_project/requests/requestHandlers"
)

// Engine wrapps everything required to request an index : an index, a request handler and an output formatter
type Engine struct {
	index          requestableIndexes.RequestableIndex
	requestHandler requestHandlers.RequestHandler
	outputFormater outputFormaters.OutputFormater
}

// NewEngine returns an Engine configured according to the input
func NewEngine(index requestableIndexes.RequestableIndex, requestType string, outputType string) *Engine {
	var requestHandler requestHandlers.RequestHandler
	switch requestType {
	case "and":
		requestHandler = requestHandlers.NewAndRequestHandler(index)
	case "binary":
		requestHandler = requestHandlers.NewBinaryRequestHandler(index)
	case "vectorial":
		requestHandler = requestHandlers.NewVectorizedRequestHandler(index)
	}
	var outputFormater outputFormaters.OutputFormater
	switch outputType {
	case "sorted":
		outputFormater = outputFormaters.NewSortDocsOutputFormater(index.GetDocIDToPath())
	case "dumb":
		outputFormater = outputFormaters.NewDumbOutputFormater()
	}
	return &Engine{
		index:          index,
		requestHandler: requestHandler,
		outputFormater: outputFormater,
	}
}

// Request outputs the response to a request to the user
func (engine *Engine) Request(request string) {
	start := time.Now()
	res := engine.requestHandler.Request(request, []string{})
	done := time.Now()
	elapsed := done.Sub(start)
	fmt.Printf("> Done computing the response to the request in %v", elapsed)
	engine.outputFormater.Output(res)
}

// GetRes returns the response to a request
func (engine *Engine) GetRes(request string, stopList []string) indexes.PostingList {
	res := *engine.requestHandler.Request(request, stopList)
	return res
}
