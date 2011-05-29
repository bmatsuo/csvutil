package main
/* 
*  File: test_main.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sun May 29 01:36:46 PDT 2011
*  Usage: test [options] 
*/
import (
    "os"
	"csv"
	//"gocsv.googlecode.com/hg"
	"github.com/kr/pretty.go"
    "log"
    //"fmt"
    //"flag"
    //"io"
)

func main() {
	reader := csv.NewReader(os.Stdin)
	reader.Trim = true
	for r := range reader.EachRow() {
		log.Printf("%# v\n", pretty.Formatter(r))
	}
	/*
	rows, err := csv.ReadAll(os.Stdin)
	if err != nil {
		panic(err.String())
	}
	log.Printf("%#\n", pretty.Formatter(rows))
	*/
}
