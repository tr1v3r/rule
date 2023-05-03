package driver

import (
	"time"
)

var _ Driver = (*JSONDriver)(nil)

// JSONDriver is a driver for JSON type rule tree
type JSONDriver struct {
	SlashPathDriver
	CommonRuleDriver
	JSONOperatorDriver
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
