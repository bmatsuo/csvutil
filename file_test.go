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
    if statErr == os.ENOENT {
        return
    }
    if statErr != nil {
        T.Errorf("Error stat'ing the test file %s; %s\n", f, statErr.String())
    }
    rmErr := os.Remove(f)
    if rmErr != nil {
        T.Error("Error removing the test file %s; %s\n", f, rmErr.String())
    }
}

// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
var(
    TestIn string = "_test-csvutil-01-i.csv"
    TestOut string = "_test-csvutil-01-o.csv"
    TestPerm uint32 = 0622
)
func TestWriteFile (T *testing.T) {
    var testFilename string = TestOut
    defer cleanTestFile(testFilename, T)
    mat, str := makeTestCSVInstance()
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

func TestReadFile (T *testing.T) {
    var testFilename string = TestOut
    defer cleanTestFile(testFilename, T)

    mat, str := makeTestCSVInstance()
    err := ioutil.WriteFile(testFilename, bytes.NewBufferString(str).Bytes(), TestPerm)
    if err != nil {
        T.Error(err)
    }
    inputMat, csvErr := ReadFile(testFilename)
    if csvErr != nil {
        T.Errorf("CSV reading error: %s", err.String())
    }
    T.Logf("\nExpected;\n'%v'\n Received:\n'%v'\n\n", mat, inputMat)
    if len(inputMat) != len(mat) {
        T.Fatal("INPUT MISMATCH; number of rows")
    }
    for i:=0 ; i<len(mat) ; i++ {
        if len(mat[i]) != len(inputMat[i]) {
            T.Errorf("INPUT MISMATCH; row %d", i)
        } else {
            for j:=0 ; j<len(mat[i]) ; j++ {
                if mat[i][j] != inputMat[i][j] {
                    T.Errorf("INPUT MISMATCH; %d, %d", i, j)
                }
            }
        }
    }
}
// END TEST1
