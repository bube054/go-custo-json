package jsonvx

import (
	"fmt"
	"testing"
)

func TestJSONQuery(t *testing.T) {
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

	parser := NewParser(data, nil)

	// parse the JSON
	node, err := parser.Parse()
	if err != nil {
		t.Fatalf("failed to parse JSON: %s", err)
	}

	rootObj, ok := AsObject(node)
	if !ok {
		t.Fatalf("expected root node to be an object, but got: %s", err.Error())
	}

	// query the "age" field
	ageNode, err := rootObj.QueryPath("age")
	if err != nil {
		t.Fatalf("failed to query 'age' field: %s", err.Error())
	}

	// assert that the age field is a number
	ageNum, ok := AsNumber(ageNode)
	if !ok {
		t.Fatalf("expected 'age' to be a number, but got: %s", err.Error())
	}

	// get the value of the number
	ageValue, err := ageNum.Value()
	if err != nil {
		// t.Fatalf("failed to convert 'age' to numeric value: %s", err.Error())
	}

	fmt.Println(ageValue) // 37

}
