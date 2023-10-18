package driver

import "time"

// Driver driver interface
type Driver interface {
	Name() string

	PathParser
	Calculator
	Modem
}

// PathParser path parser
type PathParser interface {
	// GetLevel get level from path
	GetLevel(path string) (level int)
	// GetNameByLevel get node name from path by level
	// return empty string if level is out of range
	GetNameByLevel(path string, level int) (name string)
}

// Calculator rule calculator
type Calculator interface {
	// CalcRule calc rule
	CalcRule(template string, ops ...Processor) (string, error)
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
	Process(before string) (after string, err error)

	// informatin
	Author() string
	CreatedAt() time.Time

	// Load load Processor from data
	Load([]byte) error
	// Save ...
	Save() []byte
}
