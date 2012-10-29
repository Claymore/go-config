package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

// A ParseError is returned for parsing errors.
// The first line is 1. The first column is 0.
type ParseError struct {
	Line   int   // Line where the error occurred
	Column int   // Column (rune index) where the error occurred
	Err    error // The actual error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Err)
}

// These are the errors that can be returned in ParseError.Error
var (
	ErrParse = errors.New("generic parse error")
)

type Section struct {
	Name    string
	Options map[string]string
}

func NewSection(name string) *Section {
	return &Section{
		Name:    name,
		Options: make(map[string]string),
	}
}

type Reader struct {
	r      *bufio.Reader
	field  bytes.Buffer
	line   int
	column int
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

// error creates a new ParseError based on err.
func (r *Reader) error(err error) error {
	return &ParseError{
		Line:   r.line,
		Column: r.column,
		Err:    err,
	}
}

func (r *Reader) Read() (section *Section, err error) {
	for {
		section, err = r.parseSection()
		if section != nil {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return section, nil
}

func (r *Reader) ReadAll() (sections []Section, err error) {
	for {
		section, err := r.Read()
		if err == io.EOF {
			return sections, nil
		}
		if err != nil {
			return nil, err
		}
		sections = append(sections, *section)
	}
	panic("unreachable")
}

func (r *Reader) skip(delim rune) error {
	for {
		r1, _, err := r.r.ReadRune()
		if err != nil {
			return err
		}
		if r1 == delim {
			return nil
		}
	}
	panic("unreachable")
}

func (r *Reader) parseSection() (section *Section, err error) {
	r.field.Reset()
	r1, _, err := r.r.ReadRune()
	if err != nil {
		return nil, err
	}
	if r1 == '#' {
		return nil, r.skip('\n')
	}

	if r1 == '[' {
		for {
			r1, _, err := r.r.ReadRune()
			if err != nil {
				return nil, err
			}
			if r1 == ']' {
				section = NewSection(r.field.String())
				r.skip('\n')
				break
			}
			r.field.WriteRune(r1)
		}
	}

	if section != nil {
		for {
			r1, _, err := r.r.ReadRune()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			r.r.UnreadRune()
			if r1 == '[' {
				break
			}
			hasOption, key, value, err := r.parseOption()
			if hasOption {
				section.Options[key] = value
			}
			if err != nil {
				return nil, err
			}
		}
	}

	return section, err
}

func (r *Reader) parseOption() (haveOption bool, key string, value string, err error) {
	r.field.Reset()

	var foundDelim bool
	for {
		r1, _, err := r.r.ReadRune()
		if err == io.EOF {
			value = r.field.String()
			return foundDelim, key, value, nil
		}
		if err != nil {
			return false, key, value, err
		}
		if !foundDelim {
			if r1 == '=' {
				key = r.field.String()
				r.skip(' ')
				foundDelim = true
				r.field.Reset()
				continue
			} else if r1 != ' ' {
				r.field.WriteRune(r1)
			}
		} else {
			if r1 != '\n' {
				r.field.WriteRune(r1)
			}
		}
		if r1 == '\n' {
			value = r.field.String()
			return foundDelim, key, value, nil
		}
	}
	panic("unreachable")
}
