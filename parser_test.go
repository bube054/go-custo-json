package jsonvx

import (
	"errors"
	"testing"
)

type ParserTest struct {
	msg          string
	input        []byte
	cfg          *ParserConfig
	expectedNode JSON
	expectedErr  error
}

func runJSONParserTests(t *testing.T, tests []ParserTest) {
	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {

			parser := NewParser(test.input, test.cfg)
			node, err := parser.Parse()

			if node == nil || test.expectedNode == nil {
				if node != test.expectedNode || !errors.Is(err, test.expectedErr) {
					t.Errorf("got (%v, %v), expected (%v, %v)", node, err, test.expectedNode, test.expectedErr)
				}
			} else {
				if !node.Equal(test.expectedNode) || !errors.Is(err, test.expectedErr) {
					t.Errorf("got (%v, %v), expected (%v, %v)", node, err, test.expectedNode, test.expectedErr)
				}
			}

		})
	}
}

func TestJSONParserNothing(t *testing.T) {
	var tests = []ParserTest{
		{msg: "Parse nothing", input: []byte(``), expectedNode: nil, expectedErr: ErrJSONNoContent},
	}

	runJSONParserTests(t, tests)
}

func TestJSONParserWhitespace(t *testing.T) {
	var tests = []ParserTest{
		{msg: "Parse whitespace", input: []byte(``), expectedNode: nil, expectedErr: ErrJSONNoContent},
		{msg: "Parse multiple whitespace", input: []byte(` 	 `), expectedNode: nil, expectedErr: ErrJSONNoContent},
	}

	runJSONParserTests(t, tests)
}

