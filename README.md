<div align="center">
    <p style="font-size: 50px">JSONVX</p>
  </a>
</div>

`jsonvx` is a highly configurable `JSON` `parser`, `querier`, and `formatter` for Go.

It supports both strict `JSON` (as defined by [ECMA-404](https://datatracker.ietf.org/doc/html/rfc7159)) and up to a more relaxed variant such as [JSON5](https://json5.org/). This makes it ideal for parsing user-generated data, configuration files, or legacy formats that don't fully comply with the JSON specification.

With a single `ParserConfig` struct, `jsonvx` gives you fine-grained control over how `JSON` is parsed. You can enable or disable features like:

- Hexadecimal numbers `(0xFF)`
- `NaN` and `Infinity` as numeric values
- Leading plus signs in numbers `(+42)`
- Decimal edge cases like `.5` or `5.`
- Unquoted object keys `({key: "value"})`
- Single-quoted strings `('text')`
- Newlines inside strings
- Escape characters outside the standard set
- Trailing commas in arrays or objects
- Line comments `(// comment)`
- Block comments `(/* comment */)`
- Extra whitespace in unusual places

After parsing, `jsonvx` gives you a clean abstract syntax tree `(AST)` that you can either traverse manually or query using the built-in API. Each node in the tree implements a common JSON interface, so you can safely `inspect`, `transform`, or `stringify` data as needed.

`jsonvx` is designed for flexibility and correctness — not raw performance. It prioritizes clarity and configurability over speed, making it perfect for tools, linters, formatters, and config loaders where input may vary or include non-standard extensions.

If you need full control over how JSON is interpreted and a structured way to work with the result, jsonvx is for you.

## Installing

To start using jsonvx, install Go and run `go get`:

```sh
$ go get -u github.com/bube054/jsonvx
```

## Quick Start

Get up and running with `jsonvx` in seconds. Just create a new parser, parse your JSON, and access fields using a simple path query system.

```go
package main

import (
	"fmt"

	"github.com/bube054/jsonvx"
)

func main() {
	data := []byte(`{
  "name": {"first": "Tom", "last": "Anderson"},
  "age": 37,
  "children": ["Sara", "Alex", "Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`)

	parser := jsonvx.NewParser(data, nil)

	node, err := parser.Parse()
	if err != nil {
		panic(fmt.Sprintf("failed to parse JSON: %s", err))
	}

	rootObj, ok := jsonvx.AsObject(node)
	if !ok {
		panic(fmt.Sprintf("expected root node to be an object"))
	}

	ageNode, err := rootObj.QueryPath("age")
	if err != nil {
		panic(fmt.Sprintf("failed to query 'age' field: %s", err))
	}

	ageNum, ok := jsonvx.AsNumber(ageNode)
	if !ok {
		panic(fmt.Sprintf("expected 'age' to be a number"))
	}

	ageValue, err := ageNum.Value()
	if err != nil {
		panic(fmt.Sprintf("failed to convert 'age' to numeric value: %s", err))
	}

	fmt.Println(ageValue) // 37
}
```

## Configuring The Parser
You can configure the parser using the functional options pattern, allowing you to enable relaxed JSON features individually. By default, the parser is strict (all options disabled), matching the [ECMA-404](https://datatracker.ietf.org/doc/html/rfc7159) specification. To allow non-standard or user-friendly formats (like [JSON5](https://json5.org)), pass options when creating the config:

- `AllowExtraWS`: allows extra whitespace characters that are not normally permitted by strict JSON.
  ```go
  cfg := jsonvx.NewParserConfig(WithAllowExtraWS(true))
  parser := jsonvx.NewParser([]byte{'\v', '\f', '\u0085', '\u0085'}, cfg)
  ```
- `AllowHexNumbers`: enables support for hexadecimal numeric literals (e.g., 0xFF).
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowHexNumbers(true))
	parser := jsonvx.NewParser([]byte("0x1F"), cfg)
  ```
- `AllowPointEdgeNumbers`: allows numbers like `.5` or `5.` without requiring a digit before/after the decimal point.
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowPointEdgeNumbers(true))
	parser := jsonvx.NewParser([]byte(".5"), cfg)
  ```
- `AllowInfinity`: enables the use of `Infinity` and `-Infinity` as number values.
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowInfinity(true))
	parser := jsonvx.NewParser([]byte("Infinity"), cfg)
  ```
- `AllowNaN`: allows `NaN` (Not-a-Number) as a numeric value.
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowNaN(true))
	parser := jsonvx.NewParser([]byte("NaN"), cfg)
  ```
- `AllowLeadingPlus`: permits a leading '+' in numbers (e.g., `+42`).
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowLeadingPlus(true))
	parser := jsonvx.NewParser([]byte("+99"), cfg)
  ```
- `AllowUnquoted`: enables parsing of unquoted object keys (e.g., `{foo: "bar"}`)
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowUnquoted(true))
	parser := jsonvx.NewParser([]byte(`{foo: "bar"}`), cfg)
  ```
- `AllowSingleQuotes`: allows strings to be enclosed in single quotes (' ') in addition to double quotes.
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowSingleQuotes(true))
	parser := jsonvx.NewParser([]byte(`{'name': 'Tom'}`), cfg)
  ```
- `AllowNewlineInStrings`: permits literal newlines inside string values without requiring escaping.
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowNewlineInStrings(true))
	parser := jsonvx.NewParser([]byte(`"hello
	world"`), cfg)
  ```
- `AllowOtherEscapeChars`: enables support for escape sequences other than \\, \/, \b, \n, \f, \r, \t and Unicode escapes (\uXXXX).
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowOtherEscapeChars(true))
	parser := jsonvx.NewParser([]byte(`"hello\qworld"`), cfg)
  ```
- `AllowTrailingCommaArray`: permits a trailing comma in array literals (e.g., `[1, 2, ]`).
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowTrailingCommaArray(true))
	parser := jsonvx.NewParser([]byte(`[1, 2, 3, ]`), cfg)
  ```
- `AllowTrailingCommaObject`: permits a trailing comma in object literals (e.g., `{"a": 1,}`).
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowTrailingCommaObject(true))
	parser := jsonvx.NewParser([]byte(`{"a": 1,}`), cfg)
  ```
- `AllowLineComments`: enables the use of single-line comments (// ...).
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowLineComments(true))
	parser := jsonvx.NewParser([]byte(`// comment\n123`), cfg)
  ```
- `AllowBlockComments`: enables the use of block comments (/* ... */).
  ```go
	cfg := jsonvx.NewParserConfig(WithAllowBlockComments(true))
	parser := jsonvx.NewParser([]byte(`/* comment */ 123`), cfg)
  ```
- You can also combine multiple options to create a custom configuration.
	```go
	cfg := jsonvx.NewParserConfig(WithAllowLineComments(true), WithAllowBlockComments(true))
	parser := jsonvx.NewParser([]byte(`// comment\n/* comment */ 123`), cfg)
	```

## Query Path Syntax
...

## License
[MIT](https://github.com/bube054/jsonvx/blob/main/LICENSE)