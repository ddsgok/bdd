package assert

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/ddspog/bdd/internal/shared"
)

var (
	// zeros store all different types of zeros.
	zeros = []interface{}{
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		float32(0),
		float64(0),
	}
)

// Comparison a custom function that returns true on success and false on failure
type Comparison func() bool

// PanicTestFunc defines a func that should be passed to the assert.Panics and assert.NotPanics
// methods, and represents a simple func that takes no arguments, and returns nothing.
type PanicTestFunc func()

/*
	Helper functions
*/

// ObjectsAreEqual determines if two objects are considered equal.
//
// This function does no assertion of any kind.
func ObjectsAreEqual(expected, actual interface{}) (b bool) {
	b = true

	if expected == nil || actual == nil {
		b = expected == actual
		return
	}

	if reflect.DeepEqual(expected, actual) {
		return
	}

	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)
	if expectedValue == actualValue {
		return
	}

	// Attempt comparison after type conversion
	if actualValue.Type().ConvertibleTo(expectedValue.Type()) && expectedValue == actualValue.Convert(expectedValue.Type()) {
		return
	}

	// Last ditch effort
	if fmt.Sprintf("%#v", expected) == fmt.Sprintf("%#v", actual) {
		return
	}

	b = false
	return
}

/* CallerInfo is necessary because the assert functions use the testing object
internally, causing it to print the file:line of the assert method, rather than where
the problem actually occurred in calling code.*/

// CallerInfo returns a string containing the file and line number of the assert call
// that failed.
func CallerInfo() (c string) {
	file := ""
	line := 0
	ok := false

	for i := 0; ; i++ {
		_, file, line, ok = runtime.Caller(i)
		if !ok {
			return
		}
		parts := strings.Split(file, "/")
		dir := parts[len(parts)-2]
		file = parts[len(parts)-1]
		if (dir != "assert" && dir != "mock" && dir != "require") || file == "mock_test.go" {
			break
		}
	}

	c = fmt.Sprintf("%s:%d", file, line)
	return
}

// getWhitespaceString returns a string that is long enough to overwrite the default
// output from the go testing framework.
func getWhitespaceString() (wp string) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	wp = strings.Repeat(" ", len(fmt.Sprintf("%s:%d:      ", file, line)))
	return
}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) (msg string) {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return
	}

	if len(msgAndArgs) == 1 {
		msg = msgAndArgs[0].(string)
		return
	}

	if len(msgAndArgs) > 1 {
		msg = fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
		return
	}

	return
}

// indentMessageLines Indents all lines of the message by appending
// a number of tabs to each line, in an output format compatible with
// Go's test printing (see inner comment for specifics).
func indentMessageLines(message string, tabs int) (r string) {
	outBuf := new(bytes.Buffer)

	for i, scanner := 0, bufio.NewScanner(strings.NewReader(message)); scanner.Scan(); i++ {
		if i != 0 {
			outBuf.WriteRune('\n')
		}

		for ii := 0; ii < tabs; ii++ {
			outBuf.WriteRune('\t')
			// Bizarrely, all lines except the first need one fewer tabs prepended, so deliberately advance the counter
			// by 1 prematurely.
			if ii == 0 && i > 0 {
				ii++
			}
		}

		outBuf.WriteString(scanner.Text())
	}

	r = outBuf.String()
	return
}

// Fail reports a failure through
func Fail(t shared.Tester, failureMessage string, msgAndArgs ...interface{}) (b bool) {
	message := messageFromMsgAndArgs(msgAndArgs...)

	if len(message) > 0 {
		t.Errorf("\r%s\r\tLocation:\t%s\n"+
			"\r\tError:%s\n"+
			"\r\tMessages:\t%s\n\r",
			getWhitespaceString(),
			CallerInfo(),
			indentMessageLines(failureMessage, 2),
			message)
	} else {
		t.Errorf("\r%s\r\tLocation:\t%s\n"+
			"\r\tError:%s\n\r",
			getWhitespaceString(),
			CallerInfo(),
			indentMessageLines(failureMessage, 2))
	}

	return
}

