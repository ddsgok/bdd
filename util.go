package bdd

import (
	"fmt"
	"runtime"
	"strings"
	"github.com/pkg/errors"
)

var (
	// ErrInvalidPlaceForLike received when user puts like sentence at
	// incorrect position.
	ErrInvalidPlaceForLike = errors.New("the like sentence must be called as last argument")
	// ErrWrongNumTestFuncs received when user puts more than one
	// function test on sentences.
	ErrWrongNumTestFuncs = errors.New("there's more than one func being received to test")
)

// printf is a clearer version of fmt.Sprintf.
func printf(s string, args []interface{}) (f string) {
	f = fmt.Sprintf(s, args...)
	return
}

// split args received into the testbodies received if any, and the
// set of test args. This is needed for the idea of how the sentences
// will be used. There are some possibilities:
//
// 1 - No arguments. Tests not implemented.
// 	when("a function is called")
//
// 2 - 1 argument. Like set of arguments, various tests not implemented.
// 	when("a function is called", like(s(1, 2, 3), s(2, 4, 6)))
//
// 3 - 1 argument. One simple test, only one execution for test.
// 	when("a function is called", func(it bdd.It){ /*...*/ })
//
// 4 - 2 argument. Simple tests and at last a Like set of arguments.
// On this case, the test body will be called n times, where n is the
// len of like set.
// 	when("a function is called", func(it bdd.It){ /*...*/ },
// 		like(s(1, 2, 3), s(2, 4, 6)))
func split(init Arguments, args []interface{}) (testbody testFunc, like []Arguments) {
	like = []Arguments{init}

	switch len(args) {
	case 0: // 1º poss.
		break
	case 1:
		switch args[0].(type) {
		case []Arguments: // 2º poss.
			like = args[0].([]Arguments)
		default: // 3º poss.
			testbody = newTestFunc(args[0])
		}
	default: // 4º poss.
		if _, ok := args[0].([]Arguments); ok {
			panic(ErrInvalidPlaceForLike)
		}

		if len(args) > 2 {
			panic(ErrWrongNumTestFuncs)
		}

		testbody = newTestFunc(args[0])
		like = args[1].([]Arguments)
	}

	return
}

// notImplemented is used to mark a specification that needs coding out.
func notImplemented() (fn func(Assert)) {
	fn = func(assert Assert) {
		// nothing to do here
	}
	return
}

// feature return test name, parsed to a phrase, removing Test and _ strings.
func feature() (r string) {
	pc, _, _, _ := runtime.Caller(2)
	m := fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	i := strings.LastIndex(m, ".")
	m = m[i+1:]
	m = strings.Replace(m, "Test_", "", 1)
	m = strings.Replace(m, "Test", "", 1)
	r = strings.Replace(m, "_", " ", -1)
	return
}
