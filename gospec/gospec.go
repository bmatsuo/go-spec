// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
/*
 *  Filename:    gospec.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Nov  2 17:58:46 PDT 2011
 *  Description: <no value>
 *  Usage:       gospec [options] ARGUMENT ...
 */
import (
    "fmt"
    "os"
)

func Error(err os.Error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "gospec: %s\n", err.String())
    }
}

func FatalError(err os.Error) {
    if err != nil {
        Error(err)
        os.Exit(1)
    }
}

var opt options

func main() {
    opt = parseFlags()
    root := "./spec"
    specfiles, err := SpecGoFiles(root)
    FatalError(err)
    if len(specfiles) == 0 {
        FatalError(fmt.Errorf("No spec files in %s", root))
    }

    cmd := GoTestCommand().
        SpecGoFiles(specfiles)
    if opt.Verbose {
        cmd = cmd.Verbose()
    }
    cmd = cmd.TestPattern(opt.TestPattern)
    FatalError(cmd.Run(opt.SpecPattern))
}
