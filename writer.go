package csvutil
/* 
*  File: writer.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat May 28 23:53:36 PDT 2011
*  Description: CSV writer library.
*
*   This file is part of csvutil.
*
*   csvutil is free software: you can redistribute it and/or modify
*   it under the terms of the GNU Lesser Public License as published by
*   the Free Software Foundation, either version 3 of the License, or
*   (at your option) any later version.
*
*   csvutil is distributed in the hope that it will be useful,
*   but WITHOUT ANY WARRANTY; without even the implied warranty of
*   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*   GNU Lesser Public License for more details.
*
*   You should have received a copy of the GNU Lesser Public License
*   along with csvutil.  If not, see <http://www.gnu.org/licenses/>.
 */
import (
    "os"
    "io"
    "bufio"
    "utf8"
    "strings"
)

//  A simple CSV file writer using the package bufio for effeciency.
//  But, because of this, the method Flush() must be called to ensure
//  data is written to any given io.Writer before it is closed.
type Writer struct {
    Sep int           "CSV Field seperator."
    w   io.Writer     "Base writer object."
    bw  *bufio.Writer "Base writer object."
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
func (csvw *Writer) write(p []byte) (nbytes int, err os.Error) {
    nbytes, err = csvw.bw.Write(p)
    return
}

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
    /* Do I want int64 separators?
       rune := int(c) 
       if int64(rune) != c { rune = utf8.RuneError }
    */
    w := utf8.EncodeRune(rb, csvw.Sep)
    var sep string = string(rb[0:w])
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
func (csvw *Writer) writeField(field string, ln bool) (nbytes int, err os.Error) {
    var trailChar int
    if ln {
        trailChar = '\n'
    } else {
        trailChar = csvw.Sep
    }
    // Some code modified from
    //  $GOROOT/src/pkg/fmt/print.go: func (p *pp) fmtC(c int64) @ ~317,322
    var rb []byte = make([]byte, utf8.UTFMax) // A utf8 rune buffer.
    /* Do I want int64 separators? 
        rune := int(c) 
       	if int64(rune) != c { rune = utf8.RuneError }
    */
    rbLen := utf8.EncodeRune(rb, trailChar)
    var fLen int = len(field)
    var bp []byte = make([]byte, fLen, fLen+rbLen)
    copy(bp, field)
    return csvw.write(append(bp, rb[0:rbLen]...))
}

//  Write a slice of field values with a trailing field seperator (no '\n').
//  Returns any error incurred from writing.
func (csvw *Writer) WriteFields(fields...string) (int, os.Error) {
    var n int = len(fields)
    var success int = 0
    var err os.Error
    for i := 0; i < n; i++ {
        nbytes, err := csvw.writeField(fields[i], false)
        success += nbytes
        if nbytes < len(fields[i])+utf8.RuneLen(csvw.Sep) {
            return success, err
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

//  Flush any buffered data to the underlying io.Writer.
func (csvw *Writer) Flush() os.Error {
    return csvw.bw.Flush()
}

//  Write multple CSV rows at once.
func (csvw *Writer) WriteRows(rows [][]string) (int, os.Error) {
    var success, nbytes int
    var err os.Error
    success = 0
    for i := 0; i < len(rows); i++ {
        nbytes, err = csvw.WriteRow(rows[i]...)
        success += nbytes
        if err != nil {
            return success, err
        }
    }
    return success, nil
}
