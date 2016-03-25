// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil

import (
	"bufio"
	"io"
	"strings"
)

//  readerBufferMinimumSize is the smallest csvutil will allow the
//  Reader's internal "long-line buffer" to be allocated as.
const readerBufferMinimumSize = 30

//  A reader object for CSV data utilizing the bufio package.
type Reader struct {
	*Config
	r          io.Reader     // Base reader object.
	br         *bufio.Reader // Buffering for efficiency and line reading.
	p          []byte        // A buffer for longer lines
	pi         int           // An index into the p buffer.
	lineNum    int
	pastHeader bool
}

//  Create a new reader object.
func NewReader(r io.Reader, c *Config) *Reader {
	var csvr *Reader = new(Reader)
	if csvr.Config = c; c == nil {
		csvr.Config = NewConfig()
	}
	csvr.r = r
	csvr.br = bufio.NewReader(r)
	return csvr
}

//  Create a new reader with a buffer of a specified size.
func NewReaderSize(r io.Reader, c *Config, size int) *Reader {
	var csvr *Reader = new(Reader)
	if csvr.Config = c; c == nil {
		csvr.Config = NewConfig()
	}
	csvr.r = r
	br := bufio.NewReaderSize(r, size)
	csvr.br = br
	return csvr
}

func (csvr *Reader) readLine() (string, error) {
	var (
		isPrefix = true
		piece    []byte
		err      error
	)
	for isPrefix {
		piece, isPrefix, err = csvr.br.ReadLine()
		switch err {
		case nil:
			break
		case io.EOF:
			fallthrough
		default:
			return "", err
		}
		var (
			readLen = len(piece)
			necLen  = csvr.pi + readLen
			pLen    = len(csvr.p)
		)
		if pLen == 0 {
			if pLen = readerBufferMinimumSize; pLen < necLen {
				pLen = necLen
			}
			csvr.p = make([]byte, pLen)
			csvr.pi = 0
		} else if pLen < necLen {
			if pLen = 2 * pLen; pLen < necLen {
				pLen = necLen
			}
			csvr.p = append(csvr.p, make([]byte, pLen)...)
		}
		csvr.pi += copy(csvr.p[csvr.pi:], piece)
	}
	var s = string(csvr.p[:csvr.pi])
	for i := 0; i < csvr.pi; i++ {
		csvr.p[i] = 0
	}
	csvr.pi = 0
	return s, nil
}

//  Returns the number of lines of input read by the Reader
func (csvr *Reader) LineNum() int {
	return csvr.lineNum
}

//  Attempt to read up to a new line, skipping any comment lines found in
//  the process. Return a Row object containing the fields read and any
//  error encountered.
func (csvr *Reader) ReadRow() Row {
	var (
		r    Row
		line string
	)
	// Read lines until a non-comment line is found.
	for true {
		if line, r.Error = csvr.readLine(); r.Error != nil {
			return r
		}
		csvr.lineNum++
		if !csvr.Comments {
			break
		} else if !csvr.LooksLikeComment(line) {
			break
		} else if csvr.pastHeader && !csvr.CommentsInBody {
			break
		}
	}
	csvr.pastHeader = true

	// Break the line up into fields.
	r.Fields = strings.Split(line, string(csvr.Sep))

	// Trim any unwanted characters.
	if csvr.Trim {
		for i := 0; i < len(r.Fields); i++ {
			r.Fields[i] = strings.Trim(r.Fields[i], csvr.Cutset)
		}
	}
	return r
}

//  Read rows into a preallocated buffer. Return the number of rows read,
//  and any error encountered.
func (csvr *Reader) ReadRows(rbuf [][]string) (int, error) {
	var (
		i   int
		err error
	)
	csvr.DoN(len(rbuf), func(r Row) bool {
		err = r.Error
		if r.Fields != nil {
			rbuf[i] = r.Fields
			i++
		}
		return !r.HasError()
	})
	return i, err
}

//  Reads any remaining rows of CSV data in the underlying io.Reader.
func (csvr *Reader) RemainingRows() (rows [][]string, err error) {
	return csvr.RemainingRowsSize(16)
}

//  Like csvr.RemainingRows(), but allows specification of the initial
//  row buffer capacity to avoid unnecessary reallocations.
func (csvr *Reader) RemainingRowsSize(size int) ([][]string, error) {
	var (
		err  error
		rbuf = make([][]string, 0, size)
	)
	csvr.Do(func(r Row) bool {
		err = r.Error
		//log.Printf("Scanned %v", r)
		if r.Fields != nil {
			rbuf = append(rbuf, r.Fields)
		}
		return !r.HasError()
	})
	return rbuf, err
}

//  Iteratively read the remaining rows in the reader and call f on each
//  of them. If f returns false, no more rows will be read.
func (csvr *Reader) Do(f func(Row) bool) {
	for r := csvr.ReadRow(); true; r = csvr.ReadRow() {
		if r.HasEOF() {
			//log.Printf("EOF")
			break
		}
		if !f(r) {
			//log.Printf("Break")
			break
		}
	}
}

//  Process rows from the reader like Do, but stop after processing n of
//  them. If f returns false before n rows have been process, no more rows
//  will be processed.
func (csvr *Reader) DoN(n int, f func(Row) bool) {
	var i int
	csvr.Do(func(r Row) bool {
		if i < n {
			return f(r)
		}
		return false
	})
}
