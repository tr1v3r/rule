package rule

import (
	"sort"

	"github.com/riverchu/rule/driver"
)

// Rule raw rule
type Rule struct {
	Path      string
	Operators []driver.Operator
}

// rules rules array
type rules []*Rule

func (r rules) Len() int      { return len(r) }
func (r rules) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type sorter struct {
	rules
	by by
}

// Less is part of sort.Interface.
func (s *sorter) Less(i, j int) bool { return s.by(s.rules[i], s.rules[j]) }

// By is the type of a "less" function that defines the ordering of its Rules arguments.
type by func(x, y *Rule) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by by) Sort(rules []*Rule) {
	sort.Sort(&sorter{
		rules: rules,
		by:    by, // The Sort method's receiver is the function (closure) that defines the sort order.
	})
}
