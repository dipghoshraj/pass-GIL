package main

import (
	"errors"
	"goengine/engine"
	// flatbuffers "github.com/google/flatbuffers/go"
)

const (
	HashHandler uint32 = 1
)

type handlerFn func(req *engine.Request) ([]byte, error)

var registry = map[uint32]handlerFn{
	HashHandler: handleHashRequest,
}

func dispatchFunction(input []byte) ([]byte, error) {
	req := engine.GetRootAsRequest(input, 0)

	if handler, ok := registry[req.FuncId()]; ok {
		return handler(req)
	}
	return nil, errors.New("unknown function ID")
}
