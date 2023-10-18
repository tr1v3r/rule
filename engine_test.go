package rule

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tr1v3r/rule/driver"
)

func TestBuildTree(t *testing.T) {
	var testcases = []struct {
		Rules    []*rule
		Expected string
	}{
		{Rules: []*rule{
			{path: "/a/b/c/d"},
			{path: "/a/b/c"},
			{path: "/"},
			{path: "x/y/z"},
			{path: "/a/b/m"},
		},
			Expected: `{"a":{"b":{"c":{"d":{}},"m":{}}},"x":{"y":{"z":{}}}}`},
	}

	for index, item := range testcases {
		tree, err := NewTree(&struct {
			driver.Modem
			driver.PathParser
			driver.StdCalculator
			driver.DummyDriver // just provide a method Name
		}{Modem: driver.DummyModem, PathParser: driver.SlashPathParser},
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

func InitForest(t *testing.T) Forest {
	var builders []TreeBuilder = []TreeBuilder{
		func() Tree {
			var rules = []*rule{
				{path: "/", operators: []driver.Operator{
					&driver.JSONOperator{T: "create", JSONPath: "author.first", V: "river"},
					&driver.JSONOperator{T: "create", JSONPath: "name", V: "root"},
				}},
				{path: "/a/b/c/d", operators: []driver.Operator{
					&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:a/b/c/d"},
				}},
				{path: "/a/b/c", operators: []driver.Operator{
					&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:a/b/c"},
				}},
				{path: "/a/b", operators: []driver.Operator{
					&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:a/b"},
				}},
				{path: "/x/y/z", operators: []driver.Operator{
					&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:x/y/z"},
				}},
				{path: "/a/b/m", operators: []driver.Operator{
					&driver.JSONOperator{T: "create", JSONPath: "info.path", V: "path:a/b/m"},
				}},
			}
			tree, err := NewTree(&struct {
				driver.Modem
				driver.PathParser
				driver.StdCalculator
				driver.DummyDriver // just provide a method Name
			}{Modem: driver.DummyModem, PathParser: driver.SlashPathParser},
				"json_tree_1", `{"id":1}`, rules...)
			if err != nil {
				t.Errorf("build tree fail: %s", err)
				return nil
			}
			return tree
		},
		func() Tree {
			var rules = []*rule{
				{path: "/a/b/c/d", operators: nil},
				{path: "/a/b/c", operators: nil},
				{path: "/", operators: []driver.Operator{&driver.CURLOperator{URL: "https://xxx/ping"}}},
				{path: "/x/y/z", operators: nil},
				{path: "/a/b/m", operators: nil},
			}
			tree, err := NewTree(&struct {
				driver.Modem
				driver.PathParser
				driver.StdCalculator
				driver.DummyDriver // just provide a method Name
			}{Modem: driver.DummyModem, PathParser: driver.SlashPathParser},
				"json_tree_2", `{"id":2}`, rules...)
			if err != nil {
				t.Errorf("build tree fail: %s", err)
				return nil
			}
			return tree
		},
	}
	return NewForest(builders...)
}
