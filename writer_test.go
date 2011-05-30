// See csv_test.go for more information about each test.
package csvutil

import (
    "testing"
    "bytes"
    "os"
)


// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
func TestWriteRow (T *testing.T) {
    var csvBuf []byte = make([]byte, 200)
    var bwriter *bytes.Buffer = bytes.NewBuffer(csvBuf)
    var csvw *Writer = NewWriter(bwriter)
    var csvMatrix = makeTestCSVMatrix()
    var n int = len(csvMatrix)
    var length = 0
    for i:=0 ; i<n ; i++ {
        nbytes, err := csvw.WriteFieldsln(csvMatrix[1])
        if err != nil {
            T.Errorf("Write error: %s\n", err.String())
        } else {
            T.Logf("Wrote %d bytes on row %d\n", nbytes, i)
        }
        length += nbytes
    }
    var breader *bytes.Buffer = bytes.NewBuffer(csvBuf)
    output,err := breader.ReadString('\n')
    if err != nil && err != os.EOF {
        T.Error(err.String())
    }
    if len(output) == 0 {
        T.Error("Read 0 bytes\n")
    } else {
        T.Logf("Read %d bytes from the buffer.", len(output))
    }
    var csvStr string = makeTestCSVString()
    if output != csvStr {
        T.Errorf("Unexpected output.\n\nExpected:\n'%s'\nReceived:\n'%s'\n\n",
                 csvStr, output);
    }
    T.Error("FAIL")
}
// END TEST1
