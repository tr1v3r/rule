package rule

import (
	"github.com/tr1v3r/rule/driver"
)

// NewForest build a new forest and return it
func NewForest(builders ...TreeBuilder) Forest {
	return (&forest{
		m:        make(map[string]Tree, len(builders)),
		builders: builders,
		builderM: make(map[string]TreeBuilder, len(builders)),
	}).Build()
}

// NewJSONTree build a json rule tree
func NewJSONTree[R Rule](name, template string, rules ...R) (Tree, error) {
	return NewTree(driver.NewJSONDriver(), name, template, rules...)
}

// NewLazyJSONTree build a lazy json rule tree
func NewLazyJSONTree[R Rule](name, template string, rules ...R) (Tree, error) {
	return NewLazyTree(driver.NewJSONDriver(), name, template, rules...)
}

// NewYAMLTree build a yaml rule tree
func NewYAMLTree[R Rule](name, template string, rules ...R) (Tree, error) {
	return NewTree(driver.NewYAMLDriver(), name, template, rules...)
}

// NewLazyYAMLTree build a yaml rule tree
func NewLazyYAMLTree[R Rule](name, template string, rules ...R) (Tree, error) {
	return NewLazyTree(driver.NewYAMLDriver(), name, template, rules...)
}

// NewTree build a rule tree.
func NewTree[R Rule](driver driver.Driver, name, template string, rules ...R) (Tree, error) {
	tree := newTree(driver, name, template, rules)
	if err := tree.build(rules...); err != nil {
		return nil, err
	}
	return tree, nil
}

// NewLazyTree build a lazy rule tree.
func NewLazyTree[R Rule](driver driver.Driver, name, template string, rules ...R) (Tree, error) {
	tree := newTree(driver, name, template, rules).lazy()
	if err := tree.build(rules...); err != nil {
		return nil, err
	}
	return tree, nil
}

func newTree[R Rule](diver driver.Driver, name, template string, rules []R) *tree[R] {
	return (&tree[R]{
		name: name,

		rule:     []byte(template),
		driver:   diver,
		children: make(map[string]Tree),
	})
}

func NewRule(path string, Processors ...driver.Processor) Rule { return &rule{path, Processors} }
