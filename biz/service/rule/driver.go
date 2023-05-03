package rule

import (
	"strings"
)

var _ Driver = &JSONDriver{}

// JSONDriver is a driver for JSON type rule tree
type JSONDriver struct{}

func (JSONDriver) Type() string                { return "json" }
func (d *JSONDriver) GetLevel(path string) int { return len(strings.Split(path, "/")) }
func (d *JSONDriver) GetNameByLevel(path string, level int) string {
	return strings.Split(path, "/")[level]
}
func (d *JSONDriver) CalcRule(template string, op *Rule) (string, error) {
	return "", nil
}
