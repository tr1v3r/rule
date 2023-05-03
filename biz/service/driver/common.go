package driver

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SlashPathDriver struct{}

func (d *SlashPathDriver) GetLevel(path string) int {
	if path = strings.TrimSpace(path); path != "/" && path != "" {
		return len(strings.Split(path, "/")) - 1
	}
	return 0
}
func (d *SlashPathDriver) GetNameByLevel(path string, level int) string {
	return strings.Split(path, "/")[level]
}

type CommonRuleDriver struct{}

func (d *CommonRuleDriver) CalcRule(template string, ops ...Operator) (string, error) {
	var err error
	for _, op := range ops {
		if template, err = op.Operate(template); err != nil {
			return "", fmt.Errorf("operate fail: %w", err)
		}
	}
	return template, nil
}

type JSONOperatorDriver struct{}

func (d *JSONOperatorDriver) Marshal(ops ...Operator) []byte {
	var buf = make([]json.RawMessage, 0, len(ops))
	for _, op := range ops {
		buf = append(buf, op.Save())
	}
	data, _ := json.Marshal(buf)
	return data
}
func (d *JSONOperatorDriver) Unmarshal(data []byte) ([]Operator, error) {
	var buf = make([]json.RawMessage, 0, 8)
	if err := json.Unmarshal(data, &buf); err != nil {
		return nil, fmt.Errorf("unmarshal fail: %w", err)
	}

	var ops = make([]Operator, 0, len(buf))
	for _, item := range buf {
		op := new(JSONOperator)
		if err := op.Load(item); err != nil {
			return nil, fmt.Errorf("load data fail: %w", err)
		}
		ops = append(ops, op)
	}
	return ops, nil
}
