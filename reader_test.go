// See csv_test.go for more information about each test.
/*
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
    rit := csvr.RowIterAuto()
    for row := range rit.RowsChan {
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
