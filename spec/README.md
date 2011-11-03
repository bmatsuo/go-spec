About GoSpec
=============

Package spec provides descriptive and flexible wrappers of the "testing"
package. It can be used along with command Gotest, or with the companion
command GoSpec.

GoSpec is meant to make writing comprehensive unit testing easier to manage
with Go. The "testing" package and Gotest are pretty good for out-of-the-box
functionality. But managing a large project with them is not very feasible.
Obviously inspired by Ruby's RSpec gem, GoSpec is used to write nested
specifications which test as well as describe the functionality of programs,
workflows, and objects.

Documentation
=============

Prerequisites
-------------

[Install Go](http://golang.org/doc/install.html). 

Installation
-------------

Use goinstall to install the "spec" command

    goinstall github.com/bmatsuo/go-spec/spec

Or, you can build both the package and the program yourself by cloning the repository.

    git clone https://github.com/bmatsuo/go-spec/gospec
    cd go-spec/spec
    gomake install

General Documentation
---------------------

For documentation of the "spec" package

    godoc github.com/bmatsuo/go-spec/spec

Alternatively, use a godoc http server

    godoc -http=:6060

and view the docs [here](http://localhost:6060/pkg/github.com/bmatsuo/go-spec/spec)

Author
======

Bryan Matsuo <bmatsuo@soe.ucsc.edu>

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.

Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
