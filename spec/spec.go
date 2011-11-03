// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 *  Filename:    spec.go
 *  Package:     spec
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Tue Nov  1 21:55:30 PDT 2011
 *  Description: Define SpecTest and methods Describe, They, It, Spec.
 */

/*
Package spec provides a flexible behavior-driven development (BDD) framework.
It wraps the functionality of the standard Go package "testing" to provide
descriptive and maintainable behavior specifications. Because it wraps
"testing", it can be used with command Gotest. Or, it can be used with the
companion command Gospec.

Specifications (or Specs) are defined by nesting them in a Describe call.
Spec and Describe are methods of the SpecTest type, the primary type of "spec".
The Describe method has aliases They and It. A new SpecTest is created with the
New function.


    import (
        "import"
        . "spec"
    )
    func TestFunction(T *testing.T) {
        s := New(T)
        s.Describe(`The "strconv" package`, func() {
            s.Describe("integer conversion", func(){
                s.Describe("with Itoa", func() {
                    s.It("makes strings from integers", func() {
                        s.Spec(
                            func()string{ return strconv.Itoa(123) },
                            Should, Equal, "123")
                    })
                })
                s.It("with Atoi", func() {
                    s.It("converts decimal character strings to integers", func() {
                        decconv := func()(int,error){ return strconv.Atoi("123") }
                        s.Spec(decconv, Should, Not, HaveError)
                        s.Spec(decconv, Should, Satisfy, func(x int) bool { return x == 123 })
                    })
                    s.It("can't convert hex character strings to integers", func() {
                        hexconv := func()(int,error){ return strconv.Atoi("0x123") }
                        s.Spec(hexconv , Should, Not, Equal, 0x123)
                    })
                })
            })
        })
    }

The above example defines the following tests

    The "strconv" package integer conversion with ItoA makes strings from integers
    The "strconv" package integer conversion with Atoi converts decimal character strings to integers
    The "strconv" package integer conversion with Atoi can't convert hex character strings to integers

The spec package makes use of runtime reflection, predicate functions, and deep
equality checking to be a flexible and lightweight test framework.

The Spec argument sequence has the following grammar

    VALUE [INDEX] Should [Not] FUNCTION [ARGUMENT]

The general thinking is that (element INDEX of) VALUE is an object and FUNCTION
acts as a method of VALUE with a boolean return type. The "Not" keyword
obviously negates the returned value of FUNCTION.

VALUE can be either a normal Go value (int, string, float64, struct, interface,
...). It can also be a nil-adic function with at least one return value. If
VALUE is a nil-adic function, it is called before the spec is evaluated. If
INDEX is given, that nil-adic function return value is used as the object.
The first nil-adic function return value is used an the object when no INDEX
is given.

FUNCTION can be any of the keywords

    Equal
    Satisfy
    HaveError

Equal and Satisfy both require a single argument while HaveError requires none.
Equal performs a deep equality test of the object against an argument object.
Satisfy requires a predicate function (boolean function of one argument) and
is true if the predicate is true for the object. HaveError requires the object
to be nil-adic which returns an error in its last return value. It returns true
if the function returned an error.
*/
package spec

import (
	"strings"
	"regexp"
	"fmt"
	"os"
)

var SpecPattern = os.Getenv("GOSPECPATTERN")
var specregexp *regexp.Regexp

func (t SpecTest) getSpecRegexp() {
	if len(SpecPattern) > 0 && specregexp == nil {
		var err error
		specregexp, err = regexp.Compile(SpecPattern)
		if err != nil {
			t.Fatalf("Can't compile GOSPECPATTERN %s", SpecPattern)
		}
	}
}

//  An abstraction of the type *testing.T with identical exported methods.
type Test interface {
	Log(...interface{})
	Logf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fail()
	FailNow()
	Failed() bool
}

type stringer interface {
    String() string
}

//  Returns a human-readable interpretation of a Spec sequence.
func specString(spec sequence) string {
	s := make([]string, len(spec))
	for i, v := range spec {
        switch v.value.(type) {
            case stringer:
                s[i] = v.value.(stringer).String()
		default:
			s[i] = fmt.Sprintf("%#v", v.value)
        }
	}
	return strings.Join(s, " ")
}

type pos uint8

type Quantifier uint8

const (
	All Quantifier = iota
	First
	Last
)

var quantStr = []string{
	All:   "All",
	First: "First",
	Last:  "Last",
}

func (q Quantifier) String() string { return quantStr[q] }

func errTrigger(pos string, q Quantifier) error {
	return fmt.Errorf("Bad trigger %s %s", pos, q.String())
}

type trigger struct {
	Quantifier
	fn func()
}

func popTrigger(ptr *[]trigger) (back trigger) {
	stack := *ptr
	n := len(stack)
	if n == 0 {
		return
	}
	back = stack[n-1]
	stack[n-1] = trigger{}
	*ptr = stack[:n-1]
	return
}

func popTriggerAt(i int, ptr *[]trigger) (tr trigger) {
	n := len(*ptr)
	if n == 0 {
		return
	} else if i >= n {
		return
	}
	tstack := (*ptr)[:i+1]
	tr = popTrigger(&(tstack))
	copy((*ptr)[i:], (*ptr)[i+1:])
	popTrigger(ptr)
	return

}

func (t SpecTest) Before(q Quantifier, fn func()) error {
	t.doDebug(func() {
		t.Logf("Before maketrigger")
	})
	return t.makeTrigger(&t.beforestack[t.depth-1], "Before", q, fn)
}

func (t SpecTest) After(q Quantifier, fn func()) error {
	t.doDebug(func() {
		t.Logf("After maketrigger")
	})
	return t.makeTrigger(&t.deferstack[t.depth-1], "After", q, fn)
}

