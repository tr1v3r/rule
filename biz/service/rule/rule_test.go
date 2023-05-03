package rule_test

import (
	"encoding/json"
	"testing"

	"github.com/riverchu/rule/biz/service/driver"
	"github.com/riverchu/rule/biz/service/rule"
)

func TestBuildTree(t *testing.T) {
	var rules = []*rule.Rule{
		{Path: "/a/b/c/d", Operators: nil},
		{Path: "/a/b/c", Operators: nil},
		{Path: "/", Operators: nil},
		{Path: "/x/y/z", Operators: nil},
		{Path: "/a/b/m", Operators: nil},
	}

	tree, err := rule.BuildTree("unit_test", "{}", &struct {
		driver.SlashPathDriver
		driver.CommonRuleDriver
		driver.DummyOperatorDriver
		driver.DummyDriver // just provide a method Name
	}{}, rules...)
	if err != nil {
		t.Errorf("build tree fail: %s", err)
		return
	}

	result, _ := json.MarshalIndent(tree.ShowStruct(), "", "\t")
	t.Logf("tree: %s", result)
}
