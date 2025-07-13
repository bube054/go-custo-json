package jsonvx

import (
	"fmt"
	"testing"
)

func TestJSONQuery(t *testing.T) {
	b := []byte(`[1,{"player":"palmer"},3]`)

	p := New(b, nil)
	json, err := p.Parse()

	fmt.Println("json", json)

	if err != nil {
		t.Fatalf("got err %s:", err.Error())
	}

	jsonArray, ok := json.(Array)

	if !ok {
		t.Fatalf("json is not array")
	}

	arrayItem, err := jsonArray.QueryPath("1", "player")

	fmt.Println(arrayItem, err)
}
