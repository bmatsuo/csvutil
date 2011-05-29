package csvutils
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

const (
	DEFAULT_SEP  = ','
	DEFAULT_TRIM = false
	DEFAULT_CUTSET = " \t"
)

type Reader struct {
	Sep int "Seperator character."
	Trim bool "Remove excess whitespace from field values."
	Cutset string "Set of characters to trim."
	r io.Reader "Base reader object."
	br *bufio.Reader "For reading lines."
}

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

func NewReader(r io.Reader) *Reader {
	var csvr *Reader = new(Reader).Reset()
	csvr.r = r
	csvr.br = bufio.NewReader(r)
	return csvr
}

func NewReaderSize(r io.Reader, size int) *Reader {
	/* Create a new csv.Reader with a specified buffer size. */
	var csvr *Reader = new(Reader).Reset()
	csvr.r = r
	var err os.Error
	csvr.br, err = bufio.NewReaderSize(r, size)
	if err != nil { panic(err) }
	return csvr
}

func (csvr *Reader) ReadRow() Row {
	/* Read one row of the CSV and and return an array of the fields. */
	var r Row
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
	return r
}

func (csvr *Reader) EachRow() <-chan Row {
	/* Generator function for iterating through rows. */
	var c chan Row = make(chan Row)
	var read_rows = func (c chan<- Row) {
		for true {
			var r Row = csvr.ReadRow()
			if r.Fields == nil {
				if r.Error == os.EOF {
					break
				} else {
					panic(r.Error)
				}
			}
			c <- r
		}
		close(c)
	}
	go read_rows(c)
	return c
}

func (csvr *Reader) Configure(sep int, trim bool, cutset string) *Reader {
	csvr.Sep = sep
	csvr.Trim = trim
	csvr.Cutset = cutset
	return csvr
}

func (csvr *Reader) Reset() *Reader {
	return csvr.Configure(DEFAULT_SEP, DEFAULT_TRIM, DEFAULT_CUTSET)
}

/* Comply with the reader interface. */
func (csvr*Reader) Read(b []byte) (n int, err os.Error) {
	return csvr.r.Read(b)
}
