package csvutil
/* 
*  File: reader.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat May 28 23:53:36 PDT 2011
*  Description: CSV reader library.
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
	if err != nil { panic(err) }
	return csvr
}

//  Read up to a new line and return a slice of string slices
func (csvr *Reader) ReadRow() Row {
	/* Read one row of the CSV and and return an array of the fields. */
	r := Row{Fields:nil, Error:nil}
    csvr.LastRow = r
	line, isPrefix, err := csvr.br.ReadLine()
	r.Fields = nil
	r.Error = err
	if isPrefix { panic("longline")  } // TODO fix this
	if err == os.EOF { return r }
	if err != nil { return r }
	var fields []string
	fields = strings.FieldsFunc(
			string(line),
			func (c int) bool { return c == csvr.Sep } )
	if csvr.Trim {
		for i:=0 ; i<len(fields) ; i++ {
			fields[i] = strings.Trim(fields[i], csvr.Cutset)
		}
	}
	r.Fields = fields
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
    for r := range csvr.EachRow() {
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

//  A function with a concurrent routine for iterating through all
//  remaining rows of CSV data until EOF is encountered.
//      for r := range reader.EachRow() {
//          if r.Error != nil {
//              panic(r.Error)
//          }
//          var fields []string = r.Fields
//          // Process the desired entries of "fields".
//      }
//  This method utilizes a goroutine, and has the distict possibilty
//  of 'losing' at least one row when breaking of the of the iterating
//  loop. This can be managed by explicitly handling channel returned by
//  reader.EachRow() and closing/extracting the channel after the loop.
//      ch := reader.EachRow()
//      for r := range ch {
//          if r.Error != nil {
//              panic(r.Error)
//          }
//          // Process the desired entries of "fields".
//          if r.Fields[3] == "5" {
//              close(ch)
//              break
//          }
//      }
//      chrow, ok := <- ch
//      if ok {
//          // There was an extra row in the channel
//      }
//      panicrow := reader.LastRow
//      // Continue program execution.
//  However this method is simply recommended for use only when the
//  program needs to iteratively handle all remaining content of the
//  Reader.
func (csvr *Reader) EachRow() <-chan Row {
	/* Generator function for iterating through rows. */
    c := make(chan Row)
	var read_rows = func (c chan<- Row) {
        defer func() {
            if x:=recover(); x!=nil {
                /* Do nothing. */
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
