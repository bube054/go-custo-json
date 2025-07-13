package jsonvx

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type ParserTest struct {
	msg          string
	input        []byte
	cfg          *Config
	expectedNode JSON
	expectedErr  error
}

func RunJSONParserTests(t *testing.T, tests []ParserTest) {
	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			start := time.Now()

			parser := New(test.input, test.cfg)
			node, err := parser.Parse()

			elapsed := time.Since(start)
			fmt.Printf("parsing took %.4f seconds\n", elapsed.Seconds())
			// fmt.Printf("parsing took %s\n", elapsed)

			if !reflect.DeepEqual(node, test.expectedNode) || !errors.Is(err, test.expectedErr) {
				t.Errorf("got (%v, %v), expected (%v, %v)", node, err, test.expectedNode, test.expectedErr)
			}
		})
	}
}

// please remove later
func TestJSONParserXYZ(t *testing.T) {
	var tests = []ParserTest{
		// 		{msg: "Parse anything", input: []byte(`[
		//   1,
		//   2,
		//   3,
		//   [4, 5, 6, [7, 8, [9, 10], 11], 12],
		//   13,
		//   [14, [15, 16], 17],
		//   18,
		//   19,
		//   [20, 21, [22, 23, [24, 25, 26], 27], 28],
		//   29,
		//   [30, [31, 32, [33, [34, 35], 36], 37], 38],
		//   39,
		//   40,
		//   [41, 42, [43, 44, [45, 46], [47, [48, 49, [50, 51], 52], 53]], 54],
		//   55,
		//   56,
		//   57,
		//   [58, [59, 60, [61, 62, [63, [64, [65, 66], 67], 68], 69], 70]],
		//   71,
		//   [72, 73, [74, [75, 76], [77, [78, 79]]]],
		//   80,
		//   81,
		//   [82, 83, 84, [85, [86, [87, [88, 89, [90]]]]]],
		//   91,
		//   92,
		//   93,
		//   94,
		//   95,
		//   [96, [97, [98, [99, 100]]]]
		// ]`), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewConfig(WithAllowHexNumbers(true))},
		// 		{msg: "Parse anything", input: []byte(`{
		//   "name": "Example",
		//   "age": 30,
		//   "active": true,
		//   "score": null,
		//   "tags": ["go", "json", null, true, 123],
		//   "profile": {
		//     "id": "user_123",
		//     "email": "user@example.com",
		//     "preferences": {
		//       "notifications": true,
		//       "theme": "dark",
		//       "languages": ["en", "fr", "es"]
		//     }
		//   },
		//   "metrics": {
		//     "visits": 12345,
		//     "conversion": 0.023,
		//     "history": [
		//       {
		//         "date": "2025-07-10",
		//         "value": 42
		//       },
		//       {
		//         "date": "2025-07-09",
		//         "value": 37
		//       }
		//     ]
		//   },
		//   "nested": {
		//     "a": {
		//       "b": {
		//         "c": {
		//           "d": {
		//             "e": {
		//               "flag": false,
		//               "data": [1, 2, {"x": "deep", "y": null}]
		//             }
		//           }
		//         }
		//       }
		//     }
		//   }
		// }`), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewConfig(WithAllowHexNumbers(true))},
		{msg: "Parse anything", input: []byte(mediumPayload), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewConfig(WithAllowHexNumbers(true))},
	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserNothing(t *testing.T) {
	var tests = []ParserTest{
		{msg: "Parse nothing", input: []byte(``), expectedNode: nil, expectedErr: ErrJSONNoContent},
	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserWhitespace(t *testing.T) {
	var tests = []ParserTest{
		{msg: "Parse whitespace", input: []byte(``), expectedNode: nil, expectedErr: ErrJSONNoContent},
		{msg: "Parse multiple whitespace", input: []byte(` 	 `), expectedNode: nil, expectedErr: ErrJSONNoContent},
	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserComment(t *testing.T) {
	var tests = []ParserTest{
		// {msg: "Parse single line comment", input: []byte("// line comment"), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewConfig(WithAllowLineComments(true))},
		// {msg: "Parse multiple line comments", input: []byte("// first comment\n// second comment"), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewConfig(WithAllowLineComments(true))},
		// {msg: "Parse block comment", input: []byte(`/* block comment */`), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewConfig(WithAllowBlockComments(true))},
		// {msg: "Parse multiple block comments", input: []byte(`/* first comment */ /* second comment */`), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewConfig(WithAllowBlockComments(true))},
		// {msg: "Parse primitive after line comment", input: []byte("// comment\n null"), expectedNode: Null{Token: NewToken(NULL, NONE, []byte("null"), 2, 2, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		// {msg: "Parse primitive before line comment", input: []byte("null // comment"), expectedNode: Null{Token: NewToken(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		// {msg: "Parse primitive before & after line comment", input: []byte("// comment\n null // comment \n"), expectedNode: Null{Token: NewToken(NULL, NONE, []byte("null"), 2, 2, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		// {msg: "Parse invalid multiple json values surrounded by comments", input: []byte("// comment\n null // comment \n null"), expectedNode: nil, expectedErr: ErrJSONMultipleContent, cfg: NewConfig(WithAllowLineComments(true))},
	}

	RunJSONParserTests(t, tests)
}

// func TestJSONParserNull(t *testing.T) {
// 	var tests = []ParserTest{
// 		{msg: "Parse null", input: []byte(`null`), expectedNode: Null{Token: NewToken(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil},
// 		{msg: "Parse null, with surrounding whitespace", input: []byte(` null `), expectedNode: Null{Token: NewToken(NULL, NONE, []byte("null"), 1, 2, nil)}, expectedErr: nil},
// 		{msg: "Parse null, with post line comment", input: []byte(`null // line comment`), expectedNode: Null{Token: NewToken(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
// 		{msg: "Parse null, with post block comment", input: []byte(`null /* comment */`), expectedNode: Null{Token: NewToken(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
// 		{msg: "Parse null, with pre line comment", input: []byte("// line comment \n null"), expectedNode: Null{Token: NewToken(NULL, NONE, []byte("null"), 2, 2, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
// 		{msg: "Parse null, with post block comment", input: []byte(`null /* comment */`), expectedNode: Null{Token: NewToken(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
// 		{msg: "Parse multiple nulls", input: []byte(`null null`), expectedNode: nil, expectedErr: ErrJSONMultipleContent},
// 		{msg: "Parse invalid character after null", input: []byte(`null x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar},
// 	}

// 	RunJSONParserTests(t, tests)
// }

// func TestJSONParserBoolean(t *testing.T) {
// 	var tests = []ParserTest{
// 		// parse false
// 		{msg: "Parse false", input: []byte(`false`), expectedNode: Boolean{Token: NewToken(FALSE, NONE, []byte("false"), 1, 1, nil)}, expectedErr: nil},
// 		{msg: "Parse false, with surrounding whitespace", input: []byte(` false `), expectedNode: Boolean{Token: NewToken(FALSE, NONE, []byte("false"), 1, 2, nil)}, expectedErr: nil},
// 		{msg: "Parse false, with post line comment", input: []byte(`false // line comment`), expectedNode: Boolean{Token: NewToken(FALSE, NONE, []byte("false"), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
// 		{msg: "Parse false, with post block comment", input: []byte(`false /* comment */`), expectedNode: Boolean{Token: NewToken(FALSE, NONE, []byte("false"), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
// 		{msg: "Parse multiple false's", input: []byte(`false false`), expectedNode: nil, expectedErr: ErrJSONMultipleContent},
// 		{msg: "Parse invalid character after false", input: []byte(`false x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar},

// 		// parse true
// 		{msg: "Parse true", input: []byte(`true`), expectedNode: Boolean{Token: NewToken(TRUE, NONE, []byte("true"), 1, 1, nil)}, expectedErr: nil},
// 		{msg: "Parse true, with surrounding whitespace", input: []byte(` true `), expectedNode: Boolean{Token: NewToken(TRUE, NONE, []byte("true"), 1, 2, nil)}, expectedErr: nil},
// 		{msg: "Parse true, with post line comment", input: []byte(`true // line comment`), expectedNode: Boolean{Token: NewToken(TRUE, NONE, []byte("true"), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
// 		{msg: "Parse true, with post block comment", input: []byte(`true /* comment  */`), expectedNode: Boolean{Token: NewToken(TRUE, NONE, []byte("true"), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
// 		{msg: "Parse multiple trues", input: []byte(`true true`), expectedNode: nil, expectedErr: ErrJSONMultipleContent},
// 		{msg: "Parse invalid character after true", input: []byte(`true x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar},
// 	}

// 	RunJSONParserTests(t, tests)
// }

// func TestJSONParserString(t *testing.T) {
// 	var tests = []ParserTest{
// 		// parse double quoted
// 		{msg: "Parse double quoted string", input: []byte(`"string"`), expectedNode: String{Token: NewToken(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 1, nil)}, expectedErr: nil},
// 		{msg: "Parse double quoted string, with surrounding whitespace", input: []byte(`	"string"	`), expectedNode: String{Token: NewToken(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 2, nil)}, expectedErr: nil},
// 		{msg: "Parse double quoted string, with post line comment", input: []byte(`"string" // line comment`), expectedNode: String{Token: NewToken(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
// 		{msg: "Parse double quoted string, with post block comment", input: []byte(`"string" /* comment  */`), expectedNode: String{Token: NewToken(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
// 		{msg: "Parse multiple double quoted strings", input: []byte(`"string" "string"`), expectedNode: nil, expectedErr: ErrJSONMultipleContent},
// 		{msg: "Parse invalid character after double quoted string", input: []byte(`"string" x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar},

// 		// parse single quoted
// 		{msg: "Parse single quoted string", input: []byte(`'string'`), expectedNode: String{Token: NewToken(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowSingleQuotes(true))},
// 		{msg: "Parse single quoted string, with surrounding whitespace", input: []byte(`	'string'	`), expectedNode: String{Token: NewToken(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 2, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowSingleQuotes(true))},
// 		{msg: "Parse single quoted string, with post line comment", input: []byte(`'string' // line comment`), expectedNode: String{Token: NewToken(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true), WithAllowSingleQuotes(true))},
// 		{msg: "Parse single quoted string, with post block comment", input: []byte(`'string' /* comment */`), expectedNode: String{Token: NewToken(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true), WithAllowSingleQuotes(true))},
// 		{msg: "Parse multiple single quoted strings", input: []byte(`'string' 'string'`), expectedNode: nil, expectedErr: ErrJSONMultipleContent, cfg: NewConfig(WithAllowSingleQuotes(true))},
// 		{msg: "Parse invalid character after single quoted string", input: []byte(`'string' x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar, cfg: NewConfig(WithAllowSingleQuotes(true))},
// 	}

// 	RunJSONParserTests(t, tests)
// }

// func TestJSONParserNumber(t *testing.T) {
// 	var tests = []ParserTest{
// 		// parse integer
// 		{msg: "Parse integer number", input: []byte(`123`), expectedNode: Number{Token: NewToken(NUMBER, INTEGER, []byte(`123`), 1, 1, nil)}, expectedErr: nil},
// 		{msg: "Parse integer number, with surrounding whitespace", input: []byte(`	123	`), expectedNode: Number{Token: NewToken(NUMBER, INTEGER, []byte(`123`), 1, 2, nil)}, expectedErr: nil},
// 		{msg: "Parse integer number, with post line comment", input: []byte(`123 // line comment`), expectedNode: Number{Token: NewToken(NUMBER, INTEGER, []byte(`123`), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
// 		{msg: "Parse integer number, with post block comment", input: []byte(`123 /* comment */`), expectedNode: Number{Token: NewToken(NUMBER, INTEGER, []byte(`123`), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},

// 		// DO THE REMAINING NUMBER TYPES!!!!!!!!!!!
// 	}

// 	RunJSONParserTests(t, tests)
// }

// func TestJSONParserArray(t *testing.T) {
// 	var tests = []ParserTest{
// 		{msg: "Parse empty array", input: []byte("[]"), expectedNode: Array{Items: []JSON{}}, expectedErr: nil, cfg: NewConfig(WithAllowTrailingCommaArray(false))},
// 		{msg: "Parse no trailing comma array, with no trailing comma allowed", input: []byte(`[1,2,3]`), expectedNode: Array{Items: []JSON{
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
// 		}}, expectedErr: nil},
// 		{msg: "Parse trailing comma array, with no trailing comma allowed", input: []byte(`[1,2,3,]`), expectedNode: nil, expectedErr: ErrJSONSyntax},
// 		{msg: "Parse trailing comma array, with trailing comma allowed", input: []byte(`[1,2,3,]`), expectedNode: Array{Items: []JSON{
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
// 		}}, expectedErr: nil, cfg: NewConfig(WithAllowTrailingCommaArray(true))},
// 		{msg: "Parse array, with surrounding whitespace", input: []byte(`  [1,2,3]  `), expectedNode: Array{Items: []JSON{
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 4, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 6, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 8, nil)},
// 		}}, expectedErr: nil},
// 		{msg: "Parse array, with post line comment", input: []byte(`[1,2,3] // line comment`), expectedNode: Array{Items: []JSON{
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
// 		}}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
// 		{msg: "Parse array, with post block comment", input: []byte(`[1,2,3] /* block comment */`), expectedNode: Array{Items: []JSON{
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
// 			Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
// 		}}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
// 		{msg: "Parse array of arrays", input: []byte(`[[1,2,3],[1,2,3]]`), expectedNode: Array{Items: []JSON{
// 			Array{Items: []JSON{
// 				Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 3, nil)},
// 				Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 5, nil)},
// 				Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 7, nil)},
// 			}},
// 			Array{Items: []JSON{
// 				Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 11, nil)},
// 				Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 13, nil)},
// 				Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 15, nil)},
// 			}},
// 		}}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
// 	}

// 	RunJSONParserTests(t, tests)
// }

// func TestJSONParserObject(t *testing.T) {
// 	var tests = []ParserTest{
// 		{msg: "Parse empty object", input: []byte(`
// {
//   "widget": {
//     "debug": "on",
//     "window": {
//       "title": "Sample Konfabulator Widget",
//       "name": "main_window",
//       "width": 500,
//       "height": 500
//     },
//     "image": {
//       "src": "Images/Sun.png",
//       "hOffset": 250,
//       "vOffset": 250,
//       "alignment": "center"
//     },
//     "text": {
//       "data": "Click Here",
//       "size": 36,
//       "style": "bold",
//       "vOffset": 100,
//       "alignment": "center",
//       "onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
//     }
//   }
// }
// 		`), expectedNode: Object{Properties: map[string]JSON{}}, expectedErr: nil, cfg: NewConfig(WithAllowTrailingCommaObject(false))},
// 		// {msg: "Parse empty object", input: []byte("{}"), expectedNode: Object{Properties: map[string]JSON{}}, expectedErr: nil, cfg: NewConfig(WithAllowTrailingCommaObject(false))},
// 		// {msg: "Parse no trailing comma object, with no trailing comma allowed", input: []byte(`{"key": "value"}`), expectedNode: Object{
// 		// 	Properties: map[string]JSON{"key": String{Token: NewToken(STRING, DOUBLE_QUOTED, []byte(`"value"`), 1, 9, nil)},},
// 		// }, expectedErr: nil},
// 		// {msg: "Parse trailing comma object, with no trailing comma allowed", input: []byte(`{"key": 1,}`), expectedNode: nil, expectedErr: ErrJSONSyntax},
// 		// {msg: "Parse trailing comma object, with trailing comma allowed", input: []byte(`{"key": "value"}`), expectedNode: Array{Items: []JSON{
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
// 		// }}, expectedErr: nil, cfg: NewConfig(WithAllowTrailingCommaObject(true))},
// 		// {msg: "Parse object, with surrounding whitespace", input: []byte(`  [1,2,3]  `), expectedNode: Array{Items: []JSON{
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 4, nil)},
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 6, nil)},
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 8, nil)},
// 		// }}, expectedErr: nil},
// 		// {msg: "Parse object, with post line comment", input: []byte(`[1,2,3] // line comment`), expectedNode: Array{Items: []JSON{
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
// 		// }}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
// 		// {msg: "Parse object, with post block comment", input: []byte(`[1,2,3] /* block comment */`), expectedNode: Array{Items: []JSON{
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
// 		// 	Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
// 		// }}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
// 		// {msg: "Parse object of arrays", input: []byte(`[[1,2,3],[1,2,3]]`), expectedNode: Array{Items: []JSON{
// 		// 	Array{Items: []JSON{
// 		// 		Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 3, nil)},
// 		// 		Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 5, nil)},
// 		// 		Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 7, nil)},
// 		// 	}},
// 		// 	Array{Items: []JSON{
// 		// 		Number{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 11, nil)},
// 		// 		Number{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 13, nil)},
// 		// 		Number{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 15, nil)},
// 		// 	}},
// 		// }}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
// 	}

