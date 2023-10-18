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
		Modem: &GeneralModem[*JSONProcessor]{
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

var _ Processor = (*JSONProcessor)(nil)

// JSONProcessor is a Processor for JSON type rule tree
type JSONProcessor struct {
	// P is the target path of the Processor
	P string `json:"path"`

	// T is the type of the Processor
	T string `json:"type"`
	// JSONPath is the json path of the Processor
	JSONPath string `json:"json_path"`
	// V is the value of the Processor
	V string `json:"value"`

	// A is the author of the Processor
	A string `json:"author"`
	// C is the create time of the Processor
	C time.Time `json:"created_at"`
}

func (op *JSONProcessor) Type() string         { return op.T }
func (op *JSONProcessor) Path() string         { return op.P }
func (op *JSONProcessor) Author() string       { return op.A }
func (op *JSONProcessor) CreatedAt() time.Time { return op.C }
func (op *JSONProcessor) Load(data []byte) error {
	if err := json.Unmarshal(data, op); err != nil {
		return fmt.Errorf("unmarshal fail: %w", err)
	}
	return nil
}
func (op *JSONProcessor) Save() []byte {
	data, _ := json.Marshal(op)
	return data
}
func (op *JSONProcessor) Process(before string) (after string, err error) {
	switch op.T {
	case "create", "append", "replace":
		return sjson.Set(before, op.JSONPath, op.V)
	case "set":
		return sjson.SetRaw(before, op.JSONPath, op.V)
	case "delete":
		return sjson.Delete(before, op.JSONPath)
	default:
		return "", fmt.Errorf("unknown Processor type: %s", op.T)
	}
}
