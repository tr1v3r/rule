package rule

import (
	"encoding/json"
	"time"

	"github.com/tr1v3r/rule/driver"
)

// Forest rule forest
type Forest interface {
	Register(...TreeBuilder)
	Append(...TreeBuilder) Forest

	Build() Forest
	// Refresh refresh Forest
	// interval is optional and only first value is useful when set
	// blocks when the interval is set
	Refresh(interval ...time.Duration)
	// RefreshTree refresh specified tree
	RefreshTree(name string)

	Get(name string) Tree
	Set(tree Tree)

	Info() string
}

// Tree rule tree
type Tree interface {
	Name() string

	SetRule(Rule) error
	GetRule(path string) []byte

	HasNode(path string) bool
	DelNode(path string) error
	ShowStruct() json.RawMessage
	GetProcessors() []driver.Processor
}

// Rule rule for tree
type Rule interface {
	Path() string
	Processors() []driver.Processor
}
