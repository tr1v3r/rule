package rule

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/tr1v3r/rule/driver"
)

// NewJSONTree ...
func NewJSONTree(name, template string, rules ...*Rule) (*Tree, error) {
	return NewTree(driver.NewJSONDriver(), name, template, rules...)
}

// NewYAMLTree ...
func NewYAMLTree(name, template string, rules ...*Rule) (*Tree, error) {
	return NewTree(driver.NewYAMLDriver(), name, template, rules...)
}

// NewTree build a rule tree.
func NewTree(diver driver.Driver, name, template string, rules ...*Rule) (*Tree, error) {
	tree := &Tree{
		Name: name,

		rule:     template,
		driver:   diver,
		children: make(map[string]*Tree),
	}
	for _, r := range tree.sortRule(rules) {
		if err := tree.AddRule(r); err != nil {
			return nil, fmt.Errorf("add rule error: %w", err)
		}
	}
	return tree, nil
}

// Tree is a rule tree structure.
type Tree struct {
	Name string // node name

	mu       sync.RWMutex
	children map[string]*Tree

	// current node rule
	ruleMu sync.RWMutex
	rule   string
	ops    []driver.Operator

	driver driver.Driver
	level  int
}

// AddRule add a rule node to tree or update rule node.
// make rule tree grow
func (t *Tree) AddRule(r *Rule) error {
	if level := t.driver.GetLevel(r.Path); t.level == level { // check if level matched, include root node
		return t.updateRule(r.Operators...)
	}
	return t.getChild(t.driver.GetNameByLevel(r.Path, t.level+1)).AddRule(r)
}

// GetRule get a rule from tree.
func (t *Tree) GetRule(path string) string {
	if subTree := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); subTree != nil {
		return subTree.GetRule(path)
	}
	return t.getRule()
}

// HasNode check if has node in path
func (t *Tree) HasNode(path string) bool {
	if t.driver.GetLevel(path) == t.level { // check level
		return t.Name == t.driver.GetNameByLevel(path, t.level)
	}
	if tree := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); tree != nil {
		return tree.HasNode(path)
	}
	return false
}

// DelNode delete a node from tree.
func (t *Tree) DelNode(path string) error {
	if level := t.driver.GetLevel(path); level == 0 {
		return fmt.Errorf("root node can not be deleted")
	} else if t.level+1 == level {
		return t.deleteNode(t.driver.GetNameByLevel(path, level))
	}

	if subTree := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); subTree != nil {
		return subTree.DelNode(path)
	}
	return nil
}

// ShowStruct return tree struct.
func (t *Tree) ShowStruct() json.RawMessage {
	m := make(map[string]json.RawMessage)
	for _, v := range t.getChildren() {
		m[v.Name] = v.ShowStruct()
	}
	d, _ := json.Marshal(m)
	return d
}

// GetOperators get all Operators.
func (t *Tree) GetOperators() []driver.Operator { return t.ops }

// deleteNode delete a node from tree.
func (t *Tree) deleteNode(name string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.children, name)

	return nil
}

// getChild get a child tree.
// if not found, create a new sub tree and return it
func (t *Tree) getChild(name string) (tree *Tree) {
	if tree = t.pickChild(name); tree != nil {
		return tree
	}
	tree = t.newSubTree(name)
	t.graft(tree)
	return tree
}

func (t *Tree) getChildren() (children []*Tree) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, v := range t.children {
		children = append(children, v)
	}
	return children
}

// pickChild get a child tree.
// if not found, return nil
func (t *Tree) pickChild(name string) *Tree {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.children[name]
}

// graft add a child tree on current node
func (t *Tree) graft(subTree *Tree) *Tree {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.children[subTree.Name] = subTree

	return t
}

// newSubTree create a new sub tree.
// name cannot be empty
func (t *Tree) newSubTree(name string) *Tree {
	return &Tree{
		Name: name,

		driver: t.driver,

		level:    t.level + 1,
		rule:     t.getRule(),
		children: make(map[string]*Tree),
	}
}

// updateRule parse raw rule Operator to tree node.
func (t *Tree) updateRule(ops ...driver.Operator) error {
	rule, err := t.driver.CalcRule(t.getRule(), ops...)
	if err != nil {
		return fmt.Errorf("calculate rule fail: %w", err)
	}

	t.ruleMu.Lock()
	defer t.ruleMu.Unlock()
	t.rule = rule
	t.ops = append(t.ops, ops...)

	return nil
}

// getRule return current node rule.
func (t *Tree) getRule() string {
	t.ruleMu.RLock()
	defer t.ruleMu.RUnlock()
	return t.rule
}

// sortRule sort rules by path.
func (t *Tree) sortRule(rules []*Rule) []*Rule {
	by(func(x, y *Rule) bool { return t.driver.GetLevel(x.Path) < t.driver.GetLevel(y.Path) }).Sort(rules)
	return rules
}
