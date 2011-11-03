package spec
/*
 *  Filename:    exec.go
 *  Package:     spec
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Nov  2 03:47:12 PDT 2011
 *  Description: <no value>
 */
import (
	"errors"
	"reflect"
)

func (t *SpecTest) exec(m Matcher, negated bool, args []interface{}) {
	if len(args) < 1 {
		// Serious error
		t.Errorf("Spec error: Missing Value")
		return
	}
	if negated {
		defer func() { t.passed = !t.passed }()
	}

	if n := len(args); n < m.NumIn() {
		t.Errorf("Spec error: Missing argument")
		return
	} else if n > m.NumIn() {
		t.Errorf("Spec error: Unexpected arguments %v", args[n:])
		return
	}
    t.passed, t.err = m.Matches(args)
    if t.err == nil {
        t.err = m.Error()
    }
    if t.err != nil {
        panic(t.err)
    }
    // Matcher errors messages are handled in Describe
}

func (t *SpecTest) equal(a, b interface{}) (bool, error) {
	switch a.(type) {
	case FnCall:
		a = a.(FnCall).out[0].Interface()
	}
	switch b.(type) {
	case FnCall:
		b = b.(FnCall).out[0].Interface()
	}

	t.doDebug(func() { t.Logf("%#v = %#v", a, b) })
	return reflect.DeepEqual(a, b), nil
}

func (t *SpecTest) satisfy(x, fn interface{}) (satisfies bool, err error) {
	t.doDebug(func() { t.Logf("%#v satisfies function %#v", x, fn) })
	switch x.(type) {
	case FnCall:
		x = x.(FnCall).out[0].Interface()
	default:
	}

	fnval := reflect.ValueOf(fn)

	// Check the type of fn (fn(x)bool).
	if k := fnval.Kind(); k != reflect.Func {
		return false, errors.New("Satisfy given non-function")
	}
	if typ := fnval.Type(); typ.NumIn() != 1 {
		return false, errors.New("Satisfy needs a function of one argument")
	} else if xtyp := reflect.TypeOf(x); !xtyp.AssignableTo(typ.In(0)) {
		return false, errors.New("Satisfy argument type-mismatch")
	} else if typ.NumOut() != 1 {
		return false, errors.New("Satisfy needs a predicate (func(x) bool)")
	} else if !typ.Out(0).AssignableTo(reflect.TypeOf(satisfies)) {
		return false, errors.New("Satisfy output type-mismatch")
	}
	fnout := fnval.Call([]reflect.Value{reflect.ValueOf(x)})
	satisfies = fnout[0].Interface().(bool)
	return
}

func (t *SpecTest) haveError(fn interface{}) (bool, error) {
	t.doDebug(func() { t.Logf("Function %#v has an error", fn) })

	var errval reflect.Value
	switch fn.(type) {
	case FnCall:
		fncall := fn.(FnCall)
		fntyp := fncall.fn.Type()
		if !reflect.TypeOf(errors.New("error")).AssignableTo(fntyp.Out(fntyp.NumOut() - 1)) {
			return false, errors.New("HaveError function call's last Value must be os.Error")
		}
		errval = fncall.out[len(fncall.out)-1]
	default:
		return false, errors.New("HaveError needs a function call Value")
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
	return fnerr != nil, nil
}
