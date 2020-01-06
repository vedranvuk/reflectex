package reflectex

import (
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

	type (
		Test struct {
			Name  string
			Admin bool
			Age   int
		}
	)

}
