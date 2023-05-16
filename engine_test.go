package rule_test

import (
	"encoding/json"
	"testing"

	"github.com/riverchu/rule"
	"github.com/riverchu/rule/driver"
)

func TestBuildTree(t *testing.T) {
	var rules = []*rule.Rule{
		{Path: "/a/b/c/d", Operators: nil},
		{Path: "/a/b/c", Operators: nil},
		{Path: "/", Operators: nil},
		{Path: "x/y/z", Operators: nil},
		{Path: "/a/b/m", Operators: nil},
	}

	tree, err := rule.NewTree(&struct {
		driver.PathParser
		driver.StdCalculator
		driver.DummyOperatorModem
		driver.DummyDriver // just provide a method Name
	}{PathParser: driver.SlashPathParser},
		"unit_test", "{}", rules...)
	if err != nil {
		t.Errorf("build tree fail: %s", err)
		return
	}

	result, _ := json.MarshalIndent(tree.ShowStruct(), "", "\t")
	t.Logf("tree: %s", result)
}

func TestForest_run(t *testing.T) {
	var name1 = "tree_1"
	var builder1 rule.TreeBuilder = func() (string, *rule.Tree) {
		var rules = []*rule.Rule{
			{Path: "/", Operators: []driver.Operator{
				&driver.JSONOperator{T: "create", JSONPath: "author.first", V: "river"},
				&driver.JSONOperator{T: "create", JSONPath: "name", V: "root"},
			}},
			{Path: "/a/b/c/d", Operators: []driver.Operator{
				&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:a/b/c/d"},
			}},
			{Path: "/a/b/c", Operators: []driver.Operator{
				&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:a/b/c"},
			}},
			{Path: "/a/b", Operators: []driver.Operator{
				&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:a/b"},
			}},
			{Path: "/x/y/z", Operators: []driver.Operator{
				&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:x/y/z"},
			}},
			{Path: "/a/b/m", Operators: []driver.Operator{
				&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:a/b/m"},
			}},
		}
		tree, err := rule.NewTree(&struct {
			driver.PathParser
			driver.StdCalculator
			driver.DummyOperatorModem
			driver.DummyDriver // just provide a method Name
		}{PathParser: driver.SlashPathParser},
			"json_tree_1", `{"id":1}`, rules...)
		if err != nil {
			t.Errorf("build tree fail: %s", err)
			return "", nil
		}
		return name1, tree
	}

	var name2 = "tree_2"
	var builder2 rule.TreeBuilder = func() (string, *rule.Tree) {
		var rules = []*rule.Rule{
			{Path: "/a/b/c/d", Operators: nil},
			{Path: "/a/b/c", Operators: nil},
			{Path: "/", Operators: nil},
			{Path: "/x/y/z", Operators: nil},
			{Path: "/a/b/m", Operators: nil},
		}
		tree, err := rule.NewTree(&struct {
			driver.PathParser
			driver.StdCalculator
			driver.DummyOperatorModem
			driver.DummyDriver // just provide a method Name
		}{PathParser: driver.SlashPathParser},
			"json_tree_1", `{"id":2}`, rules...)
		if err != nil {
			t.Errorf("build tree fail: %s", err)
			return "", nil
		}
		return name2, tree
	}

	f := rule.NewForest(builder1, builder2)
	t.Logf("forest info: %s", f.Info())

	rule := f.Get(name1).GetRule("a/b/c/d")
	t.Logf("got rule: %s", rule)

	rule = f.Get(name1).GetRule("a/b")
	t.Logf("got rule: %s", rule)

	rule = f.Get(name1).GetRule("a/b/x")
	t.Logf("got rule: %s", rule)
}
