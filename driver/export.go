package driver

import (
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
	Realize(rule []byte, ops ...Processor) ([]byte, error)
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
	Process(before []byte) (after []byte, err error)

	// informatin
	Author() string
	CreatedAt() time.Time

	// Load load Processor from data
	Load([]byte) error
	// Save ...
	Save() []byte
}
