package rule

import (
	"fmt"

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

// NewTileTree build a tree with tile children
func NewTileTree[R Rule](name, template string, rules ...R) (Tree, error) {
	return NewTree[R](driver.NewTileDriver(), name, template, rules...)
}

// NewLazyTileTree build a tree with tile children in lazy mode
func NewLazyTileTree[R Rule](name, template string, rules ...R) (Tree, error) {
	return NewLazyTree[R](driver.NewTileDriver(), name, template, rules...)
}

// NewTree build a rule tree.
func NewTree[R Rule](driver driver.Driver, name, template string, rules ...R) (Tree, error) {
	return buildTree(newTree[R](driver, name, template), toI(rules...)...)
}

// NewLazyTree build a lazy rule tree.
func NewLazyTree[R Rule](driver driver.Driver, name, template string, rules ...R) (Tree, error) {
	return buildTree(newTree[R](driver, name, template).lazy(), toI(rules...)...)
}

func newTree[R Rule](diver driver.Driver, name, template string) *tree {
	return (&tree{
		name: name,

		content:  []byte(template),
		driver:   diver,
		children: make(map[string]Tree),
	})
}
func buildTree(tree *tree, rules ...Rule) (Tree, error) {
	if err := tree.build(rules...); err != nil {
		return nil, fmt.Errorf("build tree fail: %w", err)
	}
	return tree, nil
}
func toI[R Rule](rules ...R) (ruleArray []Rule) {
	for _, rule := range rules {
		ruleArray = append(ruleArray, rule)
	}
	return
}

func NewRule(path string, Processors ...driver.Processor) Rule { return &rule{path, Processors} }
