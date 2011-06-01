package csvutil
/* 
*  File: iterator.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Tue May 31 08:38:58 PDT 2011
*/
import (
    "os"
)

//  An object for iterating over the rows of a Reader.
//      rit = reader.RowIterStarted()
//      for r := range reader.RowsChan {
//          if r.Error != nil {
//              rit.Break()
//              panic(r.Error)
//          }
//          var fields []string = r.Fields
//          if !fieldsOk(fields) {
//              rit.Break()
//              break
//          }
//          // Process the desired entries of "fields".
//          rit.Next()
//      }
//  This itererator is safe to break out of. For iterators meant for
//  parsing an entire stream, see the ReaderRowIteratorAuto type.
type ReaderRowIterator struct {
    stopped bool
    RowsChan <-chan Row
    control chan<- bool
}

//  A ReaderRowIteratorAuto is meant for reading until encountering the
//  end of the data stream corresponding to a Reader's underlying
//  io.Reader. 
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

//  Like (*ReaderRowIterator) Break(), but also panics afterward. This
//  helps if its possible that the panic could be recovered.
func (csvri *ReaderRowIterator) Croak(err os.Error) {
    csvri.Break()
    panic(err)
}

//  Create a new row iterator and return it.
func (csvr *Reader) RowIter() (*ReaderRowIterator) {
    ri := new(ReaderRowIterator)
    throughChan := make(chan Row)
    controlChan := make (chan bool)
    ri.RowsChan = throughChan
    ri.control = controlChan
	var read_rows = func (r chan<- Row, c <-chan bool) {
        /* Deferring may be unnecessary now (its NOT desired in this context).
            defer func() {
                x:=recover()
                if x !=nil {
                }
            } ()
        */
		for true {
            cont, ok := <-c
            if !ok || !cont {
                break
            }
            csvr.LastRow = Row{Fields:nil, Error:nil}
			var row Row = csvr.ReadRow()
			if row.HasEOF() {
				break
            }
			if row.Fields == nil {
                if row.HasError() {
				    panic(row.Error)
                }
                panic("nilfields")
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
			if row.HasEOF() {
				break
            }
			if row.Fields == nil {
                if row.HasError() {
				    panic(row.Error)
                }
                panic("nilfields")
			}
            csvr.LastRow = row
			r <- row
		}
		close(r)
	}
	go read_rows(throughChan)
	return ri
}

