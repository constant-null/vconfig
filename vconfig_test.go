package vconfig

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type Test struct {
	Str      string `vconfig:",true"`
	Sub      sub
	sub      sub
	SrtSlice []string `vconfig:"str_slice"`
}

type sub struct {
	Int       int    `vconfig:",true"`
	SomeStaff string `vconfig:"some_staff"`
}

func TestParse(t *testing.T) {
	os.Setenv("TEST_STR", "strs")
	os.Setenv("TEST_STR_SLICE", "val1 val2")
	os.Setenv("TEST_SUB_INT", "123")
	os.Setenv("TEST_SUB_SOME_STAFF", "test")
	configure()

	test := &Test{}
	err := Unmarshal(test)
	fmt.Println(err)

	fmt.Printf("%+v", test)
}

func TestParseTag(t *testing.T) {
	field1 := reflect.StructField{Name: "Test", Tag: reflect.StructTag(`vconfig:"field1,true"`)}
	tag1 := parseTag(field1, "pref")

	if tag1.Name != "pref.field1" || tag1.Required != true {
		t.Error("Parsing error")
	}

	field2 := reflect.StructField{Name: "Test", Tag: reflect.StructTag(`vconfig:"field1"`)}
	tag2 := parseTag(field2, "")

	if tag2.Name != "field1" || tag2.Required != false {
		t.Error("Parsing error", "")
	}

	field3 := reflect.StructField{Name: "Test", Tag: reflect.StructTag(`vconfig:",true"`)}
	tag3 := parseTag(field3, "pref")

	if tag3.Name != "pref.test" || tag3.Required != true {
		t.Error("Parsing error")
	}

	field4 := reflect.StructField{Name: "Test"}
	tag4 := parseTag(field4, "")

	if tag4.Name != "test" || tag4.Required != false {
		t.Error("Parsing error")
	}
}
