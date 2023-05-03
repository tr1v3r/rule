package driver_test

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDriver_Marshal(t *testing.T) {
	var buf = make([]json.RawMessage, 0, 8)
	for i := range [8]int{} {
		buf = append(buf, json.RawMessage(fmt.Sprintf(`{"index":%d}`, i)))
	}
	data, _ := json.Marshal(buf)
	t.Logf("got result: %s", string(data))

	var nextBuf = make([]json.RawMessage, 0, 8)
	if err := json.Unmarshal(data, &nextBuf); err != nil {
		t.Errorf("unmarshal fail: %s", err)
	}
	for _, v := range nextBuf {
		t.Logf("got result: %s", string(v))
	}
}
