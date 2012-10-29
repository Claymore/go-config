package config

import (
    "bytes"
    "testing"
    "fmt"
)

func TestReadAll(t *testing.T) {
    const name = "Some section"
    m := new(bytes.Buffer)
    m.WriteString(fmt.Sprintf("[%s]", name))
    reader := NewReader(m)
    sections, err := reader.ReadAll()
    if err != nil {
        t.Fatal(err)
    }
    if len(sections) != 1 {
        t.Fatalf("ReadAll should return just one section, returned %d sections", len(sections))
    }
    if sections[0].Name != name {
        t.Errorf("ReadAll should return section named %q, returned %q", name, sections[0].Name)
    }
}
