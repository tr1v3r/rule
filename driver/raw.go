package driver

import (
	"time"
)

var _ Processor = (*RawProcessor)(nil)

type RawProcessor struct {
	author    string
	createdAt time.Time

	Proc func([]byte) ([]byte, error)
}

func (op *RawProcessor) Type() string                                    { return "" }
func (op *RawProcessor) Path() string                                    { return "" }
func (op *RawProcessor) Author() string                                  { return op.author }
func (op *RawProcessor) CreatedAt() time.Time                            { return op.createdAt }
func (op *RawProcessor) Load(data []byte) error                          { return ErrSerializeNotSupport }
func (op *RawProcessor) Save() []byte                                    { return nil }
func (op *RawProcessor) Process(before []byte) (after []byte, err error) { return op.Proc(before) }
