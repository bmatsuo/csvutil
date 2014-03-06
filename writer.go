// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil

import (
	"bufio"
	"bytes"
	"io"
	"unicode/utf8"
	//"strings"
)

//  A simple CSV file writer using the package bufio for effeciency.
//  But, because of this, the method Flush() must be called to ensure
//  data is written to any given io.Writer before it is closed.
type Writer struct {
	*Config
	w  io.Writer     // Base writer object.
	bw *bufio.Writer // Buffering writer for efficiency.
}

//  Create a new CSV writer with the default field seperator and a
//  buffer of a default size.
func NewWriter(w io.Writer, c *Config) *Writer {
	csvw := new(Writer)
	if csvw.Config = c; c == nil {
		csvw.Config = NewConfig()
	}
	csvw.w = w
	csvw.bw = bufio.NewWriter(w)
	return csvw
}

//  Create a new CSV writer using a buffer of at least n bytes.
//
//      See bufio.NewWriterSize(io.Writer, int) (*bufio.NewWriter).
func NewWriterSize(w io.Writer, c *Config, n int) (*Writer, error) {
	csvw := new(Writer)
	if csvw.Config = c; c == nil {
		csvw.Config = NewConfig()
	}
	csvw.w = w
	csvw.bw = bufio.NewWriterSize(w, n)
	return csvw, nil
}

//  Write a slice of bytes to the data stream. No checking for containment
//  of the separator is done, so this file can be used to write multiple
//  fields if desired.
func (csvw *Writer) write(p []byte) (int, error) {
	return csvw.bw.Write(p)
}

//  Write a single field of CSV data. If the ln argument is true, a
//  trailing new line is printed after the field. Otherwise, when
//  the ln argument is false, a separator character is printed after
//  the field.
func (csvw *Writer) writeField(field string, ln bool) (int, error) {
	// Contains some code modified from
	//  $GOROOT/src/pkg/fmt/print.go: func (p *pp) fmtC(c int64) @ ~317,322
	var trail rune = csvw.Sep
	if ln {
		trail = '\n'
	}
	var (
		fLen = len(field)
		bp   = make([]byte, fLen+utf8.UTFMax)
	)
	copy(bp, field)
	return csvw.write(bp[:fLen+utf8.EncodeRune(bp[fLen:], trail)])
}

//  Write a slice of field values with a trailing field seperator (no '\n').
//  Returns any error incurred from writing.
func (csvw *Writer) WriteFields(fields ...string) (int, error) {
	var (
		n       = len(fields)
		success int
		err     error
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

//  Write a slice of field values with a trailing new line '\n'.
//  Returns any error incurred from writing.
func (csvw *Writer) WriteRow(fields ...string) (int, error) {
	var (
		n       = len(fields)
		success int
	)
	for i := 0; i < n; i++ {
		var EORow = i == n-1
		if nbytes, err := csvw.writeField(fields[i], EORow); err != nil {
			return success, err
		} else {
			success += nbytes
		}
	}
	return success, nil
}

//  Write a comment. Each comment string given will start on a new line. If
//  the string is contains multiple lines, comment prefixes will be
//  inserted at the beginning of each one.
func (csvw *Writer) WriteComments(comments ...string) (int, error) {
	if len(comments) == 0 {
		return 0, nil
	}

	// Break the comments into lines (w/o trailing '\n' chars)
	var lines [][]byte
	for _, c := range comments {
		var cp = make([]byte, len(c))
		copy(cp, c)
		lines = append(lines, bytes.Split(cp, []byte{'\n'})...)
	}

	// Count the total number of characters in the comments.
	var commentLen = len(lines) * (len(csvw.CommentPrefix) + 1)
	for _, cline := range lines {
		commentLen += len(cline)
	}

	// Allocate, fill, and write the comment byte slice
	var comment = make([]byte, commentLen)
	var ci int
	for _, cline := range lines {
		ci += copy(comment[ci:], csvw.CommentPrefix)
		ci += copy(comment[ci:], cline)
		ci += copy(comment[ci:], []byte{'\n'})
	}
	return csvw.write(comment[:ci])
}

//  Flush any buffered data to the underlying io.Writer.
func (csvw *Writer) Flush() error {
	return csvw.bw.Flush()
}

//  Write multple CSV rows at once.
func (csvw *Writer) WriteRows(rows [][]string) (int, error) {
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
