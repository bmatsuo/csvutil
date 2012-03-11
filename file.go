// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvutil

/* 
*  File: file.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sun May 29 23:14:48 PDT 2011
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
import (
	"io"
	"os"
)

//  Write a slice of rows (string slices) to an io.Writer object.
func Write(w io.Writer, rows [][]string) (int, error) {
	var (
		csvw        = NewWriter(w, nil)
		nbytes, err = csvw.WriteRows(rows)
	)
	if err != nil {
		return nbytes, err
	}
	return nbytes, csvw.Flush()
}

//  Write CSV data to a named file. If the file does not exist, it is
//  created. If the file exists, it is truncated upon opening. Requires
//  that file permissions be specified. Recommended permissions are 0600,
//  0622, and 0666 (6:rw, 4:w, 2:r). 
func WriteFile(filename string, perm os.FileMode, rows [][]string) (int, error) {
	var (
		out    *os.File
		nbytes int
		err    error
		mode   = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	)
	if out, err = os.OpenFile(filename, mode, perm); err != nil {
		return nbytes, err
	}
	if nbytes, err = Write(out, rows); err != nil {
		return nbytes, err
	}
	return nbytes, out.Close()
}

//  Read rows from an io.Reader until EOF is encountered.
func Read(r io.Reader) ([][]string, error) {
	var csvr = NewReader(r, nil)
	return csvr.RemainingRows()
}

//  Read a named CSV file into a new slice of new string slices.
func ReadFile(filename string) ([][]string, error) {
	var (
		in   *os.File
		rows [][]string
		err  error
	)
	if in, err = os.Open(filename); err != nil {
		return rows, err
	}
	if rows, err = Read(in); err != nil {
		return rows, err
	}
	return rows, in.Close()
}

//  Iteratively apply a function to Row objects read from an io.Reader.
func Do(r io.Reader, f func(r Row) bool) {
	var csvr = NewReader(r, nil)
	csvr.Do(f)
}

//  Iteratively apply a function to Row objects read from a named file.
func DoFile(filename string, f func(r Row) bool) error {
	var in, err = os.Open(filename)
	if err != nil {
		return err
	}
	Do(in, f)
	return in.Close()
}
