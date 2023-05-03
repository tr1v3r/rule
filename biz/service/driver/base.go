package driver

import "time"

// Driver driver interface
type Driver interface {
	Name() string

	PathDriver
	RuleDriver
	OperatorDriver
}

type PathDriver interface {
	// GetLevel get level from path
	GetLevel(path string) (level int)
	// GetNameByLevel get node name from path by level
	// return empty string if level is out of range
	GetNameByLevel(path string, level int) (name string)
}

type RuleDriver interface {
	// CalcRule calc rule
	CalcRule(template string, ops ...Operator) (string, error)
}

type OperatorDriver interface {
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
