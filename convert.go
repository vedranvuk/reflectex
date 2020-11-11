// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package reflectex

import (
	"encoding"
	"reflect"
	"strconv"
	"strings"
)

// StringToInterface converts string in to out which must be a pointer to an
// allocated memory defining a type compatible to data contained in string
// according to rules defined in description of StringToValue.
func StringToInterface(in string, out interface{}) error {
	if out == nil {
		return ErrInvalidParam
	}
	return stringToInterface(in, out)
}

func stringToInterface(in string, out interface{}) error {
	return StringToValue(in, reflect.Indirect(reflect.ValueOf(out)))
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
// Array and Slice: Values delimited by comma.
// Example: 0,1,2,3,4
//
// Map: Key=Value pairs delimited by comma.
// Example:key1=value1,key2=value2,keyN=valueN
//
// Struct: TODO
// Example:TODO
//
// Chans and func are unsupported.
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
	case reflect.Ptr:
		return StringToPointerValue(in, out)
	}
	return ErrUnsupported
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
	n, err := strconv.ParseComplex(in, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToComplex128Value converts a string to a complex128.
func StringToComplex128Value(in string, out reflect.Value) error {
	n, err := strconv.ParseComplex(in, 128)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToStringValue converts a string to a string.
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

// StringToPointerValue converts a string to a pointer.
func StringToPointerValue(in string, out reflect.Value) error {
	nv := reflect.New(out.Type().Elem())
	if err := StringToValue(in, reflect.Indirect(nv)); err != nil {
		return err
	}
	out.Set(nv)
	return nil
}
