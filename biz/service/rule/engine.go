package rule

import (
	"sort"

	"github.com/riverchu/rule/biz/service/driver"
)

type Forest struct {
	_ map[string]*Tree
}

// Rule raw rule
type Rule struct {
	Path      string
	Operators []driver.Operator
}

// RuleArray rules array
type RuleArray []*Rule

func (r RuleArray) Len() int      { return len(r) }
func (r RuleArray) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type ruleSorter struct {
	RuleArray
	by By
}

// Less is part of sort.Interface.
func (s *ruleSorter) Less(i, j int) bool { return s.by(s.RuleArray[i], s.RuleArray[j]) }

// By is the type of a "less" function that defines the ordering of its Rules arguments.
type By func(x, y *Rule) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(rules []*Rule) {
	sort.Sort(&ruleSorter{
		RuleArray: rules,
		by:        by, // The Sort method's receiver is the function (closure) that defines the sort order.
	})
}
