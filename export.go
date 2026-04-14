package rule

import (
	"time"

	"golang.org/x/time/rate"

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
	// GetValWithContext retrieves a value with runtime context.
	GetValWithContext(rc *driver.RuleContext, treeName, path string) (rule []byte, err error)

	Info() string

	// SetRateLimit sets a global rate limit for all GetVal calls as a fallback.
	SetRateLimit(r rate.Limit, burst int)
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
	// GetWithContext retrieves rule data with runtime context for dynamic construction.
	GetWithContext(rc *driver.RuleContext, path string) (rule []byte, err error)

	// Has check if has tree node for path
	Has(path string) bool
	// Del delete tree node in tree
	Del(path string) error

	// Graft graft a sub tree
	Graft(Tree)

	// ShowStruct return tree info
	ShowStruct() []byte

	// SetRateLimit sets a rate limit for Get calls on this tree.
	SetRateLimit(r rate.Limit, burst int)

	// SetFallback sets a processor to handle cases where path resolution
	// cannot find a matching child node.
	SetFallback(proc driver.Processor)
}

// Rule rule for tree
type Rule interface {
	// return path of rule
	Path() string
	// return processor on this rule
	Processors() []driver.Processor
}
