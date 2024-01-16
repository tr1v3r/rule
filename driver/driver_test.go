package driver_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/tr1v3r/rule/driver"
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
	var rule []byte

	d := driver.NewJSONDriver()

	data, err := d.Marshal([]driver.Processor{
		&driver.JSONProcessor{T: "create", JSONPath: "name.first", V: []byte("river")},
		&driver.JSONProcessor{T: "create", JSONPath: "name.last", V: []byte("chu")},
		&driver.JSONProcessor{T: "create", JSONPath: "name.last", V: []byte("Chu")},
		&driver.JSONProcessor{T: "append", JSONPath: "dear.friends.-1", V: []byte("tom")},
		&driver.JSONProcessor{T: "append", JSONPath: "dear.friends.-1", V: []byte("ken")},
		&driver.JSONProcessor{T: "set", JSONPath: "dear.family", V: []byte(`["mom","dad","bro"]`)},
		&driver.JSONProcessor{T: "create", JSONPath: "name.verbose", V: []byte("verbose")},
		&driver.JSONProcessor{T: "delete", JSONPath: "name.verbose"},
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
		rule, err = op.Process(rule)
		if err != nil {
			t.Errorf("Process fail: %s", err)
			return
		}
	}
	t.Logf("got result: %s", rule)
}

func TestYAMLProcessor(t *testing.T) {
	var rule []byte
	var err error

	rule, err = os.ReadFile("/tmp/rule.yml")
	if err != nil {
		t.Errorf("read file fail: %s", err)
		return
	}

	var ops = []driver.Processor{
		&driver.RawProcessor{Proc: func(before []byte) ([]byte, error) {
			var result any
			if err := yaml.Unmarshal([]byte(before), &result); err != nil {
				return nil, fmt.Errorf("unmarshal rule fail: %w", err)
			}
			result.(map[string]any)["unit"] = "test"
			newData, err := yaml.Marshal(result)
			return newData, err
		}},
	}
	for _, op := range ops {
		rule, err = op.Process(rule)
		if err != nil {
			t.Errorf("Process fail: %s", err)
			return
		}
	}
	t.Logf("got result: %s", rule)
}
