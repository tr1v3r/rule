package driver

import (
	"errors"
)

var (
	// ErrSerializeNotSupport not support serialize error
	ErrSerializeNotSupport = errors.New("Processor not support serialize")
)
