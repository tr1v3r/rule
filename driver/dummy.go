package driver

import (
	"time"
)

// var _ Driver = (*DummyDriver)(nil)

func NewDummyDriver() *DummyDriver {
	return &DummyDriver{
		PathParser: new(DelimiterPathParser).WithDelimiter("/"),
		Calculator: new(StdCalculator),
		Modem:      new(DummyOperatorModem),
	}
}

// DummyDriver return a dummy driver
type DummyDriver struct {
	PathParser
	Calculator
	Modem
}

func (DummyDriver) Name() string { return "dummy" }

type DummyOperatorModem struct{}

func (d *DummyOperatorModem) Marshal(ops ...Operator) []byte            { return nil }
func (d *DummyOperatorModem) Unmarshal(data []byte) ([]Operator, error) { return nil, nil }

var _ Operator = (*DummyOperator)(nil)

type DummyOperator struct{}

func (op *DummyOperator) Type() string                                    { return "" }
func (op *DummyOperator) Operate(before string) (after string, err error) { return "", nil }
func (op *DummyOperator) Author() string                                  { return "" }
func (op *DummyOperator) CreatedAt() time.Time                            { return time.Now() }
func (op *DummyOperator) Path() string                                    { return "" }
func (op *DummyOperator) Load([]byte) error                               { return nil }
func (op *DummyOperator) Save() []byte                                    { return nil }
