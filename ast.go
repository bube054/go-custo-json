package jsonvx

import (
	"fmt"
	"strings"
)

type JSON interface {
	String() string
	Literal() string
	Value() any
}

type Null struct {
	Token Token
}

func newJSONNull(token Token, cb func()) Null {
	if cb != nil {
		cb()
	}

	return Null{Token: token}
}

func (j Null) String() string {
	// return fmt.Sprintf(
	// 	"Null{Literal: %s, Value: %v}",
	// 	j.Literal(),
	// 	j.Value(),
	// )
	return "\033[1mnull\033[0m"
}

func (j Null) Literal() string {
	return string(j.Token.Literal)
}

func (j Null) Value() any {
	return j.Token.Value()
}

type Boolean struct {
	Token Token
}

func newJSONBoolean(token Token, cb func()) Boolean {
	if cb != nil {
		cb()
	}

	return Boolean{Token: token}
}

func (j Boolean) String() string {
	// return fmt.Sprintf(
	// 	"Boolean{Literal: %s, Value: %v}",
	// 	j.Literal(),
	// 	j.Value(),
	// )

	if j.Token.Kind == TRUE {
		return "\033[1mtrue\033[0m"
	} else {
		return "\033[1mfalse\033[0m"
	}
}

func (j Boolean) Literal() string {
	return string(j.Token.Literal)
}

func (j Boolean) Value() any {
	return j.Token.Value()
}

type String struct {
	Token Token
}

func newJSONString(token Token, cb func()) String {
	if cb != nil {
		cb()
	}

	return String{Token: token}
}

func (j String) String() string {
	// return fmt.Sprintf(
	// 	"String{Literal: %s, Value: %v}",
	// 	j.Literal(),
	// 	j.Value(),
	// )

	return string(j.Token.Literal)
}

func (j String) Literal() string {
	return string(j.Token.Literal)
}

func (j String) Value() any {
	return j.Token.Value()
}

type Number struct {
	Token Token
}

func newJSONNumber(token Token, cb func()) Number {
	if cb != nil {
		cb()
	}

	return Number{Token: token}
}

func (j Number) String() string {
	// return fmt.Sprintf(
	// 	"Number{Literal: %s, Value: %v}",
	// 	j.Literal(),
	// 	j.Value(),
	// )
	return string(j.Token.Literal)
}

func (j Number) Literal() string {
	return string(j.Token.Literal)
}

func (j Number) Value() any {
	return j.Token.Value()
}

type Array struct {
	Items []JSON
}

func newJSONArray(items []JSON, cb func()) Array {
	if cb != nil {
		cb()
	}

	return Array{Items: items}
}

func (j Array) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for _, item := range j.Items {
		builder.WriteString(fmt.Sprintf("%v,", item))
	}
	builder.WriteString("]")
	return builder.String()
	// var builder strings.Builder
	// builder.WriteString("Array{\n")
	// builder.WriteString(fmt.Sprintf("  Literal: %q,\n", j.Literal()))
	// builder.WriteString("  Items:\n")
	// for i, item := range j.Items {
	// 	builder.WriteString(fmt.Sprintf("    [%d]: %v\n", i, item))
	// }
	// builder.WriteString("}")
	// return builder.String()
}

func (j Array) Literal() string {
	return ""
}

func (j Array) Value() any {
	return j.Items
}

type Object struct {
	Properties map[string]JSON
}

func newJSONObject(properties map[string]JSON, cb func()) Object {
	if cb != nil {
		cb()
	}

	return Object{Properties: properties}
}

func (j Object) String() string {
	var builder strings.Builder
	builder.WriteString("Object{\n")
	builder.WriteString(fmt.Sprintf("  Literal: %q,\n", j.Literal()))
	builder.WriteString("  Properties:\n")
	for key, value := range j.Properties {
		builder.WriteString(fmt.Sprintf("    %q: %v\n", key, value))
	}
	builder.WriteString("}")
	return builder.String()
}

func (j Object) Literal() string {
	return ""
}

func (j Object) Value() any {
	return j.Properties
}
