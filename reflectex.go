// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package reflectex provides various reflect based utils.
package reflectex

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/vedranvuk/errorex"
)

var (
	ErrReflectEx    = errorex.New("reflectex")
	ErrInvalidParam = ErrReflectEx.Wrap("invalid parameter")
	ErrParse        = ErrReflectEx.Wrap("parse error")
	ErrUnsupported  = ErrReflectEx.Wrap("unsupported value")
	ErrConvert      = ErrReflectEx.WrapFormat("cannot convert '%s' to type '%s'")
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

// StringToValue converts a string to a value of v's type and returns it or an error.
func StringToValue(s string, v reflect.Value) (rv reflect.Value, err error) {
	if !v.IsValid() {
		return reflect.Value{}, ErrInvalidParam
	}
	switch v.Kind() {
	case reflect.Bool:
		var b bool
		b, err = strconv.ParseBool(s)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		n, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Float32:
		var n float64
		n, err = strconv.ParseFloat(s, 32)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Float64:
		var n float64
		n, err = strconv.ParseFloat(s, 64)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(n).Convert(v.Type())
	case reflect.String:
		rv = reflect.ValueOf(s)
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
	case reflect.Struct:

	case reflect.Ptr:
	}
	if err != nil {
		return reflect.Value{}, ErrConvert.WithArgs(s, v.Type().Name())
	}
	return
}

// StringToInterface or error.
func StringToInterface(s string, i interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(i))
	if !v.IsValid() {
		return ErrInvalidParam
	}
	var rv reflect.Value
	var err error
	switch v.Kind() {
	case reflect.Bool:
		var b bool
		b, err = strconv.ParseBool(s)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		n, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Float32:
		var n float64
		n, err = strconv.ParseFloat(s, 32)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(n).Convert(v.Type())
	case reflect.Float64:
		var n float64
		n, err = strconv.ParseFloat(s, 64)
		if err != nil {
			break
		}
		rv = reflect.ValueOf(n).Convert(v.Type())
	case reflect.String:
		rv = reflect.ValueOf(s)
	case reflect.Array:
		rv = reflect.Indirect(reflect.New(v.Type()))
		a := strings.Split(s, ",")
		for i, l := 0, v.Len(); i < l && i < len(a); i, l = i+1, l {
			if err = StringToInterface(strings.TrimSpace(a[i]), rv.Index(i).Addr().Interface()); err != nil {
				break
			}
		}
	case reflect.Slice:
		a := strings.Split(s, ",")
		rv := reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), len(a), len(a))
		for i := 0; i < len(a); i++ {
			if err = StringToInterface(a[i], rv.Index(i).Addr().Interface()); err != nil {
				break
			}
		}
	case reflect.Map:
		mt := reflect.MapOf(v.Type().Key(), v.Type().Elem())
		rv = reflect.MakeMap(mt)
		a := strings.Split(s, ",")
		for _, s := range a {
			pair := strings.Split(s, "=")
			if len(pair) != 2 {
				return ErrParse
			}
			key, err := StringToValue(pair[0], reflect.Zero(mt.Key()))
			if err != nil {
				return ErrParse
			}
			val, err := StringToValue(pair[1], reflect.Zero(mt.Elem()))
			if err != nil {
				return ErrParse
			}
			rv.SetMapIndex(key, val)
		}
	default:
		return ErrUnsupported
	}
	if rv.IsValid() {
		v.Set(rv)
	}
	return err
}
