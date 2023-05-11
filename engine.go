package rule

import "sync"

// TreeBuilder tree build method
type TreeBuilder func() (name string, tree *Tree)

// Forest rules forest
type Forest struct {
	mu sync.RWMutex
	m  map[string]*Tree

	builders []TreeBuilder
}

// Register register tree builder
func (f *Forest) Register(builder TreeBuilder) {
	f.builders = append(f.builders, builder)
}

// Build all trees in forest
func (f *Forest) Build() {
	for _, builder := range f.builders {
		f.Set(builder())
	}
}

// Get get rule tree by name
func (f *Forest) Get(name string) *Tree {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.m[name]
}

// Set set rule tree
func (f *Forest) Set(name string, tree *Tree) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.m[name] = tree
}
