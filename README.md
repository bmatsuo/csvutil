csvutil - Your CSV pocket-knife (golang)

Version
=======

    This is csvutil version 1.0_4

Synopsis
========

The csvutil package can be used to read CSV data from any io.Reader and,
write CSV data to any io.Writer. It can automatically generate CSV rows
from slices containing native Go types, other slices/arrays of native go
types, or flat structs. It can also convert CSV rows and assign values to
memory referenced by a slice of pointers.

```go
package main
import "os"
import "github.com/bmatsuo/csvutil"

type Person struct {
    Name   string
    Height float64
    Weight float64
}

func main() {
    writer := csvutil.NewWriter(os.Stdout)
    errdo := csvutil.DoFile(os.Args[1], func(r csvutil.Row) bool {
        if r.HasError() {
            panic(r.Error)
        }
        person := Person{}
        if _, errc := r.Format(&person); errc != nil {
            panic("Row is not a Person")
        }
        bmi := person.Weight / (person.Height * person.Height)
        writer.WriteRow(csvutil.FormatRow(person, bmi).Fields...)
        return true
    })
    if errdo != nil {
        panic(errdo)
    }
    writer.Flush()
}
```

Given a CSV file 'in.csv' with contents

```
alice,1.4,50
bob,2,80
chris,1.6,67
```

When the above program is called and given the path 'in.csv', the following
is printed to standard output.

```
alice,1.4,50,25.510204081632658
bob,2,80,20
chris,1.6,67,26.171874999999996
```

About
=====

The csvutil package is written to make interacting with CSV data as easy
as possible. Efficiency of its underlying functions and methods are
slightly less important factors in its design. However, that being said,
**it is important**. And, csvutil should be capable of handling both
extremely large and fairly small streams of CSV data through input and
output quite well in terms of speed and memory footprint. It should do
this while not making your code bend over backwards (more than necessary).
As libraries should never make you do.

Features
========

* Slurping/spewing CSV data. That is, reading/writing whole files or data
streams at once.

* Iteration through individual rows of a CSV data stream for a smaller
memory footprint.

* Writing of individual writing fields/rows (along with batch writing).

* Automated CSV row serialization and deserialization (formatting) for flat
data structures and types.

Todo
====

* Enhance and clean the formatting API and allow better formatting of data.

Install
=======

The easiest installation of csvutil is done through goinstall.

    goinstall github.com/bmatsuo/csvutil

Documentation
=============

The best way to read the current csvutil documentation is using
godoc.

    godoc github.com/bmatsuo/csvutil

Or better yet, you can run a godoc http server.

    godoc -http=":6060"

Then go to the url http://localhost:6060/pkg/github.com/bmatsuo/csvutil/

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.

Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
