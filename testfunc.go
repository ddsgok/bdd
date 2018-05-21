package bdd

// testFunc abstract an argument that should represent function
// received, as test function.
type testFunc struct {
	fn interface{}
}

// asWhenFunc return test function as When function.
func (tb testFunc) asWhenFunc() (wfn func(When, ...interface{})) {
	if tb.fn != nil {
		switch v := tb.fn.(type) {
		// When receiving a function without args, transform it.
		case func(When):
			wfn = func(wh When, args ...interface{}) {
				v(wh)
			}
		default:
			wfn = v.(func(When, ...interface{}))
		}
	}

	return
}

// asGoldenFunc return test function as Golden function.
func (tb testFunc) asGoldenFunc() (gfn func(When, Golden)) {
	if tb.fn != nil {
		gfn = tb.fn.(func(When, Golden))
	}

	return
}

// asItFuncs return test function as It function.
func (tb testFunc) asItFuncs() (ifn func(It, ...interface{})) {
	if tb.fn != nil {
		switch v := tb.fn.(type) {
		// When receiving a function without args, transform it.
		case func(It):
			ifn = func(it It, args ...interface{}) {
				v(it)
			}
		default:
			ifn = v.(func(It, ...interface{}))
		}
	}

	return
}

// asAssertFunc return test function as Assert function.
func (tb testFunc) asAssertFunc() (afn func(Assert, ...interface{})) {
	if tb.fn != nil {
		switch v := tb.fn.(type) {
		// When receiving a function without args, transform it.
		case func(Assert):
			afn = func(as Assert, args ...interface{}) {
				v(as)
			}
		default:
			afn = v.(func(Assert, ...interface{}))
		}
	}

	return
}

// newTestFunc creates a test func using arguments. Will check for
// errors.
func newTestFunc(args ...interface{}) (tb testFunc) {
	if len(args) > 1 {
		panic(ErrWrongNumTestFuncs)
	}

	tb = testFunc{args[0]}
	return
}
