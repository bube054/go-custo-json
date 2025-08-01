# JSONVX

`jsonvx` is a highly configurable `JSON` `parser`, `querier`, and `formatter` for Go.

It supports both strict `JSON` (as defined by [ECMA-404](https://datatracker.ietf.org/doc/html/rfc7159)) and up to a more relaxed variant such as [JSON5](https://json5.org/). This makes it ideal for parsing user-generated data, configuration files, or legacy formats that don't fully comply with the JSON specification.

With a single `ParserConfig` struct, `jsonvx` gives you fine-grained control over how `JSON` is parsed. You can enable or disable features like:

- [Hexadecimal numbers `(0xFF)`](#allowextraws)
- [`NaN` and `Infinity` as numeric values](#allownan)
- [Leading plus signs in numbers `(+42)`](#allowleadingplus)
- [Decimal edge cases like `.5` or `5.`](#allowpointedgenumbers)
- [Unquoted object keys `({key: "value"})`](#allowunquoted)
- [Single-quoted strings `('text')`](#allowsinglequotes)
- [Newlines inside strings](#allownewlineinstrings)
- [Escape characters outside the standard set](#allowotherescapechars)
- [Trailing commas in arrays or objects](#allowtrailingcommaarray)
- [Line comments `(// comment)` and Block comments `(/* comment */)`](#allowlinecomments)
- [Extra whitespace in unusual places](#allowextraws)

After parsing, `jsonvx` gives you a clean abstract syntax tree `(AST)` that you can either traverse manually or query using the built-in API. Each node in the tree implements a common `JSON interface`, so you can safely `inspect`, `transform`, or `stringify` data as needed.

`jsonvx` is designed for flexibility and correctness — not raw performance. It prioritizes clarity and configurability over speed, making it perfect for tools, linters, formatters, and config loaders where input may vary or include non-standard extensions.

If you need full control over how `JSON` is interpreted and a structured way to work with the result, `jsonvx` is for you.

## Installing

To start using `jsonvx`, install [Go](https://golang.org) and run `go get`:

```sh
$ go get -u github.com/bube054/jsonvx
```

## Quick Start

Get up and running with `jsonvx` in seconds. Just create a new parser, parse your `JSON`, and access fields using a simple path query system.

```go
package main

import (
	"fmt"

	"github.com/your/module/jsonvx"
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
		panic("expected root node to be an object")
	}

	ageNode, err := rootObj.QueryPath("age")
	if err != nil {
		panic(fmt.Sprintf("failed to query 'age' field: %s", err))
	}

	ageNum, ok := jsonvx.AsNumber(ageNode)
	if !ok {
		panic("expected 'age' to be a number")
	}

	ageValue, err := ageNum.Value()
	if err != nil {
		panic(fmt.Sprintf("failed to convert 'age' to numeric value: %s", err))
	}

	fmt.Println(ageValue) // 37
}
```

## Query Path Syntax

Access deeply nested fields in your parsed `JSON` `array` or `object` structure using the `QueryPath` method, which accepts a variadic list of strings to represent the path segments.

Using the `json` data above, we can query for specific values using the `QueryPath` method:

```go
// Get first name
strNode, _ := rootObj.QueryPath("name", "first") // => "Tom"

// get children array
arrayNode, _ := rootObj.QueryPath("children") // => ["Sara", "Alex", "Jack"]

// Get second child
strNode, _ = rootObj.QueryPath("children", "1") // => "Alex"

// Get favorite movie (with dot in key name)
strNode, _ = rootObj.QueryPath("fav.movie") // => "Deer Hunter"

// Get first friend object
objNode, _ := rootObj.QueryPath("friends", "0") // => {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]}

// Get last name of first friend
strNode, _ = rootObj.QueryPath("friends", "0", "last") // => "Murphy"

// Get second social network of second friend
strNode, _ = rootObj.QueryPath("friends", "1", "nets", "1") // => "tw"

// Get age of third friend
nuNode, _ := rootObj.QueryPath("friends", "2", "age") // => 47
```

## Configuring The Parser

You can configure the parser using the functional options pattern, allowing you to enable relaxed JSON features individually. By default, the parser is strict (all options disabled), matching the [ECMA-404](https://datatracker.ietf.org/doc/html/rfc7159) specification. To allow non-standard or user-friendly formats (like [JSON5](https://json5.org)), pass options when creating the config:

  ### `AllowExtraWS`:
  Allows extra whitespace characters that are not normally permitted by strict JSON.
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowExtraWS(true))
  parser := jsonvx.NewParser([]byte("{\"key\":\v\"value\"}"), cfg) // valid Vertical Tab
  parser := jsonvx.NewParser([]byte("{\"key\":\f\"value\"}"), cfg) // valid Form Feed
  parser := jsonvx.NewParser([]byte("{\"key\":\u0085\"value\"}"), cfg) // valid Next Line
  parser := jsonvx.NewParser([]byte("{\"key\":\u00A0\"value\"}"), cfg) // valid No-Break Space
  ```
  ### `AllowHexNumbers`:
  Enables support for hexadecimal numeric literals (e.g., 0x1F).
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowHexNumbers(true))
  parser := jsonvx.NewParser([]byte("0x1F"), cfg) // valid 31 in hex
  ```
  ### `AllowPointEdgeNumbers`:
  Allows numbers like `.5` or `5.` without requiring a digit before/after the decimal point.
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowPointEdgeNumbers(true))
  parser := jsonvx.NewParser([]byte(".5"), cfg) // valid pre decimal point number
  parser := jsonvx.NewParser([]byte("5."), cfg) // valid post decimal point number
  ```
  ### `AllowInfinity`:
  Enables the use of `Infinity` and `-Infinity` as number values.
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowInfinity(true))
  parser := jsonvx.NewParser([]byte("Infinity"), cfg) // valid positive infinity
  parser := jsonvx.NewParser([]byte("-Infinity"), cfg) // valid negative infinity
  parser := jsonvx.NewParser([]byte("+Infinity"), cfg) // valid only if AllowLeadingPlus is enabled
  ```
  ### `AllowNaN`:
  Allows `NaN` (Not-a-Number) as a numeric value.
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowNaN(true))
  parser := jsonvx.NewParser([]byte("NaN"), cfg) // valid NaN
  parser := jsonvx.NewParser([]byte("-NaN"), cfg) // valid NaN
  parser := jsonvx.NewParser([]byte("+NaN"), cfg) // valid NaN only if AllowLeadingPlus is enabled
  ```
  ### `AllowLeadingPlus`:
  Permits a leading '+' in numbers (e.g., `+42`).
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowLeadingPlus(true))
  parser := jsonvx.NewParser([]byte("+99"), cfg) // valid positive number
  ```
  ### `AllowUnquoted`:
  Enables parsing of unquoted object keys (e.g., `{foo: "bar"}`)
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowUnquoted(true))
  parser := jsonvx.NewParser([]byte(`{foo: "bar"}`), cfg) // valid only for unquoted keys and not for unquoted values
  ```
  ### `AllowSingleQuotes`:
  Allows strings to be enclosed in single quotes (' ') in addition to double quotes.
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowSingleQuotes(true))
  parser := jsonvx.NewParser([]byte(`{'name': 'Tom'}`), cfg) // valid single-quoted string
  ```
  ### `AllowNewlineInStrings`:
  Permits multiple new line characters being escaped.
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowNewlineInStrings(true))
  parser := jsonvx.NewParser([]byte(`"hello \
   world"`), cfg) // valid escaped new line
  ```
  ### `AllowOtherEscapeChars`:
  Enables support for escape sequences other than \\, \/, \b, \n, \f, \r, \t and Unicode escapes (\uXXXX).
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowOtherEscapeChars(true))
  parser := jsonvx.NewParser([]byte(`"hello\qworld"`), cfg) // valid other escape character
  ```
  ### `AllowTrailingCommaArray`:
  Permits a trailing comma in array literals (e.g., `[1, 2, ]`).
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowTrailingCommaArray(true))
  parser := jsonvx.NewParser([]byte(`[1, 2, 3, ]`), cfg) // valid array with trailing comma
  ```
  ### `AllowTrailingCommaObject`:
  Permits a trailing comma in object literals (e.g., `{"a": 1,}`).
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowTrailingCommaObject(true))
  parser := jsonvx.NewParser([]byte(`{"a": 1,}`), cfg) // valid object with trailing comma
  ```
  ### `AllowLineComments`:
  Enables the use of single-line comments (// ...).
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowLineComments(true))
  parser := jsonvx.NewParser([]byte(`// comment
  123`), cfg) // valid number after liner comment
  ```
  ### `AllowBlockComments`:
  Enables the use of block comments (/_ ... _/).
  ```go
  cfg := jsonvx.NewParserConfig(jsonvx.WithAllowBlockComments(true))
  parser := jsonvx.NewParser([]byte(`/* comment */ 123`), cfg)
  ```

## Parsing behavior
Below are examples of how to `parse` and `retrieve` values from the parsed `JSON` input using the `Parser`.

### Object

This example shows how to parse a `JSON` object using `jsonvx`, access its fields, and iterate over key-value pairs using the `ForEach` method.

Behavior

- `Objects` are stored as an `array` of key-value pairs using the `jsonvx.KeyValue` struct.
- Keys are stored as `[]byte` and `sorted` `lexicographically` for `log(n)` value retrieval.
- Each key-value pair is accessible using the `ForEach` method.
- Mixed value types (e.g., `numbers`, `strings`, `objects`) are supported.
- Duplicate keys are included, but the first key-value pair is used.

```go
parser := jsonvx.NewParser([]byte(`{"name": "Alice", "age": 30}`), jsonvx.NewParserConfig())

// Parse the JSON input
node, err := parser.Parse()
if err != nil {
	panic(fmt.Sprintf("failed to parse JSON: %s", err))
}

// Cast the parsed node to an ObjectNode
objNode, ok := jsonvx.AsObject(node)
if !ok {
	panic(fmt.Sprintf("expected root node to be an object"))
}

// Iterate over key-value pairs using ForEach
objNode.ForEach(func(key []byte, value jsonvx.JSON, object jsonvx.Object) {
	fmt.Printf("Key: %s, Value: %v\n", string(key), value)
})
```

### Array

This example demonstrates how to parse a `JSON` `array` using `jsonvx`, cast it to an `Array` node, and iterate over its elements using the built-in ForEach method.

Behavior

- `JSON` arrays are parsed as `Array` node.
- Each element in the `Array` is accessible using indexed access or `ForEach`.
- Mixed types (e.g., `numbers`, `strings`, `objects`) are supported.

```go
parser := jsonvx.NewParser([]byte(`[1, 2, 3]`), jsonvx.NewParserConfig())

// Parse the JSON input
node, err := parser.Parse()
if err != nil {
	panic(fmt.Sprintf("failed to parse JSON: %s", err))
}

// Cast the parsed node to an ArrayNode
arrNode, ok := jsonvx.AsArray(node)
if !ok {
	panic(fmt.Sprintf("expected root node to be an array"))
}

// Iterate over each item using ForEach
arrNode.ForEach(func(item jsonvx.JSON, index int, array jsonvx.Array) {
	fmt.Printf("Item %d: %v\n", index, item)
})
```

### String

This example demonstrates how to parse a JSON string value using jsonvx and retrieve its Go representation.

Behavior

- JSON string values are parsed as StringNode.
- Value() returns the raw Go string value.

```go
parser := jsonvx.NewParser([]byte(`"Hello, World!"`), jsonvx.NewParserConfig())

// Parse the JSON input
node, err := parser.Parse()
if err != nil {
	panic(fmt.Sprintf("failed to parse JSON: %s", err))
}

// Cast the parsed node to a StringNode
strNode, ok := jsonvx.AsString(node)
if !ok {
	panic(fmt.Sprintf("expected root node to be a string"))
}

// Extract the underlying Go string value
strValue, err := strNode.Value()
if err != nil {
	panic(fmt.Sprintf("failed to extract string value: %s", err))
}

fmt.Println(strValue) // Output: Hello, World!
```

### Number

This example demonstrates how to parse numeric values from JSON input using jsonvx. The parser supports multiple numeric formats including:

Behavior

- JSON numbers are parsed as a NumberNode.
- The Go value returned by .Value() is of type float64.
- Integers (e.g., 42)
- Floating-point numbers (e.g., 3.14)
- Scientific notation (e.g., 1.2e10)
- Hexadecimal (e.g., 0x1A) – if allowed by config
- Special values like Infinity, -Infinity, and NaN – if allowed by config
- Point edge numbers (e.g., .5 or 5.) – if allowed by config
- Leading plus sign (e.g., +42) – if allowed by config

```go
parser := jsonvx.NewParser([]byte("123456"), jsonvx.NewParserConfig())

// Parse the JSON input
node, err := parser.Parse()
if err != nil {
	panic(fmt.Sprintf("failed to parse JSON: %s", err))
}

// Cast the parsed node to a NumberNode
numNode, ok := jsonvx.AsNumber(node)
if !ok {
	panic(fmt.Sprintf("expected root node to be a number"))
}

// Extract the underlying Go float64 value
numValue, err := numNode.Value()
if err != nil {
	panic(fmt.Sprintf("failed to extract number value: %s", err))
}

fmt.Println(numValue) // Output: 123456
```

### Boolean

This example demonstrates how to parse a boolean (true or false) value from JSON input using the Parser, and retrieve its Go representation.

Behavior

- JSON booleans (true, false) are parsed as BooleanNode.
- The Go value returned by .Value() is of type bool.

```go
parser := jsonvx.NewParser([]byte("true"), jsonvx.NewParserConfig())

// Parse the JSON input
node, err := parser.Parse()
if err != nil {
	panic(fmt.Sprintf("failed to parse JSON: %s", err))
}

// Cast the parsed node to a BooleanNode
boolNode, ok := jsonvx.AsBoolean(node)
if !ok {
	panic(fmt.Sprintf("expected root node to be a boolean"))
}

// Extract the underlying Go value
trueValue, err := boolNode.Value()
if err != nil {
	panic(fmt.Sprintf("failed to extract boolean value: %s", err))
}

fmt.Println(trueValue) // Output: true
```

### Null

This example demonstrates how to parse a null value from a JSON string using the Parser and retrieve its Go representation.

Behavior

- A JSON null value is parsed as a NullNode.
- When accessed using Value(), a null returns nil.
- Internally and semantically, null in JSON maps to Go's nil.

```go
parser := jsonvx.NewParser([]byte("null"), jsonvx.NewParserConfig())

// Parse the JSON input
node, err := parser.Parse()
if err != nil {
	t.Fatalf("failed to parse JSON: %s", err)
}

// Cast the parsed node to a NullNode
nullNode, ok := jsonvx.AsNull(node)
if !ok {
	t.Fatalf("expected root node to be a null value, but got: %s", err.Error())
}

// Extract the underlying Go value
nilValue, _ := nullNode.Value()
fmt.Println(nilValue) // Output: <nil>
```

### Comments

This `parser` supports both line `(//)` and block `(/\* \*/)` comments in JSON-like input, but treats them as non-semantic, meaning they are ignored and not part of the `AST`.

Behavior

- `Comments` are ignored during parsing.
- A `JSON` input with only comments and no data will result in a parse error.

```go
parser := jsonvx.NewParser([]byte("/* Block Comment */"), jsonvx.NewParserConfig(
	jsonvx.WithAllowBlockComments(true), jsonvx.WithAllowLineComments(true),
))

// parse the JSON
node, err := parser.Parse()
if err != nil {
	panic(fmt.Sprintf("failed to parse JSON: %s", err))
}
```

## Maintainers

- [bube054](https://github.com/bube054) - **Attah Gbubemi David (author)**

## Other helpful projects

- [ginvalidator](https://github.com/bube054/ginvalidator)
- [validatorgo](https://github.com/bube054/validatorgo)

## License

This project is licensed under the [MIT](https://opensource.org/license/mit). See the [LICENSE](https://github.com/bube054/jsonvx/blob/master/LICENSE) file for details.
