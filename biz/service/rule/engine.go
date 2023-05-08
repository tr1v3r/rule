package rule

import "sync"

// Forest rule forest
type Forest struct {
	mu sync.RWMutex
	m  map[string]*Tree

	builders []func() (name string, tree *Tree)
}

// Build
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
