package rule

import (
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
	// GetVal get value from tree
	GetVal(treeName, path string) (rule []byte, err error)

	Info() string
}

// Tree rule tree
type Tree interface {
	// Name return tree name
	Name() string
	// Path return tree path in root tree
	// return "" when tree is root tree
	Path() string

	// Set add a rule node to tree or update rule node.
	Set(Rule) error
	// Get query rule from tree by path
	Get(path string) (rule []byte, err error)

	// Has check if has tree node for path
	Has(path string) bool
	// Del delete tree node in tree
	Del(path string) error

	// Graft graft a sub tree
	Graft(Tree)

	// ShowStruct return tree info
	ShowStruct() []byte
}

// Rule rule for tree
type Rule interface {
	// return path of rule
	Path() string
	// return processor on this rule
	Processors() []driver.Processor
}
