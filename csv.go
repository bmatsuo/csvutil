//  CSV data utilities for the Go programming language.
//
//  This is csvutil version 0.1_4
//
package csvutil
/* 
*  File: csv.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat May 28 23:53:36 PDT 2011
*/
import (
)

//  The default CSV field separator is a comma ','. By default, field
//  trimming is turned off for csvutil's I/O objects. But, when active,
//  the default cutset is " \t", that is, (ASCII) whitespace.
const (
	DEFAULT_SEP  = ','
	DEFAULT_TRIM = false
	DEFAULT_CUTSET = " \t"
)
