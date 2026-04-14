package driver_test

import (
	"strings"
	"testing"

	"github.com/tr1v3r/rule/driver"
)

func TestXMLDriver(t *testing.T) {
	d := driver.NewXMLDriver()

	data, err := d.Marshal([]driver.Processor{
		&driver.XMLProcessor{T: "create", XMLPath: "root/name/first", V: []byte("river")},
		&driver.XMLProcessor{T: "create", XMLPath: "root/name/last", V: []byte("chu")},
		&driver.XMLProcessor{T: "set", XMLPath: "root/name/last", V: []byte("Chu")},
		&driver.XMLProcessor{T: "create", XMLPath: "root/items/item", V: []byte("a")},
		&driver.XMLProcessor{T: "append", XMLPath: "root/items/item", V: []byte("b")},
		&driver.XMLProcessor{T: "delete", XMLPath: "root/name/first"},
		&driver.XMLProcessor{T: "replace", XMLPath: "root/items", V: []byte("<item>c</item><item>d</item>")},
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

	var rule []byte
	for _, op := range ops {
		rule, err = op.Process(rule)
		if err != nil {
			t.Errorf("Process fail: %s", err)
			return
		}
	}
	t.Logf("got result: %s", rule)

	result := string(rule)
	if !strings.Contains(result, "<last>Chu</last>") {
		t.Errorf("expected <last>Chu</last> in result, got: %s", result)
	}
	if strings.Contains(result, "<first>") {
		t.Errorf("expected <first> to be deleted, got: %s", result)
	}
	if !strings.Contains(result, "<item>c</item>") {
		t.Errorf("expected <item>c</item> after replace, got: %s", result)
	}
}

func TestXMLProcessor_Create(t *testing.T) {
	op := &driver.XMLProcessor{T: "create", XMLPath: "root/users/user", V: []byte("alice")}
	result, err := op.Process(nil)
	if err != nil {
		t.Errorf("Process fail: %s", err)
		return
	}
	s := string(result)
	if !strings.Contains(s, "<user>alice</user>") {
		t.Errorf("expected <user>alice</user>, got: %s", s)
	}
}

func TestXMLProcessor_Set(t *testing.T) {
	before := []byte(`<root><name>old</name></root>`)
	op := &driver.XMLProcessor{T: "set", XMLPath: "root/name", V: []byte("new")}
	result, err := op.Process(before)
	if err != nil {
		t.Errorf("Process fail: %s", err)
		return
	}
	s := string(result)
	if !strings.Contains(s, "<name>new</name>") {
		t.Errorf("expected <name>new</name>, got: %s", s)
	}
}

func TestXMLProcessor_Delete(t *testing.T) {
	before := []byte(`<root><a>1</a><b>2</b></root>`)
	op := &driver.XMLProcessor{T: "delete", XMLPath: "root/a"}
	result, err := op.Process(before)
	if err != nil {
		t.Errorf("Process fail: %s", err)
		return
	}
	s := string(result)
	if strings.Contains(s, "<a>") {
		t.Errorf("expected <a> to be deleted, got: %s", s)
	}
	if !strings.Contains(s, "<b>2</b>") {
		t.Errorf("expected <b>2</b> to remain, got: %s", s)
	}
}

func TestXMLProcessor_Replace(t *testing.T) {
	before := []byte(`<root><items><item>old</item></items></root>`)
	op := &driver.XMLProcessor{T: "replace", XMLPath: "root/items", V: []byte("<item>new1</item><item>new2</item>")}
	result, err := op.Process(before)
	if err != nil {
		t.Errorf("Process fail: %s", err)
		return
	}
	s := string(result)
	if strings.Contains(s, "old") {
		t.Errorf("expected 'old' to be replaced, got: %s", s)
	}
	if !strings.Contains(s, "<item>new1</item>") {
		t.Errorf("expected <item>new1</item>, got: %s", s)
	}
}

func TestXMLProcessor_EmptyBefore(t *testing.T) {
	op := &driver.XMLProcessor{T: "create", XMLPath: "root/hello", V: []byte("world")}
	result, err := op.Process(nil)
	if err != nil {
		t.Errorf("Process fail: %s", err)
		return
	}
	s := string(result)
	if !strings.Contains(s, "<hello>world</hello>") {
		t.Errorf("expected <hello>world</hello>, got: %s", s)
	}
}

func TestXMLProcessor_UnknownType(t *testing.T) {
	op := &driver.XMLProcessor{T: "invalid", XMLPath: "root/a"}
	_, err := op.Process(nil)
	if err == nil {
		t.Error("expected error for unknown type")
		return
	}
	if !strings.Contains(err.Error(), "unknown Processor type: invalid") {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestXMLProcessor_MarshalUnmarshal(t *testing.T) {
	d := driver.NewXMLDriver()

	original := &driver.XMLProcessor{
		T:       "create",
		XMLPath: "root/data",
		V:       []byte("test"),
		A:       "tester",
	}

	data := original.Save()

	var ops []driver.Processor
	ops, err := d.Unmarshal([]byte("[" + string(data) + "]"))
	if err != nil {
		t.Errorf("unmarshal fail: %s", err)
		return
	}
	if len(ops) != 1 {
		t.Errorf("expected 1 processor, got %d", len(ops))
		return
	}
	restored := ops[0].(*driver.XMLProcessor)
	if restored.T != original.T {
		t.Errorf("expected type %s, got %s", original.T, restored.T)
	}
	if restored.XMLPath != original.XMLPath {
		t.Errorf("expected xml_path %s, got %s", original.XMLPath, restored.XMLPath)
	}
	if string(restored.V) != string(original.V) {
		t.Errorf("expected value %s, got %s", original.V, restored.V)
	}
}
