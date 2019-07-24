package bdd

import (
	"github.com/ddsgok/bdd/internal/common"
)

// Arguments defines a set of arguments, to run on Given, When or It sentences.
type Arguments []interface{}

// When defines the action or event when Given a specific context.
type When func(when string, args ...interface{})

// It defines the specification of When something happens.
type It func(title string, args ...interface{})

// Assert defines the action of asserting things during test.
type Assert = common.Assert

// Golden defines an object to access test input and output through
// various test cases.
type Golden = common.Golden

// S return a new set of arguments, given on the function.
func S(args ...interface{}) (s Arguments) {
	s = args
	return
}

// Like defines a set of environments to be run on a sentence like
// Given, When and It. It receives a list of sets of arguments, and
// those arguments will be used to conduct table-driven tests using
// this BDD framework.
func Like(sets ...Arguments) (sa []Arguments) {
	sa = sets
	return
}
