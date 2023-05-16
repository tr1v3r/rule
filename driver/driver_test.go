package driver_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/riverchu/rule/driver"
)

func TestSlashPathParser(t *testing.T) {
	var testcases = []struct {
		path  string
		level int
	}{
		{"", 0},
		{"/", 0},
		{"/a/", 1},
		{"/a", 1},
		{"a/", 1},
		{"a", 1},
		{"/a/b/", 2},
		{"a/b/", 2},
		{"/a/b", 2},
		{"a/b", 2},
		{"/a/b/c/", 3},
		{"a/b/c/", 3},
		{"/a/b/c", 3},
		{"a/b/c", 3},
	}

	for _, item := range testcases {
		if l := driver.SlashPathParser.GetLevel(item.path); l != item.level {
			t.Errorf("get level expected %d, got %d", item.level, l)
		}
	}
}

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

func TestJSONDriver(t *testing.T) {
	var rule string

	d := driver.NewJSONDriver()

	data, err := d.Marshal([]driver.Operator{
		&driver.JSONOperator{T: "create", JSONPath: "name.first", V: "river"},
		&driver.JSONOperator{T: "create", JSONPath: "name.last", V: "chu"},
		&driver.JSONOperator{T: "create", JSONPath: "name.last", V: "Chu"},
		&driver.JSONOperator{T: "append", JSONPath: "dear.friends.-1", V: "tom"},
		&driver.JSONOperator{T: "append", JSONPath: "dear.friends.-1", V: "ken"},
		&driver.JSONOperator{T: "set", JSONPath: "dear.family", V: `["mom","dad","bro"]`},
		&driver.JSONOperator{T: "create", JSONPath: "name.verbose", V: "verbose"},
		&driver.JSONOperator{T: "delete", JSONPath: "name.verbose"},
	}...)
	if err != nil {
		t.Errorf("marshal fail: %s", err)
		return
	}
	ops, err := d.Unmarshal(data)
	if err != nil {
		t.Errorf("unmarshal fail: %s", err)
		return
	}
	for _, op := range ops {
		rule, err = op.Operate(rule)
		if err != nil {
			t.Errorf("operate fail: %s", err)
			return
		}
	}
	t.Logf("got result: %s", rule)
}
