package flowt

import (
	"encoding/json"
	"testing"

	"github.com/Jeffail/gabs"
)

func TestStringEqualsAString(t *testing.T) {
	input, _ := gabs.ParseJSON([]byte(`{"a": {"b": {"c": "abc"}}}`))
	var ch Choice
	data := `{"Variable": "$.a.b.c", "StringEquals": "abc"}`
	if err := json.Unmarshal([]byte(data), &ch); err != nil {
		t.FailNow()
	}

	if ch.IsSatisfied(input) != true {
		t.FailNow()
	}
}

func TestStringEqualsEmptyString(t *testing.T) {
	input, _ := gabs.ParseJSON([]byte(`{"a": {"b": {"c": ""}}}`))
	var ch Choice
	data := `{"Variable": "$.a.b.c", "StringEquals": ""}`
	if err := json.Unmarshal([]byte(data), &ch); err != nil {
		t.FailNow()
	}

	if ch.IsSatisfied(input) != true {
		t.FailNow()
	}
}

func TestStringNotEqualsAString(t *testing.T) {
	input, _ := gabs.ParseJSON([]byte(`{"a": {"b": {"c": "xxx"}}}`))
	var ch Choice
	data := `{"Variable": "$.a.b.c", "StringEquals": "abc"}`
	if err := json.Unmarshal([]byte(data), &ch); err != nil {
		t.FailNow()
	}

	if ch.IsSatisfied(input) != false {
		t.FailNow()
	}
}

func TestStringEqualsButNotAString(t *testing.T) {
	input, _ := gabs.ParseJSON([]byte(`{"a": {"b": {"c": 55}}}`))
	var ch Choice
	data := `{"Variable": "$.a.b.c", "StringEquals": "abc"}`
	if err := json.Unmarshal([]byte(data), &ch); err != nil {
		t.FailNow()
	}

	if ch.IsSatisfied(input) != false {
		t.FailNow()
	}
}

func TestIsPresent(t *testing.T) {
	input, _ := gabs.ParseJSON([]byte(`{"a": {"b": {"c": "abc"}}}`))
	var ch Choice
	data := `{"Variable": "$.a.b.c", "IsPresent": true}`
	if err := json.Unmarshal([]byte(data), &ch); err != nil {
		t.FailNow()
	}

	if ch.IsSatisfied(input) != true {
		t.FailNow()
	}
}

func TestIsNotPresent(t *testing.T) {
	input, _ := gabs.ParseJSON([]byte(`{"a": {"b": {}}}`))
	var ch Choice
	data := `{"Variable": "$.a.b.c", "IsPresent": false}`
	if err := json.Unmarshal([]byte(data), &ch); err != nil {
		t.FailNow()
	}

	if ch.IsSatisfied(input) != true {
		t.FailNow()
	}
}

func TestNotStringEqualsAEqualString(t *testing.T) {
	input, _ := gabs.ParseJSON([]byte(`{"a": {"b": {"c": "abc"}}}`))
	var ch Choice
	data := `{"Not": {"Variable": "$.a.b.c", "StringEquals": "abc"}}`
	if err := json.Unmarshal([]byte(data), &ch); err != nil {
		t.FailNow()
	}

	if ch.IsSatisfied(input) == true {
		t.FailNow()
	}
}

func TestNotStringEqualsAString(t *testing.T) {
	input, _ := gabs.ParseJSON([]byte(`{"a": {"b": {"c": "xxx"}}}`))
	var ch Choice
	data := `{"Not": {"Variable": "$.a.b.c", "StringEquals": "abc"}}`
	if err := json.Unmarshal([]byte(data), &ch); err != nil {
		t.FailNow()
	}

	if ch.IsSatisfied(input) == false {
		t.FailNow()
	}
}

func TestStringEqualsAStringPath(t *testing.T) {
	input, _ := gabs.ParseJSON([]byte(`{"x": "abc", "a": {"b": {"c": "abc"}}}`))
	var ch Choice
	data := `{"Variable": "$.a.b.c", "StringEqualsPath": "$.x"}`
	if err := json.Unmarshal([]byte(data), &ch); err != nil {
		t.FailNow()
	}

	if ch.IsSatisfied(input) != true {
		t.FailNow()
	}
}
