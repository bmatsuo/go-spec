// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 *  Filename:    spectutorial.go
 *  Package:     spectutorial
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Tue Nov  1 23:13:13 PDT 2011
 *  Description: <no value>
 */

// The tutorial command does ....
package main

import (
    "exec"
    "fmt"
    "os"
)

func Fatalf(err os.Error) {
    fmt.Fprintf(os.Stderr, "tutorial: ")
    fmt.Fprint(os.Stderr, err.String())
    os.Exit(1)
}

func main() {
    cmd := exec.Command("gospec", "-v", "spec")

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err := cmd.Run()
    if err != nil {
        Fatalf(err)
    }
}
