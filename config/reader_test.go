package config

import (
	"bytes"
	"fmt"
	"testing"
)

func TestReadAll(t *testing.T) {
	const name = "Some section"
	const key = "SomeOption"
	const value = "z = x + y"
	m := new(bytes.Buffer)
	m.WriteString(fmt.Sprintf("[%s]\n", name))
	m.WriteString(fmt.Sprintf("%s = %s # Some comment\n", key, value))
	m.WriteString("\n\n#[Another section]")
	reader := NewReader(m)
	sections, err := reader.ReadAll()
	fmt.Println(sections)
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
		t.Fatalf("ReadAll should return a section with just 1 option, returned %d options", len(sections[name]))
	}
	if actualValue, ok := sections[name][key]; !ok || actualValue != value {
		t.Errorf("ReadAll should return a section with option named %q and value %q, returned value %q", key, value, actualValue)
	}
}
