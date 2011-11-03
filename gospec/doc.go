/*
 *  Filename:    doc.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Nov  2 17:58:46 PDT 2011
 *  Description: Documentation for gospec
 *  Usage:       godoc github.com/bmatsuo/gospec
 */

/*
Gospec is a light wrapper around Gotest, which helps structure tests for
writing behaviour-driven tests with the "spec" package.

Gospec requires a directory structure for tests, which are separate from
source files. There is a root directory which will be called `spec/` here.

    spec/
        backend/
            types_spec.go
            interface_spec.go
            ...
        frontend/
            types_spec.go
            interface_spec.go
            ...
        ...

Gospec finds all `*_spec.go` files rooted at `spec/` runs them with Gotest
(using the -file flag).

Gospec interacts with the "spec" package by setting the GOSPECPATTERN in the
environment of the spawned Gotest process. This regular expression can select
which Specs to execute by matching against their context (test) name.

Additionally, the standard Gotest method of selecting tests by matching their
function name works. This supercedes Spec selection.

Usage:

The general gospec command syntax is

    gospec [options] [-v] [ROOT [PATTERN ...]]

Arguments:

The `ROOT` argument specifies a directory other than `./spec/` to look for
spec files.  The `PATTERN` arguments define separate regular expressions to
match against Spec contexts before running. When given, all the `PATTERN`
arguments are joined with a "|" and the value replaces the value of the flag
`-spec`.

Options:

    -root="./spec"  Directory containing spec files.

    -spec=".*"      Regexp matching Spec contexts.

    -test=".*"      Regexp matching test names (gotest -test.run).

    -v=false        Verbose program output.

*/
package documentation
