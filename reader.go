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
    //"fmt"
    //"log"
)

//  readerBufferMinimumSize is the smallest csvutil will allow the
//  Reader's internal "long-line buffer" to be allocated as.
const readerBufferMinimumSize = 30

//  A reader object for CSV data utilizing the bufio package.
type Reader struct {
    Sep     int           // Field separator character.
    Trim    bool          // Remove excess whitespace from field values.
    Cutset  string        // Set of characters to trim.
    r       io.Reader     // Base reader object.
    br      *bufio.Reader // For reading lines.
    p       []byte        // A buffer for longer lines
    pi      int           // An index into the p buffer.
}

//  Create a new reader object.
func NewReader(r io.Reader) *Reader {
    var csvr *Reader = new(Reader).Reset()
    csvr.r = r
    csvr.br = bufio.NewReader(r)
    return csvr
}

//  Create a new reader with a buffer of a specified size.
func NewReaderSize(r io.Reader, size int) *Reader {
    var csvr *Reader = new(Reader).Reset()
    csvr.r = r
    var err os.Error
    csvr.br, err = bufio.NewReaderSize(r, size)
    if err != nil {
        panic(err)
    }
    return csvr
}

func (csvr *Reader) readLine() (string, os.Error) {
    var (
        isPrefix = true
        piece    []byte
        err      os.Error
    )
    for isPrefix {
        piece, isPrefix, err = csvr.br.ReadLine()
        switch err {
        case nil:
            break
        case os.EOF:
            fallthrough
        default:
            return "", err
        }
        var (
            readLen = len(piece)
            necLen = csvr.pi+readLen
            pLen = len(csvr.p)
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
            var p = make([]byte, pLen)
            copy(p, csvr.p[:csvr.pi])
        }
        csvr.pi += copy(csvr.p[csvr.pi:], piece)
    }
    //var p = make([]byte, csvr.pi)
    //copy(p, csvr.p[:csvr.pi])
    //return p, nil
    var s = string(csvr.p[:csvr.pi])
    for i := 0; i < csvr.pi; i++ {
        csvr.p[i] = 0
    }
    csvr.pi = 0
    return s, nil
}

//  Attempt to read up to a new line. Return a Row object containing
//  the fields read and any error encountered.
func (csvr *Reader) ReadRow() Row {
    /* Read one row of the CSV and and return an array of the fields. */
    var (
        r    Row
        line string
    )
    if line, r.Error = csvr.readLine(); r.Error != nil {
        return r
    }
    r.Fields = strings.FieldsFunc(
        line, func(c int) bool { return c == csvr.Sep })
    if csvr.Trim {
        for i := 0; i < len(r.Fields); i++ {
            r.Fields[i] = strings.Trim(r.Fields[i], csvr.Cutset)
        }
    }
    return r
}


//  Read rows into a preallocated buffer. Return the number of rows read,
//  and any error encountered.
func (csvr *Reader) ReadRows(rbuf [][]string) (int, os.Error) {
    var (
        i   int
        err os.Error
    )
    csvr.DoN(len(rbuf), func(r Row)bool {
        err = r.Error
        if r.Fields != nil {
            rbuf[i] = r.Fields
            i++
        }
        return !r.HasError()
    } )
    return i, err
}

//  Reads any remaining rows of CSV data in the underlying io.Reader.
func (csvr *Reader) RemainingRows() (rows [][]string, err os.Error) {
    return csvr.RemainingRowsSize(16)
}

//  Like csvr.RemainingRows(), but allows specification of the initial
//  row buffer capacity to avoid unnecessary reallocations.
func (csvr *Reader) RemainingRowsSize(size int) ([][]string, os.Error) {
    var (
        err  os.Error
        rbuf = make([][]string, 0, size)
    )
    csvr.Do(func(r Row) bool {
        err = r.Error
        //log.Printf("Scanned %v", r)
        if r.Fields != nil {
            rbuf = append(rbuf, r.Fields)
        }
        return !r.HasError()
    } )
    return rbuf, err
}

//  Iteratively read the remaining rows in the reader and call f on each
//  of them. If f returns false, no more rows will be read.
func (csvr *Reader) Do(f func(Row) bool) {
    for r := csvr.ReadRow() ; true ; r = csvr.ReadRow() {
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
    csvr.Do( func(r Row) bool {
        if i < n {
            return f(r)
        }
        return false
    } )
}

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
    return csvr.Configure(DefaultSep, DefaultTrim, DefaultCutset)
}

/* Comply with the reader interface. */
/*
func (csvr*Reader) Read(b []byte) (n int, err os.Error) {
	return csvr.r.Read(b)
}
*/
