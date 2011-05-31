package csvutil
/* 
*  File: file.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sun May 29 23:14:48 PDT 2011
*/
import (
    "os"
)

//  Write CSV data to a named file. If the file does not exist, it is
//  created. If the file exists, it is truncated upon opening. Requires
//  that file permissions be specified. Recommended permissions are 0600,
//  0622, and 0666 (6:rw, 4:w, 2:r). 
func WriteFile(filename string, perm uint32,rows [][]string) (nbytes int, err os.Error) {
    var (
        out *os.File
        csvw *Writer
    )
    nbytes = 0
    out, err = os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, perm)
    if err != nil {
        return nbytes, err
    }
    csvw = NewWriter(out)
    nbytes, err = csvw.WriteRows(rows)
    if err != nil {
        return nbytes, err
    }
    err = out.Close()
    return nbytes, err
}

//  Read a named CSV file into a new slice of new string slices.
func ReadFile(filename string) (rows [][]string, err os.Error) {
    var (
        in *os.File
        csvr *Reader
    )
    in, err = os.Open(filename)
    if err != nil {
        return rows, err
    }
    csvr = NewReader(in)
    rows, err = csvr.RemainingRows()
    if err != nil {
        return rows, err
    }
    err = in.Close()
    return rows, err
}
