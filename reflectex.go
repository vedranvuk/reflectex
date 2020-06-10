// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package reflectex provides various reflect based utils.
package reflectex

import (
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

// LazyStructCopy copies values from src fields that have
// a coresponding field in dst to that field in dst.
// Fields must have same name and type. Tags are ignored.
// src and dest must be of struct type and addressable.
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

// StringToValue sets v's value to a value parsed from s.
// s must be convertible to v or an error is returned.
//
// TODO Describe syntax.
func StringToValue(s string, v reflect.Value) (err error) {
	if !v.IsValid() {
		return ErrInvalidParam
	}
	var parsedval reflect.Value
	switch v.Kind() {
	case reflect.Bool:
		var b bool
		b, err = strconv.ParseBool(s)
		if err != nil {
			break
		}
		parsedval = reflect.ValueOf(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		n, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			break
		}
		parsedval = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			break
		}
		parsedval = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Float32:
		var n float64
		n, err = strconv.ParseFloat(s, 32)
		if err != nil {
			break
		}
		parsedval = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Float64:
		var n float64
		n, err = strconv.ParseFloat(s, 64)
		if err != nil {
			break
		}
		parsedval = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Complex64:
		// TODO
	case reflect.Complex128:
		// TODO
	case reflect.String:
		parsedval = reflect.ValueOf(s)
	case reflect.Array:
		parsedval = reflect.Indirect(reflect.New(v.Type()))
		a := strings.Split(s, ",")
		for i, l := 0, v.Len(); i < l && i < len(a); i++ {
			if err = StringToValue(strings.TrimSpace(a[i]), parsedval.Index(i)); err != nil {
				break
			}
		}
	case reflect.Slice:
		a := strings.Split(s, ",")
		parsedval = reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), len(a), len(a))
		for i := 0; i < len(a); i++ {
			if err = StringToValue(a[i], parsedval.Index(i)); err != nil {
				break
			}
		}
	case reflect.Map:
		mt := reflect.MapOf(v.Type().Key(), v.Type().Elem())
		parsedval = reflect.MakeMap(mt)
		a := strings.Split(s, ",")
		for _, s := range a {
			pair := strings.Split(s, "=")
			if len(pair) != 2 {
				err = ErrParse
				break
			}
			key := reflect.Indirect(reflect.New(mt.Key()))
			if err = StringToValue(pair[0], key); err != nil {
				break
			}
			val := reflect.Indirect(reflect.New(mt.Elem()))
			if err = StringToValue(pair[1], val); err != nil {
				break
			}
			parsedval.SetMapIndex(key, val)
		}
	case reflect.Struct:
		// TODO
	case reflect.Func:
		// TODO
	case reflect.Chan:
		// TODO
	default:
		err = ErrUnsupported
	}
	if parsedval.IsValid() {
		v.Set(parsedval)
	}
	return err
}

// StringToInterface or error.
func StringToInterface(s string, i interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(i))
	return StringToValue(s, v)
}
