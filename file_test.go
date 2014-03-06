// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// See csv_test.go for more information about each test.
package csvutil

import (
    "testing"
    "os"
    "io/ioutil"
    "bytes"
)

func cleanTestFile(f string, T *testing.T) {
    _, statErr := os.Stat(f)
    if os.IsNotExist(statErr) {
        return
    }
    if statErr != nil {
        T.Errorf("Error stat'ing the test file %s; %s\n", f, statErr.Error())
    }
    rmErr := os.Remove(f)
    if rmErr != nil {
        T.Error("Error removing the test file %s; %s\n", f, rmErr.Error())
    }
}

// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
var (
    TestIn   string = "_test-csvutil-01-i.csv"
    TestOut  string = "_test-csvutil-01-o.csv"
    TestPerm os.FileMode = 0622
)

func TestWriteFile(T *testing.T) {
    var testFilename string = TestOut
    defer cleanTestFile(testFilename, T)
    mat, str := csvTestInstance1()
    nbytes, err := WriteFile(testFilename, TestPerm, mat)
    if err != nil {
        T.Error(err)
    }
    if nbytes == 0 {
        T.Error("Wrote 0 bytes.")
        return
    }
    T.Logf("Wrote %d bytes.\n", nbytes)
    var outputString []byte
    outputString, err = ioutil.ReadFile(testFilename)
    if err != nil {
        T.Errorf("Error reading the test output %s for verification", testFilename)
    }
    T.Logf("\nExpected:\n'%s'\nReceived:\n'%s'\n\n", outputString, str)
    if string(outputString) != str {
        T.Error("OUTPUT MISMATCH")
    }
}

func TestReadFile(T *testing.T) {
    var testFilename string = TestOut
    defer cleanTestFile(testFilename, T)

    mat, str := csvTestInstance1()
    err := ioutil.WriteFile(testFilename, bytes.NewBufferString(str).Bytes(), TestPerm)
    if err != nil {
        T.Error(err)
    }
    inputMat, csvErr := ReadFile(testFilename)
    if csvErr != nil {
        T.Errorf("CSV reading error: %s", err.Error())
    }
    T.Logf("\nExpected;\n'%v'\n Received:\n'%v'\n\n", mat, inputMat)
    if len(inputMat) != len(mat) {
        T.Fatal("INPUT MISMATCH; number of rows")
    }
    for i := 0; i < len(mat); i++ {
        if len(mat[i]) != len(inputMat[i]) {
            T.Errorf("INPUT MISMATCH; row %d", i)
        } else {
            for j := 0; j < len(mat[i]); j++ {
                if mat[i][j] != inputMat[i][j] {
                    T.Errorf("INPUT MISMATCH; %d, %d", i, j)
                }
            }
        }
    }
}
// END TEST1
