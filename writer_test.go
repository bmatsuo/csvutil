// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// See csv_test.go for more information about each test.
package csvutil

import (
    "testing"
    "bytes"
)


// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
func TestWriteRow(T *testing.T) {
    //var csvBuf []byte = make([]byte,0 , 200)
    var bwriter *bytes.Buffer = bytes.NewBufferString("")
    var csvw *Writer = NewWriter(bwriter, nil)
    var csvMatrix = makeTestCSVMatrix()
    var n int = len(csvMatrix)
    var length = 0
    for i := 0; i < n; i++ {
        nbytes, err := csvw.WriteRow(csvMatrix[i]...)
        if err != nil {
            T.Errorf("Write error: %s\n", err.String())
        }
        errFlush := csvw.Flush()
        if errFlush != nil {
            T.Logf("Wrote %d bytes on row %d\n", nbytes, i)
        }
        length += nbytes
    }
    flushErr := csvw.Flush()
    if flushErr != nil {
        T.Errorf("Error flushing output; %v\n", flushErr)
    }
    var output string = bwriter.String()
    if len(output) == 0 {
        T.Error("Read 0 bytes\n")
    } else {
        T.Logf("Read %d bytes from the buffer.", len(output))
    }
    var csvStr string = makeTestCSVString()
    if output != csvStr {
        T.Errorf("Unexpected output.\n\nExpected:\n'%s'\nReceived:\n'%s'\n\n",
            csvStr, output)
    }
}
// END TEST1
