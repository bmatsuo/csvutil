package csvutil

/*
 *  Filename:    config_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Tue Jul 12 01:56:03 PDT 2011
 *  Description: 
 *  Usage:       gotest
 */
import (
	"testing"
	"unicode/utf8"
)

//  Some rediculously dumb tests of Config methods, which are very simple.
func TestConfig(T *testing.T) {
	var config = NewConfig()

	// Test comment detection.
	config.CommentPrefix = "//"
	if !config.LooksLikeComment("// This should be a comment.\n") {
		T.Error("Did not correctly identify a // comment")
	}
	if config.LooksLikeComment("/ This, is not, a comment\n") {
		T.Error("Incorrectly labeled something a // comment")
	}

	// Test seperator detection.
	config.Sep = '\t'
	str := "\t"
	c, n := utf8.DecodeRuneInString(str)
	if c == utf8.RuneError && n == 1 {
		T.Errorf("Could not decode rune in string %q", str)
	}
	if !config.IsSep(c) {
		T.Error("Did not correctly identify a \\t separator")
	}
	if config.IsSep(52) {
		T.Error("Incorrectly labelled 52 a \\t separator")
	}
}
