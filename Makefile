# Modified the basic makefiles referred to from the
# Go home page.
# 		-- Bryan Matsuo [bmatsuo@soe.ucsc.edu]
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=csvutil
GOFILES=\
	row.go\
	reader.go\
	writer.go\
	config.go\
	file.go\
	csvutil.go\

include $(GOROOT)/src/Make.pkg
