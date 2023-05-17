package rule

import (
	"fmt"
	"sync"
)

// TreeBuilder tree build method
type TreeBuilder func() (name string, tree *Tree)

// NewForest build a new forest and return it
func NewForest(builders ...TreeBuilder) *Forest {
	return (&Forest{m: make(map[string]*Tree, 16), builders: builders}).Build()
}

// Forest rules forest
type Forest struct {
	mu sync.RWMutex
	m  map[string]*Tree

	bMu      sync.RWMutex
	builders []TreeBuilder
}

// Register register tree builder
func (f *Forest) Register(builders ...TreeBuilder) {
	f.addBuilders(builders...)
}

// Build all trees in forest
func (f *Forest) Build() *Forest {
	for _, builder := range f.getBuilders() {
		f.Set(builder())
	}
	return f
}

// Append append tree and builder to forest
func (f *Forest) Append(builders ...TreeBuilder) *Forest {
	for _, builder := range f.getBuilders() {
		f.Register(builder)
		f.Set(builder())
	}
	return f
}

// Get get rule tree by name
func (f *Forest) Get(name string) *Tree {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.m[name]
}

// Set set rule tree
func (f *Forest) Set(name string, tree *Tree) {
	if tree == nil {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.m == nil {
		f.m = make(map[string]*Tree, 16)
	}
	f.m[name] = tree
}

// Info ...
func (f *Forest) Info() string {
	return fmt.Sprintf("forest got %d tree: %s", f.count(), f.names())
}

func (f *Forest) count() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.m)
}

func (f *Forest) names() (names []string) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for name := range f.m {
		names = append(names, name)
	}
	return
}

func (f *Forest) getBuilders() []TreeBuilder {
	f.bMu.RLock()
	defer f.bMu.RUnlock()
	return f.builders
}
func (f *Forest) addBuilders(builders ...TreeBuilder) {
	f.bMu.Lock()
	defer f.bMu.Unlock()
	f.builders = append(f.builders, builders...)
}
