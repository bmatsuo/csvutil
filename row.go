// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil

/* 
*  File: row.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Wed Jun  1 16:48:20 PDT 2011
*  Description: Row related types and methods.
 */
import (
	"errors"
	"io"
	//"fmt"
	//"log"
	"reflect"
	"strconv"
)

// See strconv.Ftoa32()
var (
	FloatFmt  byte = 'g'
	FloatPrec int  = -1
)

//  A simple row structure for rows read by a csvutil.Reader that
//  encapsulates any read error enountered along with any data read
//  prior to encountering an error.
type Row struct {
	Fields []string "CSV row field data"
	Error  error    "Error encountered reading"
}

//  A wrapper for the test r.Error == os.EOF
func (r Row) HasEOF() bool {
	return r.Error == io.EOF
}

//  A wrapper for the test r.Error != nil
func (r Row) HasError() bool {
	return r.Error != nil
}

var (
	ErrorIndex         = errors.New("Not enough fields to format")
	ErrorStruct        = errors.New("Cannot format unreferenced structs")
	ErrorUnimplemented = errors.New("Unimplemented field type")
	ErrorFieldType     = errors.New("Field type incompatible.")
	ErrorNonPointer    = errors.New("Target is not a pointer.")
	ErrorCantSet       = errors.New("Cannot set value.")
)

func (r Row) formatReflectValue(i int, x reflect.Value) (int, error) {
	if i >= len(r.Fields) {
		return 0, ErrorIndex
	}
	if !x.CanSet() {
		return 0, ErrorCantSet
	}
	var (
		assigned int
		errc     error
		kind     = x.Kind()
	)
	switch kind {
	// Format pointers to standard types.
	case reflect.String:
		x.SetString(r.Fields[i])
		assigned++
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		var vint int64
		vint, errc = strconv.ParseInt(r.Fields[i], 10, 64)
		if errc == nil {
			x.SetInt(vint)
			assigned++
		}
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		var vuint uint64
		vuint, errc = strconv.ParseUint(r.Fields[i], 10, 64)
		if errc == nil {
			x.SetUint(vuint)
			assigned++
		}
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		var vfloat float64
		vfloat, errc = strconv.ParseFloat(r.Fields[i], 64)
		if errc == nil {
			x.SetFloat(vfloat)
			assigned++
		}
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		errc = ErrorUnimplemented
	case reflect.Bool:
		var vbool bool
		vbool, errc = strconv.ParseBool(r.Fields[i])
		if errc == nil {
			x.SetBool(vbool)
			assigned++
		}
	default:
		errc = ErrorFieldType
	}
	return assigned, errc
}

func (r Row) formatValue(i int, x interface{}) (int, error) {
	// TODO add complex types
	if i >= len(r.Fields) {
		return 0, ErrorIndex
	}
	var (
		assigned = 0
		errc     error
		n        int
		value    = reflect.ValueOf(x)
	)
	if !value.IsValid() {
		return 0, ErrorFieldType
	}
	//var t = value.Type()
	var kind = value.Kind()
	switch kind {
	case reflect.Ptr:
		//log.Print("PtrType")
		break
	case reflect.Array:
		//log.Print("ArrayType")
		fallthrough
	case reflect.Slice:
		//log.Print("SliceType")
		n = value.Len()
		for j := 0; j < n; j++ {
			var vj = value.Index(j)
			rvasgn, rverr := r.formatReflectValue(i+j, vj)
			assigned += rvasgn
			if rverr != nil {
				return assigned, rverr
			}
		}
		return assigned, errc
	case reflect.Map:
		//log.Print("MapType")
		return 0, ErrorUnimplemented
	default:
		return 0, ErrorFieldType
	}
	var (
		eVal  = value.Elem()
		eType = eVal.Type()
		eKind = eType.Kind()
	)
	switch eKind {
	// Format pointers to standard types.
	case reflect.Struct:
		switch kind {
		case reflect.Ptr:
			n = eVal.NumField()
			for j := 0; j < n; j++ {
				var vj = eVal.Field(j)
				rvasgn, rverr := r.formatReflectValue(i+j, vj)
				assigned += rvasgn
				if rverr != nil {
					return assigned, rverr
				}
			}
		default:
			errc = ErrorStruct
		}
	case reflect.Array:
		//log.Print("ArrayType")
		fallthrough
	case reflect.Slice:
		//log.Print("SliceType")
		n = eVal.Len()
		for j := 0; j < n; j++ {
			var vj = eVal.Index(j)
			rvasgn, rverr := r.formatReflectValue(i+j, vj)
			assigned += rvasgn
			if rverr != nil {
				return assigned, rverr
			}
		}
		return assigned, errc
	case reflect.Map:
		//log.Print("MapType")
		return 0, ErrorUnimplemented
	default:
		assigned, errc = r.formatReflectValue(i, eVal)
	}
	return assigned, errc
}

