package rule_test

import (
	"testing"

	"github.com/tr1v3r/stream"

	"github.com/tr1v3r/rule"
	"github.com/tr1v3r/rule/driver"
)

func TestBuildTree_curl_single(t *testing.T) {
	tree, err := rule.NewYAMLTree("qq", "", rule.NewRule("/", &driver.CURLProcessor{URL: "https://qq.com"}))
	if err != nil {
		t.Errorf("build tree fail: %s", err)
	}
	rule, _ := tree.Get("")
	t.Logf("get rule by curl: %s", rule)
}

func TestBuildForest_curl_tree(t *testing.T) {
	f := rule.NewForest(func() rule.Tree {
		tree, err := rule.NewYAMLTree("qq", "", rule.NewRule("/", &driver.CURLProcessor{URL: "https://qq.com"}))
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
		tree, _ := rule.NewYAMLTree("url", "", rule.NewRule("/", &driver.CURLProcessor{URL: url}))
		return tree
	}).Collect(func(trees ...any) any {
		var treesArray []rule.Tree
		for _, tree := range trees {
			treesArray = append(treesArray, tree.(rule.Tree))
		}
		return treesArray
	}).([]rule.Tree)

	for i, tree := range trees {
		rule, _ := tree.Get("")
		t.Logf("got tree [%d] %s: %s", i, tree.Name(), rule)
	}
}
