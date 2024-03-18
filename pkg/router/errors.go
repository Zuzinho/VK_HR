package router

import (
	"fmt"
)

type NoHandlerError struct {
	pair Pair
}

func newNoHandlerError(pair Pair) NoHandlerError {
	return NoHandlerError{
		pair: pair,
	}
}

func (err NoHandlerError) Error() string {
	return fmt.Sprintf("no handler for path '%s' and method '%s'", err.pair.path, err.pair.method)
}

var (
	NoHandlerErr = NoHandlerError{}
)
