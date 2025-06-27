package gocustojson

import (
	"fmt"
	"strings"
)

// Config defines the encoding and decoding behavior for the JSON parser.
//
// By default, all fields are false (strict mode). Enabling individual fields allows the parser
// to accept features that are not allowed by the standard [ECMA-404] specification but may appear
// in relaxed formats like [JSON5] or user-generated JSON.
//
// [ECMA-404]: https://datatracker.ietf.org/doc/html/rfc7159
// [JSON5]: https://json5.org/
type Config struct {
	AllowExtraWS bool // AllowExtraWS allows extra whitespace in places not normally permitted by strict JSON.

	AllowHexNumbers       bool // AllowHexNumbers enables support for hexadecimal numeric literals (e.g., 0xFF).
	AllowPointEdgeNumbers bool // AllowPointEdgeNumbers allows numbers like `.5` or `5.` without requiring a digit before/after the decimal point.

	AllowInfinity    bool // AllowInfinity enables the use of `Infinity` and `-Infinity` as number values.
	AllowNaN         bool // AllowNaN allows `NaN` (Not-a-Number) as a numeric value.
	AllowLeadingPlus bool // AllowLeadingPlus permits a leading '+' in numbers (e.g., `+42`).

	AllowUnquoted         bool // AllowUnquoted enables parsing of unquoted object keys (e.g., `{foo: "bar"}`).
	AllowSingleQuotes     bool // AllowSingleQuotes allows strings to be enclosed in single quotes (' ') in addition to double quotes.
	AllowNewlineInStrings bool // AllowNewlineInStrings permits literal newlines inside string values without requiring escaping.
	AllowOtherEscapeChars bool // AllowOtherEscapeChars enables support for escape sequences other than \\, \/, \b, \n, \f, \r, \t and Unicode escapes (\uXXXX).

	AllowTrailingCommaArray  bool // AllowTrailingCommaArray permits a trailing comma in array literals (e.g., `[1, 2, ]`).
	AllowTrailingCommaObject bool // AllowTrailingCommaObject permits a trailing comma in object literals (e.g., `{"a": 1,}`).

	AllowLineComments  bool // AllowLineComments enables the use of single-line comments (// ...).
	AllowBlockComments bool // AllowBlockComments enables the use of block comments (/* ... */).
}

// NewConfig creates a new Config instance, optionally applying one or more configuration options.
// Options are applied in the order provided.
func NewConfig(opts ...func(*Config)) *Config {
	cfg := &Config{}

	for _, o := range opts {
		o(cfg)
	}

	return cfg
}

// String returns a formatted string representing all configuration options in the Config.
// Each field is listed with its corresponding boolean value.
func (c *Config) String() string {
	var b strings.Builder
	b.WriteString("Config{\n")
	b.WriteString(fmt.Sprintf("  AllowExtraWS: %v,\n", c.AllowExtraWS))
	b.WriteString(fmt.Sprintf("  AllowHexNumbers: %v,\n", c.AllowHexNumbers))
	b.WriteString(fmt.Sprintf("  AllowPointEdgeNumbers: %v,\n", c.AllowPointEdgeNumbers))
	b.WriteString(fmt.Sprintf("  AllowInfinity: %v,\n", c.AllowInfinity))
	b.WriteString(fmt.Sprintf("  AllowNaN: %v,\n", c.AllowNaN))
	b.WriteString(fmt.Sprintf("  AllowLeadingPlus: %v,\n", c.AllowLeadingPlus))
	b.WriteString(fmt.Sprintf("  AllowUnquoted: %v,\n", c.AllowUnquoted))
	b.WriteString(fmt.Sprintf("  AllowSingleQuotes: %v,\n", c.AllowSingleQuotes))
	b.WriteString(fmt.Sprintf("  AllowNewlineInStrings: %v,\n", c.AllowNewlineInStrings))
	b.WriteString(fmt.Sprintf("  AllowOtherEscapeChars: %v,\n", c.AllowOtherEscapeChars))
	b.WriteString(fmt.Sprintf("  AllowTrailingCommaArray: %v,\n", c.AllowTrailingCommaArray))
	b.WriteString(fmt.Sprintf("  AllowTrailingCommaObject: %v,\n", c.AllowTrailingCommaObject))
	b.WriteString(fmt.Sprintf("  AllowLineComments: %v,\n", c.AllowLineComments))
	b.WriteString(fmt.Sprintf("  AllowBlockComments: %v,\n", c.AllowBlockComments))
	b.WriteString("}")
	return b.String()
}


// WithAllowExtraWS is the functional option setters for the AllowExtraWS flag.
func WithAllowExtraWS(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowExtraWS = allow
	}
}

// Functional option setters for the AllowHexNumbers flag.
func WithAllowHexNumbers(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowHexNumbers = allow
	}
}

// Functional option setters for the AllowHexNumbers flag.
func WithAllowPointEdgeNumbers(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowPointEdgeNumbers = allow
	}
}

func WithAllowInfinity(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowInfinity = allow
	}
}

func WithAllowNaN(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowNaN = allow
	}
}

func WithAllowLeadingPlus(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowLeadingPlus = allow
	}
}

func WithAllowUnquoted(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowUnquoted = allow
	}
}

func WithAllowSingleQuotes(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowSingleQuotes = allow
	}
}

func WithAllowNewlineInStrings(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowNewlineInStrings = allow
	}
}

func WithAllowOtherEscapeChars(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowOtherEscapeChars = allow
	}
}

func WithAllowTrailingCommaArray(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowTrailingCommaArray = allow
	}
}

func WithAllowTrailingCommaObject(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowTrailingCommaObject = allow
	}
}

func WithAllowLineComments(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowLineComments = allow
	}
}

func WithAllowBlockComments(allow bool) func(*Config) {
	return func(c *Config) {
		c.AllowBlockComments = allow
	}
}
