package rule

import (
	"encoding/json"

	"github.com/tr1v3r/rule/driver"
)

// Forest rule forest
type Forest interface {
	Register(...TreeBuilder)
	Build() Forest
	Append(...TreeBuilder) Forest
	Get(name string) Tree
	Set(tree Tree)
	Info() string
}

// Tree rule tree
type Tree interface {
	Name() string

	SetRule(Rule) error
	GetRule(path string) string

	HasNode(path string) bool
	DelNode(path string) error
	ShowStruct() json.RawMessage
	GetOperators() []driver.Operator
}

// Rule rule for tree
type Rule interface {
	Path() string
	Operators() []driver.Operator
}
