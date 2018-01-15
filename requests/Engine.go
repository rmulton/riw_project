package requests

import (
	"../indexes"
)

type Engine struct {
	index indexes.RequestableIndex
	requestHandler requestHandler
	outputFormater outputFormater
}

func NewEngine(index indexes.RequestableIndex) *Engine {
	// index := IndexFromFolder(folder)
	requestHandler := NewVectorizedRequestHandler(index)
	outputFormater := NewSortDocsOutputFormater()
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

