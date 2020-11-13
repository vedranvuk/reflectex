// Copyright 2020 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package reflectex provides various reflect based utils.
package reflectex

import (
	"fmt"
	"reflect"
	"sort"

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

// Show shows a reflect.Value.
func Show(v reflect.Value) {
	fmt.Printf(`Type:    %s 
Kind:    %s 
IsValid: %t 
IsZero:  %t 
CanAddr: %t
`, v.Type(), v.Kind(), v.IsValid(), v.IsZero(), v.CanAddr())
}
