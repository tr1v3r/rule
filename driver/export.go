package driver

import (
	"context"
	"time"
)

// Driver driver interface
type Driver interface {
	Name() string

	PathParser
	Realizer
	Modem
}

// PathParser path parser
type PathParser interface {
	// GetLevel get level from path
	GetLevel(path string) (level int)
	// GetNameByLevel get node name from path by level
	// return empty string if level is out of range
	GetNameByLevel(path string, level int) (name string)
	// AppendPath append path
	AppendPath(path, name string) (newPath string)
}

// Realizer realize rule
type Realizer interface {
	// Realize realize rule
	Realize(rc *RuleContext, rule []byte, ops ...Processor) ([]byte, error)
}

// Modem Processors modem
type Modem interface {
	// ProcessorsForSave get Processors data for save
	Marshal(...Processor) ([]byte, error)

	// LoadProcessors load Processors from data
	Unmarshal(data []byte) ([]Processor, error)
}

// Processor rule processor
type Processor interface {
	// Path return target tree path, not necessary
	Path() string

	// Type return processor type
	Type() string
	// Process do process rule
	Process(rc *RuleContext, before []byte) (after []byte, err error)

	// informatin
	Author() string
	CreatedAt() time.Time

	// Load load Processor from data
	Load([]byte) error
	// Save ...
	Save() []byte
}

// RuleContext carries runtime information for dynamic rule construction.
type RuleContext struct {
	context.Context
	// TreePath is the path of the current tree node being processed.
	TreePath string
	// Params holds key-value pairs from the request, for template interpolation etc.
	Params map[string]string
	// ParentContent holds the realized content of the parent node.
	ParentContent []byte
}

func (rc *RuleContext) Ctx() context.Context { return rc.Context }
