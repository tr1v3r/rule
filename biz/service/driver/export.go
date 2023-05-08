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
	CalcRule(template string, ops ...Operator) (string, error)
}

// Modem operators modem
type Modem interface {
	// OperatorsForSave get operators data for save
	Marshal(...Operator) []byte

	// LoadOperators load operators from data
	Unmarshal(data []byte) ([]Operator, error)
}

// Operator rule operator
type Operator interface {
	Type() string
	Operate(before string) (after string, err error)

	// informatin
	Author() string
	CreatedAt() time.Time
	// Path return operate path, not necessary
	Path() string

	// Load load operator from data
	Load([]byte) error
	// Save ...
	Save() []byte
}
