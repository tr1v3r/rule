package driver

import (
	"fmt"
	"strings"
)

var _ PathParser = (*DelimiterPathParser)(nil)

// SlashPathParser slash path parser
// equal to DelimiterPathParser{delimiter: "/"}
var SlashPathParser = new(DelimiterPathParser).WithDelimiter("/")

// DelimiterPathParser delimiter path parser
type DelimiterPathParser struct {
	delimiter string
}

func (d DelimiterPathParser) WithDelimiter(delimiter string) *DelimiterPathParser {
	d.delimiter = delimiter
	return &d
}
func (d *DelimiterPathParser) GetLevel(path string) int {
	if path = strings.TrimSpace(path); path != d.delimiter && path != "" {
		return len(strings.Split(path, d.delimiter)) - 1
	}
	return 0
}
func (d *DelimiterPathParser) GetNameByLevel(path string, level int) string {
	return strings.Split(path, d.delimiter)[level]
}

var _ Calculator = (*StdCalculator)(nil)

// StdCalculator standard rule driver
type StdCalculator struct{}

// CalcRule calculate rule
func (d *StdCalculator) CalcRule(template string, ops ...Operator) (string, error) {
	var err error
	for _, op := range ops {
		if template, err = op.Operate(template); err != nil {
			return "", fmt.Errorf("operate fail: %w", err)
		}
	}
	return template, nil
}
