// csvutils - CSV file utilities for the Go programming language.
//
// VERSION
//	
//  This is csutils version iteration 0.0_2
//
// SYNOPSIS
//
//  The "csvutils" package can be used to read CSV files from a io.Reader.
//
//      reader := csvutils.NewReader(os.Stdin)
//      reader.Trim = true
//      for r := range reader.EachRow() {
//          if r.Error == nil {
//              break
//          }
//          var fields []string = r.Fields
//          // ...
//          // Process the CSV row fields.
//          // ...
//      }
package csvutils
/* 
*  File: csv.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat May 28 23:53:36 PDT 2011
*/
import (
)
