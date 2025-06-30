// Package gocustojson provides configurable options for parsing of JSON.
//
// It supports parsing a broad range of JSON syntax variants — from strict [ECMA-404-compliant JSON]
// to more permissive formats like [JSON5]. The parser behavior can be customized via the Config struct,
// which exposes fine-grained toggles for non-standard features such as comments, trailing commas,
// unquoted keys, single-quoted strings, and more.
//
// [ECMA-404-compliant JSON]: https://datatracker.ietf.org/doc/html/rfc7159
// [JSON5]: https://json5.org/
package gocustojson
