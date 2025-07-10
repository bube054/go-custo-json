package jsonvx

import (
	"fmt"
	"strings"
)

type JSON interface {
	String() string
	Query(string) any
	// Literal() string
	// Value() any
}

// ///////////////////////////
type Null struct {
	Token *Token
}

func newJSONNull(token *Token, cb func()) Null {
	if cb != nil {
		cb()
	}

	return Null{Token: token}
}

func (n Null) String() string {
	return "\033[1mnull\033[0m"
}

func (n Null) Query(str string) any {
	return n.Token
}

/////////////////////////////

// ///////////////////////////
type Boolean struct {
	Token *Token
}

func newJSONBoolean(token *Token, cb func()) Boolean {
	if cb != nil {
		cb()
	}

	return Boolean{Token: token}
}

func (b Boolean) String() string {
	if b.Token.SubKind == TRUE {
		return "\033[1mtrue\033[0m"
	} else {
		return "\033[1mfalse\033[0m"
	}
}

func (b Boolean) Query(str string) any {
	return b.Token
}

/////////////////////////////

// ///////////////////////////
type String struct {
	Token *Token
}

func newJSONString(token *Token, cb func()) String {
	if cb != nil {
		cb()
	}

	return String{Token: token}
}

func (s String) String() string {
	return (s.Token.Value()).(string)
}

func (s String) Query(str string) any {
	return s.Token
}

/////////////////////////////

// ///////////////////////////
type Number struct {
	Token *Token
}

func newJSONNumber(token *Token, cb func()) Number {
	if cb != nil {
		cb()
	}

	return Number{Token: token}
}

func (n Number) String() string {
	return string((n.Token.Literal))
}

func (b Number) Query(str string) any {
	return b.Token
}

/////////////////////////////

// ///////////////////////////
type Array struct {
	Items []JSON
}

func newJSONArray(items []JSON, cb func()) Array {
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

func (a Array) Query(str string) any {
	return a.Items
}

/////////////////////////////

// ///////////////////////////
type Object struct {
	Properties map[string]JSON
}

func newJSONObject(properties map[string]JSON, cb func()) Object {
	if cb != nil {
		cb()
	}

	return Object{Properties: properties}
}

func (o Object) String() string {
	var builder strings.Builder
	builder.WriteString("{")
	count := 1
	length := len(o.Properties)
	for key, value := range o.Properties {
		if count == length {
			builder.WriteString(fmt.Sprintf("%v: %v", key, value))
		} else {
			builder.WriteString(fmt.Sprintf("%v: %v,", key, value))
		}
		count++
	}
	builder.WriteString("}")
	return builder.String()
}

func (o Object) Query(str string) any {
	return o.Properties
}

/////////////////////////////
