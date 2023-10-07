package rule

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/tr1v3r/rule/driver"
)

var _ Tree = new(tree[Rule])

// tree is a rule tree structure.
type tree[R Rule] struct {
	name string // node name

	mu       sync.RWMutex
	children map[string]Tree

	// current node rule
	ruleMu sync.RWMutex
	rule   string
	ops    []driver.Operator

	// for subtree
	driver driver.Driver
	level  int
}

func (t *tree[R]) build(rules ...R) error {
	for _, r := range sortRule(t.driver, rules) {
		if err := t.SetRule(r); err != nil {
			return fmt.Errorf("add rule error: %w", err)
		}
	}
	return nil
}

func (t *tree[R]) Name() string { return t.name }

// SetRule add a rule node to tree or update rule node.
// make rule tree grow
func (t *tree[R]) SetRule(r Rule) error {
	if level := t.driver.GetLevel(r.Path()); t.level == level { // check if level matched, include root node
		return t.updateRule(r.Operators()...)
	}
	return t.getChild(t.driver.GetNameByLevel(r.Path(), t.level+1)).SetRule(r)
}

// GetRule get a rule from tree.
func (t *tree[R]) GetRule(path string) string {
	if subTree := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); subTree != nil {
		return subTree.GetRule(path)
	}
	return t.getRule()
}

// HasNode check if has node in path
func (t *tree[R]) HasNode(path string) bool {
	if t.driver.GetLevel(path) == t.level { // check level
		return t.Name() == t.driver.GetNameByLevel(path, t.level)
	}
	if tree := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); tree != nil {
		return tree.HasNode(path)
	}
	return false
}

// DelNode delete a node from tree.
func (t *tree[R]) DelNode(path string) error {
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
func (t *tree[R]) ShowStruct() json.RawMessage {
	m := make(map[string]json.RawMessage)
	for _, v := range t.getChildren() {
		m[v.Name()] = v.ShowStruct()
	}
	d, _ := json.Marshal(m)
	return d
}

// GetOperators get all Operators.
func (t *tree[R]) GetOperators() []driver.Operator { return t.ops }

// deleteNode delete a node from tree.
func (t *tree[R]) deleteNode(name string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.children, name)

	return nil
}

// getChild get a child tree.
// if not found, create a new sub tree and return it
func (t *tree[R]) getChild(name string) (tree Tree) {
	if tree = t.pickChild(name); tree != nil {
		return tree
	}
	tree = t.newSubTree(name)
	t.graft(tree)
	return tree
}

func (t *tree[R]) getChildren() (children []Tree) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, v := range t.children {
		children = append(children, v)
	}
	return children
}

// pickChild get a child tree.
// if not found, return nil
func (t *tree[R]) pickChild(name string) Tree {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.children[name]
}

// graft add a child tree on current node
func (t *tree[R]) graft(subTree Tree) Tree {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.children[subTree.Name()] = subTree

	return t
}

// newSubTree create a new sub tree.
// name cannot be empty
func (t *tree[R]) newSubTree(name string) Tree {
	return &tree[R]{
		name: name,

		driver: t.driver,

		level:    t.level + 1,
		rule:     t.getRule(),
		children: make(map[string]Tree),
	}
}

// updateRule parse raw rule Operator to tree node.
func (t *tree[R]) updateRule(ops ...driver.Operator) error {
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
func (t *tree[R]) getRule() string {
	t.ruleMu.RLock()
	defer t.ruleMu.RUnlock()
	return t.rule
}

// sortRule sort rules by path.
func sortRule[R Rule](driver driver.Driver, rules []R) []R {
	by[R](func(x, y R) bool { return driver.GetLevel(x.Path()) < driver.GetLevel(y.Path()) }).Sort(rules)
	return rules
}
