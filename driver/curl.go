package driver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tr1v3r/pkg/fetch"
)

var _ Operator = (*CURLOperator)(nil)

// CURLOperator
type CURLOperator struct {
	// P is the target path of the operator
	P string `json:"path"`

	// URL the target url
	URL string `json:"url"`
	// Method the method to call URL
	Method string `json:"method"`
	// Body post with data
	Body   []byte `json:"body"`
	Header map[string][]string

	// A is the author of the operator
	A string `json:"author"`
	// C is the create time of the operator
	C time.Time `json:"created_at"`
}

func (op *CURLOperator) Type() string         { return "curl" }
func (op *CURLOperator) Path() string         { return op.P }
func (op *CURLOperator) Author() string       { return op.A }
func (op *CURLOperator) CreatedAt() time.Time { return op.C }
func (op *CURLOperator) Load(data []byte) error {
	if err := json.Unmarshal(data, op); err != nil {
		return fmt.Errorf("unmarshal fail: %w", err)
	}
	return nil
}
func (op *CURLOperator) Save() []byte {
	data, _ := json.Marshal(op)
	return data
}
func (op *CURLOperator) Operate(_ string) (string, error) {
	method := strings.ToUpper(strings.TrimSpace(op.Method))
	if method == "" {
		method = "GET"
	}

	_, content, _, err := fetch.DoRequestWithOptions(method, op.URL,
		[]fetch.RequestOption{fetch.WithHeaders(op.Header)}, bytes.NewReader(op.Body))
	if err != nil {
		return "", fmt.Errorf("request url %s fail: %w", op.URL, err)
	}
	return string(content), nil
}
