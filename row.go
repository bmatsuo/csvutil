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
    "os"
    //"fmt"
    //"log"
    "strconv"
    "reflect"
)

//  A simple row structure for rows read by a csvutil.Reader that
//  encapsulates any read error enountered along with any data read
//  prior to encountering an error.
type Row struct {
    Fields []string "CSV row field data"
    Error  os.Error "Error encountered reading"
}

//  A wrapper for the test r.Error == os.EOF
func (r Row) HasEOF() bool {
    return r.Error == os.EOF
}

//  A wrapper for the test r.Error != nil
func (r Row) HasError() bool {
    return r.Error != nil
}

var (
    ErrorIndex          = os.NewError("Not enough fields to format")
    ErrorStruct         = os.NewError("Cannot format unreferenced structs")
    ErrorUnimplemented  = os.NewError("Unimplemented field type")
    ErrorFieldType      = os.NewError("Field type incompatible.")
    ErrorNonPointer     = os.NewError("Target is not a pointer.")
    ErrorCantSet        = os.NewError("Cannot set value.")
)

func (r Row)formatReflectValue(i int, x reflect.Value) (int, os.Error) {
    if i >= len(r.Fields) {
        return 0, ErrorIndex
    }
    if !x.CanSet() {
        return 0, ErrorCantSet
    }
    var (
        assigned int
        errc     os.Error
        kind = x.Kind()
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
        vint, errc = strconv.Atoi64(r.Fields[i])
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
        vuint, errc = strconv.Atoui64(r.Fields[i])
        if errc == nil {
            x.SetUint(vuint)
            assigned++
        }
    case reflect.Float32:
        fallthrough
    case reflect.Float64:
        var vfloat float64
        vfloat, errc = strconv.Atof64(r.Fields[i])
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
        vbool, errc = strconv.Atob(r.Fields[i])
        if errc == nil {
            x.SetBool(vbool)
            assigned++
        }
    default:
        errc = ErrorFieldType
    }
    return assigned, errc
}

func (r Row)formatValue(i int, x interface{}) (int, os.Error) {
    // TODO add complex types
    if i >= len(r.Fields) {
        return 0, ErrorIndex
    }
    var (
        assigned = 0
        errc os.Error
        n int
        value = reflect.ValueOf(x)
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
        for j := 0 ; j < n ; j++ {
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
    var elemValue = value.Elem()
    var elemType = elemValue.Type()
    var elemKind = elemType.Kind()
    switch elemKind {
    // Format pointers to standard types.
    case reflect.Struct:
        switch kind {
        case reflect.Ptr:
            n = elemValue.NumField()
            for j := 0 ; j < n ; j++ {
                var vj = elemValue.Field(j)
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
        n = elemValue.Len()
        for j := 0 ; j < n ; j++ {
            var vj = elemValue.Index(j)
            rvasgn, rverr := r.formatReflectValue(i+j, vj)
            assigned += rvasgn
            if rverr != nil {
                return assigned, rverr
            }
        }
        return assigned, errc
    default:
        assigned, errc = r.formatReflectValue(i, elemValue)
    }
    return assigned, errc
}

//  *Unimplemented*. Should iteratively take values from the argument list
//  and assign to them successive fields from the row object. Should
//  return then number of row fields assigned to argument (fields) and any
//  error that occurred.
func (r Row) Format(x...interface{}) (int, os.Error) {
    var (
        assigned int
        vasg     int
        err os.Error
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
