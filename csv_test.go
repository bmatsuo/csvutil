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
var (
    TestMatrix1 = [][]string{
        []string{"field1", "field2", "field3"},
        []string{"Ben Franklin", "3.704", "10"},
        []string{"Tom Jefferson", "5.7", "15"}}
)
func csvTestString1() string {
    var testfields [][]string = TestMatrix1
    var rows []string = make([]string, 4)
    rows[0] = strings.Join(testfields[0], ",")
    rows[1] = strings.Join(testfields[1], ",")
    rows[2] = strings.Join(testfields[2], ",")
    rows[3] = ""
    return strings.Join(rows, "\n")
}
func csvTestInstance1() ([][]string, string) {
    return TestMatrix1, csvTestString1()
}

// END TEST1

//  TEST2 - 3x3 matrix w/ tab separators, w/o excess whitespace. And with
//  leading '#' comments. 
var (
    TestMatrix2 = [][]string{
        []string{"field1", "field2", "field3"},
        []string{"Ben Franklin", "3.704", "10"},
        []string{"Tom Jefferson", "5.7", "15"}}
)
func csvTestString2() string {
    var testfields [][]string = TestMatrix2
    var rows []string = make([]string, 5)
    rows[0] = "# This is a comment string"
    rows[0] = "# This another comment string"
    rows[1] = strings.Join(testfields[0], "\t")
    rows[2] = strings.Join(testfields[1], "\t")
    rows[3] = strings.Join(testfields[2], "\t")
    rows[4] = ""
    return strings.Join(rows, "\n")
}
func csvTestInstance2() ([][]string, string) {
    return TestMatrix2, csvTestString2()
}
// END TEST2
