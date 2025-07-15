package jsonvx

import (
	"fmt"
	"testing"
)

func TestJSONQuery(t *testing.T) {
	b := []byte(`{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`)

	p := New(b, nil)
	json, err := p.Parse()

	fmt.Println("json", json)

	if err != nil {
		t.Fatalf("got err %s:", err.Error())
	}

	jsonObject, ok := json.(*Object)
	// _ = ok

	if !ok {
		t.Fatalf("json is not array")
	}

	// arrayItem, err := jsonObject.QueryPath("name", "last")
	arrayItem, err := jsonObject.QueryPath(`fav.movie`)

	fmt.Println(arrayItem, err)
}
