package jsonvx

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotString  = errors.New("value is not a JSON string")
	ErrNotNumber  = errors.New("value is not a JSON number")
	ErrNotBoolean = errors.New("value is not a JSON boolean")
	ErrNotNull    = errors.New("value is not null")
)

type JSON interface {
	String() string
}

// Null type starts here
type Null struct {
	Token *Token
}

func newNull(Token *Token, cb func()) Null {
	if cb != nil {
		cb()
	}

	return Null{Token: Token}
}

func (n Null) String() string {
	return "\033[1mnull\033[0m"
}

func (n Null) Value() (any, error) {
	if n.Token == nil {
		return nil, ErrNotNull
	}

	return nil, nil
}

func AsNull(j JSON) (*Null, bool) {
	null, ok := j.(Null)
	return &null, ok
}

// Boolean type starts here
type Boolean struct {
	Token *Token
}

func newBoolean(Token *Token, cb func()) Boolean {
	if cb != nil {
		cb()
	}

	return Boolean{Token: Token}
}

func (b Boolean) String() string {
	if b.Token.SubKind == TRUE {
		return "\033[1mtrue\033[0m"
	} else {
		return "\033[1mfalse\033[0m"
	}
}

func (b Boolean) Value() (bool, error) {
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

func AsBoolean(j JSON) (*Boolean, bool) {
	boolean, ok := j.(Boolean)
	return &boolean, ok
}

// String type starts here
type String struct {
	Token *Token
}

func newString(Token *Token, cb func()) String {
	if cb != nil {
		cb()
	}

	return String{Token: Token}
}

func (s String) String() string {
	return (s.Token.Value()).(string)
}

func (s String) Value() (string, error) {
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

func AsString(j JSON) (*String, bool) {
	str, ok := j.(String)
	return &str, ok
}

// Number type starts here
type Number struct {
	Token *Token
}

func newNumber(Token *Token, cb func()) Number {
	if cb != nil {
		cb()
	}

	return Number{Token: Token}
}

func (n Number) String() string {
	return string((n.Token.Literal))
}

func (n Number) Value() (float64, error) {
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
	default:
		return 0, ErrNotNumber
	}
}

func AsNumber(j JSON) (*Number, bool) {
	number, ok := j.(Number)
	return &number, ok
}

// Array type starts here
type Array struct {
	Items []JSON
}

func newArray(items []JSON, cb func()) Array {
	if cb != nil {
		cb()
	}

	return Array{Items: items}
}

func (a Array) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for i, item := range a.Items {
		if i == len(a.Items)-1 {
			builder.WriteString(fmt.Sprintf("%v", item))
		} else {
			builder.WriteString(fmt.Sprintf("%v,", item))
		}
	}
	builder.WriteString("]")
	return builder.String()
}

func (a Array) Get(path string) (JSON, error) {
	return nil, nil
}

func AsArray(j JSON) (*Array, bool) {
	arr, ok := j.(Array)
	return &arr, ok
}

// Object type starts here
type KeyValue struct {
	key   []byte
	value JSON
}

type Object struct {
	Properties []KeyValue
}

func newObject(properties []KeyValue, cb func()) Object {
	if cb != nil {
		cb()
	}

	return Object{Properties: properties}
}

func (o Object) String() string {
	var builder strings.Builder
	builder.WriteString("{")
	length := len(o.Properties)
	for ind, kv := range o.Properties {
		if ind == length-1 {
			builder.WriteString(fmt.Sprintf("%s: %v", kv.key, kv.value))
		} else {
			builder.WriteString(fmt.Sprintf("%s: %v,", kv.key, kv.value))
		}
	}
	builder.WriteString("}")
	return builder.String()
}

func (o Object) Get(path string) (JSON, error) {
	return nil, nil
}

func AsObject(j JSON) (*Object, bool) {
	obj, ok := j.(Object)
	return &obj, ok
}
