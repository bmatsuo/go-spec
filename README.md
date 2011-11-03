About GoSpec
=============

GoSpec is another interpretation of RSpec for the Go language.

GoSpec provides a package "spec" for writing descriptive tests in Go. It also
provides a command line program `gospec` for running these tests. Doing these
things, GoSpec creates a thing wrapper over Go's "testing" package and `gotest`
respectively.

Documentation
=============

Prerequisites
-------------

You must have Go installed (http://golang.org/). 

Installation
-------------

See `spec/` and `gospec/`

Examples
--------

This is a "tutorial" project in `examples/tutorial/`. It's a program that does
nothing but run rspec on itself. But the tests describe and show the
capabilities of GoSpec. To run it

    cd examples/tutorial
    gomake

You can also look at the `gospec/` directory. The `gospec` program has its
tests written using the "spec" package.

    cd gospec
    gospec -v


General Documentation
---------------------

See `spec/` and `gospec/`

Author
======

Bryan Matsuo <bmatsuo@soe.ucsc.edu>

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.

Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
