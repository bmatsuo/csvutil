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
    "bytes"
)


// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
func TestWriteRow(T *testing.T) {
    //var csvBuf []byte = make([]byte,0 , 200)
    var bwriter *bytes.Buffer = bytes.NewBufferString("")
    var csvw *Writer = NewWriter(bwriter)
    var csvMatrix = makeTestCSVMatrix()
    var n int = len(csvMatrix)
    var length = 0
    for i := 0; i < n; i++ {
        nbytes, err := csvw.WriteFieldsln(csvMatrix[i])
        if err != nil {
            T.Errorf("Write error: %s\n", err.String())
        } else {
            T.Logf("Wrote %d bytes on row %d\n", nbytes, i)
        }
        length += nbytes
    }
    flushErr := csvw.Flush()
    if flushErr != nil {
        T.Errorf("Error flushing output; %v\n", flushErr)
    }
    var output string = string(bwriter.Bytes())
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
