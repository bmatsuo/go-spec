// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main
/*
 *  Filename:    main_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Nov  2 17:58:46 PDT 2011
 *  Description: <no value>
 *  Usage:       gotest
 */
import (
    "testing"
    "strings"
    . "spec"
)

func TestGospec(T *testing.T) {
    s := NewSpecTest(T)
    s.Describe("GoSpec", func() {
        s.Describe("test files", func() {
            var specfiles []string
            getspecs := func() []string { return specfiles }
            numspecs := func() int { return len(specfiles) }
            s.Before(All, func() { specfiles, _ = SpecGoFiles("./spec") })
            s.They("are in directory ./spec", func() {
                s.Spec(
                    numspecs,
                    Should, Satisfy,
                    func(x int) bool { return x > 0 })
            })
            s.They("have names ending in _spec.go", func() {
                s.Spec(
                    getspecs,
                    Should, Satisfy,
                    func(paths []string) bool {
                        for _, spec := range paths {
                            if !strings.HasSuffix(spec, "_spec.go") {
                                return false
                            }
                        }
                        return true
                    })
            })
        })
    })
}
