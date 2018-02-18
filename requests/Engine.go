package requests

import (
	"github.com/rmulton/riw_project/indexes/requestable"
	"github.com/rmulton/riw_project/requests/requestHandlers"
	"github.com/rmulton/riw_project/requests/outputFormaters"
)

type Engine struct {
	index requestable.RequestableIndex
	requestHandler requestHandlers.RequestHandler
	outputFormater outputFormaters.OutputFormater
}

func NewEngine(index indexes.RequestableIndex, requestType string, outputType string) *Engine {
	var requestHandler requestHandlers.RequestHandler
	switch requestType {
	case "and":
		requestHandler = requestHandlers.NewAndRequestHandler(index)
	case "binary":
		requestHandler = requestHandlers.NewBinaryRequestHandler(index)
	case "vectorial":
		requestHandler = requestHandlers.NewVectorizedRequestHandler(index)
	}
	var outputFormater outputFormater
	switch outputType {
	case "sorted":
		outputFormater = outputFormaters.NewSortDocsOutputFormater(index.GetDocIDToPath())
	case "dumb":
		outputFormater = outputFormaters.NewDumbOutputFormater()
	}
	return &Engine{
		index: index,
		requestHandler: requestHandler,
		outputFormater: outputFormater,
	}
}

func (engine *Engine) Request (request string) {
	res := engine.requestHandler.Request(request)
	engine.outputFormater.Output(res)
}

