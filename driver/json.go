package driver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/sjson"
)

// check interface
var _ Driver = (*JSONDriver)(nil)

// NewJSONDriver create a new json driver
func NewJSONDriver() *JSONDriver {
	return &JSONDriver{
		PathParser: new(DelimiterPathParser).WithDelimiter("/"),
		Calculator: new(StdCalculator),
		Modem: &GeneralModem[*JSONOperator]{
			Marshaler:   json.Marshal,
			Unmarshaler: json.Unmarshal,
		},
	}
}

// JSONDriver is a driver for JSON type rule tree
type JSONDriver struct {
	PathParser
	Calculator
	Modem
}

// Name return driver name
func (JSONDriver) Name() string { return "json" }

var _ Operator = (*JSONOperator)(nil)

// JSONOperator is a operator for JSON type rule tree
type JSONOperator struct {
	// P is the target path of the operator
	P string `json:"path"`

	// T is the type of the operator
	T string `json:"type"`
	// JSONPath is the json path of the operator
	JSONPath string `json:"json_path"`
	// V is the value of the operator
	V string `json:"value"`

	// A is the author of the operator
	A string `json:"author"`
	// C is the create time of the operator
	C time.Time `json:"created_at"`
}

func (op *JSONOperator) Type() string         { return op.T }
func (op *JSONOperator) Path() string         { return op.P }
func (op *JSONOperator) Author() string       { return op.A }
func (op *JSONOperator) CreatedAt() time.Time { return op.C }
func (op *JSONOperator) Load(data []byte) error {
	if err := json.Unmarshal(data, op); err != nil {
		return fmt.Errorf("unmarshal fail: %w", err)
	}
	return nil
}
func (op *JSONOperator) Save() []byte {
	data, _ := json.Marshal(op)
	return data
}
func (op *JSONOperator) Operate(before string) (after string, err error) {
	switch op.T {
	case "create", "append", "replace":
		return sjson.Set(before, op.JSONPath, op.V)
	case "set":
		return sjson.SetRaw(before, op.JSONPath, op.V)
	case "delete":
		return sjson.Delete(before, op.JSONPath)
	default:
		return "", fmt.Errorf("unknown operator type: %s", op.T)
	}
}
