// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil
/* 
*  File: row.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Wed Jun  1 16:48:20 PDT 2011
*  Description: Row related types and methods.
 */
import (
    "os"
)

//  A simple row structure for rows read by a csvutil.Reader that
//  encapsulates any read error enountered along with any data read
//  prior to encountering an error.
type Row struct {
    Fields []string "CSV row field data"
    Error  os.Error "Error encountered reading"
}

//  A wrapper for the test r.Error == os.EOF
func (r Row) HasEOF() bool {
    return r.Error == os.EOF
}

//  A wrapper for the test r.Error != nil
func (r Row) HasError() bool {
    return r.Error != nil
}
