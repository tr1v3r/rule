package driver

import "errors"

var (
	// ErrSerializeNotSupport not support serialize error
	ErrSerializeNotSupport = errors.New("operator not support serialize")
)
