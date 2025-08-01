// Package jsonvx provides a highly configurable parser, querier, and formatter for JSON-like data.
//
// jsonvx is designed to support both strict ECMA-404-compliant JSON and a wide range of non-standard,
// relaxed variants — including JSON5 and other formats commonly found in real-world data. Its core
// philosophy is flexibility: developers can fine-tune nearly every aspect of JSON parsing behavior
// through the ParserConfig struct, enabling or disabling specific features to suit various parsing needs.
//
// Key features:
//
//   - Parsing:
//
//   - Strict mode (standards-compliant) and relaxed mode (e.g., unquoted keys, single-quoted strings)
//
//   - Support for hexadecimal numbers, `Infinity`, `NaN`, and leading plus signs
//
//   - Support for trailing commas, line/block comments, and non-standard escape sequences
//
//   - Configurable whitespace handling and edge-case number formats
//
//   - Querying:
//
//   - Traverse deeply nested arrays and objects using path-based access
//
//   - Query scalar values with meaningful error reporting
//
//   - Type-safe access to JSON values via Go interfaces
//
//   - Formatting:
//
//   - Render parsed JSON back into readable, syntax-highlighted string representations
//
//   - Useful for diagnostics, debugging, or building REPL-like tools
//
// jsonvx is especially useful when working with JSON data from dynamic sources,
// legacy APIs, or configuration files that bend or break the formal rules of JSON.
//
// See the ParserConfig struct for all available parsing options.
//
// Specification references:
//   - ECMA-404: https://datatracker.ietf.org/doc/html/rfc7159
//   - JSON5: https://json5.org/
package jsonvx

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Common errors for type assertion and query path resolution.
var (
	ErrNotString  = errors.New("value is not a JSON string")
	ErrNotNumber  = errors.New("value is not a JSON number")
	ErrNotBoolean = errors.New("value is not a JSON boolean")
	ErrNotNull    = errors.New("value is not null")

	ErrInvalidQueryKey   = errors.New("invalid query key")
	ErrExpectedIndex     = errors.New("invalid query key, expected integer index")
	ErrIndexOutOfRange   = errors.New("index out of range")
	ErrEmptyArray        = errors.New("array is empty")
	ErrInvalidJSONType   = errors.New("invalid JSON type")
	ErrQueryExceedsDepth = errors.New("query exceeds depth for scalar value")
)

// JSON is a common interface implemented by all JSON types (Null, Boolean, etc.).
type JSON interface {
	fmt.Stringer
	Equal(JSON) bool
}

// Null represents a JSON null value.
type Null struct {
	Token *Token
}

// newNull creates a new *Null value, optionally invoking a callback.
func newNull(Token *Token, cb func()) *Null {
	if cb != nil {
		cb()
	}

	return &Null{Token: Token}
}

func (n *Null) String() string {
	return "\033[1mnull\033[0m"
}

// Value returns nil if the token is valid, or an error otherwise.
func (n *Null) Value() (any, error) {
	if n.Token == nil {
		return nil, ErrNotNull
	}

	return nil, nil
}

func (n *Null) Equal(n2 JSON) bool {
	if n == nil || n2 == nil {
		return n == n2
	}

	other, ok := AsNull(n2)

	if !ok {
		return false
	}

	return n.Token.Equal(other.Token)
}

// AsNull safely casts a JSON to a *Null.
func AsNull(j JSON) (*Null, bool) {
	null, ok := j.(*Null)
	return null, ok
}

// Boolean represents a JSON boolean (true or false).
type Boolean struct {
	Token *Token
}

// newBoolean creates a new *Boolean value, optionally invoking a callback
func newBoolean(Token *Token, cb func()) *Boolean {
	if cb != nil {
		cb()
	}

	return &Boolean{Token: Token}
}

func (b *Boolean) String() string {
	if b.Token.SubKind == TRUE {
		return "\033[1mtrue\033[0m"
	} else {
		return "\033[1mfalse\033[0m"
	}
}

// Value returns the boolean value or an error if the token is invalid.
func (b *Boolean) Value() (bool, error) {
	if b.Token == nil {
		return false, ErrNotBoolean
	}

	val := b.Token.Value()

	if val == nil {
		return false, ErrNotBoolean
	}

	boolVal, ok := val.(bool)

	if !ok {
		return false, ErrNotBoolean
	}

	return boolVal, nil
}

func (b *Boolean) Equal(b2 JSON) bool {
	if b == nil || b2 == nil {
		return b == b2
	}

	other, ok := AsBoolean(b2)

	if !ok {
		return false
	}

	return b.Token.Equal(other.Token)
}

