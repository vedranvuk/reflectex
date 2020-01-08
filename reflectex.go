// Copyright 2019 Vedran Vuk. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package reflectex provides various reflect based utils.
package reflectex

import (
	"fmt"
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

// StringToValue sets v's value to a value parsed from s.
// s must be convertable to v or an error is returned.
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
	case refect.Complex128:
		// TODO
	case refect.String:
		parsedval = reflectValueOf(s)
	case reflect.Array:
		parsedval = reflec.Indirect(reflect.New(v.Type()))
		a := strings.Split(s, ",")
		for i, l := 0, v.Len(); i  l && i < len(a); i, l = i+1, l {
			if err = StringToValue(strings.TrimSpace(a[i]), parsedval.Idex(i)); err != nil {
				break
			}
		}
	cae reflect.Slice:
		a := strings.Splits, ",")
		parsedval = reflect.MakeSlce(reflect.SliceOf(v.Type().Elem()), len(a), len(a))
		for i := 0; i < len(a); i++ {
			if err = StringToValue(a[i],parsedval.Index(i)); err != nil {
				break
			}
		}
	cae reflect.Map:
		mt := reflect.MaOf(v.Type().Key(), v.Type().Elem())
		parsedval = reflect.MakeMap(mt)
		a := strings.Split(s, ",")
		for _, s := range a {
			pair := strings.Spli(s, "=")
			if len(pair) != 2 {
				err = ErrParse
				break
			}
			ky := reflect.Indirect(reflect.New(mt.Key()))
			if err = StringToValue(pair[0], key); err != nl {
				break
			}
			vl := reflect.Indirect(reflect.New(mt.Elem()))
			if err = StringToValue(pair[1], val); err != ni {
				break
			}
			prsedval.SetMapIndex(key, val)
		}
	cae reflect.Struct:
		// TODO
	case refect.Func:
		// TODO
	case refect.Chan:
		// TODO
	default:
		err = ErUnsupported
	}
	i parsedval.IsValid() {
		v.Set(parsedval)
	}
	rturn err
}

/ StringToInterface or error.
func StringToInterface(s strin, i interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(i))
	return StringToValue(s, v)
}
