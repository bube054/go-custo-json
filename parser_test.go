package gocustojson

import (
	"fmt"
	"testing"
)

func TestJSONParser(t *testing.T) {
	// b := []byte(`null`)
	// b := []byte(`true`)
	// b := []byte(`false`)
	// b := []byte(`"abc"`)
	// b := []byte(`abc`)
	// b := []byte(`1245`)
	// b := []byte(` [ null , true , false , "abc" , 418 , ] `)
	// b := []byte(` [ null , true , false , "abc" , 418 , [1] ] `)
	b := []byte(`{"key1":"value1","key2":"value2","nullkey": null,"truthy": true,"falsy":false,"numint":12345,"num_float":123.456,"arr":[ null , true , false , "abc" , 418 , ],}`)

	p := New(b, NewConfig(WithAllowUnquoted(true)))
	json, err := p.Parse()

	if err != nil {
		t.Errorf("Got: %v", err)
	} else {
		fmt.Println(json)
	}

}
