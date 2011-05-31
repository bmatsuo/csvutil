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

// Write a slice of field values with a trailing field seperator and no '\n'.
// Returns any error incurred from writing.
func (csvw *Writer) WriteFields(fields []string) (int, os.Error) {
    var n int = len(fields)
    var success int = 0
    var err os.Error
    for i:=0 ; i<n ; i++ {
        nbytes, err := fmt.Fprintf(csvw.w,"%s%c",fields[i], csvw.Sep)
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
    var nbytes int
    var err os.Error
    if n > 1 {
        success, err = csvw.WriteFields(fields[:n-1])
        if err != nil {
            return success, err
        }
    }
    if n >= 1 {
        nbytes, err = fmt.Fprint(csvw.w, fields[n-1])
        success += nbytes
        if nbytes != len(fields[n-1]) {
            return success, err
        }
    }
    nbytes, err = fmt.Fprintf(csvw.w, "%c", '\n')
    success += nbytes
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
