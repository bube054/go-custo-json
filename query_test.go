package jsonvx

import (
	"errors"
	"testing"
)

type QueryTest struct {
	msg          string
	input        []byte
	cfg          *ParserConfig
	queryPaths   []string
	expectedNode JSON
	expectedErr  error
}

func runJSONQueryTests(t *testing.T, tests []QueryTest) {
	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {

			parser := NewParser(test.input, test.cfg)
			node, _ := parser.Parse()

			var result JSON
			var err error

			switch val := node.(type) {
			case *Array:
				result, err = val.QueryPath(test.queryPaths...)
			case *Object:
				result, err = val.QueryPath(test.queryPaths...)

			default:
				t.Errorf("JSON %v must be queryable", node)
			}

			if result == nil || test.expectedNode == nil {
				if result != test.expectedNode || !errors.Is(err, test.expectedErr) {
					t.Errorf("got (%v, %v), expected (%v, %v)", result, err, test.expectedNode, test.expectedErr)
				}
			} else {
				if !result.Equal(test.expectedNode) || !errors.Is(err, test.expectedErr) {

					t.Errorf("got (%v, %v), expected (%v, %v)", result, err, test.expectedNode, test.expectedErr)
				}
			}

		})
	}
}

func TestJSONQuery(t *testing.T) {
	var data = []byte(`{
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

	var tests = []QueryTest{
		{
			msg:          "Query string `name.first`",
			input:        data,
			cfg:          nil,
			queryPaths:   []string{"name", "first"},
			expectedNode: &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"Tom"`), 2, 23, nil)},
			expectedErr:  nil,
		},
		{
			msg:        "Query array `children`",
			input:      data,
			cfg:        nil,
			queryPaths: []string{"children"},
			expectedNode: &Array{
				Items: []JSON{
					&String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"Sara"`), 4, 18, nil)},
					&String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"Alex"`), 4, 26, nil)},
					&String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"Jack"`), 4, 34, nil)},
				},
			},
			expectedErr: nil,
		},
		{
			msg:          "Query number `friends.2.age`",
			input:        data,
			cfg:          nil,
			queryPaths:   []string{"friends", "2", "age"},
			expectedNode: &Number{Token: newTokenPtr(NUMBER, INTEGER, []byte(`47`), 9, 50, nil)},
			expectedErr:  nil,
		},
		{
			msg:        "Query object `name`",
			input:      data,
			cfg:        nil,
			queryPaths: []string{"name"},
			expectedNode: &Object{
				Properties: []KeyValue{
					newKeyValue([]byte(`first`), &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"Tom"`), 2, 23, nil)}),
					newKeyValue([]byte(`last`), &String{Token: newTokenPtr(STRING, DOUBLE_QUOTED, []byte(`"Anderson"`), 2, 38, nil)}),
				},
			},
			expectedErr: nil,
		},
	}

	runJSONQueryTests(t, tests)
}
