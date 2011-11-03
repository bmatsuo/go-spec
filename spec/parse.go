package spec
/*
 *  Filename:    parse.go
 *  Package:     spec
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Wed Nov  2 03:43:00 PDT 2011
 *  Description: <no value>
 */
import (
	"errors"
	"reflect"
)
//  Syntactic sugar for Spec sequences. See Spec.
type Sugar uint8

const Should Sugar = 0

func (s Sugar) String() string {
	if s == Should {
		return "Should"
	}
	panic("Bad Sugar")
}

//  Functions over Spec Values that can be used in Spec sequences. See Spec.
type Function uint

const (
	Equal Function = iota
	Satisfy
	HaveError
	Not
)

var fnNumArg = []int{
	Equal:     1,
	Satisfy:   1,
	HaveError: 0,
	Not:       1,
}

var fnString = []string{
	Equal:     "Equal",
	Satisfy:   "Satisfy",
	HaveError: "HaveError",
	Not:       "Not",
}

//  Return the number of arguments taken by Function fn. This is the number
//  of arguments following the function name in the Spec sequence.
func (fn Function) NumArg() int { return fnNumArg[fn] }

//  Return the Function's name as a string.
func (fn Function) String() string { return fnString[fn] }

//  A call to a nil-adic function. Output values can be accessed.
type FnCall struct {
	fn  reflect.Value
	out []reflect.Value
}

type token uint

const (
	tSugar token = iota
	tFunction
	tValue
)

var tokenStr = []string{
	tSugar:    "Sugar",
	tFunction: "Function",
	tValue:    "Value",
}

func (t token) String() string   { return tokenStr[t] }
func (t token) IsSugar() bool    { return t == tSugar }
func (t token) IsFunction() bool { return t == tFunction }
func (t token) IsValue() bool    { return t == tValue }

func tokenOf(v interface{}) (t token) {
	switch v.(type) {
	case Function:
		t = tFunction
	case Sugar:
		t = tSugar
	default:
		t = tValue
	}
	return
}

func tokenize(sequence []interface{}) []token {
	tokens := make([]token, len(sequence))
	for i := range sequence {
		tokens[i] = tokenOf(sequence[i])
	}
	return tokens
}

type kind uint

const (
	kNative kind = iota
	kFnCall
)

type elem struct {
	token
	kind
	value interface{}
}

type sequence []elem

var (
	ErrBadSugar           = errors.New("Bad Sugar")
	ErrBadFunction        = errors.New("Bad Value")
	ErrMissingSugar       = errors.New("Missing Value")
	ErrMissingFunction    = errors.New("Missing Value")
	ErrMissingValue       = errors.New("Missing Value")
	ErrUnexpectedSugar    = errors.New("Unexpected Sugar")
	ErrUnexpectedFunction = errors.New("Unexpected Function")
	ErrUnexpectedValue    = errors.New("Unexpected Value")
)

func (t *SpecTest) scan(values []interface{}) (seq sequence, err error) {
	if len(values) == 0 {
		return nil, ErrMissingValue
	}

	seq = make(sequence, len(values))
	for i, t := range tokenize(values) {
		seq[i] = elem{t, kNative, values[i]}
	}
	return
}

