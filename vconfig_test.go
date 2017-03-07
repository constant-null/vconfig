package vconfig

import (
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

type Test struct {
	Bool     bool
	Float64  float64
	Str      string `vconfig:",true"`
	Sub      sub
	usub     sub
	SrtSlice []string `vconfig:"str_slice"`
}

type sub struct {
	Int       int    `vconfig:",true"`
	SomeStaff string `vconfig:"some_staff"`
}

func TestMain(m *testing.M) {
	viper.SetConfigFile("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("TEST")

	os.Setenv("TEST_BOOL", "true")
	os.Setenv("TEST_FLOAT64", "1.01")
	os.Setenv("TEST_STR", "strs")
	os.Setenv("TEST_STR_SLICE", "val1 val2")
	os.Setenv("TEST_SUB_INT", "123")
	os.Setenv("TEST_SUB_SOME_STAFF", "test")

	retCode := m.Run()

	os.Clearenv()
	os.Exit(retCode)
}

func TestUmarshal(t *testing.T) {
	test := &Test{}
	err := Unmarshal(test)

	if err != nil {
		t.Errorf("Error while unmarshaling config %s", err)
	}

	expected := Test{
		true,
		1.01,
		"strs",
		sub{123, "test"},
		sub{},
		[]string{"val1", "val2"},
	}

	if !reflect.DeepEqual(expected, *test) {
		t.Fail()
	}
}

func TestUnmarshalErr(t *testing.T) {
	err := Unmarshal(Test{})
	if err == nil {
		t.Error("An error should occure")
	}

	os.Unsetenv("TEST_STR")

	err = Unmarshal(&Test{})
	if err == nil {
		t.Error("An error should occure")
	}
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
