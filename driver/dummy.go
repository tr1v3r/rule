package driver

import (
	"time"
)

var _ Driver = (*DummyDriver)(nil)

var DummyModem = &GeneralModem[Operator]{
	Marshaler:   func(any) ([]byte, error) { return nil, nil },
	Unmarshaler: func([]byte, any) error { return nil },
}

// NewDummyDriver ...
func NewDummyDriver() *DummyDriver {
	return &DummyDriver{
		PathParser: new(DelimiterPathParser).WithDelimiter("/"),
		Calculator: new(StdCalculator),
		Modem:      DummyModem,
	}
}

// DummyDriver return a dummy driver
type DummyDriver struct {
	PathParser
	Calculator
	Modem
}

func (DummyDriver) Name() string { return "dummy" }

var _ Operator = (*DummyOperator)(nil)

type DummyOperator struct{}

func (op *DummyOperator) Type() string                                    { return "dummy" }
func (op *DummyOperator) Path() string                                    { return "" }
func (op *DummyOperator) Operate(before string) (after string, err error) { return "", nil }
func (op *DummyOperator) Author() string                                  { return "dummy" }
func (op *DummyOperator) CreatedAt() time.Time                            { return time.Now() }
func (op *DummyOperator) Load([]byte) error                               { return nil }
func (op *DummyOperator) Save() []byte                                    { return nil }
