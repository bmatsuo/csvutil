package csvutil
/* 
*  File: reader.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat May 28 23:53:36 PDT 2011
*  Description: CSV reader library.
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
	"io"
	"bufio"
	"strings"
	"os"
)

//  A reader object for CSV data utilizing the bufio package.
type Reader struct {
	Sep int "Field separator character."
	Trim bool "Remove excess whitespace from field values."
	Cutset string "Set of characters to trim."
    LastRow Row "The last row read by the Reader."
	r io.Reader "Base reader object."
	br *bufio.Reader "For reading lines."
    bufSize int "Initial size of internal line buffers."
}

//  A simple row structure for rows read by a csvutil.Reader that
//  encapsulates any read error enountered along with any data read
//  prior to encountering an error.
type Row struct {
	Fields []string "CSV row field data"
	Error os.Error "Error encountered reading"
}

func mkrow(fields []string, err os.Error) Row {
	var r Row
	r.Fields = fields
	r.Error = err
	return r
}

//  An object for iterating over the rows of a Reader.
//      rit = reader.RowIterAuto()
//      for r := range reader.RowsChan {
//          if r.Error != nil {
//              panic(r.Error)
//          }
//          var fields []string = r.Fields
//          // Process the desired entries of "fields".
//          rit.Next()
//      }
//  The (*Reader) RowIterAuto() (*RowReaderIterator) method creates a
//  new iterator object and immediately reads a row from the Reader.
//  If this behavior is not desired, use the underlying alternate method
//  (*Reader) RowIter() (*RowReaderIterator).
type ReaderRowIterator struct {
    stopped bool
    RowsChan <-chan Row
    control chan<- bool
}

//  A ReaderRowIterator meant for reading until the end of the
//  data stream of the corresponding Reader's io.Reader. 
//      rit = reader.RowIterAuto()
//      for r := range reader.RowsChan {
//          if r.Error != nil {
//              panic(r.Error)
//          }
//          var fields []string = r.Fields
//          // Process the desired entries of "fields".
//      }
//  This iteration using a range statement as above, it is not safe to
//  break out of the loop. This generally causes the 'loss' of at least
//  one row.
//
//  For iterating rows in a way such that the iteration can be stopped
//  safely, use ReaderRowIterator objects instead.
type ReaderRowIteratorAuto struct {
    stopped bool
    RowsChan <-chan Row
}

//  Tell the iterator to get another row from the Reader.
func (csvri *ReaderRowIterator) Next() {
    if csvri.stopped {
        panic("stopped")
    }
    csvri.control <- true
}

//  Tell the iterator to stop fetching rows and exit its goroutine.
//  Calling the (*ReaderRowIterator) Break() method is not necessary,
//  but to avoid doing so will cause the iterating goroutine to sleep
//  for the duration of the program.
func (csvri *ReaderRowIterator) Break() {
    if csvri.stopped {
        return
    }
    close(csvri.control)
    csvri.stopped = true
}

//  Create a new reader object.
func NewReader(r io.Reader) *Reader {
	var csvr *Reader = new(Reader).Reset()
	csvr.r = r
	csvr.br = bufio.NewReader(r)
    csvr.bufSize = 80
	return csvr
}

//  Create a new reader with a buffer of a specified size.
func NewReaderSize(r io.Reader, size int) *Reader {
	var csvr *Reader = new(Reader).Reset()
	csvr.r = r
	var err os.Error
	csvr.br, err = bufio.NewReaderSize(r, size)
	if err != nil { panic(err) }
    csvr.bufSize = size
	return csvr
}

// Create a new row iterator and return it.
func (csvr *Reader) RowIter() (*ReaderRowIterator) {
    ri := new(ReaderRowIterator)
    throughChan := make(chan Row)
    controlChan := make (chan bool)
    ri.RowsChan = throughChan
    ri.control = controlChan
	var read_rows = func (r chan<- Row, c <-chan bool) {
        defer func() {
            if x:=recover(); x!=nil {
                /* Do nothing. */
            }
        } ()
		for true {
            cont, ok := <-c
            if !ok || !cont {
                break
            }
            csvr.LastRow = Row{Fields:nil, Error:nil}
			var row Row = csvr.ReadRow()
			if row.Fields == nil {
				if row.Error == os.EOF {
					break
				} else {
					panic(row.Error)
				}
			}
            csvr.LastRow = row
			r <- row
		}
		close(r)
	}
	go read_rows(throughChan, controlChan)
	return ri
}
//  For convenience, return a new ReaderRowIterator rit that has
//  already already been the target of (*ReaderRowIterator) Next().
func (csvr *Reader) RowIterStarted() (rit *ReaderRowIterator) {
    rit = csvr.RowIter()
    rit.Next()
    return rit
}

// Create a new ReaderRowIteratorAuto object and return it.
func (csvr *Reader) RowIterAuto() (*ReaderRowIteratorAuto) {
    ri := new(ReaderRowIteratorAuto)
    throughChan := make(chan Row)
    ri.RowsChan = throughChan
	var read_rows = func (r chan<- Row) {
        defer func() {
            if x:=recover(); x!=nil {
                /* Do nothing. */
            }
        } ()
		for true {
            csvr.LastRow = Row{Fields:nil, Error:nil}
			var row Row = csvr.ReadRow()
			if row.Fields == nil {
				if row.Error == os.EOF {
					break
				} else {
					panic(row.Error)
				}
			}
            csvr.LastRow = row
			r <- row
		}
		close(r)
	}
	go read_rows(throughChan)
	return ri
}


