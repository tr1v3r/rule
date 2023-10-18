package driver_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tr1v3r/rule/driver"
	"gopkg.in/yaml.v3"
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

	data, err := d.Marshal([]driver.Processor{
		&driver.JSONProcessor{T: "create", JSONPath: "name.first", V: "river"},
		&driver.JSONProcessor{T: "create", JSONPath: "name.last", V: "chu"},
		&driver.JSONProcessor{T: "create", JSONPath: "name.last", V: "Chu"},
		&driver.JSONProcessor{T: "append", JSONPath: "dear.friends.-1", V: "tom"},
		&driver.JSONProcessor{T: "append", JSONPath: "dear.friends.-1", V: "ken"},
		&driver.JSONProcessor{T: "set", JSONPath: "dear.family", V: `["mom","dad","bro"]`},
		&driver.JSONProcessor{T: "create", JSONPath: "name.verbose", V: "verbose"},
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
	var rule string

	f, err := os.ReadFile("/tmp/rule.yml")
	if err != nil {
		t.Errorf("read file fail: %s", err)
		return
	}
	rule = string(f)

	var ops = []driver.Processor{
		&driver.RawProcessor{Proc: func(before string) (string, error) {
			var result any
			if err := yaml.Unmarshal([]byte(before), &result); err != nil {
				return "", fmt.Errorf("unmarshal rule fail: %w", err)
			}
			result.(map[string]any)["unit"] = "test"
			newData, err := yaml.Marshal(result)
			return string(newData), err
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
