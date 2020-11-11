// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package reflectex

import "testing"

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
