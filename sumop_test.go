package bdd

import (
	"fmt"
)

// TestInt is a named type after int for test purposes.
type TestInt int

// Sum is a function to sum a TestInt to int of TestInt types.
func (t TestInt) Sum(a interface{}) (n TestInt) {
	switch val := a.(type) {
	case int:
		n = t + TestInt(val)
	case TestInt:
		n = t + val
	default:
		n = TestInt(0)
	}
	return
}

// TestSumOp is a test structure to suite test a struct.
type TestSumOp struct {
	LastResultAsString string  `json:"last_result"`
	Handicap           TestInt `json:"handicap"`
}

// Sum is a test method to suite test a struct.
func (s *TestSumOp) Sum(a, b int) (x int) {
	x = int(TestInt(a).Sum(b).Sum(s.Handicap))
	s.LastResultAsString = fmt.Sprintf("%v", x)
	return
}

// NewTestSumOp creates a new TestSumOp with a handicap h.
func NewTestSumOp(h int) (t *TestSumOp) {
	t = &TestSumOp{
		Handicap: TestInt(h),
	}
	return
}