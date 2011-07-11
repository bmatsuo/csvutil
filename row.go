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
    ErrorIndex         = os.NewError("Not enough fields to format")
    ErrorStruct        = os.NewError("Cannot format unreferenced structs")
    ErrorUnimplemented = os.NewError("Unimplemented field type")
    ErrorFieldType     = os.NewError("Field type incompatible.")
    ErrorNonPointer    = os.NewError("Target is not a pointer.")
    ErrorCantSet       = os.NewError("Cannot set value.")
)

func (r Row) formatReflectValue(i int, x reflect.Value) (int, os.Error) {
    if i >= len(r.Fields) {
        return 0, ErrorIndex
    }
    if !x.CanSet() {
        return 0, ErrorCantSet
    }
    var (
        assigned int
        errc     os.Error
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

func (r Row) formatValue(i int, x interface{}) (int, os.Error) {
    // TODO add complex types
    if i >= len(r.Fields) {
        return 0, ErrorIndex
    }
    var (
        assigned = 0
        errc     os.Error
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

//  *Unimplemented*. Should iteratively take values from the argument list
//  and assign to them successive fields from the row object. Should
//  return then number of row fields assigned to argument (fields) and any
//  error that occurred.
func (r Row) Format(x ...interface{}) (int, os.Error) {
    var (
        assigned int
        vasg     int
        err      os.Error
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

func formatReflectValue(x reflect.Value) (string, os.Error) {
    /*
       if !x.CanSet() {
           return "", ErrorCantSet
       }
    */
    var (
        errc os.Error
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
        return strconv.Itoa64(int64(x.Interface().(int8))), nil
    case reflect.Int16:
        return strconv.Itoa64(int64(x.Interface().(int16))), nil
    case reflect.Int32:
        return strconv.Itoa64(int64(x.Interface().(int32))), nil
    case reflect.Int64:
        return strconv.Itoa64(x.Interface().(int64)), nil
    case reflect.Uint:
        return strconv.Uitoa(x.Interface().(uint)), nil
    case reflect.Uint8:
        return strconv.Uitoa64(uint64(x.Interface().(uint8))), nil
    case reflect.Uint16:
        return strconv.Uitoa64(uint64(x.Interface().(uint16))), nil
    case reflect.Uint32:
        return strconv.Uitoa64(uint64(x.Interface().(uint32))), nil
    case reflect.Uint64:
        return strconv.Uitoa64(x.Interface().(uint64)), nil
    case reflect.Float32:
        return strconv.Ftoa32(x.Interface().(float32), FloatFmt, FloatPrec), nil
    case reflect.Float64:
        return strconv.Ftoa64(x.Interface().(float64), FloatFmt, FloatPrec), nil
    case reflect.Complex64:
        fallthrough
    case reflect.Complex128:
        errc = ErrorUnimplemented
    case reflect.Bool:
        return strconv.Btoa(x.Interface().(bool)), nil
    default:
        errc = ErrorFieldType
    }
    return "", errc
}

func formatValue(x interface{}) ([]string, os.Error) {
    // TODO add complex types
    var (
        formatted    = make([]string, 0, 1)
        appendwhenok = func(s string, e os.Error) os.Error {
            if e == nil {
                formatted = append(formatted, s)
            }
            return e
        }
        errc  os.Error
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

func FormatRow(x ...interface{}) Row {
    var (
        err          os.Error
        formatted    = make([]string, 0, len(x))
        appendwhenok = func(s []string, e os.Error) os.Error {
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