// AsBoolean safely casts a JSON to a *Boolean.
func AsBoolean(j JSON) (*Boolean, bool) {
	boolean, ok := j.(*Boolean)
	return boolean, ok
}

// String represents a JSON string.
type String struct {
	Token *Token
}

// newString creates a new *String value, optionally invoking a callback
func newString(Token *Token, cb func()) *String {
	if cb != nil {
		cb()
	}

	return &String{Token: Token}
}

func (s *String) String() string {
	return (s.Token.Value()).(string)
}

// Value returns the string or an error if invalid.
func (s *String) Value() (string, error) {
	if s.Token == nil {
		return "", ErrNotString
	}

	val := s.Token.Value()

	if val == nil {
		return "", ErrNotString
	}

	strVal, ok := val.(string)

	if !ok {
		return "", ErrNotString
	}

	return strVal, nil
}

func (s *String) Equal(s2 JSON) bool {
	if s == nil || s2 == nil {
		return s == s2
	}

	other, ok := AsString(s)

	if !ok {
		return false
	}

	return s.Token.Equal(other.Token)
}

// AsString safely casts a JSON to a *String.
func AsString(j JSON) (*String, bool) {
	str, ok := j.(*String)
	return str, ok
}

// Number represents a JSON number (integer, float, hex, etc.).
type Number struct {
	Token *Token
}

// newNumber creates a new *Number value, optionally invoking a callback
func newNumber(Token *Token, cb func()) *Number {
	if cb != nil {
		cb()
	}

	return &Number{Token: Token}
}

func (n *Number) String() string {
	return string((n.Token.Literal))
}

// Value attempts to parse the number as float64, depending on its subtype.
func (n *Number) Value() (float64, error) {
	if n.Token == nil {
		return 0, ErrNotNumber
	}

	switch n.Token.SubKind {
	case INTEGER, HEX:
		numVal, err := ToInt(n.Token.Literal)
		if err != nil {
			return 0, ErrNotNumber
		}
		return float64(numVal), nil
	case FLOAT, SCI_NOT:
		numVal, err := ToFloat(n.Token.Literal)
		if err != nil {
			return 0, ErrNotNumber
		}
		return numVal, nil
	case INF:
		if bytes.Equal(n.Token.Literal, []byte("Infinity")) ||
			bytes.Equal(n.Token.Literal, []byte("+Infinity")) {
			return math.Inf(1), nil
		}
		if bytes.Equal(n.Token.Literal, []byte("-Infinity")) {
			return math.Inf(-1), nil
		}
		return 0, ErrNotNumber
	case NaN:
		if bytes.Equal(n.Token.Literal, []byte("NaN")) ||
			bytes.Equal(n.Token.Literal, []byte("+NaN")) ||
			bytes.Equal(n.Token.Literal, []byte("-NaN")) {
			return math.NaN(), nil
		}
		return 0, ErrNotNumber

	default:
		return 0, ErrNotNumber
	}
}

func (n *Number) Equal(n2 JSON) bool {
	if n == nil || n2 == nil {
		return n == n2
	}

	other, ok := AsNumber(n2)

	if !ok {
		return false
	}

	return n.Token.Equal(other.Token)
}

// AsNumber safely casts a JSON to a *Number.
func AsNumber(j JSON) (*Number, bool) {
	number, ok := j.(*Number)
	return number, ok
}

// Array represents a JSON array.
type Array struct {
	Items []JSON
}

// newArray creates a new *Array value, optionally invoking a callback
func newArray(items []JSON, cb func()) *Array {
	if cb != nil {
		cb()
	}

	return &Array{Items: items}
}

func (a *Array) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, item := range a.Items {
		if i == len(a.Items)-1 {
			builder.WriteString(fmt.Sprintf("%v", item))
		} else {
			builder.WriteString(fmt.Sprintf("%v, ", item))
		}
	}
	builder.WriteString("]")
	return builder.String()
}

func (a *Array) Len() int {
	return len(a.Items)
}

// QueryPath retrieves a nested item using a slice of string indices.
func (a *Array) QueryPath(paths ...string) (JSON, error) {
	if len(paths) == 0 {
		return a, nil
	}

	indexStr := paths[0]
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrExpectedIndex, indexStr)
	}

	if a.Len() == 0 {
		return nil, ErrEmptyArray
	}

	if index < 0 || index >= a.Len() {
		return nil, fmt.Errorf("%w: %d", ErrIndexOutOfRange, index)
	}

	item := a.Items[index]

	rest := paths[1:]

	switch val := item.(type) {
	case *Null, *Boolean, *Number, *String:
		if len(rest) > 0 {
			return nil, ErrQueryExceedsDepth
		}

		return val, nil
	case *Array:
		return val.QueryPath(rest...)
	case *Object:
		return val.QueryPath(rest...)
	default:
		return nil, ErrInvalidJSONType
	}
}

