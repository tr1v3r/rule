package ivy_test

import (
	"testing"

	"github.com/tr1v3r/stream"

	"github.com/tr1v3r/ivy"
	"github.com/tr1v3r/ivy/driver"
)

func TestBuildTree_curl_single(t *testing.T) {
	tree, err := ivy.NewYAMLTree("qq", "", ivy.NewDirective("/", &driver.CURLProcessor{URL: "https://qq.com"}))
	if err != nil {
		t.Errorf("build tree fail: %s", err)
	}
	rule, _ := tree.Get("")
	t.Logf("get rule by curl: %s", rule)
}

func TestBuildForest_curl_tree(t *testing.T) {
	f := ivy.NewForest(func() ivy.Tree {
		tree, err := ivy.NewYAMLTree("qq", "", ivy.NewDirective("/", &driver.CURLProcessor{URL: "https://qq.com"}))
		if err != nil {
			t.Errorf("build tree fail: %s", err)
			return nil
		}
		rule, _ := tree.Get("")
		t.Logf("get rule by curl: %s", rule)
		return tree
	})
	rule, _ := f.Get("qq").Get("/")
	t.Logf("get rule from forest: %s", rule)
}

func TestBuildForest_stream(t *testing.T) {
	trees := stream.SliceOf("https://qq.com", "https://163.com").Parallel(64).Convert(func(url string) any {
		tree, _ := ivy.NewYAMLTree("url", "", ivy.NewDirective("/", &driver.CURLProcessor{URL: url}))
		return tree
	}).Collect(func(trees ...any) any {
		var treesArray []ivy.Tree
		for _, tree := range trees {
			treesArray = append(treesArray, tree.(ivy.Tree))
		}
		return treesArray
	}).([]ivy.Tree)

	for i, tree := range trees {
		rule, _ := tree.Get("")
		t.Logf("got tree [%d] %s: %s", i, tree.Name(), rule)
	}
}

func TestTileTree(t *testing.T) {
	var proc = func(content string) func(*driver.RuleContext, []byte) ([]byte, error) {
		return func(_ *driver.RuleContext, _ []byte) ([]byte, error) {
			return []byte(content), nil
		}
	}

	// tree, err := ivy.NewTileTree("test_tile_tree", "template",
	tree, err := ivy.NewLazyTileTree("test_tile_tree", "template",
		ivy.NewDirective("/abc", &driver.RawProcessor{Proc: proc("content1")}),
		ivy.NewDirective("/123", &driver.RawProcessor{Proc: proc("content2")}),
		ivy.NewDirective("/test", &driver.RawProcessor{Proc: proc("content3")}),
		ivy.NewDirective("/@@@", &driver.RawProcessor{Proc: proc("content4")}),
	)
	if err != nil {
		t.Errorf("build tile tree fail: %s", err)
	}

	var data []byte

	if data, err = tree.Get("/abc"); err != nil {
		t.Errorf("get /abc fail: %s", err)
	}
	if data, err = tree.Get("/123"); err != nil {
		t.Errorf("get /123 fail: %s", err)
	}
	if data, err = tree.Get("/test"); err != nil {
		t.Errorf("get /test fail: %s", err)
	}
	if data, err = tree.Get("/@@@"); err != nil {
		t.Errorf("get /@@@ fail: %s", err)
	}

	_ = data
}