// Implements asserts that an object is implemented by the specified interface.
//
//    assert.Implements(t, (*MyInterface)(nil), new(MyObject), "MyObject")
func Implements(t shared.Tester, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) (b bool) {
	interfaceType := reflect.TypeOf(interfaceObject).Elem()

	if !reflect.TypeOf(object).Implements(interfaceType) {
		b = Fail(t, fmt.Sprintf("Object must implement %v", interfaceType), msgAndArgs...)
		return
	}

	b = true
	return
}

// IsType asserts that the specified objects are of the same type.
func IsType(t shared.Tester, expectedType interface{}, object interface{}, msgAndArgs ...interface{}) (b bool) {
	if !ObjectsAreEqual(reflect.TypeOf(object), reflect.TypeOf(expectedType)) {
		b = Fail(t, fmt.Sprintf("Object expected to be of type %v, but was %v", reflect.TypeOf(expectedType), reflect.TypeOf(object)), msgAndArgs...)
		return
	}

	b = true
	return
}

// Equal asserts that two objects are equal.
//
//    assert.Equal(t, 123, 123, "123 and 123 should be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func Equal(t shared.Tester, expected, actual interface{}, msgAndArgs ...interface{}) (b bool) {
	if !ObjectsAreEqual(expected, actual) {
		b = Fail(t, fmt.Sprintf("Not equal: %#v (expected)\n"+
			"        != %#v (actual)", expected, actual), msgAndArgs...)
		return
	}

	b = true
	return
}

// Exactly asserts that two objects are equal is value and type.
//
//    assert.Exactly(t, int32(123), int64(123), "123 and 123 should NOT be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func Exactly(t shared.Tester, expected, actual interface{}, msgAndArgs ...interface{}) (b bool) {
	if aType, bType := reflect.TypeOf(expected), reflect.TypeOf(actual); aType != bType {
		b = Fail(t, "Types expected to match exactly", "%v != %v", aType, bType)
	} else {
		b = Equal(t, expected, actual, msgAndArgs...)
	}

	return
}

// NotNil asserts that the specified object is not nil.
//
//    assert.NotNil(t, err, "err should be something")
//
// Returns whether the assertion was successful (true) or not (false).
func NotNil(t shared.Tester, object interface{}, msgAndArgs ...interface{}) (success bool) {
	success = true

	if object == nil {
		success = false
	} else {
		value := reflect.ValueOf(object)
		kind := value.Kind()
		if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
			success = false
		}
	}

	if !success {
		Fail(t, "Expected not to be nil.", msgAndArgs...)
	}

	return
}

// isNil checks if a specified object is nil or not, without Failing.
func isNil(object interface{}) (b bool) {
	if object == nil {
		b = true
	} else {
		value := reflect.ValueOf(object)
		kind := value.Kind()
		if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
			b = true
		}
	}

	return
}

// Nil asserts that the specified object is nil.
//
//    assert.Nil(t, err, "err should be nothing")
//
// Returns whether the assertion was successful (true) or not (false).
func Nil(t shared.Tester, object interface{}, msgAndArgs ...interface{}) (b bool) {
	if isNil(object) {
		b = true
	} else {
		b = Fail(t, fmt.Sprintf("Expected nil, but got: %#v", object), msgAndArgs...)
	}

	return
}

// isEmpty gets whether the specified object is considered empty or not.
func isEmpty(object interface{}) bool {

	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	for _, v := range zeros {
		if object == v {
			return true
		}
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Map:
		fallthrough
	case reflect.Slice, reflect.Chan:
		{
			return objValue.Len() == 0
		}
	case reflect.Ptr:
		{
			switch object.(type) {
			case *time.Time:
				return object.(*time.Time).IsZero()
			default:
				return false
			}
		}
	}
	return false
}

// Empty asserts that the specified object is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
// assert.Empty(t, obj)
//
// Returns whether the assertion was successful (true) or not (false).
func Empty(t shared.Tester, object interface{}, msgAndArgs ...interface{}) bool {

	pass := isEmpty(object)
	if !pass {
		Fail(t, fmt.Sprintf("Should be empty, but was %v", object), msgAndArgs...)
	}

	return pass

}

// NotEmpty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
// if assert.NotEmpty(t, obj) {
//   assert.Equal(t, "two", obj[1])
// }
//
// Returns whether the assertion was successful (true) or not (false).
func NotEmpty(t shared.Tester, object interface{}, msgAndArgs ...interface{}) bool {

	pass := !isEmpty(object)
	if !pass {
		Fail(t, fmt.Sprintf("Should NOT be empty, but was %v", object), msgAndArgs...)
	}

	return pass

}

