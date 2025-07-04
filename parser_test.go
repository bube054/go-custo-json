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
		{msg: "Parse nothing, with multiple whitespace", input: []byte(` 	 `), expectedNode: nil, expectedErr: ErrJSONNoContent},
		{msg: "Parse nothing, with line comment", input: []byte(`// line comment`), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Parse nothing, with block comment", input: []byte(`/* block comment */`), expectedNode: nil, expectedErr: ErrJSONNoContent, cfg: NewConfig(WithAllowBlockComments(true))},
	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserNull(t *testing.T) {
	var tests = []ParserTest{
		{msg: "Parse null", input: []byte(`null`), expectedNode: JSONNull{token: NewToken(NULL, NONE, []byte("null"), 1, 0, nil)}, expectedErr: nil},
		{msg: "Parse null, with surrounding whitespace", input: []byte(` null `), expectedNode: JSONNull{token: NewToken(NULL, NONE, []byte("null"), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse null, with line comment", input: []byte(`null // line comment`), expectedNode: JSONNull{token: NewToken(NULL, NONE, []byte("null"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Parse null, with block comment", input: []byte(`null /*
		*/`), expectedNode: JSONNull{token: NewToken(NULL, NONE, []byte("null"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
	}

	RunJSONParserTests(t, tests)
}

func TestJSONParserBoolean(t *testing.T) {
	var tests = []ParserTest{
		// parse false
		{msg: "Parse false", input: []byte(`false`), expectedNode: JSONBoolean{token: NewToken(FALSE, NONE, []byte("false"), 1, 0, nil)}, expectedErr: nil},
		{msg: "Parse false, with surrounding whitespace", input: []byte(` false `), expectedNode: JSONBoolean{token: NewToken(FALSE, NONE, []byte("false"), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse false, with line comment", input: []byte(`false // line comment`), expectedNode: JSONBoolean{token: NewToken(FALSE, NONE, []byte("false"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Parse false, with block comment", input: []byte(`false /*
		// */`), expectedNode: JSONBoolean{token: NewToken(FALSE, NONE, []byte("false"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},

		// parse true
		{msg: "Parse true", input: []byte(`true`), expectedNode: JSONBoolean{token: NewToken(TRUE, NONE, []byte("true"), 1, 0, nil)}, expectedErr: nil},
		{msg: "Parse true, with surrounding whitespace", input: []byte(` true `), expectedNode: JSONBoolean{token: NewToken(TRUE, NONE, []byte("true"), 1, 1, nil)}, expectedErr: nil},
		{msg: "Parse true, with line comment", input: []byte(`true // line comment`), expectedNode: JSONBoolean{token: NewToken(TRUE, NONE, []byte("true"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowLineComments(true))},
		{msg: "Parse true, with block comment", input: []byte(`true /*
		// */`), expectedNode: JSONBoolean{token: NewToken(TRUE, NONE, []byte("true"), 1, 0, nil)}, expectedErr: nil, cfg: NewConfig(WithAllowBlockComments(true))},
	}

	RunJSONParserTests(t, tests)
}
