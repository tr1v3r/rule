package driver

import (
	"encoding/json"
	"time"
)

var _ Driver = (*YAMLDriver)(nil)

func NewYAMLDriver() *YAMLDriver {
	return &YAMLDriver{
		PathParser: new(DelimiterPathParser).WithDelimiter("/"),
		Calculator: new(StdCalculator),
		Modem: &GeneralModem[*YAMLOperator]{
			Marshaler:   json.Marshal,
			Unmarshaler: json.Unmarshal,
		},
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

type YAMLOperator struct {
	T string `json:"type"`
	V string `json:"value"`
}

func (op *YAMLOperator) Type() string           { return op.T }
func (op *YAMLOperator) Path() string           { return "" }
func (op *YAMLOperator) Author() string         { return "" }
func (op *YAMLOperator) CreatedAt() time.Time   { return time.Now() }
func (op *YAMLOperator) Load(data []byte) error { return json.Unmarshal(data, op) }
func (op *YAMLOperator) Save() []byte {
	data, _ := json.Marshal(op)
	return data
}
func (op *YAMLOperator) Operate(before string) (after string, err error) {
	// var result any
	// if err := yaml.Unmarshal([]byte(before), &result); err != nil {
	// 	return "", fmt.Errorf("unmarshal yaml fail: %w", err)
	// }
	switch op.T {
	case "append":
		return before + op.T, nil
	}
	return before, nil
}
