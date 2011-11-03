package spec
/*
 *  Filename:    exec.go
 *  Package:     spec
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Nov  2 03:47:12 PDT 2011
 *  Description: <no value>
 */
import (
    "reflect"
    "os"
)

func (t *SpecTest) exec(fn Function, negated bool, args []interface{}) (pass bool, err os.Error) {
    if len(args) < 1 {
        // Serious error
        t.Errorf("Spec error: Missing Value")
        return
    }
    if negated {
        defer func() { pass = !pass }()
    }
    val := args[0]

    if n:= len(args)-1; n < fn.NumArg() {
        t.Errorf("Spec error: Missing argument")
        return
    } else if n > fn.NumArg() {
        t.Errorf("Spec error: Unexpected arguments %v", args[n:])
        return
    }

    // Handle the Function fn.
    switch fn {
    case Equal:
        pass, err = t.equal(val, args[1])
    case Satisfy:
        pass, err = t.satisfy(val, args[1])
    case HaveError:
        pass, err = t.haveError(val)
    }
    return
}

func (t *SpecTest) equal(a, b interface{}) (bool, os.Error) {
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

func (t *SpecTest) satisfy(x, fn interface{}) (satisfies bool, err os.Error) {
    t.doDebug(func() { t.Logf("%#v satisfies function %#v", x, fn) })
    switch x.(type) {
    case FnCall:
        x = x.(FnCall).out[0].Interface()
    default:
    }

    fnval := reflect.ValueOf(fn)

    // Check the type of fn (fn(x)bool).
    if k := fnval.Kind(); k != reflect.Func {
        return false, os.NewError("Satisfy given non-function")
    }
    if typ := fnval.Type(); typ.NumIn() != 1 {
        return false, os.NewError("Satisfy needs a function of one argument")
    } else if xtyp := reflect.TypeOf(x); !xtyp.AssignableTo(typ.In(0)) {
        return false, os.NewError("Satisfy argument type-mismatch")
    } else if typ.NumOut() != 1 {
        return false, os.NewError("Satisfy needs a predicate (func(x) bool)")
    } else if !typ.Out(0).AssignableTo(reflect.TypeOf(satisfies)) {
        return false, os.NewError("Satisfy output type-mismatch")
    }
    fnout := fnval.Call([]reflect.Value{reflect.ValueOf(x)})
    satisfies = fnout[0].Interface().(bool)
    return
}

func (t *SpecTest) haveError(fn interface{}) (bool, os.Error) {
    t.doDebug(func() { t.Logf("Function %#v has an error", fn) })

    var errval reflect.Value
    switch fn.(type) {
    case FnCall:
        fncall := fn.(FnCall)
        fntyp := fncall.fn.Type()
        if !reflect.TypeOf(os.NewError("error")).AssignableTo(fntyp.Out(fntyp.NumOut() - 1)) {
            return false, os.NewError("HaveError function call's last Value must be os.Error")
        }
        errval = fncall.out[len(fncall.out)-1]
    default:
        return false, os.NewError("HaveError needs a function call Value")
    }

    var fnerr os.Error
    switch v := errval.Interface(); v.(type) {
    case nil:
        fnerr = nil
    case os.Error:
        fnerr = v.(os.Error)
    default:
        return false, os.NewError("Function call can not have error")
    }
    return fnerr != nil, nil
}
