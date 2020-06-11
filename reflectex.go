// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package reflectex provides various reflect based utils.
package reflectex

import (
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
)

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
// Returned value is a struct value or nil in case of an error.
func FilterStruct(in interface{}, filter ...string) interface{} {
	v := reflect.Indirect(reflect.ValueOf(in))
	if !v.IsValid() {
		return nil
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	fields := []reflect.StructField{}
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanSet() {
			continue
		}
		sort.Strings(filter)
		pos := sort.SearchStrings(filter, v.Type().Field(i).Name)
		if pos < len(filter) && filter[pos] == v.Type().Field(i).Name {
			continue
		}
		fields = append(fields, v.Type().Field(i))
	}
	structType := reflect.StructOf(fields)
	structVal := reflect.New(structType)

	if err := LazyStructCopy(v.Interface(), structVal.Interface()); err != nil {
		return nil
	}
	return structVal.Interface()
}

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

// StringToBool converts a string to a bool.
func StringToBool(in string, out reflect.Value) error {
	b, err := strconv.ParseBool(in)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(b))
	return nil
}

// StringToInt converts a string to a int of any width.
func StringToInt(in string, out reflect.Value) error {
	n, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToUint converts a string to an uint of any width.
func StringToUint(in string, out reflect.Value) error {
	n, err := strconv.ParseUint(in, 10, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToFloat32 converts a string to a float32.
func StringToFloat32(in string, out reflect.Value) error {
	n, err := strconv.ParseFloat(in, 32)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToFloat64 converts a string to a float64.
func StringToFloat64(in string, out reflect.Value) error {
	n, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return err
	}
	out.Set(reflect.ValueOf(n).Convert(out.Type()))
	return nil
}

// StringToComplex64 converts a string to a complex64.
func StringToComplex64(in string, out reflect.Value) error {
	// TODO Implement StringToComplex64
	return nil
}

// StringToComplex128 converts a string to a complex128.
func StringToComplex128(in string, out reflect.Value) error {
	// TODO Implement StringToComplex128
	return nil
}

// StringToString converts a string to a string.
// A real eye opener.
func StringToString(in string, out reflect.Value) error {
	out.Set(reflect.ValueOf(in))
	return nil
}

// StringToArray converts a string to an array.
func StringToArray(in string, out reflect.Value) error {
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

// StringToSlice converts a string to a slice.
func StringToSlice(in string, out reflect.Value) error {
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

// StringToMap converts a string to a map.
func StringToMap(in string, out reflect.Value) error {
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

// StringToStruct converts a string to a struct.
func StringToStruct(in string, out reflect.Value) error {
	// TODO Implement StringToStruct
	return nil
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
		return StringToBool(in, out)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return StringToInt(in, out)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return StringToUint(in, out)
	case reflect.Float32:
		return StringToFloat32(in, out)
	case reflect.Float64:
		return StringToFloat64(in, out)
	case reflect.Complex64:
		return StringToComplex64(in, out)
	case reflect.Complex128:
		return StringToComplex128(in, out)
	case reflect.String:
		return StringToString(in, out)
	case reflect.Array:
		return StringToArray(in, out)
	case reflect.Slice:
		return StringToSlice(in, out)
	case reflect.Map:
		return StringToMap(in, out)
	case reflect.Struct:
		return StringToStruct(in, out)
	}
	return ErrUnsupported
}

// StringToInterface converts string in to out which must be a pointer to an
// allocated memory defining a type compatible to data contained in string
// according to rules defined in description of StringToInterface.
func StringToInterface(in string, out interface{}) error {
	return StringToValue(in, reflect.ValueOf(out))
}
