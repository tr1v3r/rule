package rule

import (
	"encoding/json"
	"strings"
	"testing"
)

var _ Driver = &DummyDriver{}

type DummyDriver struct{}

func (d *DummyDriver) GetLevel(path string) int { return len(strings.Split(path, "/")) - 1 }
func (d *DummyDriver) GetNameByLevel(path string, level int) string {
	return strings.Split(path, "/")[level]
}
func (d *DummyDriver) CalcRule(template string, op *Rule) (string, error) {
	return "", nil
}

func TestBuildTree(t *testing.T) {
	var rules = []*Rule{
		{Path: "/a/b/c/d"},
		{Path: "/a/b/c"},
		{Path: ""},
		{Path: "/x/y/z"},
		{Path: "/a/b/m"},
	}

	tree := BuildTree("unit_test", "dummy", &DummyDriver{}, rules...)

	result, _ := json.MarshalIndent(tree.ShowStruct(), "", "\t")
	t.Logf("tree: %s", result)
}