func (t *SpecTest) parse(seq sequence) (fn Function, negated bool, args []interface{}, err error) {
	var (
		i, k   int         // Index,   Gobbled
		v1, v2 interface{} // Object,  Argument
	)

	t.doDebug(func() { t.Log(specString(seq[i:])) })
	// Parse the object Value to be "spec'ed"
	k, v1, err = t.parseValue(seq)
	i += k
	if err != nil {
		return
	}

	t.doDebug(func() { t.Log(specString(seq[i:])) })
	// Look for Should separating value and function.
	switch seq[i].token {
	case tSugar:
		if seq[i].value.(Sugar) != Should {
			err = ErrBadSugar
		}
	case tFunction:
		err = ErrMissingSugar
	default:
		err = ErrUnexpectedValue
	}
	i++

	t.doDebug(func() { t.Log(specString(seq[i:])) })
	// Look for function negation, Not.
	negated = false
	switch seq[i].token {
	case tFunction:
		if seq[i].value.(Function) == Not {
			negated = true
			i++
		}
	default:
	}

	if negated {
		t.doDebug(func() { t.Log(specString(seq[i:])) })
	}
	// Look for Function.
	k, fn, err = t.parseFunction(seq[i:])
	i += k
	if err != nil {
		return
	}

	args = make([]interface{}, 1, 2)
	args[0] = v1
	if fn.NumArg() > 0 {
		// Look for a Function argument if necessary.
		if fn.NumArg() > 1 {
			err = errors.New("Needs too many arguments")
			return
		}
		t.doDebug(func() { t.Log(specString(seq[i:])) })
		k, v2, err = t.parseValue(seq[i:])
		i += k
		if err != nil {
			return
		}
		if i < len(seq) {
			err = errors.New("Excess specification pieces")
		}
		args = append(args, v2)
	}
	return
}
/*
func (t *SpecTest) parseSpec(specpieces []interface{}) (fn Function, negated bool, args []interface{}, err os.Error) {
    var (
        i, k   int          // Index,   Gobbled
        v1, v2 interface{}  // Object,  Argument
    )
    if len(specpieces) == 0 {
        err = os.NewError("No pieces")
        return
    }

    t.doDebug(func() { t.Log(specString(specpieces[i:])) })
    // Parse the object Value to be "spec'ed"
    k, v1, err = t.parseValue(specpieces)
    i += k
    if err != nil {
        return
    }

    t.doDebug(func() { t.Log(specString(specpieces[i:])) })
    // Look for Should separating value and function.
    switch specpieces[i].(type) {
    case Sugar:
        if specpieces[i].(Sugar) != Should {
            err = os.NewError("Bad Sugar")
        }
    case Function:
        err = os.NewError("Missing sugar")
    default:
        err = os.NewError("Unexpected value")
    }
    i++

    t.doDebug(func() { t.Log(specString(specpieces[i:])) })
    // Look for function negation, Not.
    negated = false
    switch specpieces[i].(type) {
    case Function:
        if specpieces[i].(Function) == Not {
            negated = true
            i++
        }
    default:
    }

    if negated {
        t.doDebug(func() { t.Log(specString(specpieces[i:])) })
    }
    // Look for Function.
    k, fn, err = t.parseFunction(specpieces[i:])
    i += k
    if err != nil {
        return
    }

    args = make([]interface{}, 1, 2)
    args[0] = v1
    if fn.NumArg() > 0 {
        // Look for a Function argument if necessary.
        if fn.NumArg() > 1 {
            err = os.NewError("Needs too many arguments")
            return
        }
        t.doDebug(func() { t.Log(specString(specpieces[i:])) })
        k, v2, err = t.parseValue(specpieces[i:])
        i += k
        if err != nil {
            return
        }
        if i < len(specpieces) {
            err = os.NewError("Excess specification pieces")
        }
        args = append(args, v2)
    }
    return
}
*/

func (t *SpecTest) parseValue(seq sequence) (i int, v interface{}, err error) {
	if len(seq) == 0 {
		panic("empty")
	}

	// Gobble pieces of the Value.
	var valpieces sequence
	done := false
	for !done && i < len(seq) {
		var piece elem = seq[i]
		switch piece.token {
		case tSugar:
			if i == 0 {
				err = ErrMissingValue
				return
			}
			// Piece Should terminates the value.
			switch piece.value.(Sugar) {
			case Should:
			default:
				err = ErrBadSugar
				return
			}
			i--
			done = true
		case tFunction:
			err = ErrMissingSugar
			return
		default:
			valpieces = append(valpieces, piece)
		}
		i++
	}

	// The object of the Spec Function.
	v = valpieces[0].value
	valpieces[0].kind = kNative

	// Evaluate nil-adic function calls.
	switch v.(type) {
	case FnCall:
	default:
		fntyp := reflect.TypeOf(v)
		if k := fntyp.Kind(); k != reflect.Func {
			break
		}
		if fntyp.NumIn() != 0 {
			break
		}
		if fntyp.NumOut() == 0 {
			err = errors.New("Value-less function")
			return
		}
		valpieces[0].kind = kFnCall
		fnval := reflect.ValueOf(v)
		v = FnCall{fnval, fnval.Call([]reflect.Value{})}
	}

	// Return when no indexing values are given.
	if len(valpieces) <= 1 {
		return
	}

	if len(valpieces) > 2 {
		err = errors.New("Too many Value pieces")
		return
	}

	// Create a slice of reflect.Values to index
	var vals []reflect.Value
	switch v.(type) {
	case FnCall:
		vals = v.(FnCall).out
	default:
		// Can't index the Value.
		err = errors.New("Too many Value pieces")
		return
	}

	// Index the Value to create another Value.
	index := valpieces[1]
	switch index.value.(type) {
	case int:
		j := index.value.(int)
		if j >= len(vals) {
			err = errors.New("Index out of range")
			return
		}
		v = vals[j].Interface()
	default:
		err = errors.New("Missing 'int' Value index")
	}
	return
}

func (t *SpecTest) parseFunction(seq sequence) (i int, fn Function, err error) {
	if len(seq) == 0 {
		err = errors.New("empty")
		return
	}

	switch piece := seq[0]; piece.token {
	case tFunction:
		if fn = piece.value.(Function); fn == Not {
			// Function Not is handled before.
			err = errors.New("Double negative")
		}
	case tSugar:
		err = errors.New("Unexpected Sugar")
	default:
		err = errors.New("Missing Function")
	}
	i++
	return
}
