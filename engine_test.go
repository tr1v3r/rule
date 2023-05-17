package rule_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/riverchu/rule"
	"github.com/riverchu/rule/driver"
)

func TestBuildTree(t *testing.T) {
	var testcases = []struct {
		Rules    []*rule.Rule
		Expected string
	}{
		{Rules: []*rule.Rule{
			{Path: "/a/b/c/d"},
			{Path: "/a/b/c"},
			{Path: "/"},
			{Path: "x/y/z"},
			{Path: "/a/b/m"},
		},
			Expected: `{"a":{"b":{"c":{"d":{}},"m":{}}},"x":{"y":{"z":{}}}}`},
	}

	for index, item := range testcases {
		tree, err := rule.NewTree(&struct {
			driver.PathParser
			driver.StdCalculator
			driver.DummyOperatorModem
			driver.DummyDriver // just provide a method Name
		}{PathParser: driver.SlashPathParser},
			"unit_test_"+fmt.Sprint(index), "{}", item.Rules...)
		if err != nil {
			t.Errorf("build tree fail: %s", err)
			return
		}

		if result := string(tree.ShowStruct()); item.Expected != result {
			t.Errorf("build tree fail, expect: %s\ngot: %s", item.Expected, result)

			printS, _ := json.MarshalIndent(tree.ShowStruct(), "", "\t")
			t.Logf("got tree struct: %s", printS)
		}
	}
}

func TestForest_run(t *testing.T) {
	var testcases = []struct {
		Name       string
		TargetPath string
		Expected   string
	}{
		{"tree_1", "a/b", `{"id":1,"author":{"first":"river"},"name":"root","info":{"path":"path:a/b"}}`},
		{"tree_1", "a/b/c/d", `{"id":1,"author":{"first":"river"},"name":"root","info":{"path":"path:a/b/c/d"}}`},
		{"tree_1", "a/b/x", `{"id":1,"author":{"first":"river"},"name":"root","info":{"path":"path:a/b"}}`},
		{"tree_2", "x/y/z", `{"code":200,"msg":"pong"}`},
	}

	f := InitForest(t)
	t.Logf("forest info: %s", f.Info())

	for _, item := range testcases {
		rule := f.Get(item.Name).GetRule(item.TargetPath)
		if rule != item.Expected {
			t.Errorf("check rule fail, expect: %s\ngot: %s", item.Expected, rule)
		}
		t.Logf("got rule: %s", rule)
	}
}

func InitForest(t *testing.T) *rule.Forest {
	var builders []rule.TreeBuilder = []rule.TreeBuilder{
		func() (string, *rule.Tree) {
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
				"json_tree", `{"id":1}`, rules...)
			if err != nil {
				t.Errorf("build tree fail: %s", err)
				return "", nil
			}
			return "tree_1", tree
		},
		func() (string, *rule.Tree) {
			var rules = []*rule.Rule{
				{Path: "/a/b/c/d", Operators: nil},
				{Path: "/a/b/c", Operators: nil},
				{Path: "/", Operators: []driver.Operator{&driver.CURLOperator{URL: "https://xxx/ping"}}},
				{Path: "/x/y/z", Operators: nil},
				{Path: "/a/b/m", Operators: nil},
			}
			tree, err := rule.NewTree(&struct {
				driver.PathParser
				driver.StdCalculator
				driver.DummyOperatorModem
				driver.DummyDriver // just provide a method Name
			}{PathParser: driver.SlashPathParser},
				"json_tree", `{"id":2}`, rules...)
			if err != nil {
				t.Errorf("build tree fail: %s", err)
				return "", nil
			}
			return "tree_2", tree
		},
	}
	return rule.NewForest(builders...)
}