func TestJSONParserComment(t *testing.T) {
	var tests = []ParserTest{
		{msg: "Parse single line comment", input: []byte("// line comment"), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse multiple line comments", input: []byte("// first comment\n// second comment"), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse block comment", input: []byte(`/* block comment */`), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewParserConfig(WithAllowBlockComments(true))},
		{msg: "Parse multiple block comments", input: []byte(`/* first comment */ /* second comment */`), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewParserConfig(WithAllowBlockComments(true))},
		{msg: "Parse primitive after line comment", input: []byte("// comment\n null"), expectedNode: &Null{Token: newTokenPtr(NULL, NONE, []byte("null"), 2, 2, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse primitive before line comment", input: []byte("null // comment"), expectedNode: &Null{Token: newTokenPtr(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse primitive before & after line comment", input: []byte("// comment\n null // comment \n"), expectedNode: &Null{Token: newTokenPtr(NULL, NONE, []byte("null"), 2, 2, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse invalid multiple json values surrounded by comments", input: []byte("// comment\n null // comment \n null"), expectedNode: nil, expectedErr: ErrJSONMultipleContent, cfg: NewParserConfig(WithAllowLineComments(true))},
	}

	runJSONParserTests(t, tests)
}

func TestJSONParserNull(t *testing.T) {
	var tests = []ParserTest{
		{msg: "Parse null", input: []byte(`null`), expectedNode: &Null{Token: newTokenPtr(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse null, with surrounding whitespace", input: []byte(` null `), expectedNode: &Null{Token: newTokenPtr(NULL, NONE, []byte("null"), 1, 2, nil)}, expectedErr: nil},
		{msg: "Parse null, with post line comment", input: []byte(`null // line comment`), expectedNode: &Null{Token: newTokenPtr(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse null, with post block comment", input: []byte(`null /* comment */`), expectedNode: &Null{Token: newTokenPtr(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowBlockComments(true))},
		{msg: "Parse null, with pre line comment", input: []byte("// line comment \n null"), expectedNode: &Null{Token: newTokenPtr(NULL, NONE, []byte("null"), 2, 2, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse null, with post block comment", input: []byte(`null /* comment */`), expectedNode: &Null{Token: newTokenPtr(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowBlockComments(true))},
		{msg: "Parse multiple nulls", input: []byte(`null null`), expectedNode: nil, expectedErr: ErrJSONMultipleContent},
		{msg: "Parse invalid character after null", input: []byte(`null x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar},
	}

	runJSONParserTests(t, tests)
}

func TestJSONParserBoolean(t *testing.T) {
	var tests = []ParserTest{
		// parse false
		{msg: "Parse false", input: []byte(`false`), expectedNode: &Boolean{Token: newTokenPtr(BOOLEAN, FALSE, []byte("false"), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse false, with surrounding whitespace", input: []byte(` false `), expectedNode: &Boolean{Token: newTokenPtr(BOOLEAN, FALSE, []byte("false"), 1, 2, nil)}, expectedErr: nil},
		{msg: "Parse false, with post line comment", input: []byte(`false // line comment`), expectedNode: &Boolean{Token: newTokenPtr(BOOLEAN, FALSE, []byte("false"), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse false, with post block comment", input: []byte(`false /* comment */`), expectedNode: &Boolean{Token: newTokenPtr(BOOLEAN, FALSE, []byte("false"), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowBlockComments(true))},
		{msg: "Parse multiple false's", input: []byte(`false false`), expectedNode: nil, expectedErr: ErrJSONMultipleContent},
		{msg: "Parse invalid character after false", input: []byte(`false x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar},

		// parse true
		{msg: "Parse true", input: []byte(`true`), expectedNode: &Boolean{Token: newTokenPtr(BOOLEAN, TRUE, []byte("true"), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse true, with surrounding whitespace", input: []byte(` true `), expectedNode: &Boolean{Token: newTokenPtr(BOOLEAN, TRUE, []byte("true"), 1, 2, nil)}, expectedErr: nil},
		{msg: "Parse true, with post line comment", input: []byte(`true // line comment`), expectedNode: &Boolean{Token: newTokenPtr(BOOLEAN, TRUE, []byte("true"), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse true, with post block comment", input: []byte(`true /* comment  */`), expectedNode: &Boolean{Token: newTokenPtr(BOOLEAN, TRUE, []byte("true"), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowBlockComments(true))},
		{msg: "Parse multiple trues", input: []byte(`true true`), expectedNode: nil, expectedErr: ErrJSONMultipleContent},
		{msg: "Parse invalid character after true", input: []byte(`true x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar},
	}

	runJSONParserTests(t, tests)
}

func TestJSONParserString(t *testing.T) {
	var tests = []ParserTest{
		// parse double quoted
		{msg: "Parse double quoted string", input: []byte(`"string"`), expectedNode: &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse double quoted string, with surrounding whitespace", input: []byte(`	"string"	`), expectedNode: &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 2, nil)}, expectedErr: nil},
		{msg: "Parse double quoted string, with post line comment", input: []byte(`"string" // line comment`), expectedNode: &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse double quoted string, with post block comment", input: []byte(`"string" /* comment  */`), expectedNode: &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"string"`), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowBlockComments(true))},
		{msg: "Parse multiple double quoted strings", input: []byte(`"string" "string"`), expectedNode: nil, expectedErr: ErrJSONMultipleContent},
		{msg: "Parse invalid character after double quoted string", input: []byte(`"string" x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar},

		// parse single quoted
		{msg: "Parse single quoted string", input: []byte(`'string'`), expectedNode: &String{Token: newTokenPtr(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowSingleQuotes(true))},
		{msg: "Parse single quoted string, with surrounding whitespace", input: []byte(`	'string'	`), expectedNode: &String{Token: newTokenPtr(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 2, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowSingleQuotes(true))},
		{msg: "Parse single quoted string, with post line comment", input: []byte(`'string' // line comment`), expectedNode: &String{Token: newTokenPtr(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true), WithAllowSingleQuotes(true))},
		{msg: "Parse single quoted string, with post block comment", input: []byte(`'string' /* comment */`), expectedNode: &String{Token: newTokenPtr(STRING, SINGLE_QUOTED, []byte(`'string'`), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowBlockComments(true), WithAllowSingleQuotes(true))},
		{msg: "Parse multiple single quoted strings", input: []byte(`'string' 'string'`), expectedNode: nil, expectedErr: ErrJSONMultipleContent, cfg: NewParserConfig(WithAllowSingleQuotes(true))},
		{msg: "Parse invalid character after single quoted string", input: []byte(`'string' x`), expectedNode: nil, expectedErr: ErrJSONUnexpectedChar, cfg: NewParserConfig(WithAllowSingleQuotes(true))},
	}

	runJSONParserTests(t, tests)
}

func TestJSONParserNumber(t *testing.T) {
	var tests = []ParserTest{
		// parse integer
		{msg: "Parse integer number", input: []byte(`123`), expectedNode: &Number{Token: newTokenPtr(NUMBER, INTEGER, []byte(`123`), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse integer number, with surrounding whitespace", input: []byte(`	123	`), expectedNode: &Number{Token: newTokenPtr(NUMBER, INTEGER, []byte(`123`), 1, 2, nil)}, expectedErr: nil},
		{msg: "Parse integer number, with post line comment", input: []byte(`123 // line comment`), expectedNode: &Number{Token: newTokenPtr(NUMBER, INTEGER, []byte(`123`), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse integer number, with post block comment", input: []byte(`123 /* comment */`), expectedNode: &Number{Token: newTokenPtr(NUMBER, INTEGER, []byte(`123`), 1, 1, nil)}, expectedErr: nil, cfg: NewParserConfig(WithAllowBlockComments(true))},

		// DO THE REMAINING NUMBER TYPES!!!!!!!!!!!
	}

	runJSONParserTests(t, tests)
}

func TestJSONParserArray(t *testing.T) {
	var tests = []ParserTest{
		{msg: "Parse empty array", input: []byte("[]"), expectedNode: &Array{Items: []JSON{}}, expectedErr: nil, cfg: NewParserConfig(WithAllowTrailingCommaArray(false))},
		{msg: "Parse no trailing comma array, with no trailing comma allowed", input: []byte(`[1,2,3]`), expectedNode: &Array{Items: []JSON{
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
		}}, expectedErr: nil},
		{msg: "Parse trailing comma array, with no trailing comma allowed", input: []byte(`[1,2,3,]`), expectedNode: nil, expectedErr: ErrJSONSyntax},
		{msg: "Parse trailing comma array, with trailing comma allowed", input: []byte(`[1,2,3,]`), expectedNode: &Array{Items: []JSON{
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
		}}, expectedErr: nil, cfg: NewParserConfig(WithAllowTrailingCommaArray(true))},
		{msg: "Parse array, with surrounding whitespace", input: []byte(`  [1,2,3]  `), expectedNode: &Array{Items: []JSON{
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("1"), 1, 4, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("2"), 1, 6, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("3"), 1, 8, nil)},
		}}, expectedErr: nil},
		{msg: "Parse array, with post line comment", input: []byte(`[1,2,3] // line comment`), expectedNode: &Array{Items: []JSON{
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
		}}, expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true))},
		{msg: "Parse array, with post block comment", input: []byte(`[1,2,3] /* block comment */`), expectedNode: &Array{Items: []JSON{
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("1"), 1, 2, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("2"), 1, 4, nil)},
			&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("3"), 1, 6, nil)},
		}}, expectedErr: nil, cfg: NewParserConfig(WithAllowBlockComments(true))},
		{msg: "Parse array of arrays", input: []byte(`[[1,2,3],[1,2,3]]`), expectedNode: &Array{Items: []JSON{
			&Array{Items: []JSON{
				&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("1"), 1, 3, nil)},
				&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("2"), 1, 5, nil)},
				&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("3"), 1, 7, nil)},
			}},
			&Array{Items: []JSON{
				&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("1"), 1, 11, nil)},
				&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("2"), 1, 13, nil)},
				&Number{Token: newTokenPtr(NUMBER, INTEGER, []byte("3"), 1, 15, nil)},
			}},
		}}, expectedErr: nil, cfg: nil},
	}

	runJSONParserTests(t, tests)
}

func TestJSONParserObject(t *testing.T) {
	var tests = []ParserTest{
		{msg: "Parse empty object", input: []byte("{}"), expectedNode: &Object{Properties: []KeyValue{}}, expectedErr: nil, cfg: NewParserConfig(WithAllowTrailingCommaObject(false))},
		{msg: "Parse no trailing comma object, with no trailing comma allowed", input: []byte(`{"key": "value"}`), expectedNode: &Object{
			Properties: []KeyValue{
				newKeyValue([]byte(`key`), &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"value"`), 1, 9, nil)}),
			}},
			expectedErr: nil,
		},
		{msg: "Parse trailing comma object, with no trailing comma allowed", input: []byte(`{"key": 1,}`), expectedNode: nil, expectedErr: ErrJSONSyntax},
		{msg: "Parse trailing comma object, with trailing comma allowed", input: []byte(`{"key": "value",}`), expectedNode: &Object{
			Properties: []KeyValue{
				newKeyValue([]byte(`key`), &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"value"`), 1, 9, nil)}),
			}},
			expectedErr: nil, cfg: NewParserConfig(WithAllowTrailingCommaObject(true)),
		},
		{msg: "Parse object, with surrounding whitespace", input: []byte(` 	 {"key": "value"} 	 `), expectedNode: &Object{
			Properties: []KeyValue{
				newKeyValue([]byte(`key`), &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"value"`), 1, 12, nil)}),
			}},
			expectedErr: nil, cfg: nil,
		},
		{msg: "Parse object, with post line comment", input: []byte(`{"key": "value"} // line comment`), expectedNode: &Object{
			Properties: []KeyValue{
				newKeyValue([]byte(`key`), &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"value"`), 1, 9, nil)}),
			}},
			expectedErr: nil, cfg: NewParserConfig(WithAllowLineComments(true)),
		},
		{msg: "Parse object, with post block comment", input: []byte(`{"key": "value"} /* block comment */`), expectedNode: &Object{
			Properties: []KeyValue{
				newKeyValue([]byte(`key`), &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"value"`), 1, 9, nil)}),
			}},
			expectedErr: nil, cfg: NewParserConfig(WithAllowBlockComments(true)),
		},
		{msg: "Parse object of objects", input: []byte(`{"key": {"key2": "value2"}}`), expectedNode: &Object{
			Properties: []KeyValue{
				newKeyValue([]byte(`key`), &Object{
					Properties: []KeyValue{
						newKeyValue([]byte(`key2`), &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"value2"`), 1, 18, nil)}),
					},
				}),
			}},
			expectedErr: nil, cfg: nil,
		},
	}
	runJSONParserTests(t, tests)
}
