package csvutil
/* 
*  File: writer.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat May 28 23:53:36 PDT 2011
*  Description: CSV writer library.
*/
import (
    "os"
    "io"
    "utf8"
    "fmt"
    "bytes"
    "strings"
)

// A simple CSV file writer.
type Writer struct {
    Sep int "CSV Field seperator."
    w io.Writer "Base writer object."
}

//  Create a new CSV writer with the default field seperator
func NewWriter(w io.Writer) *Writer {
    csvw := new(Writer)
    csvw.Sep = DEFAULT_SEP
    csvw.w = w
    return csvw
}

//  Write a slice of bytes to the data stream. No checking for containment
//  of the separator is done, so this file can be used to write multiple
//  fields if desired.
func (csvw *Writer) Write(p []byte) (nbytes int, err os.Error) {
    nbytes, err = csvw.w.Write(p)
    return
}

//  Write a string to the data stream. No checking for containment of
//  the separator is done, so this file can be used to write multiple
//  fields if desired. 
func (csvw *Writer) WriteString(str string) (nbytes int, err os.Error) {
    return csvw.Write(bytes.NewBufferString(str).Bytes())
}

//  Attempt to write a string to underlying io.Writer, but panic if a
//  separator character found in the stream.
func (csvw *Writer) WriteStringSafe(str string) (nbytes int, err os.Error) {
    var sep string = fmt.Sprintf("%c", csvw.Sep)
    var i int = strings.Index(str, sep)
    if i != -1 {
        panic("sepfound")
    }
    return csvw.WriteString(str)
}

//  Write a single field of CSV data. If the ln argument is true, a
//  trailing new line is printed after the field. Otherwise, when
//  the ln argument is false, a separator character is printed after
//  the field.
func (csvw *Writer) WriteField(field string, ln bool) (nbytes int, err os.Error) {
    var trailChar int
    if ln {
        trailChar = '\n'
    } else {
        trailChar = csvw.Sep
    }
    return fmt.Fprintf(csvw.w,"%s%c",field, trailChar)
}

//  Write a slice of field values with a trailing field seperator and no '\n'.
//  Returns any error incurred from writing.
func (csvw *Writer) WriteFields(fields []string) (int, os.Error) {
    var n int = len(fields)
    var success int = 0
    var err os.Error
    for i:=0 ; i<n ; i++ {
        nbytes, err := csvw.WriteField(fields[i], false)
        success += nbytes
        if nbytes < len(fields[i]) + utf8.RuneLen(csvw.Sep) {
            return success, err
        }
    }
    return success, err
}

// Write a slice of field values with a trailing new line '\n'.
// Returns any error incurred from writing.
func (csvw *Writer) WriteFieldsln(fields []string) (int, os.Error) {
    var n int = len(fields)
    var success int = 0
    var err os.Error
    for i:=0 ; i<n ; i++ {
        var onLastField bool = i == n-1
        nbytes, err := csvw.WriteField(fields[i], onLastField)
        success += nbytes

        var trail int = csvw.Sep
        if onLastField {
            trail = '\n'
        }

        if nbytes < len(fields[i]) + utf8.RuneLen(trail) {
            return success, err
        }
    }
    return success, err
}

// Write multple CSV rows at once.
func (csvw *Writer) WriteRows(rows [][]string) (int, os.Error) {
    var success, nbytes int
    var err os.Error
    success = 0
    for i:=0 ; i<len(rows) ; i++ {
        nbytes, err = csvw.WriteFieldsln(rows[i])
        success += nbytes
        if err != nil {
            return success, err
        }
    }
    return success, err
}
