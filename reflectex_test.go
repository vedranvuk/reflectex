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
	if !reflect.DeepEqual(out, &struct{ Age int }{0}) {
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

func TestStringToPointerValue(t *testing.T) {

	in := "69"
	var val *int
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if *val != 69 {
		t.Fatal("StringToValue(pointer) failed")
	}

}

func TestStringToDeepPointerValue(t *testing.T) {

	in := "69"
	var val ***int
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if ***val != 69 {
		t.Fatal("StringToValue(pointer) failed")
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

func TestCompareInterfaceBool(t *testing.T) {
	if CompareInterfaces(false, false) != 0 {
		t.Fatal("TestCompareInterfaceBool failed.")
	}
	if CompareInterfaces(true, false) != 1 {
		t.Fatal("TestCompareInterfaceBool failed.")
	}
	if CompareInterfaces(false, true) != -1 {
		t.Fatal("TestCompareInterfaceBool failed.")
	}
}

func TestCompareInterfaceInt(t *testing.T) {
	if CompareInterfaces(2, 2) != 0 {
		t.Fatal("TestCompareInterfaceInt failed.")
	}
	if CompareInterfaces(0, 2) != -1 {
		t.Fatal("TestCompareInterfaceInt failed.")
	}
	if CompareInterfaces(2, 0) != 1 {
		t.Fatal("TestCompareInterfaceInt failed.")
	}
}

func TestCompareInterfaceUint(t *testing.T) {
	if CompareInterfaces(uint(2), uint(2)) != 0 {
		t.Fatal("TestCompareInterfaceUint failed.")
	}
	if CompareInterfaces(uint(0), uint(2)) != -1 {
		t.Fatal("TestCompareInterfaceUint failed.")
	}
	if CompareInterfaces(uint(2), uint(0)) != 1 {
		t.Fatal("TestCompareInterfaceUint failed.")
	}
}

func TestCompareInterfaceFloat(t *testing.T) {
	if CompareInterfaces(3.14, 3.14) != 0 {
		t.Fatal("TestCompareInterfaceFloat failed.")
	}
	if CompareInterfaces(3.14, 6.28) != -1 {
		t.Fatal("TestCompareInterfaceFloat failed.")
	}
	if CompareInterfaces(6.28, 3.14) != 1 {
		t.Fatal("TestCompareInterfaceFloat failed.")
	}
}

func TestCompareInterfaceComplex(t *testing.T) {
	if CompareInterfaces(3+4i, 3+4i) != 0 {
		t.Fatal("TestCompareInterfaceComplex failed.")
	}
	if CompareInterfaces(3+4i, 3+5i) != -1 {
		t.Fatal("TestCompareInterfaceComplex failed.")
	}
	if CompareInterfaces(3+5i, 3+4i) != 1 {
		t.Fatal("TestCompareInterfaceComplex failed.")
	}
}

func TestCompareInterfaceArraySlice(t *testing.T) {
	a := []int{0, 1, 2, 3, 4}
	b := []int{9, 8, 7, 6, 5}
	c := [5]int{0, 1, 2, 3, 4}

	if CompareInterfaces(a, a) != 0 {
		t.Fatal("TestCompareInterfaceArraySlice failed.")
	}
	if CompareInterfaces(a, b) != -1 {
		t.Fatal("TestCompareInterfaceArraySlice failed.")
	}
	if CompareInterfaces(b, a) != 1 {
		t.Fatal("TestCompareInterfaceArraySlice failed.")
	}
	if CompareInterfaces(a, c) != 1 {
		t.Fatal("TestCompareInterfaceArraySlice failed.")
	}
	if CompareInterfaces(c, c) != 0 {
		t.Fatal("TestCompareInterfaceArraySlice failed.")
	}
}

func TestCompareInterfaceMap(t *testing.T) {

	a := map[string]int{"0": 0, "1": 1, "2": 2}
	b := map[string]int{"1": 1, "2": 2, "3": 3}
	c := map[int]string{0: "0", 1: "1", 2: "2"}

	if CompareInterfaces(a, a) != 0 {
		t.Fatal("TestCompareInterfaceMap failed.")
	}
	if CompareInterfaces(a, b) != -1 {
		t.Fatal("TestCompareInterfaceMap failed.")
	}
	if CompareInterfaces(a, c) != 1 {
		t.Fatal("TestCompareInterfaceMap failed.")
	}
}

func TestCompareInterfaceString(t *testing.T) {

	a := "one"
	b := "two"

	if CompareInterfaces(a, a) != 0 {
		t.Fatal("TestCompareInterfaceString failed.")
	}

	if CompareInterfaces(a, b) != -1 {
		t.Fatal("TestCompareInterfaceString failed.")
	}

	if CompareInterfaces(b, a) != 1 {
		t.Fatal("TestCompareInterfaceString failed.")
	}
}

func TestCompareInterfaceStruct(t *testing.T) {

	type A struct {
		Name   string
		Age    int
		hidden bool
	}

	type B struct {
		Age  int
		Name string
	}

	type C struct {
		FOO string
		BAR int
	}

	type D struct {
		Field interface{}
	}

	if CompareInterfaces(&A{}, &A{}) != 0 {
		t.Fatal("TestCompareInterfaceStruct failed")
	}

	if CompareInterfaces(&A{}, &B{}) != 0 {
		t.Fatal("TestCompareInterfaceStruct failed")
	}

	if CompareInterfaces(&A{}, &C{}) != -1 {
		t.Fatal("TestCompareInterfaceStruct failed")
	}

	if CompareInterfaces(&C{}, &A{}) != 1 {
		t.Fatal("TestCompareInterfaceStruct failed")
	}

	if CompareInterfaces(&C{"Foo", 42}, &C{"Foo", 43}) != -1 {
		t.Fatal("TestCompareInterfaceStruct failed")
	}

	if CompareInterfaces(&C{"Bar", 1337}, &C{"Bar", 1337}) != 0 {
		t.Fatal("TestCompareInterfaceStruct failed")
	}

	if CompareInterfaces(&D{42}, &D{42}) != 0 {
		t.Fatal("TestCompareInterfaceStruct failed")
	}
}

func TestCompareInterfacesInterface(t *testing.T) {

	var a interface{}
	var b interface{}

	a = "abc"
	b = "def"

	if CompareInterfaces(a, a) != 0 {
		t.Fatal("TestCompareInterfacesInterface failed.")
	}
	if CompareInterfaces(a, b) != -1 {
		t.Fatal("TestCompareInterfacesInterface failed.")
	}
	if CompareInterfaces(b, a) != 1 {
		t.Fatal("TestCompareInterfacesInterface failed.")
	}
}

func TestCompareInterfacesPointer(t *testing.T) {

	a := int(1)
	b := int(2)

	pa := &a
	pb := &b
	var pc *int = nil

	if CompareInterfaces(pa, pa) != 0 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
	if CompareInterfaces(pa, pb) != -1 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
	if CompareInterfaces(pb, pa) != 1 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
	if CompareInterfaces(pa, pc) != 1 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
	if CompareInterfaces(pc, pc) != 0 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
}

func TestCompareInterfacesUnsafepointer(t *testing.T) {

	a := uintptr(1000)
	b := uintptr(2000)

	pa := &a
	pb := &b
	var pc uintptr = 0

	if CompareInterfaces(pa, pa) != 0 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
	if CompareInterfaces(pa, pb) != -1 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
	if CompareInterfaces(pb, pa) != 1 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
	if CompareInterfaces(pa, pc) != 1 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
	if CompareInterfaces(pc, pc) != 0 {
		t.Fatal("TestCompareInterfacesPointer failed.")
	}
}

func BenchmarkStructPartialEqual(b *testing.B) {

	type TestA struct {
		AField0 string
		AField1 int
		AField2 uint
		AField3 float32
		AField4 float64
		AField5 complex64
		AField6 complex128
		AField7 rune
		AField8 bool
		Field9  []byte
	}

	type TestB struct {
		BField0 string
		BField1 int
		BField2 uint
		BField3 float32
		BField4 float64
		BField5 complex64
		BField6 complex128
		BField7 rune
		BField8 bool
		Field9  []byte
	}

	for i := 0; i < b.N; i++ {
		StructPartialEqual(&TestA{}, &TestB{})
	}
}

func BenchmarkLazyStructCopy(b *testing.B) {

	b.StopTimer()

	type TestA struct {
		Field0 string
		Field1 int
		Field2 uint
		Field3 float32
		Field4 float64
		Field5 complex64
		Field6 complex128
		Field7 rune
		Field8 bool
		Field9 []byte
	}

	type TestB struct {
		Field0 string
		Field1 int
		Field2 uint
		Field3 float32
		Field4 float64
		Field5 complex64
		Field6 complex128
		Field7 rune
		Field8 bool
		Field9 []byte
	}

	testA := &TestA{"one", 2, 3, 4.0, 5.0, 6.0i, 7.0i, '8', true, []byte("nein")}
	testB := &TestB{}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		LazyStructCopy(testA, testB)
	}
}

func BenchmarkFilterStructCopy(b *testing.B) {

	type Test struct {
		Field0 string
		Field1 int
		Field2 uint
		Field3 float32
		Field4 float64
		Field5 complex64
		Field6 complex128
		Field7 rune
		Field8 bool
		Field9 []byte
	}

	for i := 0; i < b.N; i++ {
		FilterStruct(&Test{})
	}
}

func BenchmarkFilterStructFilter(b *testing.B) {

	type Test struct {
		Field0 string
		Field1 int
		Field2 uint
		Field3 float32
		Field4 float64
		Field5 complex64
		Field6 complex128
		Field7 rune
		Field8 bool
		Field9 []byte
	}

	for i := 0; i < b.N; i++ {
		FilterStruct(&Test{}, "Field0", "Field1", "Field2", "Field3", "Field4", "Field5", "Field6", "Field7", "Field8", "Field9")
	}
}

func BenchmarkCompareInterfaces(b *testing.B) {

	b.StopTimer()

	type Test struct {
		Field0 string
		Field1 int
		Field2 uint
		Field3 float32
		Field4 float64
		Field5 complex64
		Field6 complex128
		Field7 rune
		Field8 bool
		Field9 []byte
	}

	test := &Test{"one", 2, 3, 4.0, 5.0, 6.0i, 7.0i, '8', true, []byte("nein")}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		CompareInterfaces(test, test)
	}
}
