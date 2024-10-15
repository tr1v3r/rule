package rule

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/tr1v3r/rule/driver"
)

var _ Tree = (*tree)(nil)

// tree is a rule tree structure.
type tree struct {
	name string // node name
	path string // nod path

	mu       sync.RWMutex
	children map[string]Tree

	// for subtree
	driver driver.Driver
	level  int

	// current node rule
	contentMu sync.RWMutex
	content   []byte

	// procs processor array
	// only set when tree build, only concurrent reads, so mutex is verbose
	procs []driver.Processor

	// Lazy Mode:
	// In Lazy Mode, tree nodes are not created or calculated during initialization.
	// Only the root node exists initially, and other nodes are dynamically initialized
	// and calculated when accessed. This mode employs lazy evaluation, saving memory
	// and computation resources, particularly useful in scenarios where only a subset
	// of the nodes are accessed.
	lazyMode bool
	// Instant Mode:
	// In Instant Mode, every time a node is accessed, its data is recalculated
	// and the entire path from the root to the accessed node is refreshed.
	// Even if the nodes have been previously created or calculated, they are
	// forcefully recalculated to ensure up-to-date data. This mode emphasizes
	// real-time computation, ideal for scenarios requiring frequent updates
	// and data consistency.
	instantMode bool

	procMu   sync.RWMutex
	realized bool
}

func (t *tree) lazy() *tree {
	t.lazyMode = true
	return t
}

func (t *tree) instant() *tree {
	t.instantMode = true
	return t
}

func (t *tree) build(rules ...Rule) error {
	for _, r := range byLevel(t.driver, rules) {
		if err := t.Set(r); err != nil {
			return fmt.Errorf("apply rule fail: %w", err)
		}
	}
	return nil
}

func (t *tree) Name() string { return t.name }
func (t *tree) Path() string { return t.path }

func (t *tree) Set(r Rule) error {
	if level := t.driver.GetLevel(r.Path()); t.level == level { // check if level matched, include root node
		return t.apply(r.Processors()...)
	}
	return t.getChild(t.driver.GetNameByLevel(r.Path(), t.level+1)).Set(r)
}

func (t *tree) Get(path string) ([]byte, error) {
	if t == nil {
		return nil, ErrNotExistsTree
	}

	if err := t.realize(t.procs); err != nil {
		return nil, fmt.Errorf("realize rule on %s fail: %w", t.Path(), err)
	}

	if child := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); child != nil {
		if child, ok := child.(*tree); ok {
			child.inherit(t)
		}
		return child.Get(path)
	}
	return t.get(), nil
}

// inherit set content by parent's content after check mode and realization
func (t *tree) inherit(parent *tree) {
	if t.lazyMode && t.needRealize() {
		t.set(parent.get())
	}
}

// Has check if has node in path
func (t *tree) Has(path string) bool {
	if t.driver.GetLevel(path) == t.level { // check level
		return t.Name() == t.driver.GetNameByLevel(path, t.level)
	}
	if tree := t.pickChild(t.driver.GetNameByLevel(path, t.level+1)); tree != nil {
		return tree.Has(path)
	}
	return false
}

// Del delete a node from tree.
func (t *tree) Del(path string) error {
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
func (t *tree) ShowStruct() []byte {
	m := make(map[string]json.RawMessage)
	for _, v := range t.getChildren() {
		m[v.Name()] = v.ShowStruct()
	}
	d, _ := json.Marshal(m)
	return d
}

// deleteNode delete a node from tree.
func (t *tree) deleteNode(name string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.children, name)

	return nil
}

// getChild get a child tree.
// if not found, create a new sub tree and return it
func (t *tree) getChild(name string) (tree Tree) {
	if tree = t.pickChild(name); tree != nil {
		return tree
	}
	tree = t.newSubTree(name)
	t.Graft(tree)
	return tree
}

func (t *tree) getChildren() (children []Tree) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, v := range t.children {
		children = append(children, v)
	}
	return children
}

// pickChild get a child tree.
// if not found, return nil
func (t *tree) pickChild(name string) Tree {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.children[name]
}

// Graft graft a sub tree
func (t *tree) Graft(child Tree) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.children[child.Name()] = child
}

// newSubTree create a new sub tree.
// name cannot be empty
func (t *tree) newSubTree(name string) Tree {
	return &tree{
		name: name,
		path: t.driver.AppendPath(t.path, name),

		driver:      t.driver,
		lazyMode:    t.lazyMode,
		instantMode: t.instantMode,

		level:    t.level + 1,
		content:  t.get(),
		children: make(map[string]Tree),
	}
}

// updateRule parse raw rule Processor to tree node.
func (t *tree) apply(procs ...driver.Processor) error {
	if t.lazyMode {
		t.procs = procs
		return nil
	}
	return t.realize(procs)
}

func (t *tree) realize(procs []driver.Processor) error {
	t.procMu.Lock()
	defer t.procMu.Unlock()
	if !t.instantMode && t.realized {
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

func (t *tree) set(rule []byte) {
	t.contentMu.Lock()
	defer t.contentMu.Unlock()
	t.content = rule
}

// get return current node rule.
func (t *tree) get() (rule []byte) {
	t.contentMu.RLock()
	defer t.contentMu.RUnlock()
	return t.content
}

func (t *tree) needRealize() bool {
	t.procMu.RLock()
	defer t.procMu.RUnlock()
	return !t.realized
}

// byLevel sort rules by path level
func byLevel[R Rule](driver driver.Driver, rules []R) []R {
	by[R](func(x, y R) bool { return driver.GetLevel(x.Path()) < driver.GetLevel(y.Path()) }).Sort(rules)
	return rules
}
