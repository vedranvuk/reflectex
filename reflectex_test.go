// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package reflectex

import (
	"reflect"
	"testing"
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
