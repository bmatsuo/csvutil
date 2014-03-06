// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// See csv_test.go for more information about each test.
package csvutil

import (
	"testing"
)

// TEST1 - Simple 3x3 matrix w/ comma separators and w/o excess whitespace.
func TestWriteRow(T *testing.T) {
	var csvw, buff = BufferWriter(nil)
	var csvMatrix = TestMatrix1
	var n int = len(csvMatrix)
	var length = 0
	for i := 0; i < n; i++ {
		nbytes, err := csvw.WriteRow(csvMatrix[i]...)
		if err != nil {
			T.Errorf("Write error: %s\n", err)
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
	var output string = buff.String()
	if len(output) == 0 {
		T.Error("Read 0 bytes\n")
	} else {
		T.Logf("Read %d bytes from the buffer.", len(output))
	}
	var csvStr string = csvTestString1()
	if output != csvStr {
		T.Errorf("Unexpected output.\n\nExpected:\n'%s'\nReceived:\n'%s'\n\n",
			csvStr, output)
	}
}

// END TEST1

func TestWriterComments(T *testing.T) {
	var config = NewConfig()
	config.Sep = '\t'
	var (
		matrix       = TestMatrix2
		comments     = TestMatrix2Comments
		verification = csvTestString2()
		writer, buff = BufferWriter(config)
	)
	writer.WriteComments(comments...)
	writer.WriteRows(matrix)
	writer.Flush()
	var output = buff.String()
	if output != verification {
		T.Errorf("Error writing comments\n\n'%s'\n'%s'", verification, output)
	}
}
