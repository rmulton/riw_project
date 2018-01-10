package requests

import (
)

type Engine struct {
	index *Index
	requestHandler requestHandler
	outputFormater outputFormater
}

func NewEngine(folder string) *Engine {
	index := NewIndex(folder)
	requestHandler := NewVectorizedRequestHandler()
	outputFormater := NewSortDocsOutputFormater()
	return &Engine{
		index: index,
		requestHandler: requestHandler,
		outputFormater: outputFormater,
	}
}

func (engine *Engine) Request (request string) {
	res := engine.requestHandler.request(request, engine.index)
	engine.outputFormater.output(res)
}