// getLen try to get length of object.
// return (false, 0) if impossible.
func getLen(x interface{}) (ok bool, length int) {
	v := reflect.ValueOf(x)
	defer func() {
		if e := recover(); e != nil {
			ok = false
		}
	}()
	return true, v.Len()
}

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
//
//    assert.Len(t, mySlice, 3, "The size of slice is not 3")
//
// Returns whether the assertion was successful (true) or not (false).
func Len(t shared.Tester, object interface{}, length int, msgAndArgs ...interface{}) bool {
	ok, l := getLen(object)
	if !ok {
		return Fail(t, fmt.Sprintf("\"%s\" could not be applied builtin len()", object), msgAndArgs...)
	}

	if l != length {
		return Fail(t, fmt.Sprintf("\"%s\" should have %d item(s), but have %d", object, length, l), msgAndArgs...)
	}
	return true
}

// True asserts that the specified value is true.
//
//    assert.True(t, myBool, "myBool should be true")
//
// Returns whether the assertion was successful (true) or not (false).
func True(t shared.Tester, value bool, msgAndArgs ...interface{}) bool {

	if value != true {
		return Fail(t, "Should be true", msgAndArgs...)
	}

	return true

}

// False asserts that the specified value is true.
//
//    assert.False(t, myBool, "myBool should be false")
//
// Returns whether the assertion was successful (true) or not (false).
func False(t shared.Tester, value bool, msgAndArgs ...interface{}) bool {

	if value != false {
		return Fail(t, "Should be false", msgAndArgs...)
	}

	return true

}

// NotEqual asserts that the specified values are NOT equal.
//
//    assert.NotEqual(t, obj1, obj2, "two objects shouldn't be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func NotEqual(t shared.Tester, expected, actual interface{}, msgAndArgs ...interface{}) bool {

	if ObjectsAreEqual(expected, actual) {
		return Fail(t, "Should not be equal", msgAndArgs...)
	}

	return true

}

// Contains asserts that the specified string contains the specified substring.
//
//    assert.Contains(t, "Hello World", "World", "But 'Hello World' does contain 'World'")
//
// Returns whether the assertion was successful (true) or not (false).
func Contains(t shared.Tester, s, contains string, msgAndArgs ...interface{}) bool {

	if !strings.Contains(s, contains) {
		return Fail(t, fmt.Sprintf("\"%s\" does not contain \"%s\"", s, contains), msgAndArgs...)
	}

	return true

}

// NotContains asserts that the specified string does NOT contain the specified substring.
//
//    assert.NotContains(t, "Hello World", "Earth", "But 'Hello World' does NOT contain 'Earth'")
//
// Returns whether the assertion was successful (true) or not (false).
func NotContains(t shared.Tester, s, contains string, msgAndArgs ...interface{}) bool {

	if strings.Contains(s, contains) {
		return Fail(t, fmt.Sprintf("\"%s\" should not contain \"%s\"", s, contains), msgAndArgs...)
	}

	return true

}

