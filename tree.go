package rule

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/tr1v3r/rule/driver"
)

var _ Tree = (*tree[Rule])(nil)

// tree is a rule tree structure.
type tree[R Rule] struct {
	name string // node name
	path string // nod path

	mu       sync.RWMutex
	children map[string]Tree

	// current node rule
	ruleMu sync.RWMutex
	rule   []byte

	// procs processor array
	// only set when tree build, only concurrent reads, so mutex is verbose
	procs []driver.Processor

	// for subtree
	driver driver.Driver
	level  int

	lazyMode bool

	procMu   sync.RWMutex
	realized bool
}

func (t *tree[R]) lazy() *tree[R] {
	t.lazyMode = true
	return t
}

func (t *tree[R]) build(rules ...R) error {
	for _, r := range byLevel(t.driver, rules) {
		if err := t.Set(r); err != nil {
			return fmt.Errorf("apply rule fail: %w", err)
		}
	}
	return nil
}

func (t *tree[R]) Name() string { return t.name }
func (t *tree[R]) Path() string { return t.path }

func (t *tree[R]) Set(r Rule) error {
	if level := t.driver.GetLevel(r.Path()); t.level == level { // check if level matched, include root node
		return t.apply(r.Processors()...)
	}
	return t.getChild(t.driver.GetNameByLevel(r.Path(), t.level+1)).Set(r)
}

func (t *tree[R]) Get(path string) ([]byte, error) {
	if t == nil {
		return nil, ErrNilTree
	}

	if t.lazyMode && t.needRealize() {
		if err := t.realize(t.procs); err != nil {
			return nil, fmt.Errorf("realize rule on %s fail: %w", t.Path(), err)
		}
	}

	if child := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); child != nil {
		if t.lazyMode {
			if tree, ok := child.(*tree[R]); ok && tree.needRealize() {
				tree.set(t.get())
			}
		}
		return child.Get(path)
	}
	return t.get(), nil
}

// Has check if has node in path
func (t *tree[R]) Has(path string) bool {
	if t.driver.GetLevel(path) == t.level { // check level
		return t.Name() == t.driver.GetNameByLevel(path, t.level)
	}
	if tree := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); tree != nil {
		return tree.Has(path)
	}
	return false
}

// Del delete a node from tree.
func (t *tree[R]) Del(path string) error {
	if level := t.driver.GetLevel(path); level == 0 {
		return fmt.Errorf("root node can not be deleted")
	} else if t.level+1 == level {
		return t.deleteNode(t.driver.GetNameByLevel(path, level))
	}

	if child := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); child != nil {
		return child.Del(path)
	}
	return nil
}

// ShowStruct return tree struct.
func (t *tree[R]) ShowStruct() []byte {
	m := make(map[string]json.RawMessage)
	for _, v := range t.getChildren() {
		m[v.Name()] = v.ShowStruct()
	}
	d, _ := json.Marshal(m)
	return d
}

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
	t.Graft(tree)
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

// Graft graft a sub tree
func (t *tree[R]) Graft(child Tree) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.children[child.Name()] = child
}

// newSubTree create a new sub tree.
// name cannot be empty
func (t *tree[R]) newSubTree(name string) Tree {
	return &tree[R]{
		name: name,
		path: t.driver.AppendPath(t.path, name),

		driver:   t.driver,
		lazyMode: t.lazyMode,

		level:    t.level + 1,
		rule:     t.get(),
		children: make(map[string]Tree),
	}
}

// updateRule parse raw rule Processor to tree node.
func (t *tree[R]) apply(procs ...driver.Processor) error {
	t.procs = procs

	if t.lazyMode {
		return nil
	}
	return t.realize(procs)
}

func (t *tree[R]) realize(procs []driver.Processor) error {
	t.procMu.Lock()
	defer t.procMu.Unlock()
	if t.realized {
		return nil
	}

	rule, err := t.driver.Realize(t.get(), procs...)
	if err != nil {
		return fmt.Errorf("realize rule fail: %w", err)
	}
	t.set(rule)

	t.realized = true
	return nil
}

func (t *tree[R]) set(rule []byte) {
	t.ruleMu.Lock()
	defer t.ruleMu.Unlock()
	t.rule = rule
}

// get return current node rule.
func (t *tree[R]) get() (rule []byte) {
	t.ruleMu.RLock()
	defer t.ruleMu.RUnlock()
	return t.rule
}

func (t *tree[R]) needRealize() bool {
	t.procMu.RLock()
	defer t.procMu.RUnlock()
	return t.realized
}

// byLevel sort rules by path level
func byLevel[R Rule](driver driver.Driver, rules []R) []R {
	by[R](func(x, y R) bool { return driver.GetLevel(x.Path()) < driver.GetLevel(y.Path()) }).Sort(rules)
	return rules
}
