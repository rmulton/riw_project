package requests

import (
	"../indexes"
)

type Engine struct {
	index indexes.RequestableIndex
	requestHandler requestHandler
	outputFormater outputFormater
}

func NewEngine(index indexes.RequestableIndex, requestType string, outputType string) *Engine {
	var requestHandler requestHandler
	switch requestType {
	case "and":
		requestHandler = NewAndRequestHandler(index)
	case "binary":
		requestHandler = NewBinaryRequestHandler(index)
	case "vectorial":
		requestHandler = NewVectorizedRequestHandler(index)
	}
	var outputFormater outputFormater
	switch outputType {
	case "sorted":
		// TODO: get saved folder path in argument or config
		outputFormater = NewSortDocsOutputFormater("./saved")
	case "dumb":
		outputFormater = NewDumbOutputFormater()
	}
	return &Engine{
		index: index,
		requestHandler: requestHandler,
		outputFormater: outputFormater,
	}
}

func (engine *Engine) Request (request string) {
	res := engine.requestHandler.request(request)
	engine.outputFormater.output(res)
}

