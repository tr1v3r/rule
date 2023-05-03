package rule

import "sort"

type Forest struct {
	m map[string]*Tree
}

type Driver interface {
	// GetLevel get level from path
	GetLevel(path string) (level int)

	// GetNameByLevel get node name from path by level
	// return empty string if level is out of range
	GetNameByLevel(path string, level int) (name string)

	// CalcRule calc rule
	CalcRule(template string, op *Rule) (string, error)
}

// Rule raw rule
type Rule struct {
	Path    string
	Operate any
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

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(x, y *Rule) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(rules []*Rule) {
	sort.Sort(&ruleSorter{
		RuleArray: rules,
		by:        by, // The Sort method's receiver is the function (closure) that defines the sort order.
	})
}
