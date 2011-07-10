//  CSV data utilities for the Go programming language.
//
//  This is csvutil version 0.2_51
//
package csvutil
/* 
*  File: csv.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat May 28 23:53:36 PDT 2011
*
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
import ()

//  The default CSV field separator is a comma ','. By default, field
//  trimming is turned off for csvutil's I/O objects. But, when active,
//  the default cutset is " \t", that is, (ASCII) whitespace.
const (
    DefaultSep    = ','
    DefaultTrim   = false
    DefaultCutset = " \t"
)
