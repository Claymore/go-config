// Package config reads INI configuration files.
//
// A configaration file contains zero or more sections of zero or more options
// per section. Each section is separated by the newline character. Options are
// also separated by the newline character.
//
// Carriage returns before newline characters are silently removed.
//
// Blank lines and lines with only whitespace characters are ignored. Lines
// beginning with comment characters ('#' and ';') are also ignored.
//
// Each section must have a header wrapped in square brackets. Any option
// appearing on the next line after a section header will belong to this section.
//
// An options consists of a name and a value separated with ':' or '=' characters.
// Leading and trailing spaces will be trimmed from options names. There might be
// options without a value.

package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// A ParseError is returned for parsing errors.
// The first line is 1. The first column is 1.
type ParseError struct {
	Line   int   // Line where the error occurred
	Column int   // Column (rune index) where the error occurred
	Err    error // The actual error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Err)
}

// These are the errors that can be returned in ParseError.Error.
var (
	ErrParse              = errors.New("generic parse error")
	ErrEmptySectionHeader = errors.New("empty section header")
)

// A Reader reads sections of options from a configation file.
type Reader struct {
	r              *bufio.Reader
	field          bytes.Buffer
	line           int
	column         int
	currentSection string
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r:              bufio.NewReader(r),
		currentSection: "default",
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

// readRune reads one rune from r, folding \r\n to \n and keeping track
// of how far into the line we have read. r.column will point to the start
// of this rune, not the end of this rune.
func (r *Reader) readRune() (rune, error) {
	r1, _, err := r.r.ReadRune()

	// Handle \r\n here.  We make the simplifying assumption that
	// anytime \r is followed by \n that it can be folded to \n.
	// We will not detect text which contains both \r\n and bare \n.
	if r1 == '\r' {
		r1, _, err = r.r.ReadRune()
		if err == nil {
			if r1 != '\n' {
				r.r.UnreadRune()
				r1 = '\r'
			}
		}
	}
	r.column++
	return r1, err
}

// unreadRune puts the last rune read from r back.
func (r *Reader) unreadRune() {
	r.r.UnreadRune()
	r.column--
}

// ReadAll reads all the sections from r.
// Each section is a map.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read until EOF, it does not treat end of file as an error to be
// reported.
func (r *Reader) ReadAll() (sections map[string]map[string]string, err error) {
	sections = make(map[string]map[string]string)
	for {
		r.line++
		r.column = 0
		r1, err := r.readRune()
		if err == io.EOF {
			return sections, nil
		} else if err != nil {
			return nil, err
		}

		switch r1 {
		case '#', ';':
			r.skip('\n')
		case '[':
			section, err := r.parseHeader()
			if err != nil {
				return nil, err
			}
			if _, ok := sections[section]; !ok {
				sections[section] = make(map[string]string)
			}
			r.currentSection = section
		default:
			r.unreadRune()
			key, value, err := r.parseOption()
			if err != nil {
				return nil, err
			}
			key = strings.TrimSpace(key)

			if len(key) != 0 {
				if r.currentSection == "default" {
					if _, ok := sections["default"]; !ok {
						sections["default"] = make(map[string]string)
					}
				}
				sections[r.currentSection][key] = value
			}
		}
	}
	panic("unreachable")
}

// skip reads runes up to and including the rune delim or until error.
func (r *Reader) skip(delim rune) error {
	for {
		r1, err := r.readRune()
		if err != nil {
			return err
		}
		if r1 == delim {
			return nil
		}
	}
	panic("unreachable")
}

func (r *Reader) parseHeader() (section string, err error) {
	r.field.Reset()
	for {
		r1, err := r.readRune()
		if err != nil {
			return section, err
		}

		switch r1 {
		case '#', ';':
			return section, r.error(ErrParse)
		case ']':
			section = r.field.String()
			if section == "" {
				return section, r.error(ErrEmptySectionHeader)
			}
			err = r.skip('\n')
			if err == nil || err == io.EOF {
				return section, nil
			} else {
				return section, err
			}
		default:
			r.field.WriteRune(r1)
		}
	}
	panic("unreachable")
}

func (r *Reader) parseOption() (key string, value string, err error) {
	r.field.Reset()
	var (
		lastRune   rune
		foundDelim bool
	)
	for {
		r1, err := r.readRune()
		if err == io.EOF {
			value = r.field.String()
			return key, value, nil
		}
		if err != nil {
			return key, value, err
		}

		if (lastRune == 0 || lastRune == ' ') && (r1 == '#' || r1 == ';') {
			value = r.field.String()
			return key, value[:len(value)-1], r.skip('\n')
		}

		switch r1 {
		case '=', ':':
			if !foundDelim {
				key = r.field.String()
				r.skip(' ')
				foundDelim = true
				r.field.Reset()
			} else {
				r.field.WriteRune(r1)
				lastRune = r1
			}
		case '\n':
			value = r.field.String()
			return key, value, nil
		default:
			r.field.WriteRune(r1)
			lastRune = r1
		}
	}
	panic("unreachable")
}
