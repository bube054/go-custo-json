package jsonvx

import (
	"fmt"
	"strings"
)

type JSONNode interface {
	String() string
	Literal() string
	Value() any
}

type JSONNull struct {
	Token Token
}

func newJSONNull(token Token, cb func()) JSONNull {
	if cb != nil {
		cb()
	}

	return JSONNull{Token: token}
}

func (j JSONNull) String() string {
	// return fmt.Sprintf(
	// 	"JSONNull{Literal: %s, Value: %v}",
	// 	j.Literal(),
	// 	j.Value(),
	// )
	return "null"
}

func (j JSONNull) Literal() string {
	return string(j.Token.Literal)
}

func (j JSONNull) Value() any {
	return j.Token.Value()
}

type JSONBoolean struct {
	Token Token
}

func newJSONBoolean(token Token, cb func()) JSONBoolean {
	if cb != nil {
		cb()
	}

	return JSONBoolean{Token: token}
}

func (j JSONBoolean) String() string {
	// return fmt.Sprintf(
	// 	"JSONBoolean{Literal: %s, Value: %v}",
	// 	j.Literal(),
	// 	j.Value(),
	// )

	if j.Token.Kind == TRUE {
		return "true"
	} else {
		return "falSE"
	}
}

func (j JSONBoolean) Literal() string {
	return string(j.Token.Literal)
}

func (j JSONBoolean) Value() any {
	return j.Token.Value()
}

type JSONString struct {
	Token Token
}

func newJSONString(token Token, cb func()) JSONString {
	if cb != nil {
		cb()
	}

	return JSONString{Token: token}
}

func (j JSONString) String() string {
	// return fmt.Sprintf(
	// 	"JSONString{Literal: %s, Value: %v}",
	// 	j.Literal(),
	// 	j.Value(),
	// )

	return string(j.Token.Literal)
}

func (j JSONString) Literal() string {
	return string(j.Token.Literal)
}

func (j JSONString) Value() any {
	return j.Token.Value()
}

type JSONNumber struct {
	Token Token
}

func newJSONNumber(token Token, cb func()) JSONNumber {
	if cb != nil {
		cb()
	}

	return JSONNumber{Token: token}
}

func (j JSONNumber) String() string {
	// return fmt.Sprintf(
	// 	"JSONNumber{Literal: %s, Value: %v}",
	// 	j.Literal(),
	// 	j.Value(),
	// )
	return string(j.Token.Literal)
}

func (j JSONNumber) Literal() string {
	return string(j.Token.Literal)
}

func (j JSONNumber) Value() any {
	return j.Token.Value()
}

type JSONArray struct {
	Items []JSONNode
}

func newJSONArray(items []JSONNode, cb func()) JSONArray {
	if cb != nil {
		cb()
	}

	return JSONArray{Items: items}
}

func (j JSONArray) String() string {
	var builder strings.Builder
	builder.WriteString("[")
	for _, item := range j.Items {
		builder.WriteString(fmt.Sprintf("%v,", item))
	}
	builder.WriteString("]")
	return builder.String()
	// var builder strings.Builder
	// builder.WriteString("JSONArray{\n")
	// builder.WriteString(fmt.Sprintf("  Literal: %q,\n", j.Literal()))
	// builder.WriteString("  Items:\n")
	// for i, item := range j.Items {
	// 	builder.WriteString(fmt.Sprintf("    [%d]: %v\n", i, item))
	// }
	// builder.WriteString("}")
	// return builder.String()
}

func (j JSONArray) Literal() string {
	return ""
}

func (j JSONArray) Value() any {
	return j.Items
}

type JSONObject struct {
	Properties map[string]JSONNode
}

func (j JSONObject) String() string {
	var builder strings.Builder
	builder.WriteString("JSONObject{\n")
	builder.WriteString(fmt.Sprintf("  Literal: %q,\n", j.Literal()))
	builder.WriteString("  Properties:\n")
	for key, value := range j.Properties {
		builder.WriteString(fmt.Sprintf("    %q: %v\n", key, value))
	}
	builder.WriteString("}")
	return builder.String()
}

func (j JSONObject) Literal() string {
	return ""
}

func (j JSONObject) Value() any {
	return j.Properties
}
