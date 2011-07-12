package csvutil
/*
 *  Filename:    config_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Tue Jul 12 01:56:03 PDT 2011
 *  Description: 
 *  Usage:       gotest
 */
import (
    "utf8"
    "testing"
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
    var str = "\t"
    var utf8str = utf8.NewString(str)
    if !config.IsSep(utf8str.At(0)) {
        T.Error("Did not correctly identify a \\t separator")
    }
    if config.IsSep(52) {
        T.Error("Incorrectly labelled 52 a \\t separator")
    }
}
