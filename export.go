package ivy

import (
	"time"

	"golang.org/x/time/rate"

	"github.com/tr1v3r/ivy/driver"
)

// Forest manages a collection of trees.
type Forest interface {
	Register(...TreeBuilder)
	Append(...TreeBuilder) Forest

	Build() Forest
	// Refresh refreshes all trees.
	// interval is optional and only first value is useful when set.
	// Blocks when the interval is set.
	Refresh(interval ...time.Duration)
	// RefreshTree refreshes the specified tree.
	RefreshTree(name string)

	Get(name string) Tree
	Set(tree Tree)
	// GetVal retrieves a value from the named tree at the given path.
	GetVal(treeName, path string) (val []byte, err error)
	// GetValWithContext retrieves a value with runtime context.
	GetValWithContext(rc *driver.RuleContext, treeName, path string) (val []byte, err error)

	Info() string

	// SetRateLimit sets a global rate limit for all GetVal calls as a fallback.
	SetRateLimit(r rate.Limit, burst int)

	// SetDefaultContext sets the default RuleContext on all trees in the forest.
	SetDefaultContext(rc *driver.RuleContext)
}

// Tree is a hierarchical node structure with path-based access.
type Tree interface {
	// Name returns the tree name.
	Name() string
	// Path returns the tree path in the root tree.
	// Returns "" when the tree is the root tree.
	Path() string

	// Set adds or updates a directive on the tree.
	Set(Directive) error
	// Get retrieves the value at the given path.
	Get(path string) (val []byte, err error)
	// GetWithContext retrieves a value with runtime context for dynamic construction.
	GetWithContext(rc *driver.RuleContext, path string) (val []byte, err error)

	// Has checks if a node exists at the given path.
	Has(path string) bool
	// Del deletes a node at the given path.
	Del(path string) error

	// Graft attaches a sub-tree.
	Graft(Tree)

	// ShowStruct returns the tree structure as JSON.
	ShowStruct() []byte

	// SetRateLimit sets a rate limit for Get calls on this tree.
	SetRateLimit(r rate.Limit, burst int)

	// SetFallback sets a processor to handle cases where path resolution
	// cannot find a matching child node.
	SetFallback(proc driver.Processor)

	// SetDefaultContext sets the default RuleContext used by realize when
	// no request-scoped context is provided (e.g. via Get or during build).
	SetDefaultContext(rc *driver.RuleContext)
}

// Directive defines a path and the processors to apply at that path.
type Directive interface {
	// Path returns the path this directive applies to.
	Path() string
	// Processors returns the processors for this directive.
	Processors() []driver.Processor
}
