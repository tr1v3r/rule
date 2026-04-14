package rule

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

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

func TestLazyCacheTree_TTLExpiry(t *testing.T) {
	var realizeCount int32

	// Use a processor that produces different output each time to detect re-realization
	incrementProcessor := &driver.RawProcessor{
		Proc: func(before []byte) ([]byte, error) {
			n := atomic.AddInt32(&realizeCount, 1)
			return fmt.Appendf(nil, `{"realize":%d}`, n), nil
		},
	}

	tree, err := NewLazyCacheTree(
		&struct {
			driver.Modem
			driver.PathParser
			driver.StdRealizer
			driver.DummyDriver
		}{Modem: driver.DummyModem, PathParser: driver.SlashPathParser},
		"cache_test", `{}`, 50*time.Millisecond,
		NewRule("/", incrementProcessor),
	)
	if err != nil {
		t.Fatalf("build tree fail: %s", err)
	}

	// First Get — triggers realization
	result1, err := tree.Get("/")
	if err != nil {
		t.Fatalf("first get fail: %s", err)
	}
	if atomic.LoadInt32(&realizeCount) != 1 {
		t.Fatalf("expected realizeCount=1, got %d", realizeCount)
	}

	// Second Get immediately — should use cache
	result2, err := tree.Get("/")
	if err != nil {
		t.Fatalf("second get fail: %s", err)
	}
	if string(result1) != string(result2) {
		t.Errorf("cache miss on immediate second get: %s vs %s", result1, result2)
	}
	if atomic.LoadInt32(&realizeCount) != 1 {
		t.Errorf("expected realizeCount=1 (cached), got %d", realizeCount)
	}

	// Wait for TTL to expire
	time.Sleep(60 * time.Millisecond)

	// Third Get — should re-realize with different result
	result3, err := tree.Get("/")
	if err != nil {
		t.Fatalf("third get fail: %s", err)
	}
	if atomic.LoadInt32(&realizeCount) != 2 {
		t.Errorf("expected realizeCount=2 (re-realized), got %d", realizeCount)
	}
	if string(result3) == string(result1) {
		t.Errorf("expected different result after TTL expiry")
	}
	t.Logf("cached: %s, after TTL: %s", result1, result3)
}

func TestLazyCacheTree_ZeroTTL(t *testing.T) {
	// Zero TTL should behave like standard lazy — cache forever
	tree, err := NewLazyCacheTree(
		&struct {
			driver.Modem
			driver.PathParser
			driver.StdRealizer
			driver.DummyDriver
		}{Modem: driver.DummyModem, PathParser: driver.SlashPathParser},
		"zero_ttl_test", `{"v":0}`, 0,
		NewRule("/", &driver.JSONProcessor{T: "create", JSONPath: "v", V: []byte("1")}),
	)
	if err != nil {
		t.Fatalf("build tree fail: %s", err)
	}

	result1, _ := tree.Get("/")
	time.Sleep(20 * time.Millisecond)
	result2, _ := tree.Get("/")
	if string(result1) != string(result2) {
		t.Errorf("zero TTL should cache forever, got different results: %s vs %s", result1, result2)
	}
}