//  Iteratively take values from the argument list and assigns to them
//  successive fields from the row object. Returns the number of row fields
//  assigned to arguments and any error that occurred.
func (r Row) Format(x ...interface{}) (int, error) {
	var (
		assigned int
		vasg     int
		err      error
	)
	for _, elm := range x {
		vasg, err = r.formatValue(assigned, elm)
		assigned += vasg
		if err != nil {
			return assigned, err
		}
	}
	return assigned, err
}

func formatReflectValue(x reflect.Value) (string, error) {
	/*
	   if !x.CanSet() {
	       return "", ErrorCantSet
	   }
	*/
	var (
		errc error
		kind = x.Kind()
		//vintstr   string
	)
	switch kind {
	// Format pointers to standard types.
	case reflect.String:
		return x.Interface().(string), nil
	case reflect.Int:
		return strconv.Itoa(x.Interface().(int)), nil
	case reflect.Int8:
		return strconv.FormatInt(int64(x.Interface().(int8)), 10), nil
	case reflect.Int16:
		return strconv.FormatInt(int64(x.Interface().(int16)), 10), nil
	case reflect.Int32:
		return strconv.FormatInt(int64(x.Interface().(int32)), 10), nil
	case reflect.Int64:
		return strconv.FormatInt(x.Interface().(int64), 10), nil
	case reflect.Uint:
		return strconv.FormatUint(uint64(x.Interface().(uint)), 10), nil
	case reflect.Uint8:
		return strconv.FormatUint(uint64(x.Interface().(uint8)), 10), nil
	case reflect.Uint16:
		return strconv.FormatUint(uint64(x.Interface().(uint16)), 10), nil
	case reflect.Uint32:
		return strconv.FormatUint(uint64(x.Interface().(uint32)), 10), nil
	case reflect.Uint64:
		return strconv.FormatUint(x.Interface().(uint64), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(float64(x.Interface().(float32)), FloatFmt, FloatPrec, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(x.Interface().(float64), FloatFmt, FloatPrec, 64), nil
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		errc = ErrorUnimplemented
	case reflect.Bool:
		return strconv.FormatBool(x.Interface().(bool)), nil
	default:
		errc = ErrorFieldType
	}
	return "", errc
}

func formatValue(x interface{}) ([]string, error) {
	// TODO add complex types
	var (
		formatted    = make([]string, 0, 1)
		appendwhenok = func(s string, e error) error {
			if e == nil {
				formatted = append(formatted, s)
			}
			return e
		}
		errc  error
		n     int
		value = reflect.ValueOf(x)
	)
	if !value.IsValid() {
		return formatted, ErrorFieldType
	}
	//var t = value.Type()
	var kind = value.Kind()
	switch kind {
	case reflect.Ptr:
		//log.Print("PtrType")
		break
	case reflect.Struct:
		n = value.NumField()
		for j := 0; j < n; j++ {
			var vj = value.Field(j)
			errc = appendwhenok(formatReflectValue(vj))
			if errc != nil {
				break
			}
		}
		return formatted, errc
	case reflect.Array:
		//log.Print("ArrayType")
		fallthrough
	case reflect.Slice:
		//log.Print("SliceType")
		n = value.Len()
		for j := 0; j < n; j++ {
			var vj = value.Index(j)
			errc = appendwhenok(formatReflectValue(vj))
			if errc != nil {
				break
			}
		}
		return formatted, errc
	case reflect.Map:
		//log.Print("MapType")
		return formatted, ErrorUnimplemented
	default:
		errc = appendwhenok(formatReflectValue(value))
		return formatted, errc
	}
	var (
		eVal  = value.Elem()
		eType = eVal.Type()
		eKind = eType.Kind()
	)
	switch eKind {
	// Format pointers to standard types.
	case reflect.Struct:
		switch kind {
		case reflect.Ptr:
			n = eVal.NumField()
			for j := 0; j < n; j++ {
				var vj = eVal.Field(j)
				errc = appendwhenok(formatReflectValue(vj))
				if errc != nil {
					break
				}
			}
			return formatted, errc
		default:
			errc = ErrorStruct
		}
	case reflect.Array:
		//log.Print("ArrayType")
		fallthrough
	case reflect.Slice:
		//log.Print("SliceType")
		n = eVal.Len()
		for j := 0; j < n; j++ {
			var vj = eVal.Index(j)
			errc = appendwhenok(formatReflectValue(vj))
			if errc != nil {
				break
			}
		}
		return formatted, errc
	case reflect.Map:
		//log.Print("MapType")
		return formatted, ErrorUnimplemented
	default:
		errc = appendwhenok(formatReflectValue(eVal))
	}
	return nil, errc
}

//  Iteratively take values from the argument list and formats them (or
//  their elements/fields) as a (list of) string(s). Returns a Row object
//  that contains the formatted arguments, as well as any error that
//  occured.
func FormatRow(x ...interface{}) Row {
	var (
		err          error
		formatted    = make([]string, 0, len(x))
		appendwhenok = func(s []string, e error) error {
			if e == nil {
				formatted = append(formatted, s...)
			}
			return e
		}
	)
	for _, elm := range x {
		err = appendwhenok(formatValue(elm))
		if err != nil {
			break
		}
	}
	return Row{formatted, err}
}