// Condition uses a Comparison to assert a complex condition.
func Condition(t shared.Tester, comp Comparison, msgAndArgs ...interface{}) bool {
	result := comp()
	if !result {
		Fail(t, "Condition failed!", msgAndArgs...)
	}
	return result
}

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f PanicTestFunc) (bool, interface{}) {

	didPanic := false
	var message interface{}
	func() {

		defer func() {
			if message = recover(); message != nil {
				didPanic = true
			}
		}()

		// call the target function
		f()

	}()

	return didPanic, message

}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panics(t, func(){
//     GoCrazy()
//   }, "Calling GoCrazy() should panic")
//
// Returns whether the assertion was successful (true) or not (false).
func Panics(t shared.Tester, f PanicTestFunc, msgAndArgs ...interface{}) bool {

	if funcDidPanic, panicValue := didPanic(f); !funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should panic\n\r\tPanic value:\t%v", f, panicValue), msgAndArgs...)
	}

	return true
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   assert.NotPanics(t, func(){
//     RemainCalm()
//   }, "Calling RemainCalm() should NOT panic")
//
// Returns whether the assertion was successful (true) or not (false).
func NotPanics(t shared.Tester, f PanicTestFunc, msgAndArgs ...interface{}) bool {

	if funcDidPanic, panicValue := didPanic(f); funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should not panic\n\r\tPanic value:\t%v", f, panicValue), msgAndArgs...)
	}

	return true
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
//   assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second, "The difference should not be more than 10s")
//
// Returns whether the assertion was successful (true) or not (false).
func WithinDuration(t shared.Tester, expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) bool {

	dt := expected.Sub(actual)
	if dt < -delta || dt > delta {
		return Fail(t, fmt.Sprintf("Max difference between %v and %v allowed is %v, but difference was %v", expected, actual, delta, dt), msgAndArgs...)
	}

	return true
}

// toFloat converts arg to float64 value.
func toFloat(x interface{}) (float64, bool) {
	var xf float64
	xok := true

	switch xn := x.(type) {
	case uint8:
		xf = float64(xn)
	case uint16:
		xf = float64(xn)
	case uint32:
		xf = float64(xn)
	case uint64:
		xf = float64(xn)
	case int:
		xf = float64(xn)
	case int8:
		xf = float64(xn)
	case int16:
		xf = float64(xn)
	case int32:
		xf = float64(xn)
	case int64:
		xf = float64(xn)
	case float32:
		xf = float64(xn)
	case float64:
		xf = float64(xn)
	default:
		xok = false
	}

	return xf, xok
}

// InDelta asserts that the two numerals are within delta of each other.
//
// 	 assert.InDelta(t, math.Pi, (22 / 7.0), 0.01)
//
// Returns whether the assertion was successful (true) or not (false).
func InDelta(t shared.Tester, expected, actual interface{}, delta float64, msgAndArgs ...interface{}) bool {

	af, aok := toFloat(expected)
	bf, bok := toFloat(actual)

	if !aok || !bok {
		return Fail(t, fmt.Sprintf("Parameters must be numerical"), msgAndArgs...)
	}

	dt := af - bf
	if dt < -delta || dt > delta {
		return Fail(t, fmt.Sprintf("Max difference between %v and %v allowed is %v, but difference was %v", expected, actual, delta, dt), msgAndArgs...)
	}

	return true
}

// min(|expected|, |actual|) * epsilon
func calcEpsilonDelta(expected, actual interface{}, epsilon float64) float64 {
	af, aok := toFloat(expected)
	bf, bok := toFloat(actual)

	if !aok || !bok {
		// invalid input
		return 0
	}

	if af < 0 {
		af = -af
	}
	if bf < 0 {
		bf = -bf
	}
	var delta float64
	if af < bf {
		delta = af * epsilon
	} else {
		delta = bf * epsilon
	}
	return delta
}

// InEpsilon asserts that expected and actual have a relative error less than epsilon
//
// Returns whether the assertion was successful (true) or not (false).
func InEpsilon(t shared.Tester, expected, actual interface{}, epsilon float64, msgAndArgs ...interface{}) bool {
	delta := calcEpsilonDelta(expected, actual, epsilon)

	return InDelta(t, expected, actual, delta, msgAndArgs...)
}

/*
	Errors
*/

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.NoError(t, err) {
//	   assert.Equal(t, actualObj, expectedObj)
//   }
//
// Returns whether the assertion was successful (true) or not (false).
func NoError(t shared.Tester, err error, msgAndArgs ...interface{}) bool {
	if isNil(err) {
		return true
	}

	return Fail(t, fmt.Sprintf("No error is expected but got %v", err), msgAndArgs...)
}

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Error(t, err, "An error was expected") {
//	   assert.Equal(t, err, expectedError)
//   }
//
// Returns whether the assertion was successful (true) or not (false).
func Error(t shared.Tester, err error, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)
	return NotNil(t, err, "An error is expected but got nil. %s", message)

}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   if assert.Error(t, err, "An error was expected") {
//	   assert.Equal(t, err, expectedError)
//   }
//
// Returns whether the assertion was successful (true) or not (false).
func EqualError(t shared.Tester, theError error, errString string, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)
	if !NotNil(t, theError, "An error is expected but got nil. %s", message) {
		return false
	}
	s := "An error with value \"%s\" is expected but got \"%s\". %s"
	return Equal(t, theError.Error(), errString,
		s, errString, theError.Error(), message)
}
