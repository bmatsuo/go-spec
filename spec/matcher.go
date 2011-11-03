// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spec

/*  Filename:    matcher.go
 *  Author:      Bryan Matsuo <bryan.matsuo@gmail.com>
 *  Created:     Thu Nov  3 05:19:41 PDT 2011
 *  Description: 
 */

import (
	"reflect"
	"errors"
	"fmt"
)

var boolType = reflect.TypeOf(true)
var errorType = reflect.TypeOf(errors.New("ERROR"))

//  The default set of Spec matchers.
var (
	MEqual     = MatcherMust(NewMatcher(matcherEqual))
	MSatisfy   = MatcherMust(NewMatcher(matcherSatisfy))
	MHaveError = MatcherMust(NewMatcher(matcherHaveError))
    MPanic     = MatcherMust(NewMatcher(matcherPanic))
)

//  If x is not a FnCall, return x. Otherwise, return the first return value
//  of x.
func valueOfSpecValue(x interface{}) (y interface{}) {
	switch x.(type) {
	case FnCall:
		y = x.(FnCall).out[0].Interface()
    default:
        y = x
	}
    return
}

//  Use reflect.DeepEqual to test two values' equality.
func matcherEqual(a, b interface{}) (pass bool, err error) {
	//t.doDebug(func() { t.Logf("%#v = %#v", a, b) })
	pass = reflect.DeepEqual(
        valueOfSpecValue(a),
        valueOfSpecValue(b))
    return
}

func matcherSatisfy(x, fn interface{}) (pass bool, err error) {
	//t.doDebug(func() { t.Logf("%#v satisfies function %#v", x, fn) })
	fnval := reflect.ValueOf(fn)
	// Check the type of fn (fn(x)bool).
	if k := fnval.Kind(); k != reflect.Func {
		return false, errors.New("Satisfy given non-function")
	}
    x = valueOfSpecValue(x)
	if typ := fnval.Type(); typ.NumIn() != 1 {
        // error: fn must accept a single value.
		return false, errors.New("Satisfy needs a function of one argument")
	} else if xtyp := reflect.TypeOf(x); !xtyp.AssignableTo(typ.In(0)) {
        // error: fn must accept x
		return false, errors.New("Satisfy argument type-mismatch")
	} else if typ.NumOut() != 1 {
        // error: fn must return one value
		return false, errors.New("Satisfy needs a predicate (func(x) bool)")
	} else if !typ.Out(0).AssignableTo(reflect.TypeOf(boolType)) {
        // error: fn must return bool
		return false, errors.New("Satisfy output type-mismatch")
	}
	fnout := fnval.Call([]reflect.Value{reflect.ValueOf(x)})
	pass = fnout[0].Bool()
	return
}

func matcherHaveError(fn interface{}) (pass bool, err error) {
	//t.doDebug(func() { t.Logf("Function %#v has an error", fn) })

	var errval reflect.Value
	switch fn.(type) {
	case FnCall:
		fncall := fn.(FnCall)
		fntyp := fncall.fn.Type()
		if !errorType.AssignableTo(fntyp.Out(fntyp.NumOut() - 1)) {
            // error: fn's last return value error.
			return false, errors.New("HaveError function call's last Value must be os.Error")
		}
		errval = fncall.out[len(fncall.out)-1]
	default:
        err = errors.New("HaveError needs a function call Value")
		return
	}

	var fnerr error
	switch v := errval.Interface(); v.(type) {
	case nil:
		fnerr = nil
	case error:
		fnerr = v.(error)
	default:
		return false, errors.New("Function call can not have error")
	}
    pass = fnerr != nil
	return
}

func matcherPanic(fn interface{}) (pass bool, err error) {
	//t.doDebug(func() { t.Logf("Function %#v has an error", fn) })

	switch fn.(type) {
	case FnCall:
		fncall := fn.(FnCall)
		pass = fncall.panicv != nil
		if pass {
			return
		}
	default:
        err = errors.New("HaveError needs a function call Value")
	}
	return
}

type Matcher interface {
	// Run the matcher against arguments.
	Matches(args ...interface{}) bool
	// Return any error encountered executing the matcher.
	Error() error
	// Return the number of arguments needed by the matcher
	NumIn() int
}

type match struct {
	fn  reflect.Value // A bool function of at least one argument
	typ reflect.Type  // A type with kind reflect.Func
	err error         // An error encountered when running err (panic / bug)
}

type errpanic struct {
	v interface{}
}

func (ep errpanic) Error() string {
	return fmt.Sprintf("runtime panic: %v", ep.v)
}

func (m match) exec(args []reflect.Value) (pass bool) {
	defer func() {
		if e := recover(); e != nil {
			m.err = errpanic{e}
		}
	}()
	out := m.fn.Call(args)
	pass, m.err = out[0].Bool(), out[1].Interface().(error)
	return
}

func (m match) Matches(args ...interface{}) bool {
	n := len(args)
	// Check the arguments.
	if n != m.NumIn() {
		m.err = errors.New("wrong number of arguments")
	}
	// Turn interfaces into reflect.Values and call the matcher.
	vals := make([]reflect.Value, n)
	for i := range args {
		vals[i] = reflect.ValueOf(args[i])
	}
	return m.exec(vals)
}
func (m match) Error() error { return m.err }
func (m match) NumIn() int   { return m.typ.NumIn() }

//  Create a new Matcher object from function fn. Function fn must take
//  at least one argument and return exactly one bool.
func NewMatcher(fn interface{}) (Matcher, error) {
	m := new(match)
	m.fn = reflect.ValueOf(fn)
	m.typ = m.fn.Type()

    // Check the kind of fn.
	if k := m.typ.Kind(); k != reflect.Func {
		return m, errors.New("matcher not a function")
	}

    // Check the number of inputs on fn
	if numin := m.typ.NumIn(); numin == 0 {
		return m, errors.New("nil-adic matcher")
	}

    // Check the number of outputs on fn
	if numout := m.typ.NumOut(); numout < 2 {
		return m, errors.New("not enough matcher return values")
	} else if numout > 2 {
		return m, errors.New("too many matcher return values")
	}

    // Check the types of fn's outputs
	if bout := m.typ.Out(0); !bout.AssignableTo(boolType) {
		return m, errors.New("matcher with non-bool return")
	}
	if !errorType.AssignableTo(m.typ.Out(1)) {
		return m, errors.New("matcher with non-error second return value")
	}

	return m, nil
}

func MatcherMust(m Matcher, err error) Matcher {
	if err != nil {
		panic(err)
	}
	return m
}
