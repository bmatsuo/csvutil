// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file doesn't actually do any tests.
// But it has helper functions/data for testing functions.
package csvutil

import (
    "strings"
)


// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
//  This is a simple test of the readers ability to return something of
//  ideal input format in the proper internal form. the dimensions and
//  values of the resulting parsed matrix are verified against the [][]string
//  matrix which created the CSV data.
func makeTestCSVMatrix() [][]string {
    var testfields [][]string = make([][]string, 3)
    for i := 0; i < 3; i++ {
        testfields[i] = make([]string, 3)
    }
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
func makeTestCSVString() string {
    var testfields [][]string = makeTestCSVMatrix()
    var rows []string = make([]string, 4)
    rows[0] = strings.Join(testfields[0], ",")
    rows[1] = strings.Join(testfields[1], ",")
    rows[2] = strings.Join(testfields[2], ",")
    rows[3] = ""
    return strings.Join(rows, "\n")
}
func makeTestCSVInstance() ([][]string, string) {
    return makeTestCSVMatrix(), makeTestCSVString()
}
// END TEST1
