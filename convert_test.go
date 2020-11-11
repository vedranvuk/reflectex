// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package reflectex

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

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

func BenchmarkStringToValueBool(b *testing.B) {
	val := false
	in := "true"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
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

func BenchmarkStringToValueInt(b *testing.B) {
	val := 0
	in := "-42"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
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

func BenchmarkStringToValueUint(b *testing.B) {
	val := 0
	in := "1337"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
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
func BenchmarkStringToValueFloat32(b *testing.B) {
	val := float32(0.0)
	in := "3.14"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
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

func BenchmarkStringToValueFloat64(b *testing.B) {
	val := float64(0.0)
	in := "3.14"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
	}
}

func TestStringToValueComplex64(t *testing.T) {
	val := complex64(complex(0, 0))
	in := "3.14+10i"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if val != 3.14+10i {
		t.Fatal("StringToValue(complex64) failed")
	}
}

func BenchmarkStringToValueComplex64(b *testing.B) {
	val := complex64(complex(0, 0))
	in := "3.14+10i"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
	}
}

func TestStringToValueComplex128(t *testing.T) {
	val := complex128(complex(0, 0))
	in := "3.14+10i"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if val != 3.14+10i {
		t.Fatal("StringToValue(complex128) failed")
	}
}

func BenchmarkStringToValueComplex128(b *testing.B) {
	val := complex128(complex(0, 0))
	in := "3.14+10i"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
	}
}

func TestStringToValueString(t *testing.T) {
	val := string("")
	in := "foobar"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if val != "foobar" {
		t.Fatal("StringToValue(string) failed")
	}
}

func BenchmarkStringToValueString(b *testing.B) {
	val := string("")
	in := "foobar"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
	}
}

func TestStringToValueArray(t *testing.T) {
	val := [3]int{0, 0, 0}
	in := "1,2,3"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if val != [3]int{1, 2, 3} {
		t.Fatal("StringToValue(array) failed")
	}
}

func BenchmarkStringToValueArray(b *testing.B) {
	val := [3]int{0, 0, 0}
	in := "1,2,3"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
	}
}

func TestStringToValueSlice(t *testing.T) {
	val := []byte{}
	in := "1,2,3"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(val, []byte{1, 2, 3}) != 0 {
		t.Fatal("StringToValue(slice) failed")
	}
}

func BenchmarkStringToValueSlice(b *testing.B) {
	val := []int{0, 0, 0}
	in := "1,2,3"
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
	}
}

func TestStringToValueMap(t *testing.T) {
	val := map[string]int{}
	in := "1=1,2=2,3=3"
	out := reflect.Indirect(reflect.ValueOf(&val))
	if err := StringToValue(in, out); err != nil {
		t.Fatal(err)
	}
	expect := map[string]int{"1": 1, "2": 2, "3": 3}
	for valk, valv := range val {
		if expv, ok := expect[valk]; !ok {
			t.Fatal("StringToValue(map) failed")
		} else {
			if expv != valv {
				t.Fatal("StringToValue(map) failed")
			}
		}
	}
}

func BenchmarkStringToValueMap(b *testing.B) {
	val := map[string]int{}
	in := "1=1,2=2,3=3"
	out := reflect.Indirect(reflect.ValueOf(&val))
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
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

func BenchmarkStringToPointerValue(b *testing.B) {
	in := "69"
	var val *int
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
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
		t.Fatal("StringToValue(deeppointer) failed")
	}
}

func BenchmarkStringToDeepPointerValue(b *testing.B) {
	in := "69"
	var val ***int
	out := reflect.Indirect(reflect.ValueOf(&val))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringToValue(in, out)
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
