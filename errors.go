package rule

import (
	"errors"
)

var (
	// ErrNotExistsTree tree not exists
	ErrNotExistsTree = errors.New("tree not exists")
	// ErrRateLimited rate limited
	ErrRateLimited = errors.New("rate limited")
)
