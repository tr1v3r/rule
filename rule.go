package rule

import (
	"sort"

	"github.com/tr1v3r/rule/driver"
)

// rule raw rule for node
type rule struct {
	path      string
	operators []driver.Operator
}

func (r *rule) Path() string                 { return r.path }
func (r *rule) Operators() []driver.Operator { return r.operators }

// rules rules array
type rules[R Rule] []R

func (r rules[R]) Len() int      { return len(r) }
func (r rules[R]) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type sorter[R Rule] struct {
	rules[R]
	by by[R]
}

// Less is part of sort.Interface.
func (s *sorter[R]) Less(i, j int) bool { return s.by(s.rules[i], s.rules[j]) }

// By is the type of a "less" function that defines the ordering of its Rules arguments.
type by[R Rule] func(x, y R) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by by[R]) Sort(rules []R) {
	sort.Sort(&sorter[R]{
		rules: rules,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	})
}
