package config

import (
	"bytes"
	"fmt"
	"testing"
)

func TestParseValidSectionHeader(t *testing.T) {
	const name = "Some section"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("[%s]\n", name))
	reader := NewReader(m)
	sections, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 {
		t.Fatalf("ReadAll should return just one section, returned %d sections", len(sections))
	}
	if _, ok := sections[name]; !ok {
		t.Errorf("ReadAll should return section named %q", name)
	}
}

func TestParseInvalidSectionHeader(t *testing.T) {
	m := new(bytes.Buffer)
	m.WriteString("[Some section\n")
	reader := NewReader(m)
	_, err := reader.ReadAll()
	if err.(*ParseError).Err != ErrParse {
		t.Errorf("ReadAll should return %q error, returned %q", ErrParse, err)
	}
}

func TestParseEmptySectionHeader(t *testing.T) {
	m := new(bytes.Buffer)
	m.WriteString("[]")
	reader := NewReader(m)
	_, err := reader.ReadAll()
	if err.(*ParseError).Err != ErrEmptySectionHeader {
		t.Errorf("ReadAll should return %q error, returned %q", ErrEmptySectionHeader, err)
	}
}

func TestDefaultParseOption(t *testing.T) {
	const name = "default"
	const key = "SomeOption"
	const value = "z = x + y"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("%s = %s # Some comment\n", key, value))
	reader := NewReader(m)
	sections, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 {
		t.Fatalf("ReadAll should return just one section, returned %d sections", len(sections))
	}
	if _, ok := sections[name]; !ok {
		t.Errorf("ReadAll should return section named %q", name)
	}
	if len(sections[name]) != 1 {
		t.Fatalf("ReadAll should return a section with just one option, returned %d options", len(sections[name]))
	}
	if actualValue, ok := sections[name][key]; !ok || actualValue != value {
		t.Errorf("ReadAll should return a section with option named %q and value %q, returned value %q", key, value, actualValue)
	}
}

func TestNondefaultParseOption(t *testing.T) {
	const name = "Some section"
	const key = "SomeOption"
	const value = "z = x + y"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("[%s]\n", name))
	m.WriteString(fmt.Sprintf("%s = %s # Some comment\n", key, value))
	reader := NewReader(m)
	sections, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 1 {
		t.Fatalf("ReadAll should return just one section, returned %d sections", len(sections))
	}
	if _, ok := sections[name]; !ok {
		t.Errorf("ReadAll should return section named %q", name)
	}
	if len(sections[name]) != 1 {
		t.Fatalf("ReadAll should return a section with just one option, returned %d options", len(sections[name]))
	}
	if actualValue, ok := sections[name][key]; !ok || actualValue != value {
		t.Errorf("ReadAll should return a section with option named %q and value %q, returned value %q", key, value, actualValue)
	}
}
