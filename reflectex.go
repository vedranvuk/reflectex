// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package reflectex provides various reflect related utils.
package reflectex

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/vedranvuk/errorex"
)

var (
	ErrReflectEx        = errorex.New("reflectex")
	ErrInvalidParam     = ErrReflectEx.Wrap("invalid parameter")
	ErrUnsupportedValue = ErrReflectEx.Wrap("invalid parameter")
)

// LazyStructCopy copies src fields that have a coresponding dst field.
// Fields must have same name and type. Tags are ignored.
// src and dest must be of struct type and addressable.
func LazyStructCopy(src, dst interface{}) error {
	srcv := reflect.Indirect(reflect.ValueOf(src))
	dstv := reflect.Indirect(reflect.ValueOf(dst))
	if srcv.Kind() != reflect.Struct || dstv.Kind() != reflect.Struct {
		return ErrInvalidParam
	}
	for i := 0; i < srcv.NumField(); i++ {
		tgt := dstv.FieldByName(srcv.Type().Field(i).Name)
		if !tgt.IsValid() {
			continue
		}
		if tgt.Kind() != srcv.Field(i).Kind() {
			continue
		}
		tgt.Set(srcv.Field(i))
	}
	return nil
}

// StructPartialEqual compares two structs and tells if there is at least
// one field in both that match both by name and type.
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

// StringToValue converts a string to a value of v's type and returns it or an error.
func StringToValue(s string, v reflect.Value) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Bool:
		if b, e := strconv.ParseBool(s); e == nil {
			return reflect.ValueOf(b), nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, e := strconv.ParseInt(s, 10, 64); e == nil {
			return reflect.ValueOf(n).Convert(v.Type()), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if n, e := strconv.ParseUint(s, 10, 64); e == nil {
			return reflect.ValueOf(n).Convert(v.Type()), nil
		}
	case reflect.Float32:
		if n, e := strconv.ParseFloat(s, 32); e == nil {
			return reflect.ValueOf(n).Convert(v.Type()), nil
		}
	case reflect.Float64:
		if n, e := strconv.ParseFloat(s, 64); e == nil {
			return reflect.ValueOf(n).Convert(v.Type()), nil
		}
	case reflect.String:
		return reflect.ValueOf(s), nil
	case reflect.Array:
		av := reflect.Indirect(v)
		a := strings.Split(s, ",")
		for i, l := 0, v.Len(); i < l && i < len(a); i, l = i+1, l {
			rv, re := StringToValue(strings.TrimSpace(a[i]), v.Index(i))
			if re == nil {
				av.Index(i).Set(rv)
			}
		}
		return av, nil
	case reflect.Slice:
		st := v.Type().Elem()
		a := strings.Split(s, ";")
		sv := reflect.MakeSlice(reflect.SliceOf(st), len(a), len(a))
		for i := 0; i < len(a); i++ {
			rv, re := StringToValue(a[i], sv.Index(0))
			if re == nil {
				sv.Index(i).Set(rv)
			}
		}
		return sv, nil
	case reflect.Map:
		mt := reflect.MapOf(v.Type().Key(), v.Type().Elem())
		mv := reflect.MakeMap(mt)
		a := strings.Split(s, ";")
		for _, s := range a {
			pair := strings.Split(s, "=")
			if len(pair) != 2 {
				continue
			}
			key, err := StringToValue(pair[0], reflect.Zero(mt.Key()))
			if err != nil {
				continue
			}
			val, err := StringToValue(pair[1], reflect.Zero(mt.Elem()))
			if err != nil {
				continue
			}
			mv.SetMapIndex(key, val)
		}
		return mv, nil
	}

	return reflect.Value{}, ErrUnsupportedValue
}
