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

Get searches json for the specified path. A path is in dot syntax, such as "name.last" or "age". When the value is found it's returned immediately.

```go
package main

import "github.com/tidwall/gjson"

const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

func main() {
	value := gjson.Get(json, "name.last")
	println(value.String())
}
```

This will print:

```
Prichard
```

_There's also [GetBytes](#working-with-bytes) for working with JSON byte slices._

## Path Syntax

Below is a quick overview of the path syntax, for more complete information please
check out [GJSON Syntax](SYNTAX.md).

A path is a series of keys separated by a dot.
A key may contain special wildcard characters '\*' and '?'.
To access an array value use the index as the key.
To get the number of elements in an array or to access a child path, use the '#' character.
The dot and wildcard characters can be escaped with '\\'.

```json
{
  "name": { "first": "Tom", "last": "Anderson" },
  "age": 37,
  "children": ["Sara", "Alex", "Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {
      "first": "Dale",
      "last": "Murphy",
      "age": 44,
      "nets": ["ig", "fb", "tw"]
    },
    { "first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"] },
    { "first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"] }
  ]
}
```

```
"name.last"          >> "Anderson"
"age"                >> 37
"children"           >> ["Sara","Alex","Jack"]
"children.#"         >> 3
"children.1"         >> "Alex"
"child*.2"           >> "Jack"
"c?ildren.0"         >> "Sara"
"fav\.movie"         >> "Deer Hunter"
"friends.#.first"    >> ["Dale","Roger","Jane"]
"friends.1.last"     >> "Craig"
```

You can also query an array for the first match by using `#(...)`, or find all
matches with `#(...)#`. Queries support the `==`, `!=`, `<`, `<=`, `>`, `>=`
comparison operators and the simple pattern matching `%` (like) and `!%`
(not like) operators.

```
friends.#(last=="Murphy").first    >> "Dale"
friends.#(last=="Murphy")#.first   >> ["Dale","Jane"]
friends.#(age>45)#.last            >> ["Craig","Murphy"]
friends.#(first%"D*").last         >> "Murphy"
friends.#(first!%"D*").last        >> "Craig"
friends.#(nets.#(=="fb"))#.first   >> ["Dale","Roger"]
```

_Please note that prior to v1.3.0, queries used the `#[...]` brackets. This was
changed in v1.3.0 as to avoid confusion with the new
[multipath](SYNTAX.md#multipaths) syntax. For backwards compatibility,
`#[...]` will continue to work until the next major release._

<!-- goos: windows
goarch: amd64
pkg: github/bube054/jsonvx
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkJSONParserMediumPayload-8   	   10527	    116476 ns/op	  132784 B/op	     134 allocs/op
PASS
ok  	github/bube054/jsonvx	2.694s -->

goos: windows
goarch: amd64
pkg: github/bube054/jsonvx
cpu: Intel(R) Core(TM) i7-7700HQ CPU @ 2.80GHz
BenchmarkJSONParserMediumPayload-8   	   12444	     90460 ns/op	  123402 B/op	      34 allocs/op
PASS
ok  	github/bube054/jsonvx	2.703s

BenchmarkJSONParserMediumPayload-8         15084             76120 ns/op          123400 B/op         34 allocs/op