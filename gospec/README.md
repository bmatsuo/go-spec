About Gospec
============

Gospec is a lightweight wrapper for gotest (for "spec" package users). It
spearches for spec files (`*_spec.go`) in the directory `./spec/`.  It
executes these files with the help of the standard Gotest command.

Using Gospec does not depend on the ["spec" package](https://github.com/bmatsuo/go-spec/tree/master/spec).
As long as the naming and directory structure is followed any "testing" files
can be executed by Gospec.

Go-Spec allows for spec.Specs to be run selectively by matching the context they
belong to against a regular expression by using the environment variable
GOSPECPATTERN, which is used by the "spec" package.

Documentation
=============
Installation
------------

Install Gospec with Goinstall

    goinstall github.com/bmatsuo/go-spec/gospec

Or, you can build program yourself by cloning the repository.

    git clone https://github.com/bmatsuo/go-spec/gospec
    cd go-spec/gospec
    gomake install

Usage
-----

The general gospec command syntax is

    gospec [options] [-v] [ROOT [PATTERN ...]]

Arguments
---------

The `ROOT` argument specifies a directory other than `./spec/` to look for
spec files.  The `PATTERN` arguments define separate regular expressions to
match against Spec contexts before running. When given, all the `PATTERN`
arguments are joined with a "|" and the value replaces the value of the flag
`-spec`.

Options
-------

    -root="./spec"  Directory containing spec files.

    -spec=".*"      Regexp matching Spec contexts.

    -test=".*"      Regexp matching test names (gotest -test.run).

    -v=false        Verbose program output.

More documentation
------------------

For full documentation of `gospec`

    godoc github.com/bmatsuo/go-spec/gospec

Alternatively, use a godoc http server

    godoc -http=:6060

and view the docs [here](http://localhost:6060/pkg/github.com/bmatsuo/go-spec/gospec)


Author
======

Bryan Matsuo <bmatsuo@soe.ucsc.edu>

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.

Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
