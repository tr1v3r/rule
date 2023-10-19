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
			driver.StdRealizer
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
		rule, err := f.Get(item.Name).Get(item.TargetPath)
		if err != nil {
			t.Errorf("query rule %s on %s fail: %s", item.TargetPath, item.Name, err)
			continue
		}
		if string(rule) != item.Expected {
			t.Errorf("check rule fail, expect: %s\ngot: %s", item.Expected, string(rule))
			continue
		}
		t.Logf("got rule: %s", string(rule))
	}
}

func InitForest(t *testing.T) Forest {
	var builders []TreeBuilder = []TreeBuilder{
		func() Tree {
			var rules = []*rule{
				{path: "/", processors: []driver.Processor{
					&driver.JSONProcessor{T: "create", JSONPath: "author.first", V: []byte("river")},
					&driver.JSONProcessor{T: "create", JSONPath: "name", V: []byte("root")},
				}},
				{path: "/a/b/c/d", processors: []driver.Processor{
					&driver.JSONProcessor{T: "create", JSONPath: "info.path", V: []byte("path:a/b/c/d")},
				}},
				{path: "/a/b/c", processors: []driver.Processor{
					&driver.JSONProcessor{T: "create", JSONPath: "info.path", V: []byte("path:a/b/c")},
				}},
				{path: "/a/b", processors: []driver.Processor{
					&driver.JSONProcessor{T: "create", JSONPath: "info.path", V: []byte("path:a/b")},
				}},
				{path: "/x/y/z", processors: []driver.Processor{
					&driver.JSONProcessor{T: "create", JSONPath: "info.path", V: []byte("path:x/y/z")},
				}},
				{path: "/a/b/m", processors: []driver.Processor{
					&driver.JSONProcessor{T: "create", JSONPath: "info.path", V: []byte("path:a/b/m")},
				}},
			}
			tree, err := NewTree(&struct {
				driver.Modem
				driver.PathParser
				driver.StdRealizer
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
				{path: "/a/b/c/d", processors: nil},
				{path: "/a/b/c", processors: nil},
				{path: "/", processors: []driver.Processor{&driver.CURLProcessor{URL: "https://xxx/ping"}}},
				{path: "/x/y/z", processors: nil},
				{path: "/a/b/m", processors: nil},
			}
			tree, err := NewTree(&struct {
				driver.Modem
				driver.PathParser
				driver.StdRealizer
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
