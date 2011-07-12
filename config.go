package csvutil
/*
 *  Filename:    config.go
 *  Package:     csvutil
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Tue Jul 12 01:56:03 PDT 2011
 *  Description: Define the configuration type for Readers and Writers.
 */
import ()

//  A configuration structure that can be shared between a Reader and Writer.
type Config struct {
    // General configuration
    //  Field seperator
    Sep int
    //  Trim leading/trailing whitespace in fields.
    Trim bool
    //  Characters to trim from fields.
    Cutset string
    //  Prefix for comment lines.
    CommentPrefix string

    // Reader specific config
    //  Are comments allowed in the input.
    Comments bool
    //  Comments can appear in the body (Comments must be true).
    CommentsInBody bool
}

//  The default configuration is used for Readers and Writers when none is
//  given.
var (
    DefaultConfig = &Config{
        Sep: ',', Trim: false, Cutset: " \t", CommentPrefix: "#",
        Comments: false, CommentsInBody: false}
)

//  Return a freshly allocated Config that is initialized to DefaultConfig.
func NewConfig() *Config {
    var c = new(Config)
    *c = *DefaultConfig
    return c
}

func (c *Config) LooksLikeComment(line string) bool {
    return line[:len(c.CommentPrefix)] == c.CommentPrefix
}

func (c *Config) IsSep(rune int) bool {
    return rune == c.Sep
}