// 	RunJSONParserTests(t, tests)
// }

func BenchmarkJSONParserAny(b *testing.B) {
	json := []byte(`[
  1,
  2,
  3,
  [4, 5, 6, [7, 8, [9, 10], 11], 12],
  13,
  [14, [15, 16], 17],
  18,
  19,
  [20, 21, [22, 23, [24, 25, 26], 27], 28],
  29,
  [30, [31, 32, [33, [34, 35], 36], 37], 38],
  39,
  40,
  [41, 42, [43, 44, [45, 46], [47, [48, 49, [50, 51], 52], 53]], 54],
  55,
  56,
  57,
  [58, [59, 60, [61, 62, [63, [64, [65, 66], 67], 68], 69], 70]],
  71,
  [72, 73, [74, [75, 76], [77, [78, 79]]]],
  80,
  81,
  [82, 83, 84, [85, [86, [87, [88, 89, [90]]]]]],
  91,
  92,
  93,
  94,
  95,
  [96, [97, [98, [99, 100]]]]
]`)
	p := New(json, nil)
	for i := 0; i < b.N; i++ {
		p.Parse()
	}
}

func BenchmarkJSONParserSmallPayload(b *testing.B) {
	json := []byte(smallPayload)
	p := New(json, nil)
	for i := 0; i < b.N; i++ {
		p.Parse()
	}
}

