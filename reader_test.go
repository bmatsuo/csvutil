package csvutil

import (
    "testing"
    "strings"
)


// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
//  This is a simple test of the readers ability to return something of
//  ideal input format in the proper internal form. the dimensions and
//  values of the resulting parsed matrix are verified against the [][]string
//  matrix which created the CSV data.
func makeTestCSVMatrix() [][]string {
    var testfields [][]string = make([][]string, 3);
    for i:=0 ; i<3 ; i++ { testfields[i] = make([]string, 3) }
    testfields[0][0] = "field1"
    testfields[0][1] = "field2"
    testfields[0][2] = "field3"
    testfields[1][0] = "Ben Franklin"
    testfields[1][1] = "3.704"
    testfields[1][2] = "10"
    testfields[2][0] = "Tom Jefferson"
    testfields[2][1] = "5.7"
    testfields[2][2] = "15"
    return testfields
}
func makeTestCSVString () string {
    var testfields [][]string = makeTestCSVMatrix()
    var rows []string = make([]string, 3);
    rows[0] = strings.Join(testfields[0], ",")
    rows[1] = strings.Join(testfields[1], ",")
    rows[2] = strings.Join(testfields[2], ",")
    return strings.Join(rows, "\n")
}
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
