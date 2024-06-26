package rule

import (
	"fmt"
	"sync"
	"time"

	"github.com/tr1v3r/pkg/guard"
	"github.com/tr1v3r/pkg/log"
)

var _ Forest = (*forest)(nil)

// TreeBuilder tree build method
type TreeBuilder func() Tree

// forest rules forest
type forest struct {
	mu sync.RWMutex
	m  map[string]Tree

	bMu      sync.RWMutex
	builders []TreeBuilder
	builderM map[string]TreeBuilder
}

// Register register tree builder
func (f *forest) Register(builders ...TreeBuilder) { f.appendBuilders(builders...) }

// BindTreeBuilder bind tree and builder
func (f *forest) BindTreeBuilder(name string, builder TreeBuilder) {
	f.bMu.Lock()
	defer f.bMu.Unlock()
	f.builderM[name] = builder
}

// Refresh refresh rule forest
func (f *forest) Refresh(interval ...time.Duration) {
	if len(interval) == 0 {
		f.Build()
		return
	}
	for range time.Tick(interval[0]) { // nolint
		f.Build()
	}
}

// RefreshTree refresh tree
func (f *forest) RefreshTree(name string) {
	if build := f.getBuilder(name); build != nil {
		f.Set(build())
	}
}

// Build all trees in forest
func (f *forest) Build() Forest {
	for _, build := range f.getBuilders() {
		if tree := build(); tree != nil {
			f.Set(tree)
			f.BindTreeBuilder(tree.Name(), build)
		}
	}
	return f
}

// Append append tree and builder to forest
func (f *forest) Append(builders ...TreeBuilder) Forest {
	for _, builder := range f.getBuilders() {
		f.Register(builder)
		f.Set(builder())
	}
	return f
}

// Get get rule tree by name
func (f *forest) Get(name string) Tree {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.m[name]
}

// Set set rule tree
func (f *forest) Set(tree Tree) {
	if tree == nil {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.m == nil {
		f.m = make(map[string]Tree, 16)
	}
	f.m[tree.Name()] = tree
}

// GetVal get value from tree
func (f *forest) GetVal(treeName, path string) (rule []byte, err error) {
	tree := f.Get(treeName)
	if tree == nil {
		return nil, ErrNotExistsTree
	}
	return tree.Get(path)
}

// Info ...
func (f *forest) Info() string {
	return fmt.Sprintf("forest got %d tree: %s", f.count(), f.names())
}

func (f *forest) count() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.m)
}

func (f *forest) names() (names []string) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for name := range f.m {
		names = append(names, name)
	}
	return
}

func (f *forest) getBuilders() []TreeBuilder {
	f.bMu.RLock()
	defer f.bMu.RUnlock()
	return f.builders
}
func (f *forest) appendBuilders(builders ...TreeBuilder) {
	var wrappedBuilders = make([]TreeBuilder, len(builders))
	for i, builder := range builders {
		wrappedBuilders[i] = func() Tree {
			defer func() {
				if e := recover(); e != nil {
					log.Error("build tree panic: %s, stack: %s", e, guard.CatchStack())
				}
			}()
			return builder()
		}
	}

	f.bMu.Lock()
	defer f.bMu.Unlock()
	f.builders = append(f.builders, wrappedBuilders...)
}
func (f *forest) getBuilder(name string) TreeBuilder {
	f.bMu.RLock()
	defer f.bMu.RUnlock()
	return f.builderM[name]
}
