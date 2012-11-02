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

func TestParseInvalidSectionHeader1(t *testing.T) {
	m := new(bytes.Buffer)
	m.WriteString("[Some section")
	reader := NewReader(m)
	_, err := reader.ReadAll()
	if err.(*ParseError).Err != ErrInvalidSectionHeader {
		t.Errorf("ReadAll should return %q error, returned %q", ErrInvalidSectionHeader, err)
	}
}

func TestParseInvalidSectionHeader2(t *testing.T) {
	m := new(bytes.Buffer)
	m.WriteString("[Some section[\n")
	reader := NewReader(m)
	_, err := reader.ReadAll()
	if err.(*ParseError).Err != ErrInvalidSectionHeader {
		t.Errorf("ReadAll should return %q error, returned %q", ErrInvalidSectionHeader, err)
	}
}

func TestParseInvalidSectionHeader3(t *testing.T) {
	m := new(bytes.Buffer)
	m.WriteString("[Some section]\n[Oops\noption = value")
	reader := NewReader(m)
	_, err := reader.ReadAll()
	if err.(*ParseError).Err != ErrInvalidSectionHeader {
		t.Errorf("ReadAll should return %q error, returned %q", ErrInvalidSectionHeader, err)
	}
}

func TestParseEmptySectionHeader(t *testing.T) {
	m := new(bytes.Buffer)
	m.WriteString("[]")
	reader := NewReader(m)
	_, err := reader.ReadAll()
	if err.(*ParseError).Err != ErrInvalidSectionHeader {
		t.Errorf("ReadAll should return %q error, returned %q", ErrInvalidSectionHeader, err)
	}
}

func TestParseEmptyLines(t *testing.T) {
	const section1 = "Section 1"
	const section2 = "Section 2"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("\n[%s]\n\n[%s]", section1, section2))
	reader := NewReader(m)
	sections, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 2 {
		t.Fatalf("ReadAll should return two sections, returned %d sections", len(sections))
	}
	if _, ok := sections[section1]; !ok {
		t.Errorf("ReadAll should return section named %q", section1)
	}
	if _, ok := sections[section2]; !ok {
		t.Errorf("ReadAll should return section named %q", section2)
	}
}

func TestParseCommentLines(t *testing.T) {
	const section1 = "Section 1"
	const section2 = "Section 2"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("\n[%s]\n#[Comment section]\n\t#[Another comment section]\n[%s]", section1, section2))
	reader := NewReader(m)
	sections, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(sections) != 2 {
		t.Fatalf("ReadAll should return two sections, returned %d sections", len(sections))
	}
	if _, ok := sections[section1]; !ok {
		t.Errorf("ReadAll should return section named %q", section1)
	}
	if _, ok := sections[section2]; !ok {
		t.Errorf("ReadAll should return section named %q", section2)
	}
}

func TestDefaultParseOption(t *testing.T) {
	const name = "default"
	const key = "Some option"
	const value = "some value"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("%s = %s\n", key, value))
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

func TestParseOption1(t *testing.T) {
	const name = "Some section"
	const key = "SomeOption"
	const value = "z = x + y"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("[%s]\n", name))
	m.WriteString(fmt.Sprintf("%s = %s\n", key, value))
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

func TestParseOption2(t *testing.T) {
	const name = "Some section"
	const key = "SomeOption"
	const value = "z = x + y"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("[%s]\n", name))
	m.WriteString(fmt.Sprintf("%s: %s\n", key, value))
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

func TestParseOption3(t *testing.T) {
	const name = "Some section"
	const key = "SomeOption"
	const value = "z = x + y"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("[%s]\n", name))
	m.WriteString(fmt.Sprintf("\t\t  %s: %s\n", key, value))
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

func TestParseOptionWithComment(t *testing.T) {
	const name = "Some section"
	const key = "SomeOption"
	const value = "z = x + y"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("[%s]\n", name))
	m.WriteString(fmt.Sprintf("%s: %s # Some comment\n", key, value))
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

func TestParseOptionWithoutAComment(t *testing.T) {
	const name = "Some section"
	const key = "SomeOption"
	const value = "z = x + y#Not a comment"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("[%s]\n", name))
	m.WriteString(fmt.Sprintf("%s: %s\n", key, value))
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
