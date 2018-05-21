package bdd

import "testing"

var (
	// sentences will store the management functions for sentences on
	// this package.
	sentences = newSentencesManager()
)

// SentencesManager allows the user of this bdd package, to choose
// which style of tests to use.
type SentencesManager interface {
	Given() func(*testing.T, string, ...interface{})
	Golden() func(*testing.T, string, ...interface{})
	All() (func(*testing.T, string, ...interface{}), func(...Arguments) []Arguments, func(...interface{}) Arguments)
}

// sentencesManagement will contains all functions getters: Given, Like,
// S and GoldenGiven.
type sentencesManagement struct {}

// Given returns the Given function, to be named by user.
func (sm *sentencesManagement) Given() (fn func(*testing.T, string, ...interface{})) {
	fn = Given
	return
}

// Golden returns the GivenWithGolden function, to be named by user.
func (sm *sentencesManagement) Golden() (fn func(*testing.T, string, ...interface{})) {
	fn = GivenWithGolden
	return
}

// All returns the set of sentences Give, Like and S to be named by
// user.
func (sm *sentencesManagement) All() (given func(*testing.T, string, ...interface{}), like func(...Arguments) []Arguments, s func(...interface{}) Arguments) {
	given = Given
	like = Like
	s = S
	return
}

// newSentencesManager creates a empty sentences manager.
func newSentencesManager() (sm SentencesManager) {
	sm = &sentencesManagement{}
	return
}

// Sentences return the manager for sentences, with options for user to
// choose desired testing style.
func Sentences() (sm SentencesManager) {
	sm = sentences
	return
}