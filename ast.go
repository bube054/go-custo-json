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
	token Token
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
	return string(j.token.Literal)
}

func (j JSONNull) Value() any {
	return j.token.Value()
}

type JSONBoolean struct {
	token Token
}

func (j JSONBoolean) String() string {
	return fmt.Sprintf(
		"JSONBoolean{Literal: %s, Value: %v}",
		j.Literal(),
		j.Value(),
	)
}

func (j JSONBoolean) Literal() string {
	return string(j.token.Literal)
}

func (j JSONBoolean) Value() any {
	return j.token.Value()

}

type JSONString struct {
	token Token
}

func (j JSONString) String() string {
	return fmt.Sprintf(
		"JSONString{Literal: %s, Value: %v}",
		j.Literal(),
		j.Value(),
	)
}

func (j JSONString) Literal() string {
	return string(j.token.Literal)
}

func (j JSONString) Value() any {
	return j.token.Value()
}

type JSONNumber struct {
	token Token
}

func (j JSONNumber) String() string {
	return fmt.Sprintf(
		"JSONNumber{Literal: %s, Value: %v}",
		j.Literal(),
		j.Value(),
	)
}

func (j JSONNumber) Literal() string {
	return string(j.token.Literal)
}

func (j JSONNumber) Value() any {
	return j.token.Value()
}

type JSONArray struct {
	items []JSONNode
}

func (j JSONArray) String() string {
	var builder strings.Builder
	builder.WriteString("JSONArray{\n")
	builder.WriteString(fmt.Sprintf("  Literal: %q,\n", j.Literal()))
	builder.WriteString("  Items:\n")
	for i, item := range j.items {
		builder.WriteString(fmt.Sprintf("    [%d]: %v\n", i, item))
	}
	builder.WriteString("}")
	return builder.String()
}

func (j JSONArray) Literal() string {
	return ""
}

func (j JSONArray) Value() any {
	return j.items
}

type JSONObject struct {
	properties map[string]JSONNode
}

func (j JSONObject) String() string {
	var builder strings.Builder
	builder.WriteString("JSONObject{\n")
	builder.WriteString(fmt.Sprintf("  Literal: %q,\n", j.Literal()))
	builder.WriteString("  Properties:\n")
	for key, value := range j.properties {
		builder.WriteString(fmt.Sprintf("    %q: %v\n", key, value))
	}
	builder.WriteString("}")
	return builder.String()
}

func (j JSONObject) Literal() string {
	return ""
}

func (j JSONObject) Value() any {
	return j.properties
}
