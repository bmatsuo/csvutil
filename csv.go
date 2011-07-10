// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//  CSV data utilities for the Go programming language.
package csvutil
import ()

//  The default CSV field separator is a comma ','. By default, field
//  trimming is turned off for csvutil's I/O objects. But, when active,
//  the default cutset is " \t", that is, (ASCII) whitespace.
const (
    DefaultSep    = ','
    DefaultTrim   = false
    DefaultCutset = " \t"
)
