// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package reflectex

import (
	"reflect"
	"testing"
	"time"
)

func TestStructPartialEqual(t *testing.T) {

	type (
		Src struct {
			FieldA string
			FieldB int
			FieldC bool
		}

		Dst struct {
			FieldA string
			FieldB float64
		}
	)

	src := &Src{"FieldA", 42, true}
	dst := &Dst{"DoNotWant", 3.14}

	if !StructPartialEqual(src, dst) {
		t.Fatal("StructPartialEqual failed")
	}

}

func TestLazyStructCopy(t *testing.T) {

	type (
		Src struct {
			FieldA string
			FieldB int
			FieldC bool
		}

		Dst struct {
			FieldA string
			FieldB float64
		}
	)

	src := &Src{"FieldA", 42, true}
	dst := &Dst{"DoNotWant", 3.14}
	if err := LazyStructCopy(src, dst); err != nil {
		t.Fatalf("LazyStructCopy failed: %v", err)
	}

	if src.FieldA != dst.FieldA {
		t.Fatal("LazyStructCopy failed.")
	}
}

func TestFilterStruct(t *testing.T) {

	type Test struct {
		Name    string
		Surname string
		Age     int
		nope    bool
	}

	in := &Test{"Foo", "Bar", 42, true}
	out := FilterStruct(in, "Name", "Surname")
	if !reflect.DeepEqual(out, &struct{ Age int }{42}) {
		t.Fatal("FilterStruct failed")
	}

}

func TestStringToValueTextUnmarshaler(t *testing.T) {

	val := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	in := now.Format(time.RFC3339Nano)
	out := reflect.ValueOf(&val)
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if !val.Equal(now) {
		t.Logf("in: %#v\n", now)
		t.Logf("out: %#v\n", val)
		t.Fatal("StringToValue(TextUnmarshaler) failed")
	}

}

func TestStringToValueBool(t *testing.T) {

	val := false
	in := "true"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if !val {
		t.Fatal("StringToValue(bool) failed")
	}

}

func TestStringToValueInt(t *testing.T) {

	val := 0
	in := "-42"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if val != -42 {
		t.Fatal("StringToValue(int) failed")
	}

}

func TestStringToValueUint(t *testing.T) {

	val := 0
	in := "1337"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if val != 1337 {
		t.Fatal("StringToValue(uint) failed")
	}

}

func TestStringToValueFloat32(t *testing.T) {

	val := float32(0.0)
	in := "3.14"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if val != 3.14 {
		t.Fatal("StringToValue(float32) failed")
	}

}

func TestStringToValueFloat64(t *testing.T) {

	val := float64(0.0)
	in := "3.14"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if val != 3.14 {
		t.Fatal("StringToValue(float64) failed")
	}

}

func TestStringToInterface(t *testing.T) {

	s := ""
	if err := StringToInterface("string", &s); err != nil {
		t.Fatal("string", err)
	}
	if s != "string" {
		t.Fatalf("StringToInterface(string) failed: want '%s', got '%s'", "string", s)
	}

	a := [5]int{0, 1, 2, 3, 4}
	if err := StringToInterface("9,8,7,6,5", &a); err != nil {
		t.Fatal("array", err)
	}
	if a != [5]int{9, 8, 7, 6, 5} {
		t.Fatalf("StringToInterface(array) failed: want '%s', got '%v'", "[9 8 7 6 5]", a)
	}

	sl := []string{"one", "two", "three"}
	if err := StringToInterface("red, green, blue", &sl); err != nil {
		t.Fatal("slice", err)
	}

	m := map[string]string{
		"apple":      "green",
		"banana":     "yellow",
		"grapefruit": "red",
	}
	if err := StringToInterface("allice=small,julie=petite,annie=fat(ish)", &m); err != nil {
		t.Fatal("map", err)
	}
}
