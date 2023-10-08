package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/tr1v3r/rule/driver"
)

func Test_Rule(t *testing.T) {
	var items = []RuleDataItem{
		{Path: "/", Operators: []struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}{
			{Type: "curl", Data: (&driver.CURLOperator{URL: "https://r1v3.com/ping", C: time.Now()}).Save()},
		}},
	}

	data, err := json.Marshal(items)
	if err != nil {
		t.Errorf("marshal rules fail: %s", err)
		return
	}
	t.Logf("got rules data: %s", data)
}