// ArrayCallback defines the function signature for iterating over items in a JSON array.
// - item: the current element
// - index: the index of the element
// - array: the array being iterated
type ArrayCallback func(item JSON, index int, array *Array)

// ForEach calls the given callback for each item in the array.
// It provides the item, its index, and the array itself.
func (a *Array) ForEach(cb ArrayCallback) {
	for i, item := range a.Items {
		cb(item, i, a)
	}
}

func (a *Array) Equal(a2 JSON) bool {
	if a == nil || a2 == nil {
		return a == a2
	}

	other, ok := AsArray(a2)

	if !ok {
		return false
	}

	if a.Len() != other.Len() {
		return false
	}

	for i, item := range a.Items {
		item2 := other.Items[i]

		if !item.Equal(item2) {
			return false
		}

	}

	return true
}

// AsArray safely casts a JSON to a *Array.
func AsArray(j JSON) (*Array, bool) {
	arr, ok := j.(*Array)
	return arr, ok
}

// KeyValue represents a key-value pair in a JSON object.
type KeyValue struct {
	key   []byte
	value JSON
}

func (kv *KeyValue) Equal(kv2 *KeyValue) bool {
	if kv == nil || kv2 == nil {
		return kv == kv2
	}

	if !bytes.Equal(kv.key, kv2.key) {
		return false
	}

	return kv.value.Equal(kv2.value)
}

// Object represents a JSON object.
type Object struct {
	Properties []KeyValue
}

// newObject creates a new *Object value, optionally invoking a callback
func newObject(properties []KeyValue, cb func()) *Object {
	if cb != nil {
		cb()
	}

	return &Object{Properties: properties}
}

func (o *Object) String() string {
	var builder strings.Builder
	builder.WriteString("{")
	length := len(o.Properties)
	for ind, kv := range o.Properties {
		if ind == length-1 {
			builder.WriteString(fmt.Sprintf("%s: %v", kv.key, kv.value))
		} else {
			builder.WriteString(fmt.Sprintf("%s: %v, ", kv.key, kv.value))
		}
	}
	builder.WriteString("}")
	return builder.String()
}

func (o *Object) Len() int {
	return len(o.Properties)
}

// QueryPath retrieves a nested value via key-based path traversal.
func (o *Object) QueryPath(paths ...string) (JSON, error) {
	if len(paths) == 0 {
		return o, nil
	}

	keyStr := paths[0]
	keyBytes := []byte(keyStr)

	index := sort.Search(len(o.Properties), func(i int) bool {
		return bytes.Compare(o.Properties[i].key, keyBytes) >= 0
	})

	if index >= o.Len() || !bytes.Equal(o.Properties[index].key, keyBytes) {
		return nil, fmt.Errorf("key not found: %q", keyStr)
	}

	item := o.Properties[index]
	rest := paths[1:]

	switch val := item.value.(type) {
	case *Null, *Boolean, *Number, *String:
		if len(rest) > 0 {
			return nil, ErrQueryExceedsDepth
		}
		return val, nil
	case *Array:
		return val.QueryPath(rest...)
	case *Object:
		return val.QueryPath(rest...)
	default:
		return nil, ErrInvalidJSONType
	}
}

// ObjectCallback defines the function signature for iterating over properties in a JSON object.
// - key: the property's key as a byte slice
// - value: the property's value
// - object: the object being iterated
type ObjectCallback func(key []byte, value JSON, object *Object)

// ForEach calls the given callback for each key-value pair in the object.
// It provides the key, value, and the object itself.
func (o *Object) ForEach(cb ObjectCallback) {
	for _, prop := range o.Properties {
		cb(prop.key, prop.value, o)
	}
}

func (o *Object) Equal(o2 JSON) bool {
	other, ok := AsObject(o2)

	if !ok {
		return false
	}

	if o.Len() != other.Len() {
		return false
	}

	for i, prop := range o.Properties {
		prop2 := other.Properties[i]

		if !prop.Equal(&prop2) {
			return false
		}
	}

	return true
}

// AsObject safely casts a JSON to a *Object.
func AsObject(j JSON) (*Object, bool) {
	obj, ok := j.(*Object)
	return obj, ok
}
