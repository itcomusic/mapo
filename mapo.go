// Package maporder provides a map that maintains the order of insertion.
package mapo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

const (
	delimObjectOpen  = '{'
	delimObjectClose = '}'
	delimArrayOpen   = '['
	delimArrayClose  = ']'
)

type Map struct {
	value map[string]any
	key   []string
}

// New creates a new Map.
func New() *Map {
	return &Map{value: make(map[string]any)}
}

// NewWithSize creates a new Map with the given capacity.
func NewWithSize(cap int) *Map {
	return &Map{
		value: make(map[string]any, cap),
		key:   make([]string, 0, cap),
	}
}

// Set sets the key to value.
//
// If the key already exists, the order is not updated.
func (m *Map) Set(key string, value any) {
	if m.value == nil {
		m.value = make(map[string]any)
	}

	if _, ok := m.value[key]; !ok {
		m.key = append(m.key, key)
	}
	m.value[key] = value
}

// Get returns the value associated with the key.
func (m *Map) Get(key string) (any, bool) {
	if m.value == nil {
		return nil, false
	}
	v, ok := m.value[key]
	return v, ok
}

// Keys returns the keys in the map.
func (m *Map) Keys() []string {
	return m.key
}

// Delete deletes the key.
func (m *Map) Delete(key string) {
	delete(m.value, key)
	for i := range m.key {
		if m.key[i] == key {
			m.key = append(m.key[:i], m.key[i+1:]...)
			return
		}
	}
}

func (m *Map) UnmarshalJSON(b []byte) error {
	if m.value == nil {
		m.value = make(map[string]any)
	}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()

	// must open with a delim token '{'
	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != delimObjectOpen {
		return fmt.Errorf("expect object open with %q", delimObjectOpen)
	}

	err = m.parseObject(dec)
	if err != nil {
		return err
	}

	t, err = dec.Token()
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("expect end of object but got more token: %T: %v: %w", t, t, err)
	}
	return nil
}

func (m *Map) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte(delimObjectOpen)

	encoder := json.NewEncoder(&buf)
	for i, key := range m.key {
		if i > 0 {
			buf.WriteByte(',')
		}

		// add key
		if err := encoder.Encode(key); err != nil {
			return nil, err
		}
		buf.WriteByte(':')

		// add value
		if err := encoder.Encode(m.value[key]); err != nil {
			return nil, err
		}
	}

	buf.WriteByte(delimObjectClose)
	return buf.Bytes(), nil
}

func (m *Map) parseObject(dec *json.Decoder) (err error) {
	var t json.Token
	for dec.More() {
		t, err = dec.Token()
		if err != nil {
			return err
		}

		key, ok := t.(string)
		if !ok {
			return fmt.Errorf("expecting key should be always a string: %T: %v", t, t)
		}

		t, err = dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		var value any
		value, err = handleDelim(t, dec)
		if err != nil {
			return err
		}

		m.key = append(m.key, key)
		m.value[key] = value
	}

	t, err = dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != delimObjectClose {
		return fmt.Errorf("expect object close with %q", delimObjectClose)
	}
	return nil
}

func parseArray(dec *json.Decoder) ([]any, error) {
	var t json.Token
	res := make([]any, 0)
	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return nil, err
		}

		var value any
		value, err = handleDelim(t, dec)
		if err != nil {
			return nil, err
		}
		res = append(res, value)
	}

	t, err := dec.Token()
	if err != nil {
		return nil, err
	}

	if delim, ok := t.(json.Delim); !ok || delim != delimArrayClose {
		return nil, fmt.Errorf("expect array close with %q", delimArrayClose)
	}
	return res, nil
}

func handleDelim(t json.Token, dec *json.Decoder) (any, error) {
	if delim, ok := t.(json.Delim); ok {
		switch delim {
		case delimObjectOpen:
			om2 := New()
			if err := om2.parseObject(dec); err != nil {
				return nil, err
			}
			return om2, nil

		case delimArrayOpen:
			value, err := parseArray(dec)
			if err != nil {
				return nil, err
			}
			return value, nil

		default:
			return nil, fmt.Errorf("unexpected delimiter: %q", delim)
		}
	}
	return t, nil
}