func (t SpecTest) makeTrigger(stack *[]trigger, pos string, q Quantifier, fn func()) error {
	if q == Last && pos == "Before" {
		return errTrigger(pos, q)
	}
	s := *stack
	s = append(s, trigger{q, fn})
	t.doDebug(func() {
		t.Logf("triggers %#v", s)
	})
	*stack = s
	return nil
}

//  The primary object of the spec package. Describe tests using the Describe,
//  It, and They methods. Write individual tests using the Spec methods.
type SpecTest struct {
	Test
	depth       int
	spec        sequence
	runspec     bool
	passed      bool
	ranspec     bool
	err         error
	beforestack [][]trigger
	deferstack  [][]trigger
	descstack   []string
	debug       bool
}

//  Create a new SpecTest. Call this function at the begining of your test functions.
//      import (
//          "testing"
//          . "spec"
//      )
//      func TestObject(T *testing.T) {
//          s := New(T)
//          s.Describe("My object", func() {
//              ...
//          })
//      }
func New(T Test) *SpecTest { return &SpecTest{Test: T, descstack: nil, debug: false} }

//  Execute a function if t.debug is true.
func (t *SpecTest) doDebug(fn func()) {
	if t.debug {
		fn()
	}
}

//  Return a string describing the current tests being executed by t.
func (t *SpecTest) String() string { return strings.Join(t.descstack, " ") }

//  Begin a block that describes a given thing. Can be called again from the
//  does function to describe more specific elements of that thing.
func (t *SpecTest) Describe(thing string, does func()) {
	t.getSpecRegexp()

	t.descstack = append(t.descstack, thing)
	t.beforestack = append(t.beforestack, nil)
	t.deferstack = append(t.deferstack, nil)
	t.depth++

	oldrunspec := t.runspec
	if specregexp != nil && !specregexp.MatchString(t.String()) {
		t.runspec = false
	} else if specregexp != nil {
		t.runspec = true
	}

	defer func() {
		// Clear the SpecTest when the description's scope is left.
		t.depth--
		t.descstack = t.descstack[:t.depth]
		popTrigger(&t.beforestack[t.depth])
		popTrigger(&t.deferstack[t.depth])
		t.spec = nil
		t.passed = true
		t.ranspec = false
		t.err = nil
		t.runspec = oldrunspec
	}()

	after := t.deferstack
	if k := len(after[t.depth-1]); k > 0 {
		for j := 0; j < k; j++ {
			tr := after[t.depth-1][j]
			if tr.Quantifier == Last {
				defer tr.fn()
			}
		}
	}
	t.deferstack = after

	// Do the described tests.
	does()
	if t.ranspec {
		// Compute the result of executed Spec calls.
		ok := t.passed && t.err == nil
		result := "PASS"
		if t.err != nil {
			result = "ERROR"
		} else if !t.passed {
			result = "FAIL"
		}

		// Write a message summarizing Spec calls.
		msg := fmt.Sprintf("%s: %s", t.String(), result)
		if !ok {
			msg += fmt.Sprintf("\n\t%s", specString(t.spec))
		}
		if t.err != nil {
			msg += fmt.Sprintf("\n\tError: %s", t.err.Error())
		}

		// Write the message as an error if there was a problem.
		if ok {
			t.Log(msg)
		} else {
			t.Error(msg)
		}
	}
}

//  A synonymn of Describe. It's function is meant to contain calls to Spec.
func (t *SpecTest) It(specification string, check func()) { t.Describe(specification, check) }
//  A synonymn of Describe. It's function is meant to contain calls to Spec.
func (t *SpecTest) They(specification string, check func()) { t.Describe(specification, check) }

//  Specify a relation between two objects.
//      Spec("abc", Should, Equal, "abc")
//      Spec("abc", Should, Satisfy, func(x string)bool{ return "abc" })
//      v := Value(func() (string, os.Error) { return "abc", os.NewError("Oops!"))}
//      Spec( v, Should, HaveError)
//      Spec( v, Should, Equal, "abc")
func (t *SpecTest) Spec(spec ...interface{}) {
	// Run triggers regardless of context's matching status
	before := t.beforestack
	for i := 0; i < t.depth; i++ {
		if k := len(before[i]); k > 0 {
			for j := 0; j < k; j++ {
				tr := before[i][j]
				t.doDebug(func() {
					t.Logf("firing %#v", tr)
				})
				tr.fn()
				if tr.Quantifier == First {
					popTriggerAt(j, &before[i])
					k--
					j--
				}
			}
		}
	}
	t.beforestack = before
	after := t.deferstack
	for i := 0; i < t.depth; i++ {
		if k := len(after[i]); k > 0 {
			for j := 0; j < k; j++ {
				tr := after[i][j]
				t.doDebug(func() {
					t.Logf("defering %#v", tr)
				})
				defer tr.fn()
				if tr.Quantifier == First {
					popTriggerAt(j, &after[i])
					k--
					j--
				}
			}
		}
	}
	t.deferstack = after

	// Don't run the spec if we are not in a matching context.
	if !t.runspec {
		return
	}

	var (
		seq     sequence
		m       Matcher
		negated bool
		args    []interface{}
	)

	t.doDebug(func() {
		t.Logf("Executing")
	})

	t.ranspec = true
	seq, t.err = t.scan(spec)
	if t.err != nil {
		t.spec = seq
		return
	}
	m, negated, args, t.err = t.parse(seq)
	if t.err != nil {
		t.spec = seq
		return
	}
	t.passed, t.err = t.exec(m, negated, args)
	if !t.passed || t.err != nil {
		t.spec = seq
	}
}
