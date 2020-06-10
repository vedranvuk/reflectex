package reflectex

import (
	"reflect"
	"testing"
)

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
		t.Fatal("STructPartialEqual failed")
	}

}

func TestStringToValue(t *testing.T) {

	type Data struct {
		String string
		Array  [5]int
		Slice  []string
		Map    map[string]string
	}

	type Test struct {
		Value    string
		Expected error
	}

	data := &Data{
		"string",
		[5]int{0, 1, 2, 3, 4},
		[]string{"one", "two", "three"},
		map[string]string{
			"apple":      "green",
			"banana":     "yellow",
			"grapefruit": "red",
		},
	}

	dv := reflect.Indirect(reflect.ValueOf(data))
	for i := 0; i < dv.NumField(); i++ {
		// fmt.Println(dv.Type().Field(i).Name)
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

	// TODO more tests
}
