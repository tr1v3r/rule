package rule_test

import (
	"testing"

	"github.com/tr1v3r/rule"
	"github.com/tr1v3r/rule/driver"
	"github.com/tr1v3r/stream"
)

func TestBuildTree_curl_single(t *testing.T) {
	tree, err := rule.NewYAMLTree("qq", "", rule.NewRule("/", &driver.CURLOperator{URL: "https://qq.com"}))
	if err != nil {
		t.Errorf("build tree fail: %s", err)
	}
	t.Logf("get rule by curl: %s", tree.GetRule(""))
}

func TestBuildForest_curl_tree(t *testing.T) {
	f := rule.NewForest(func() rule.Tree {
		tree, err := rule.NewYAMLTree("qq", "", rule.NewRule("/", &driver.CURLOperator{URL: "https://qq.com"}))
		if err != nil {
			t.Errorf("build tree fail: %s", err)
			return nil
		}
		t.Logf("get rule by curl: %s", tree.GetRule(""))
		return tree
	})
	t.Logf("get rule from forest: %s", f.Get("qq").GetRule("/"))
}

func TestBuildForest_stream(t *testing.T) {
	trees := stream.SliceOf("https://qq.com", "https://163.com").Parallel(64).Convert(func(url string) any {
		tree, _ := rule.NewYAMLTree("url", "", rule.NewRule("/", &driver.CURLOperator{URL: url}))
		return tree
	}).Collect(func(trees ...any) any {
		var treesArray []rule.Tree
		for _, tree := range trees {
			treesArray = append(treesArray, tree.(rule.Tree))
		}
		return treesArray
	}).([]rule.Tree)

	for i, tree := range trees {
		t.Logf("got tree [%d] %s: %s", i, tree.Name(), tree.GetRule(""))
	}
}
