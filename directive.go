package ivy

import (
	"sort"

	"github.com/tr1v3r/ivy/driver"
)

var _ Directive = (*directive)(nil)

// directive is a path + processors pair that defines a transformation on the tree.
type directive struct {
	path       string
	processors []driver.Processor
}

func (d *directive) Path() string                   { return d.path }
func (d *directive) Processors() []driver.Processor { return d.processors }

// directives is a sortable slice of Directive.
type directives[R Directive] []R

func (d directives[R]) Len() int      { return len(d) }
func (d directives[R]) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

type sorter[R Directive] struct {
	directives[R]
	by by[R]
}

// Less is part of sort.Interface.
func (s *sorter[R]) Less(i, j int) bool { return s.by(s.directives[i], s.directives[j]) }

// By is the type of a "less" function that defines the ordering of its Directive arguments.
type by[R Directive] func(x, y R) bool

// Sort sorts the argument slice according to the function.
func (b by[R]) Sort(directives []R) {
	sort.Sort(&sorter[R]{
		directives: directives,
		by:         b,
	})
}
