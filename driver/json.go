package driver

import (
	"encoding/json"
	"fmt"
	"time"
)

var _ Driver = (*JSONDriver)(nil)

func NewJSONDriver() *JSONDriver {
	return &JSONDriver{
		PathParser: new(DelimiterPathParser).WithDelimiter("/"),
		Calculator: new(StdCalculator),
		Modem:      new(JSONOperatorModem),
	}
}

// JSONDriver is a driver for JSON type rule tree
type JSONDriver struct {
	PathParser
	Calculator
	Modem
}

func (JSONDriver) Name() string { return "json" }

var _ Operator = (*JSONOperator)(nil)

type JSONOperator struct{}

func (op *JSONOperator) Type() string                                    { return "" }
func (op *JSONOperator) Operate(before string) (after string, err error) { return "", nil }
func (op *JSONOperator) Author() string                                  { return "" }
func (op *JSONOperator) CreatedAt() time.Time                            { return time.Now() }
func (op *JSONOperator) Path() string                                    { return "" }
func (op *JSONOperator) Load([]byte) error                               { return nil }
func (op *JSONOperator) Save() []byte                                    { return nil }

var _ Modem = (*JSONOperatorModem)(nil)

// JSONOperatorModem modem for json operator
type JSONOperatorModem struct{}

func (d *JSONOperatorModem) Marshal(ops ...Operator) []byte {
	var buf = make([]json.RawMessage, 0, len(ops))
	for _, op := range ops {
		buf = append(buf, op.Save())
	}
	data, _ := json.Marshal(buf)
	return data
}
func (d *JSONOperatorModem) Unmarshal(data []byte) ([]Operator, error) {
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