//  Read up to a new line and return a slice of string slices
func (csvr *Reader) ReadRow() Row {
	/* Read one row of the CSV and and return an array of the fields. */
    var(
        r Row
        line, readLine, ln []byte
        err os.Error
        isPrefix bool
    )
	r = Row{Fields:nil, Error:nil}
    csvr.LastRow = r
    isPrefix = true
    for isPrefix {
        readLine, isPrefix, err = csvr.br.ReadLine()
	    r.Fields = nil
	    r.Error = err
	    if err == os.EOF { return r }
	    if err != nil { return r }
	    if isPrefix {
            readLen := len(readLine)
            if line == nil {
                line = make([]byte, 0, 2*readLen)
            }
            ln = make([]byte, readLen)
            copy(ln, readLine)
            line = append(line, ln...)
        } else {
            // isPrefix is false here. The loop will break next iteration.
            if line == nil {
                line = readLine
            }
        }
    }
	r.Fields = strings.FieldsFunc(
			string(line),
			func (c int) bool { return c == csvr.Sep } )
	if csvr.Trim {
		for i:=0 ; i<len(r.Fields) ; i++ {
			r.Fields[i] = strings.Trim(r.Fields[i], csvr.Cutset)
		}
	}
    csvr.LastRow = r
	return r
}


//  Read rows into a preallocated buffer. Any error encountered is
//  returned. Returns the number of rows read in a single value return
//  context any errors encountered (including os.EOF).
func (csvr *Reader) ReadRows(rbuf [][]string) (int, os.Error) {
    var(
        err os.Error
        numRead int = 0
        n int = len(rbuf)
    )
    for i:=0 ; i<n ; i++ {
        r := csvr.ReadRow()
        numRead++
        if r.Error != nil {
            err = r.Error
            if r.Fields != nil {
                rbuf[i] = r.Fields
            }
            break
        }
        rbuf[i] = r.Fields
    }
    return numRead, err
}

//  Convenience methor to read at most n rows from csvr. Simple allocates
//  a row slice rs and calls csvr.ReadRows(rs). Returns the actual number
//  of rows read and any error that occurred (and halted reading).
func (csvr *Reader) ReadNRows(n int) (int, os.Error) {
    rows := make([][]string, n)
    return csvr.ReadRows(rows)
}

//  Reads any remaining rows of CSV data in the underlying io.Reader.
//  Uses resizing when a preallocated buffer of rows fills. Up to 16
//  rows can be read without any doubling occuring.
func (csvr *Reader) RemainingRows() (rows [][]string, err os.Error) {
    return csvr.RemainingRowsSize(16)
}

//  Like csvr.RemainingRows(), but allows specification of the initial
//  row buffer capacity.
func (csvr *Reader) RemainingRowsSize(size int) (rows [][]string, err os.Error) {
    err = nil
    var rbuf [][]string = make([][]string, 0, size)
    rit := csvr.RowIterAuto()
    for r := range rit.RowsChan {
        /*
        if cap(rbuf) == len(rbuf) {
            newbuf := make([][]string, len(rbuf), 2*len(rbuf))
            copy(rbuf,newbuf)
            rbuf = newbuf
        }
        */
        if r.Error != nil {
            err = r.Error
            if r.Fields != nil {
                rbuf = append(rbuf, r.Fields)
            }
            break
        }
        rbuf = append(rbuf, r.Fields)
    }
    return rbuf, err
}

/*
func (csvr *Reader) EachRow() <-chan Row {
	// Generator function for iterating through rows.
    c := make(chan Row)
	var read_rows = func (c chan<- Row) {
        defer func() {
            if x:=recover(); x!=nil {
                // Do nothing.
            }
        } ()
		for true {
            csvr.LastRow = Row{Fields:nil, Error:nil}
			var r Row = csvr.ReadRow()
			if r.Fields == nil {
				if r.Error == os.EOF {
					break
				} else {
					panic(r.Error)
				}
			}
            csvr.LastRow = r
			c <- r
		}
		close(c)
	}
	go read_rows(c)
	return c
}
*/

//  A function routine for setting all the configuration variables of a
//  csvutil.Reader in a single line.
func (csvr *Reader) Configure(sep int, trim bool, cutset string) *Reader {
	csvr.Sep = sep
	csvr.Trim = trim
	csvr.Cutset = cutset
	return csvr
}

//  Reset a Reader's configuration to the defaults. This is mostly meant
//  for internal use but is safe for general use.
func (csvr *Reader) Reset() *Reader {
	return csvr.Configure(DEFAULT_SEP, DEFAULT_TRIM, DEFAULT_CUTSET)
}

/* Comply with the reader interface. */
/*
func (csvr*Reader) Read(b []byte) (n int, err os.Error) {
	return csvr.r.Read(b)
}
*/
