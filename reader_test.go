// See csv_test.go for more information about each test.
package csvutil

import (
    "testing"
    "strings"
)


// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
func TestReadRow (T *testing.T) {
    var csvStr string = makeTestCSVString()
    var sreader *strings.Reader = strings.NewReader(csvStr)
    var csvr *Reader = NewReader(sreader)
    var n int = -1
    var rows [][]string
    var headrow Row = csvr.ReadRow()
    n = len(headrow.Fields)
    if n != 3 {
        T.Errorf("Unexpected row size %d\n", n)
    }
    rows = make([][]string, n) // Expect a square matrix.
    rows[0] = headrow.Fields
    var i int = 1
    for row := range csvr.EachRow() {
        var k int = len(row.Fields)
        if k != n {
            T.Errorf("Unexpected row size %d (!= %d)\n", k, n)
        }
        var j int = 0
        for j=0 ; j<k ; j++ {
            var field string = row.Fields[j]
            if len(field) < 1 {
                T.Error("Unexpected non-empty string\n")
            }
        }
        rows[i] = row.Fields
        i++
    }
    var test_matrix [][]string = makeTestCSVMatrix()
    var assert_val = func (i,j int) {
        if rows[i][j] != test_matrix[i][j] {
            T.Errorf("Unexpected value in (%d,%d), %s", i, j, rows[i][j])
        }
    }
    for i:=0 ; i<n ; i++ {
        for j:=0 ; j<n ; j++ {
            assert_val(i,j)
        }
    }
}
// END TEST1
