// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package reflectex provides various reflect based utils.
package reflectex

import (
	"bytes"
	"encoding"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/vedranvuk/errorex"
)

var (
	// ErrReflectEx is the base error of reflectex package.
	ErrReflectEx = errorex.New("reflectex")
	// ErrInvalidParam is returned when an invalid param is passed to a func.
	ErrInvalidParam = ErrReflectEx.Wrap("invalid parameter")
	// ErrParse is returned when a parse error occurs.
	ErrParse = ErrReflectEx.Wrap("parse error")
	// ErrUnsupported is returned when an unsupported value is encountered.
	ErrUnsupported = ErrReflectEx.Wrap("unsupported value")
	// ErrConvert is returned when a conversion is unable to complete.
	ErrConvert = ErrReflectEx.WrapFormat("cannot convert '%s' to type '%s'")

	// ErrNotImplemented help.
	ErrNotImplemented = ErrReflectEx.WrapFormat("NOT IMPLEMENTED '%s'")
)

// StructPartialEqual compares two structs and tells if there is at least
// one field in both that match both by name and type.
// Tags in both x and y are ignored.
func StructPartialEqual(x, y interface{}) bool {
	xv := reflect.Indirect(reflect.ValueOf(x))
	yv := reflect.Indirect(reflect.ValueOf(y))
	if xv.Kind() != reflect.Struct || yv.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < xv.NumField(); i++ {
		tgt := yv.FieldByName(xv.Type().Field(i).Name)
		if !tgt.IsValid() {
			continue
		}
		return true
	}
	return false
}

// LazyStructCopy copies values from src fields that have a coresponding field
// in dst to that field in dst. Fields must have same name and type. Tags are
// ignored. src and dest must be of struct type and addressable.
func LazyStructCopy(src, dst interface{}) error {
	srcv := reflect.Indirect(reflect.ValueOf(src))
	dstv := reflect.Indirect(reflect.ValueOf(dst))
	if srcv.Kind() != reflect.Struct || dstv.Kind() != reflect.Struct {
		return ErrInvalidParam
	}
	for i := 0; i < srcv.NumField(); i++ {
		name := srcv.Type().Field(i).Name
		tgt := dstv.FieldByName(name)
		if !tgt.IsValid() {
			continue
		}
		if tgt.Kind() != srcv.Field(i).Kind() {
			continue
		}
		if name == "_" {
			continue
		}
		if name[0] >= 97 && name[0] <= 122 {
			continue
		}
		tgt.Set(srcv.Field(i))
	}
	return nil
}

// FilterStruct returns a copy of in struct with specified fields removed.
// In must be a pointer to a struct or a struct value.
// Values of non-filtered fields are not copied from the source to result.
// Returned value is a struct value or nil in case of an error.
func FilterStruct(in interface{}, filter ...string) interface{} {

	v := reflect.Indirect(reflect.ValueOf(in))
	if !v.IsValid() {
		return nil
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	sort.Strings(filter)

	fields := make([]reflect.StructField, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanSet() {
			continue
		}
		pos := sort.SearchStrings(filter, v.Type().Field(i).Name)
		if pos < len(filter) && filter[pos] == v.Type().Field(i).Name {
			continue
		}
		fields = append(fields, v.Type().Field(i))
	}

	structType := reflect.StructOf(fields)
	structVal := reflect.New(structType)

	return structVal.Interface()
}

// StringToBoolValue converts a string to a bool.
func StringToBoolValue(in string, out reflect.Value) error {
	b, err := strconv.ParseBool(in)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(b))
	return nil
}