func BenchmarkJSONParserMediumPayload(b *testing.B) {
	json := []byte(mediumPayload)
	p := New(json, nil)
	for i := 0; i < b.N; i++ {
		p.Parse()
	}
}

func BenchmarkJSONParserLargePayload(b *testing.B) {
	json := []byte(`{
  "person": {
    "id": "d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
    "name": {
      "fullName": "Leonid Bugaev",
      "givenName": "Leonid",
      "familyName": "Bugaev"
    },
    "email": "leonsbox@gmail.com",
    "gender": "male",
    "location": "Saint Petersburg, Saint Petersburg, RU",
    "geo": {
      "city": "Saint Petersburg",
      "state": "Saint Petersburg",
      "country": "Russia",
      "lat": 59.9342802,
      "lng": 30.3350986
    },
    "bio": "Senior engineer at Granify.com",
    "site": "http://flickfaver.com",
    "avatar": "https://d1ts43dypk8bqh.cloudfront.net/v1/avatars/d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
    "employment": {
      "name": "www.latera.ru",
      "title": "Software Engineer",
      "domain": "gmail.com"
    },
    "facebook": {
      "handle": "leonid.bugaev"
    },
    "github": {
      "handle": "buger",
      "id": 14009,
      "avatar": "https://avatars.githubusercontent.com/u/14009?v=3",
      "company": "Granify",
      "blog": "http://leonsbox.com",
      "followers": 95,
      "following": 10
    },
    "twitter": {
      "handle": "flickfaver",
      "id": 77004410,
      "bio": null,
      "followers": 2,
      "following": 1,
      "statuses": 5,
      "favorites": 0,
      "location": "",
      "site": "http://flickfaver.com",
      "avatar": null
    },
    "linkedin": {
      "handle": "in/leonidbugaev"
    },
    "googleplus": {
      "handle": null
    },
    "angellist": {
      "handle": "leonid-bugaev",
      "id": 61541,
      "bio": "Senior engineer at Granify.com",
      "blog": "http://buger.github.com",
      "site": "http://buger.github.com",
      "followers": 41,
      "avatar": "https://d1qb2nb5cznatu.cloudfront.net/users/61541-medium_jpg?1405474390"
    },
    "klout": {
      "handle": null,
      "score": null
    },
    "foursquare": {
      "handle": null
    },
    "aboutme": {
      "handle": "leonid.bugaev",
      "bio": null,
      "avatar": null
    },
    "gravatar": {
      "handle": "buger",
      "urls": [

      ],
      "avatar": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
      "avatars": [
        {
          "url": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
          "type": "thumbnail"
        }
      ]
    },
    "fuzzy": false
  },
  "company": null
}`)
	p := New(json, nil)
	for i := 0; i < b.N; i++ {
		p.Parse()
	}
}
