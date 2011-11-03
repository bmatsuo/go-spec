// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spectutorial
/*
 *  Filename:    main_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Tue Nov  1 22:33:31 PDT 2011
 *  Description: <no value>
 *  Usage:       gotest
 */
import (
    "testing"
    . "spec"
    "os"
)

type TestStruct struct {
    a, b, c string
}

func TestSpecTutorial(T *testing.T) {
    s := New(T)
    s.Describe("A SpecTest", func() {

        s.Describe("Describe method", func() {
            s.Describe("message", func() {
                msg := "A SpecTest Describe method message is in a simple hierarchy"
                s.It("is in a simple hierarchy", func() {
                    s.Spec(s.String(), Should, Equal, msg)
                })
                s.It("changes depending on the scope", func() {
                    s.Spec(s.String(), Should, Not, Equal, msg)
                })
            })
            s.Describe("function", func() {
                s.It("can make multiple Spec method calls", func() {
                    s.Spec(1, Should, Equal, 1)
                    s.Spec(2, Should, Equal, 2)
                })
                s.Describe("Trigger", func() {
                    x := 0
                    getx := func()int{ return x }
                    s.Before(All, func() { x++ })
                    s.After(All, func() { x-- })
                    s.Describe("on All nested Specs", func() {
                        s.Before(All, func() { x++ })
                        s.After(All, func() { x-- })
                        s.It("can run before every nested Spec", func() {
                            s.Spec(getx, Should, Equal, 2)
                        })
                        s.It("can run after every nested Spec", func() {
                            s.Spec(getx, Should, Equal, 2)
                        })
                    })

                    s.Before(First, func() { x++ })
                    s.After(First, func() { x-- })
                    s.Describe("on the First nested Spec", func() {
                        s.It("can run before the next nested Spec", func() {
                            s.Spec(getx, Should, Equal, 2)
                        })
                        s.It("can run after the next nested Spec", func() {
                            s.Spec(getx, Should, Equal, 1)
                        })
                        s.It("is not run on subsequent nested Specs", func() {
                            s.Spec(getx, Should, Equal, 1)
                        })
                    })
                    s.Describe("outside of its scope", func() {
                        s.It("will not run", func() {
                            s.Spec(getx, Should, Equal, 1)
                        })
                    })
                })
            })
        })

        s.Describe("Spec method", func() {
            s.Describe("call", func() {
                s.It("can determine the equality two element Values", func() {
                    s.Spec(1, Should, Equal, 1)
                })
                s.It("can determine the inequality of two element Values", func() {
                    s.Spec(1, Should, Not, Equal, 2)
                })
                s.It("can verify predicate func of a Value returns true", func() {
                    s.Spec(false, Should, Satisfy, func(x bool) bool { return !x })
                })
            })

            s.Describe("argument", func() {
                s.Describe("Value", func() {
                    s.It("can be a native Go type", func() {
                        s.Spec(1, Should, Equal, 1)
                        s.Spec("abc", Should, Not, Equal, 1+3i)
                        s.Spec(false, Should, Not, Equal, true)
                    })
                    s.It("can be a struct", func() {
                        s.Spec(TestStruct{"1", "2", "3"}, Should, Equal, TestStruct{"1", "2", "3"})
                    })
                    s.It("gets called as a nil-adic func", func() {
                        s.Spec(3, Should, Equal, func() int { return 3 })
                    })
                    s.Describe("called as a nil-adic func", func() {
                        s.It("defaults to the first return value", func() {
                            s.Spec(
                                func() (bool, os.Error) { return true, nil },
                                Should, Equal, true)
                        })
                        s.It("is indexed with a following integer", func() {
                            s.Spec(
                                func() (int, os.Error) { return 45, nil },
                                0, Should, Equal, 45)
                        })
                        s.Describe("used as on Object", func() {
                            s.It("can checked for errors", func() {
                                s.Spec(
                                    func() (bool, os.Error) { return true, os.NewError("blah") },
                                    Should, HaveError)
                                s.Spec(
                                    func() (bool, os.Error) { return true, nil },
                                    Should, Not, HaveError)
                            })
                        })
                    })
                    s.Describe("can be a slice", func() {
                    })
                })
            })
        })
    })
}
