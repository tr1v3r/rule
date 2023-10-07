package driver

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

var _ Driver = (*YAMLDriver)(nil)

func NewYAMLDriver() *JSONDriver {
	return &JSONDriver{
		PathParser: new(DelimiterPathParser).WithDelimiter("/"),
		Calculator: new(StdCalculator),
		Modem:      new(YAMLOperatorModem),
	}
}

// YAMLDriver is a driver for YAML type rule tree
type YAMLDriver struct {
	PathParser
	Calculator
	Modem
}

func (YAMLDriver) Name() string { return "yaml" }

var _ Operator = (*YAMLOperator)(nil)

type YAMLOperator struct{}

func (op *YAMLOperator) Type() string                                    { return "" }
func (op *YAMLOperator) Path() string                                    { return "" }
func (op *YAMLOperator) Operate(before string) (after string, err error) { return before, nil }
func (op *YAMLOperator) Author() string                                  { return "" }
func (op *YAMLOperator) CreatedAt() time.Time                            { return time.Now() }
func (op *YAMLOperator) Load([]byte) error                               { return nil }
func (op *YAMLOperator) Save() []byte                                    { return nil }

var _ Modem = (*YAMLOperatorModem)(nil)

// YAMLOperatorModem operator driver for yaml
type YAMLOperatorModem struct{}

func (d *YAMLOperatorModem) Marshal(ops ...Operator) ([]byte, error) {
	var buf = make([][]byte, 0, len(ops))
	for _, op := range ops {
		buf = append(buf, op.Save())
	}
	return yaml.Marshal(buf)
}
func (d *YAMLOperatorModem) Unmarshal(data []byte) ([]Operator, error) {
	var buf = make([][]byte, 0, 8)
	if err := json.Unmarshal(data, &buf); err != nil {
		return nil, fmt.Errorf("unmarshal fail: %w", err)
	}

	var ops = make([]Operator, 0, len(buf))
	for _, item := range buf {
		op := new(YAMLOperator)
		if err := op.Load(item); err != nil {
			return nil, fmt.Errorf("load data fail: %w", err)
		}
		ops = append(ops, op)
	}
	return ops, nil
}
