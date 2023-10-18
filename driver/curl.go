package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tr1v3r/pkg/fetch"
)

var _ Processor = (*CURLProcessor)(nil)

// CURLProcessor
type CURLProcessor struct {
	// P is the target path of the Processor
	P string `json:"path,omitempty"`

	// URL the target url
	URL string `json:"url"`
	// Method the method to call URL
	Method string `json:"method,omitempty"`
	// Body post with data
	Body   []byte              `json:"body,omitempty"`
	Header map[string][]string `json:"header,omitempty"`

	// A is the author of the Processor
	A string `json:"author"`
	// C is the create time of the Processor
	C time.Time `json:"created_at"`
}

func (op *CURLProcessor) Type() string         { return "curl" }
func (op *CURLProcessor) Path() string         { return op.P }
func (op *CURLProcessor) Author() string       { return op.A }
func (op *CURLProcessor) CreatedAt() time.Time { return op.C }
func (op *CURLProcessor) Load(data []byte) error {
	if err := json.Unmarshal(data, op); err != nil {
		return fmt.Errorf("unmarshal fail: %w", err)
	}
	return nil
}
func (op *CURLProcessor) Save() []byte {
	data, _ := json.Marshal(op)
	return data
}
func (op *CURLProcessor) Process(_ []byte) ([]byte, error) {
	method := strings.ToUpper(strings.TrimSpace(op.Method))
	if method == "" {
		method = "GET"
	}

	_, content, _, err := fetch.DoRequestWithOptions(method, op.URL,
		[]fetch.RequestOption{fetch.WithHeaders(op.Header)}, bytes.NewReader(op.Body))
	if err != nil {
		return nil, fmt.Errorf("request url %s fail: %w", op.URL, err)
	}
	return content, nil
}
