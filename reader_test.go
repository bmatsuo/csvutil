// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil

import (
    "testing"
)

func TestDo(T *testing.T) {
    var csvr *Reader = StringReader(csvTestString1(), nil)
    var rowlen = -1
    csvr.Do(func(r Row) bool {
        if r.HasError() {
            T.Errorf("Read error encountered: %v", r.Error)
            return false
        }
        if rowlen == -1 {
            rowlen = len(r.Fields)
        } else if rowlen != len(r.Fields) {
            T.Error("Row length error, non-rectangular.")
        }
        return true
    })
}


// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
func TestReadRow(T *testing.T) {
    T.Log("Beginning test\n")
    var csvr *Reader = StringReader(csvTestString1(), nil)
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
    csvr.Do(func(row Row) bool {
        var k int = len(row.Fields)
        if k != n {
            T.Errorf("Unexpected row size %d (!= %d)\n", k, n)
        }
        var j int = 0
        for j = 0; j < k; j++ {
            var field string = row.Fields[j]
            if len(field) < 1 {
                T.Error("Unexpected non-empty string\n")
            }
        }
        rows[i] = row.Fields
        i++
        return true
    })
    var test_matrix [][]string = TestMatrix1
    var assert_val = func(i, j int) {
        if rows[i][j] != test_matrix[i][j] {
            T.Errorf("Unexpected value in (%d,%d), %s", i, j, rows[i][j])
        }
    }
    for i := 0; i < n; i++ {
        for j := 0; j < n; j++ {
            assert_val(i, j)
        }
    }
    T.Log("Finished test\n")
}
// END TEST1

func TestComments(T *testing.T) {
    // Create the test configuration.
    var config = NewConfig()
    config.Sep = '\t'
    config.Comments = true

    // Create a Reader for the test string and parse the rows.
    var (
        reader    = StringReader(csvTestString2(), config)
        rows, err = reader.RemainingRows()
    )
    if err != nil {
        T.Error("Error:", err.Error())
    }
    if len(rows) != len(TestMatrix2) {
        T.Error("Different number of rows in parsed and original", len(rows), len(TestMatrix2))
    }
    for i, row := range rows {
        if i >= len(TestMatrix2) {
            break
        }
        for j, s := range row {
            if j >= len(TestMatrix2[i]) {
                T.Errorf("Row %d: different number of columns in parsed %d and original %d",
                    i, len(row), len(TestMatrix2[i]))
                break
            }
            if s != TestMatrix2[i][j] {
                T.Errorf("'%s' != '%s' at (%d,%d)", s, TestMatrix2[i][j], i, j)
            }
        }
    }
}
