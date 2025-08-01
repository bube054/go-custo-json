package jsonvx

import (
	"fmt"
	"strings"
)

// ParserConfig defines the encoding and decoding behavior for the JSON parser.
//
// By default, all fields are false (strict mode). Enabling individual fields allows the parser
// to accept features that are not allowed by the standard [ECMA-404] specification but may appear
// in relaxed formats like [JSON5] or user-generated JSON.
//
// [ECMA-404]: https://datatracker.ietf.org/doc/html/rfc7159
// [JSON5]: https://json5.org/
type ParserConfig struct {
	AllowExtraWS bool // AllowExtraWS allows extra whitespace characters that are not normally permitted by strict JSON.

	AllowHexNumbers       bool // AllowHexNumbers enables support for hexadecimal numeric literals (e.g., 0xFF).
	AllowPointEdgeNumbers bool // AllowPointEdgeNumbers allows numbers like `.5` or `5.` without requiring a digit before/after the decimal point.

	AllowInfinity    bool // AllowInfinity enables the use of `Infinity` and `-Infinity` as number values.
	AllowNaN         bool // AllowNaN allows `NaN` (Not-a-Number) as a numeric value.
	AllowLeadingPlus bool // AllowLeadingPlus permits a leading '+' in numbers (e.g., `+42`).

	AllowUnquoted         bool // AllowUnquoted enables parsing of unquoted object keys (e.g., `{foo: "bar"}`).
	AllowSingleQuotes     bool // AllowSingleQuotes allows strings to be enclosed in single quotes (' ') in addition to double quotes.
	AllowNewlineInStrings bool // AllowNewlineInStrings permits multiple new line characters being escaped.
	AllowOtherEscapeChars bool // AllowOtherEscapeChars enables support for escape sequences other than \\, \/, \b, \n, \f, \r, \t and Unicode escapes (\uXXXX).

	AllowTrailingCommaArray  bool // AllowTrailingCommaArray permits a trailing comma in array literals (e.g., `[1, 2, ]`).
	AllowTrailingCommaObject bool // AllowTrailingCommaObject permits a trailing comma in object literals (e.g., `{"a": 1,}`).

	AllowLineComments  bool // AllowLineComments enables the use of single-line comments (// ...).
	AllowBlockComments bool // AllowBlockComments enables the use of block comments (/* ... */).
}

// NewParserConfig creates a new ParserConfig instance, optionally applying one or more configuration options.
// Options are applied in the order provided.
func NewParserConfig(opts ...func(*ParserConfig)) *ParserConfig {
	cfg := &ParserConfig{}

	for _, o := range opts {
		o(cfg)
	}

	return cfg
}

// String returns a formatted string representing all configuration options in the ParserConfig.
// Each field is listed with its corresponding boolean value.
func (c *ParserConfig) String() string {
	var b strings.Builder
	b.WriteString("ParserConfig{\n")

	configFields := []struct {
		name  string
		value bool
	}{
		{"AllowExtraWS", c.AllowExtraWS},
		{"AllowHexNumbers", c.AllowHexNumbers},
		{"AllowPointEdgeNumbers", c.AllowPointEdgeNumbers},
		{"AllowInfinity", c.AllowInfinity},
		{"AllowNaN", c.AllowNaN},
		{"AllowLeadingPlus", c.AllowLeadingPlus},
		{"AllowUnquoted", c.AllowUnquoted},
		{"AllowSingleQuotes", c.AllowSingleQuotes},
		{"AllowNewlineInStrings", c.AllowNewlineInStrings},
		{"AllowOtherEscapeChars", c.AllowOtherEscapeChars},
		{"AllowTrailingCommaArray", c.AllowTrailingCommaArray},
		{"AllowTrailingCommaObject", c.AllowTrailingCommaObject},
		{"AllowLineComments", c.AllowLineComments},
		{"AllowBlockComments", c.AllowBlockComments},
	}

	for _, f := range configFields {
		b.WriteString(fmt.Sprintf("  %s: %v,\n", f.name, f.value))
	}

	b.WriteString("}")
	return b.String()
}

// WithAllowExtraWS is the functional option setters for the AllowExtraWS flag.
func WithAllowExtraWS(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowExtraWS = allow
	}
}

// Functional option setters for the AllowHexNumbers flag.
func WithAllowHexNumbers(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowHexNumbers = allow
	}
}

// Functional option setters for the AllowHexNumbers flag.
func WithAllowPointEdgeNumbers(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowPointEdgeNumbers = allow
	}
}

func WithAllowInfinity(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowInfinity = allow
	}
}

func WithAllowNaN(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowNaN = allow
	}
}

func WithAllowLeadingPlus(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowLeadingPlus = allow
	}
}

func WithAllowUnquoted(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowUnquoted = allow
	}
}

func WithAllowSingleQuotes(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowSingleQuotes = allow
	}
}

func WithAllowNewlineInStrings(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowNewlineInStrings = allow
	}
}

func WithAllowOtherEscapeChars(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowOtherEscapeChars = allow
	}
}

func WithAllowTrailingCommaArray(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowTrailingCommaArray = allow
	}
}

func WithAllowTrailingCommaObject(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowTrailingCommaObject = allow
	}
}

func WithAllowLineComments(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowLineComments = allow
	}
}

func WithAllowBlockComments(allow bool) func(*ParserConfig) {
	return func(c *ParserConfig) {
		c.AllowBlockComments = allow
	}
}

// JSON5Config returns a ParserConfig with all features enabled for JSON5 compatibility.
func JSON5Config() *ParserConfig {
	return &ParserConfig{
		AllowExtraWS:             true,
		AllowHexNumbers:          true,
		AllowPointEdgeNumbers:    true,
		AllowInfinity:            true,
		AllowNaN:                 true,
		AllowLeadingPlus:         true,
		AllowUnquoted:            true,
		AllowSingleQuotes:        true,
		AllowNewlineInStrings:    true,
		AllowOtherEscapeChars:    true,
		AllowTrailingCommaArray:  true,
		AllowTrailingCommaObject: true,
		AllowLineComments:        true,
		AllowBlockComments:       true,
	}
}
