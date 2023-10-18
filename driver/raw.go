package driver

import (
	"time"
)

var _ Operator = (*RawOperator)(nil)

type RawOperator struct {
	author    string
	createdAt time.Time

	Proc func(string) (string, error)
}

func (op *RawOperator) Type() string                                    { return "" }
func (op *RawOperator) Path() string                                    { return "" }
func (op *RawOperator) Author() string                                  { return op.author }
func (op *RawOperator) CreatedAt() time.Time                            { return op.createdAt }
func (op *RawOperator) Load(data []byte) error                          { return ErrSerializeNotSupport }
func (op *RawOperator) Save() []byte                                    { return nil }
func (op *RawOperator) Operate(before string) (after string, err error) { return op.Proc(before) }
