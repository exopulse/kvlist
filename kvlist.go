package kvlist

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// KeyValue holds key and value.
type KeyValue struct {
	Key   string
	Value string
}

func (kv *KeyValue) Write(p []byte) (int, error) {
	b := bytes.FieldsFunc(p, parseFunction())

	if len(b) != 1 {
		return 0, errors.New("key-value field format error")
	}

	s := string(b[0])
	eq := strings.Index(s, "=")

	if eq < 0 {
		return 0, errors.Errorf("key-value assignment error: %s", s)
	}

	v, err := strconv.Unquote(s[eq+1:])

	if err != nil {
		return 0, err
	}

	kv.Key = s[:eq]
	kv.Value = v

	return 0, nil
}

func parseFunction() func(rune) bool {
	inQuotes := false
	inEscapes := false

	return func(r rune) bool {
		if inEscapes {
			inEscapes = false

			return false
		}

		switch {
		case r == '"':
			inQuotes = !inQuotes
		case r == '\\':
			inEscapes = true
		case unicode.IsSpace(r):
			if !inQuotes {
				return true
			}
		}

		return false
	}
}

func (kv *KeyValue) Read(p []byte) (int, error) {
	s := fmt.Sprintf("%s=%s", kv.Key, strconv.QuoteToASCII(kv.Value))

	return copy(p, s), io.EOF
}

func (kv *KeyValue) String() string {
	return fmt.Sprintf("%s=%s", kv.Key, kv.Value)
}

// KeyValueList holds multiple key-value pairs. Keys are case sensitive. This list is mutable and not thread-safe.
type KeyValueList struct {
	list []*KeyValue
}

// Add adds new key-value pair. Duplicate keys are allowed.
func (l *KeyValueList) Add(key, value string) *KeyValueList {
	return l.AddKeyValue(KeyValue{key, value})
}

// AddKeyValue adds new key-value pair. Duplicate keys are allowed.
func (l *KeyValueList) AddKeyValue(kv KeyValue) *KeyValueList {
	l.list = append(l.list, &kv)

	return l
}

// Put adds new key-value pair. Existing key is removed.
func (l *KeyValueList) Put(key, value string) *KeyValueList {
	return l.PutKeyValue(KeyValue{key, value})
}

// PutKeyValue adds new key-value pair. Existing key is removed.
func (l *KeyValueList) PutKeyValue(kv KeyValue) *KeyValueList {
	for _, e := range l.list {
		if e.Key == kv.Key {
			e.Value = kv.Value
			return l
		}
	}

	l.list = append(l.list, &kv)

	return l
}

// DeleteKey deletes key-value pair by key. Only first occurrence of key is deleted.
// Method returns true if key was deleted.
func (l *KeyValueList) DeleteKey(key string) bool {
	for i, e := range l.list {
		if e.Key == key {
			l.list = append(l.list[:i], l.list[i+1:]...)

			return true
		}
	}

	return false
}

// DeleteKeys deletes all key-value pair by key. Method returns true if any key was deleted.
func (l *KeyValueList) DeleteKeys(key string) bool {
	deletedAny := false

	for true {
		deleted := l.DeleteKey(key)

		if deleted {
			deletedAny = true
		} else {
			break
		}
	}

	return deletedAny
}

// Count returns count of key-value pairs in list.
func (l *KeyValueList) Count() int {
	return len(l.list)
}

// Clear removes all items from the list.
func (l *KeyValueList) Clear() *KeyValueList {
	l.list = nil

	return l
}

// Get returns key-value pair by index.
func (l *KeyValueList) Get(index int) KeyValue {
	return *l.list[index]
}

// GetKey returns key value by key.
func (l *KeyValueList) GetKey(key string) (string, bool) {
	for _, e := range l.list {
		if e.Key == key {
			return e.Value, true
		}
	}

	return "", false
}

// ScanKey populates value receiver with key value.
func (l *KeyValueList) ScanKey(receiver *string, key string) bool {
	for _, e := range l.list {
		if e.Key == key {
			*receiver = e.Value

			return true
		}
	}

	return false
}

// GetKeys returns key-value pairs by key.
func (l *KeyValueList) GetKeys(key string) []KeyValue {
	var kvs []KeyValue

	for _, e := range l.list {
		if e.Key == key {
			kvs = append(kvs, *e)
		}
	}

	return kvs
}

// Items returns all key-value as a slice.
func (l *KeyValueList) Items() []KeyValue {
	var kvs []KeyValue

	for _, e := range l.list {
		kvs = append(kvs, *e)
	}

	return kvs
}

func (l *KeyValueList) String() string {
	return fmt.Sprintf("%s", l.list)
}

func (l *KeyValueList) Read(p []byte) (int, error) {
	buffer := new(bytes.Buffer)

	for _, e := range l.list {
		if buffer.Len() > 0 {
			_, err := buffer.WriteRune(' ')

			if err != nil {
				return 0, err
			}
		}

		if _, err := buffer.ReadFrom(e); err != nil && err != io.EOF {
			return 0, err
		}
	}

	return copy(p, buffer.Bytes()), io.EOF
}

func (l *KeyValueList) Write(p []byte) (int, error) {
	bs := bytes.FieldsFunc(p, parseFunction())

	for _, b := range bs {
		s := string(b)
		eq := strings.Index(s, "=")

		if eq < 0 {
			return 0, errors.Errorf("key-value assignment error: %s", s)
		}

		v, err := strconv.Unquote(s[eq+1:])

		if err != nil {
			return 0, err
		}

		l.list = append(l.list, &KeyValue{s[:eq], v})
	}

	return 0, nil
}

// New creates new key-value list.
func New() *KeyValueList {
	return new(KeyValueList)
}

// NewFromString creates new key-value list from string.
func NewFromString(s string) (*KeyValueList, error) {
	l := new(KeyValueList)

	if _, err := l.Write([]byte(s)); err != nil {
		return nil, err
	}

	return l, nil
}
