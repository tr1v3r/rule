package driver

import (
	"time"
)

var _ Driver = (*DummyDriver)(nil)

var DummyModem = &GeneralModem[Processor]{
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

var _ Processor = (*DummyProcessor)(nil)

type DummyProcessor struct{}

func (op *DummyProcessor) Type() string                                    { return "dummy" }
func (op *DummyProcessor) Path() string                                    { return "" }
func (op *DummyProcessor) Process(before string) (after string, err error) { return "", nil }
func (op *DummyProcessor) Author() string                                  { return "dummy" }
func (op *DummyProcessor) CreatedAt() time.Time                            { return time.Now() }
func (op *DummyProcessor) Load([]byte) error                               { return nil }
func (op *DummyProcessor) Save() []byte                                    { return nil }
