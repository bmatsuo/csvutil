// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil
import (
    "os"
    "io"
    "bufio"
    "utf8"
    //"strings"
)

//  A simple CSV file writer using the package bufio for effeciency.
//  But, because of this, the method Flush() must be called to ensure
//  data is written to any given io.Writer before it is closed.
type Writer struct {
    Sep int           // CSV field seperator.
    w   io.Writer     // Base writer object.
    bw  *bufio.Writer // Buffering writer for efficiency.
}

//  Create a new CSV writer with the default field seperator and a
//  buffer of a default size.
func NewWriter(w io.Writer) *Writer {
    csvw := new(Writer)
    csvw.Sep = DefaultSep
    csvw.w = w
    csvw.bw = bufio.NewWriter(w)
    return csvw
}

//  Create a new CSV writer using a buffer of at least n bytes.
//
//      See bufio.NewWriterSize(io.Writer, int) (*bufio.NewWriter).
func NewWriterSize(w io.Writer, n int) (*Writer, os.Error) {
    csvw := new(Writer)
    csvw.Sep = DefaultSep
    csvw.w = w
    var bufErr os.Error
    csvw.bw, bufErr = bufio.NewWriterSize(w, n)
    return csvw, bufErr
}

//  Write a slice of bytes to the data stream. No checking for containment
//  of the separator is done, so this file can be used to write multiple
//  fields if desired.
func (csvw *Writer) write(p []byte) (int, os.Error) {
    return csvw.bw.Write(p)
}
/*
//  Write a string to the data stream. No checking for containment of
//  the separator is done, so this file can be used to write multiple
//  fields if desired. 
func (csvw *Writer) WriteString(str string) (nbytes int, err os.Error) {
    var b []byte = make([]byte, len(str))
    copy(b, str)
    return csvw.write(b)
}

//  Attempt to write a string to underlying io.Writer, but panic if a
//  separator character found in the stream.
func (csvw *Writer) WriteStringSafe(str string) (nbytes int, err os.Error) {
    // Some code modified from
    //  $GOROOT/src/pkg/fmt/print.go: func (p *pp) fmtC(c int64) @ ~317,322
    var rb []byte = make([]byte, utf8.UTFMax)
    //Do I want int64 separators?
    //rune := int(c) 
    //if int64(rune) != c { rune = utf8.RuneError }
    w := utf8.EncodeRune(rb, csvw.Sep)
    var sep string = string(rb[0:w])
    var i int = strings.Index(str, sep)
    if i != -1 {
        panic("sepfound")
    }
    return csvw.WriteString(str)
}
*/

//  Write a single field of CSV data. If the ln argument is true, a
//  trailing new line is printed after the field. Otherwise, when
//  the ln argument is false, a separator character is printed after
//  the field.
func (csvw *Writer) writeField(field string, ln bool) (int, os.Error) {
    // Contains some code modified from
    //  $GOROOT/src/pkg/fmt/print.go: func (p *pp) fmtC(c int64) @ ~317,322
    var trail int = csvw.Sep
    if ln {
        trail = '\n'
    }
    var (
        fLen  = len(field)
        bp    = make([]byte, fLen+utf8.UTFMax)
    )
    copy(bp, field)
    return csvw.write(bp[:fLen+utf8.EncodeRune(bp[fLen:], trail)])
}

//  Write a slice of field values with a trailing field seperator (no '\n').
//  Returns any error incurred from writing.
func (csvw *Writer) WriteFields(fields...string) (int, os.Error) {
    var(
        n       = len(fields)
        success int
        err     os.Error
    )
    for i := 0; i < n; i++ {
        if nbytes, err := csvw.writeField(fields[i], false); err != nil {
            return success, err
        } else {
            success += nbytes
        }
    }
    return success, err
}

/*
func (csvw *Writer) WriteFieldsln(fields...string) (int, os.Error) {
    var n int = len(fields)
    var success int = 0
    var err os.Error
    for i := 0; i < n; i++ {
        var onLastField bool = i == n-1
        nbytes, err := csvw.writeField(fields[i], onLastField)
        success += nbytes

        var trail int = csvw.Sep
        if onLastField {
            trail = '\n'
        }

        if nbytes < len(fields[i])+utf8.RuneLen(trail) {
            return success, err
        }
    }
    return success, err
}
*/

// Write a slice of field values with a trailing new line '\n'.
// Returns any error incurred from writing.
func (csvw *Writer) WriteRow(fields...string) (int, os.Error) {
    var (
        n       = len(fields)
        success int
    )
    for i := 0; i < n; i++ {
        var onLastField bool = i == n-1
        if nbytes, err := csvw.writeField(fields[i], onLastField); err != nil {
            return success, err
        } else {
            success += nbytes
        }
    }
    return success, nil
}

//  Flush any buffered data to the underlying io.Writer.
func (csvw *Writer) Flush() os.Error {
    return csvw.bw.Flush()
}

//  Write multple CSV rows at once.
func (csvw *Writer) WriteRows(rows [][]string) (int, os.Error) {
    var success int
    for i := 0; i < len(rows); i++ {
        if nbytes, err := csvw.WriteRow(rows[i]...); err != nil {
            return success, err
        } else {
            success += nbytes
        }
    }
    return success, nil
}