// StringToIntValue converts a string to a int of any width.
func StringToIntValue(in string, out reflect.Value) error {
	n, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToUintValue converts a string to an uint of any width.
func StringToUintValue(in string, out reflect.Value) error {
	n, err := strconv.ParseUint(in, 10, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToFloat32Value converts a string to a float32.
func StringToFloat32Value(in string, out reflect.Value) error {
	n, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToFloat64Value converts a string to a float64.
func StringToFloat64Value(in string, out reflect.Value) error {
	n, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToComplex64Value converts a string to a complex64.
func StringToComplex64Value(in string, out reflect.Value) error {
	// TODO Implement StringToComplex64
	return ErrNotImplemented.WrapArgs("StringToComplex64Value")
}

// StringToComplex128Value converts a string to a complex128.
func StringToComplex128Value(in string, out reflect.Value) error {
	// TODO Implement StringToComplex128
	return ErrNotImplemented.WrapArgs("StringToComplex128Value")

}

// StringToStringValue converts a string to a string.
// A real eye opener.
func StringToStringValue(in string, out reflect.Value) error {
	out.Set(reflect.ValueOf(in))
	return nil
}

// StringToArrayValue converts a string to an array.
func StringToArrayValue(in string, out reflect.Value) error {
	v := reflect.Indirect(reflect.New(out.Type()))
	a := strings.Split(in, ",")
	for i, l := 0, out.Len(); i < l && i < len(a); i++ {
		if err := StringToValue(strings.TrimSpace(a[i]), v.Index(i)); err != nil {
			return err
		}
	}
	out.Set(v)
	return nil
}

// StringToSliceValue converts a string to a slice.
func StringToSliceValue(in string, out reflect.Value) error {
	a := strings.Split(in, ",")
	parsedval := reflect.MakeSlice(reflect.SliceOf(out.Type().Elem()), len(a), len(a))
	for i := 0; i < len(a); i++ {
		if err := StringToValue(a[i], parsedval.Index(i)); err != nil {
			return err
		}
	}
	out.Set(parsedval)
	return nil
}

// StringToMapValue converts a string to a map.
func StringToMapValue(in string, out reflect.Value) error {
	mt := reflect.MapOf(out.Type().Key(), out.Type().Elem())
	parsedval := reflect.MakeMap(mt)
	a := strings.Split(in, ",")
	for _, s := range a {
		pair := strings.Split(s, "=")
		if len(pair) != 2 {
			return ErrParse
		}
		key := reflect.Indirect(reflect.New(mt.Key()))
		if err := StringToValue(pair[0], key); err != nil {
			return err
		}
		val := reflect.Indirect(reflect.New(mt.Elem()))
		if err := StringToValue(pair[1], val); err != nil {
			return err
		}
		parsedval.SetMapIndex(key, val)
	}
	out.Set(parsedval)
	return nil
}

// StringToStructValue converts a string to a struct.
func StringToStructValue(in string, out reflect.Value) error {
	// TODO Implement StringToStruct
	return ErrNotImplemented.WrapArgs("StringToStructValue")
}

// StringToValue intends to set out to a value parsed from in which must be
// convertible to out or an error is returned. If out is a compound type its'
// value(s) is replaced.
//
// StringToValue tries to be a one call converter to many different value kinds
// in one unifying call, primarily for conversion of simple types such as bools
// and numbers and types which implement TextUnmarshaler. Parsing compound
// types like arrays, slices, maps and structs requires a defined syntax so a
// simple, possibly logical syntax is implemented for completeness sake, as
// described:
//
// Array and Slice: Values enclosed in square brackets, delimited by comma.
// Example: [0,1,2,3,4]
//
// Map: Key/Value pairs enclosed in square brackets, delimited by comma.
// Example:[key1=value1,key2=value2,keyN=valueN]
//
// Struct: Values enclosed in curly braces, delimited by comma.
// Example:{value1,value2,[1,2,3],[key1=value1,key2=value2],{value1,value2}}
//
// Keys and Values can be enclosed in double quotes to retain spaces and
// special characters. Inner quotes must be escaped. Basically a mini-json.
//
// Pointers, chans and func are unsupported.
//
// If an error occurs it is returned.
func StringToValue(in string, out reflect.Value) error {

	bum, ok := out.Interface().(encoding.TextUnmarshaler)
	if ok {
		if err := bum.UnmarshalText([]byte(in)); err != nil {
			return err
		}
		return nil
	}

	out = reflect.Indirect(out)
	if !out.IsValid() {
		return ErrInvalidParam
	}

	switch out.Kind() {
	case reflect.Bool:
		return StringToBoolValue(in, out)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return StringToIntValue(in, out)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return StringToUintValue(in, out)
	case reflect.Float32:
		return StringToFloat32Value(in, out)
	case reflect.Float64:
		return StringToFloat64Value(in, out)
	case reflect.Complex64:
		return StringToComplex64Value(in, out)
	case reflect.Complex128:
		return StringToComplex128Value(in, out)
	case reflect.String:
		return StringToStringValue(in, out)
	case reflect.Array:
		return StringToArrayValue(in, out)
	case reflect.Slice:
		return StringToSliceValue(in, out)
	case reflect.Map:
		return StringToMapValue(in, out)
	case reflect.Struct:
		return StringToStructValue(in, out)
	}
	return ErrUnsupported
}

// StringToInterface converts string in to out which must be a pointer to an
// allocated memory defining a type compatible to data contained in string
// according to rules defined in description of StringToValue.
func StringToInterface(in string, out interface{}) error {
	return StringToValue(in, reflect.ValueOf(out))
}

// compareKind compares a and b reflect.Kind as integer index in enumerarion and
//
// Returns a negative number (-1) if a is less than b.
// Returns a zero (0) if a is equal to b.
// Returns a positive number (1) if a is more than b.
func compareKind(a, b reflect.Kind) int {
	if int(a) > int(b) {
		return 1
	}
	if int(a) < int(b) {
		return -1
	}
	return 0
}

// CompareValues recursively compares two possibly compound values a and b for
// equality. It returns:
//
// a negative number (-1) if a is less than b.
// a zero (0) if a is equal to b.
// a positive number (1) if a is more than b.
//
// Comparison is done using magic as described below:
//
// Types are compared for reflect.Kind equality first and foremost. Comparison
// logic is taken as a and b's index in reflect.Kind enumeration.
//
// For structs, only published fields are enumerated and compared. Private
// fields do not affect comparison. A struct with less public fields returns a
// less result. Structs with equal number of fields are compared alphabetically
// ascending comparing field value kinds, names and finally values.
//
// Comparisons between two Arrays and/or slices return a less result for values
// with less dimensions. In equal dimensioned arrays or slices bytes are
// compared using bytes.Compare().
//
// Maps with less elements return a less result. Maps with equal number of
// keys are compared for same kinds, then key names and finally equal values.
// Comparison is done in ascending order after converting keys to strings/bytes.
//
// Pointer types are dereferenced do their values before comparison. This
// includes UnsafePointer and Uintptr.
//
// Complex numbers are compared as strings.
//
// Channel and func types are not supported.
//
// If an error occurs it is returned with a compare value that should be
// disregarded.
//
func CompareValues(a, b reflect.Value) int {

	if res := compareKind(a.Kind(), b.Kind()); res != 0 {
		return res
	}

	// Dereference all pointers to their concrete value.
	// TODO sync dereference.

	if a.Kind() == reflect.Ptr {
		for {
			a = a.Elem()
			if a.Kind() != reflect.Ptr {
				break
			}
		}
	}

	if b.Kind() == reflect.Ptr {
		for {
			b = b.Elem()
			if b.Kind() != reflect.Ptr {
				break
			}
		}
	}

	// TODO Dereference interfaces.

	// Compare by reflect.Kind.

	switch a.Kind() {
	case reflect.Bool:
		if a.Bool() == b.Bool() {
			return 0
		}
		if a.Bool() {
			return 1
		}
		return -1

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if res := compareKind(a.Kind(), b.Kind()); res != 0 {
			return res
		}
		if a.Int() == b.Int() {
			return 0
		}
		if a.Int() > b.Int() {
			return 1
		}
		return -1

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if res := compareKind(a.Kind(), b.Kind()); res != 0 {
			return res
		}
		if a.Uint() == b.Uint() {
			return 0
		}
		if a.Uint() > b.Uint() {
			return 1
		}
		return -1

	case reflect.Float32, reflect.Float64:
		if res := compareKind(a.Kind(), b.Kind()); res != 0 {
			return res
		}
		if a.Float() == b.Float() {
			return 0
		}
		if a.Float() > b.Float() {
			return 1
		}
		return -1
	case reflect.Complex64, reflect.Complex128:
		if res := compareKind(a.Kind(), b.Kind()); res != 0 {
			return res
		}
		if a.String() == b.String() {
			return 0
		}
		if a.String() > b.String() {
			return 1
		}
		return -1
	case reflect.Array, reflect.Slice:
		// TODO Compare indices sequentially.
		if res := compareKind(a.Kind(), b.Kind()); res != 0 {
			return res
		}
		if a.Len() == b.Len() {
			return bytes.Compare(a.Bytes(), b.Bytes())
		}
		if a.Len() > b.Len() {
			return 1
		}
		return -1

	case reflect.Map:

		// Compare lengths.
		if a.Len() > b.Len() {
			return 1
		}
		if a.Len() < b.Len() {
			return -1
		}

		// Compare keys.
		akeys := a.MapKeys()
		bkeys := b.MapKeys()
		sort.Slice(akeys, func(i, j int) bool {
			return akeys[i].String() < akeys[j].String()
		})
		sort.Slice(bkeys, func(i, j int) bool {
			return bkeys[i].String() < bkeys[j].String()
		})
		for i := 0; i < len(akeys); i++ {
			if res := compareKind(akeys[i].Kind(), bkeys[i].Kind()); res != 0 {
				return res
			}
			if res := strings.Compare(akeys[i].String(), bkeys[i].String()); res != 0 {
				return res
			}
		}

		// Compare values.
		for i := 0; i < len(akeys); i++ {
			aval := a.MapIndex(akeys[i])
			bval := b.MapIndex(bkeys[i])
			if a.Kind() != b.Kind() {
				if res := compareKind(aval.Kind(), bval.Kind()); res != 0 {
					return res
				}
				if res := CompareValues(aval, bval); res != 0 {
					return res
				}
			}
		}

	case reflect.String:
		return strings.Compare(a.String(), b.String())

	case reflect.Struct:

		// Enum public fields.
		aflds := make([]reflect.StructField, 0, a.NumField())
		bflds := make([]reflect.StructField, 0, b.NumField())
		for i := 0; i < a.NumField(); i++ {
			if !a.Field(i).CanSet() {
				continue
			}
			aflds = append(aflds, a.Type().Field(i))
		}
		for i := 0; i < b.NumField(); i++ {
			if !b.Field(i).CanSet() {
				continue
			}
			bflds = append(bflds, b.Type().Field(i))
		}

		// Compare by public field count.
		if len(aflds) > len(bflds) {
			return 1
		}
		if len(aflds) < len(bflds) {
			return -1
		}

		// Sort the fields and compare by kind and name.
		sort.Slice(aflds, func(i, j int) bool {
			return aflds[i].Name < aflds[j].Name
		})
		sort.Slice(bflds, func(i, j int) bool {
			return bflds[i].Name < bflds[j].Name
		})
		for i := 0; i < len(aflds); i++ {
			// Compare kind.
			if res := compareKind(aflds[i].Type.Kind(), bflds[i].Type.Kind()); res != 0 {
				return res
			}
			// Compare field name.
			if res := strings.Compare(aflds[i].Name, bflds[i].Name); res != 0 {
				return res
			}
			// Compare field value.
			if res := CompareValues(a.FieldByName(aflds[i].Name), b.FieldByName(bflds[i].Name)); res != 0 {
				return res
			}
		}

	case reflect.Interface:

	case reflect.Uintptr:
	case reflect.Ptr:
	case reflect.UnsafePointer:

	default:
		// Ignore, return 0.
	}

	return 0
}

// CompareInterfaces compares two interfaces for equality between the types
// contained within them. See CompareValue for details.
func CompareInterfaces(a, b interface{}) int {
	return CompareValues(reflect.ValueOf(a), reflect.ValueOf(b))
}
