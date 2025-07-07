package jsonvx

import (
	"errors"
	"reflect"
	"testing"
)

type ParserTest struct {
	msg          string
	input        []byte
	cfg          *Config
	expectedNode JSONNode
	expectedErr  error
}

func RunJSONParserTests(t *testing.T, tests []ParserTest) {
	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			parser := New(test.input, test.cfg)
			node, err := parser.Parse()

			if !reflect.DeepEqual(node, test.expectedNode) || !errors.Is(err, test.expectedErr) {
				t.Errorf("got (%v, %v), expected (%v, %v)", node, err, test.expectedNode, test.expectedErr)
			}
		})
	}
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
		// {msg: "Parse primitive after line comment", input: []byte("// comment\n null"), expectedNode: JSONNull{Token: NewToken(NULL, NONE, []byte("null"), 2, 2, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		// {msg: "Parse primitive before line comment", input: []byte("null // comment"), expectedNode: JSONNull{Token: NewToken(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		// {msg: "Parse primitive before & after line comment", input: []byte("// comment\n null // comment \n null"), expectedNode: JSONNull{Token: NewToken(NULL, NONE, []byte("null"), 2, 2, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},

	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserNull(t *testing.T) {
	var tests = []ParserTest{
		// {msg: "Parse null", input: []byte(`null`), expectedNode: JSONNull{Token: NewToken(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil},
		// {msg: "Parse null, with surrounding whitespace", input: []byte(` null `), expectedNode: JSONNull{Token: NewToken(NULL, NONE, []byte("null"), 1, 2, nil)}, expectedErr: nil},
		{msg: "Parse null, with post line comment", input: []byte(`null // line comment`), expectedNode: JSONNull{Token: NewToken(NULL, NONE, []byte("null"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		// {msg: "Parse null, with post block comment", input: []byte(`null /*
		// */`), expectedNode: JSONNull{Token: NewToken(NULL, NONE, []byte("null"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},

		// {msg: "Parse null, with pre line comment", input: []byte("// line comment \n null"), expectedNode: JSONNull{Token: NewToken(NULL, NONE, []byte("null"), 2, 2, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		// {msg: "Parse null, with post block comment", input: []byte(`null /*
		// */`), expectedNode: JSONNull{Token: NewToken(NULL, NONE, []byte("null"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserBoolean(t *testing.T) {
	var tests = []ParserTest{
		// parse false
		{msg: "Parse false", input: []byte(`false`), expectedNode: JSONBoolean{Token: NewToken(FALSE, NONE, []byte("false"), 1, 0, nil)}, expectedErr: nil},
		{msg: "Parse false, with surrounding whitespace", input: []byte(` false `), expectedNode: JSONBoolean{Token: NewToken(FALSE, NONE, []byte("false"), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse false, with post line comment", input: []byte(`false // line comment`), expectedNode: JSONBoolean{Token: NewToken(FALSE, NONE, []byte("false"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Parse false, with post block comment", input: []byte(`false /*
		 */`), expectedNode: JSONBoolean{Token: NewToken(FALSE, NONE, []byte("false"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},

		// parse true
		{msg: "Parse true", input: []byte(`true`), expectedNode: JSONBoolean{Token: NewToken(TRUE, NONE, []byte("true"), 1, 0, nil)}, expectedErr: nil},
		{msg: "Parse true, with surrounding whitespace", input: []byte(` true `), expectedNode: JSONBoolean{Token: NewToken(TRUE, NONE, []byte("true"), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse true, with post line comment", input: []byte(`true // line comment`), expectedNode: JSONBoolean{Token: NewToken(TRUE, NONE, []byte("true"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Parse true, with post block comment", input: []byte(`true /*
		 */`), expectedNode: JSONBoolean{Token: NewToken(TRUE, NONE, []byte("true"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserString(t *testing.T) {
	var tests = []ParserTest{
		// parse double quote
		{msg: "Parse double string", input: []byte(`"string"`), expectedNode: JSONString{Token: NewToken(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 0, nil)}, expectedErr: nil},
		{msg: "Parse double string, with surrounding whitespace", input: []byte(`	"string"	`), expectedNode: JSONString{Token: NewToken(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse double string, with post line comment", input: []byte(`"string" // line comment`), expectedNode: JSONString{Token: NewToken(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Parse double string, with post block comment", input: []byte(`"string" /*
		 */`), expectedNode: JSONString{Token: NewToken(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},

		// parse single quote
		{msg: "Parse single string", input: []byte(`'string'`), expectedNode: JSONString{Token: NewToken(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Parse single string, with surrounding whitespace", input: []byte(`	'string'	`), expectedNode: JSONString{Token: NewToken(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 1, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowSingleQuotes(true))},
		{msg: "Parse single string, with post line comment", input: []byte(`'string' // line comment`), expectedNode: JSONString{Token: NewToken(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true), WithAllowSingleQuotes(true))},
		{msg: "Parse single string, with post block comment", input: []byte(`'string' /*
		  */`), expectedNode: JSONString{Token: NewToken(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true), WithAllowSingleQuotes(true))},
	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserNumber(t *testing.T) {
	var tests = []ParserTest{
		// parse integer
		{msg: "Parse integer number", input: []byte(`123`), expectedNode: JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte(`123`), 1, 0, nil)}, expectedErr: nil},
		{msg: "Parse integer number, with surrounding whitespace", input: []byte(`	123	`), expectedNode: JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte(`123`), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse integer number, with post line comment", input: []byte(`123 // line comment`), expectedNode: JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte(`123`), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Parse integer number, with post block comment", input: []byte(`123 /*
		 */`), expectedNode: JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte(`123`), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},

		// DO THE REMAINING NUMBER TYPES!!!!!!!!!!!
	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserArray(t *testing.T) {
	var tests = []ParserTest{
		// parse array
		{msg: "Parse empty array", input: []byte(`[[1,2,3,],[4,5,6,],]`), expectedNode: JSONArray{Items: []JSONNode{}}, expectedErr: nil, cfg: NewConfig(WithAllowTrailingCommaArray(false))},
		// {msg: "Parse no trailing comma array, with no trailing comma allowed", input: []byte(`[1,2,3]`), expectedNode: JSONArray{Items: []JSONNode{
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 1, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 3, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 5, nil)},
		// }}, expectedErr: nil},
		// {msg: "Parse trailing comma array, with no trailing comma allowed", input: []byte(`[1,2,3,]`), expectedNode: nil, expectedErr: ErrJSONSyntax},
		// {msg: "Parse trailing comma array, with trailing comma allowed", input: []byte(`[1,2,3,]`), expectedNode: JSONArray{Items: []JSONNode{
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 1, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 3, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 5, nil)},
		// }}, expectedErr: nil,  cfg: NewConfig(WithAllowTrailingCommaArray(true)),},
		// {msg: "Parse array, with surrounding whitespace", input: []byte(`  [1,2,3]  `), expectedNode: JSONArray{Items: []JSONNode{
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 3, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 5, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 7, nil)},
		// }}, expectedErr: nil},
		// {msg: "Parse array, with post line comment", input: []byte(`[1,2,3] // line comment`), expectedNode: JSONArray{Items: []JSONNode{
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 1, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 3, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 5, nil)},
		// }}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true)),},
		// {msg: "Parse array, with post block comment", input: []byte(`[1,2,3] /* block comment */`), expectedNode: JSONArray{Items: []JSONNode{
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 1, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 3, nil)},
		// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 5, nil)},
		// }}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true)),},
		// {msg: "Parse array of arrays", input: []byte(`[[1,2,3,],[1,2,3,],]`), expectedNode: JSONArray{Items: []JSONNode{
		// {msg: "Parse array of arrays", input: []byte(`[[1,2,3,],[4,5,6,]]`), expectedNode: JSONArray{Items: []JSONNode{
			// JSONArray{Items: []JSONNode{
			// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
			// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
			// 	JSONNumber{Token: NewToken(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
			// }},
			// JSONArray{Items: []JSONNode{}},
		// }}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
	}

	RunJSONParserTests(t, tests)
}
