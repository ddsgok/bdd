package spec

import (
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/ddspog/bdd/colors"
	"github.com/ddspog/bdd/internal/common"
)

// failingLineData stores information about error, captured on assert
// sentence. It focus on describing the 3 lines, centered on error.
type failingLineData struct {
	prev     string
	content  string
	next     string
	filename string
	number   int
}

// TestSpecification holds the state of the test context for a specific specification.
type TestSpecification struct {
	T                       *testing.T
	Feature                 string
	Given                   string
	When                    string
	It                      string
	AssertFn                func(common.Assert)
	AssertionFailed         bool
	AssertionFailedMessages []string

	NotImplemented bool
}

// PrintFeature prints line informing about feature being tested.
func (spec *TestSpecification) PrintFeature() {
	if config.LastFeature != spec.Feature {
		if config.Output != OutputNone {
			fmt.Printf("%sFeature: %s%s\n", config.AnsiOfFeature, spec.Feature, colors.Reset)
		}
		config.LastFeature = spec.Feature
	}

	config.ResetLasts()
}

// PrintContext prints line informing about context being tested.
func (spec *TestSpecification) PrintContext() {
	if config.LastGiven != spec.Given {
		if config.Output != OutputNone {
			fmt.Printf("%s  Given %s%s\n", config.AnsiOfGiven, withLeftPadding(spec.Given, 2), colors.Reset)
		}
		config.LastGiven = spec.Given
	}

	config.ResetWhen()
}

// PrintWhen prints line informing about situation being tested.
func (spec *TestSpecification) PrintWhen() {
	if config.LastWhen != spec.When {
		if config.Output != OutputNone {
			fmt.Printf("%s    When %s%s\n", config.AnsiOfWhen, spec.When, colors.Reset)
		}
		config.LastWhen = spec.When
	}

	config.ResetIt()
}

// PrintIt prints line informing about verification being tested when
// successful.
func (spec *TestSpecification) PrintIt() {
	if config.Output != OutputNone {
		fmt.Printf("%s    » It %s %s\n", config.AnsiOfThen, spec.It, colors.Reset)
	}
	config.LastIt = spec.It
}

// PrintItWithError prints line informing about verification being
// tested when verification fail.
func (spec *TestSpecification) PrintItWithError() {
	if config.Output != OutputNone {
		fmt.Printf("%s    » It %s %s\n", config.AnsiOfThenWithError, spec.It, colors.Reset)
	}
	config.LastIt = spec.It
}

// PrintItNotImplemented prints line informing about verification not
// implemented.
func (spec *TestSpecification) PrintItNotImplemented() {
	if config.Output != OutputNone {
		fmt.Printf("%s    » It %s «-- NOT IMPLEMENTED%s\n", config.AnsiOfThenNotImplemented, spec.It, colors.Reset)
	}
	config.LastIt = spec.It
}

// PrintError prints text detailing how the verification failed on
// test.
func (spec *TestSpecification) PrintError(message string) {
	if failingLine, err := failingLine(); err == nil && config.Output != OutputNone {
		fmt.Printf("%s%s%s\n", config.AnsiOfExpectedError, message, colors.Reset)
		fmt.Printf("%s        in %s:%d%s\n", config.AnsiOfCode, path.Base(failingLine.filename), failingLine.number, colors.Reset)
		fmt.Printf("%s        ---------\n", config.AnsiOfCode)
		fmt.Printf("%s        %d. %s%s\n", config.AnsiOfCode, failingLine.number-1, withSoftTabs(failingLine.prev), colors.Reset)
		fmt.Printf("%s        %d. %s %s\n", config.AnsiOfCodeError, failingLine.number, failingLine.content, colors.Reset)
		fmt.Printf("%s        %d. %s%s\n", config.AnsiOfCode, failingLine.number+1, withSoftTabs(failingLine.next), colors.Reset)
		fmt.Println()
		spec.T.Fail()
		fmt.Println()
	}
}

// Run handles contextual printing and some delegation
// to the Assert's implementation for error handling
func (spec *TestSpecification) Run() {

	// execute the Assertion
	spec.AssertFn(config.assertFn(spec))

	// if there was no error (which handles its own printing),
	// print the spec here.
	if spec.NotImplemented {
		spec.PrintItNotImplemented()
	} else if !spec.AssertionFailed {
		spec.PrintIt()
	}

	spec.AssertionFailed = false
}

func New(t *testing.T, feat, given string) (sp *TestSpecification) {
	sp = &TestSpecification{
		T:       t,
		Feature: feat,
		Given:   given,
	}
	return
}

// failingLine returns information about current failing line on test.
func failingLine() (fl failingLineData, err error) {
	fl = failingLineData{}

	// this entire func is now a hack because of where it is being called,
	// which is now one caller higher.  previously it was being called in the
	// Expect struct which had the right caller info.  but now, it is being
	// called after the Assertion has been executed to print details to the
	// string.

	_, filename, ln, _ := runtime.Caller(6)

	// this is really hacky, need to find a way of not using magic numbers for runtime.Caller
	// If we are not in a test file, we must still be inside this package,
	// so we need to go up one more stack frame to get to the test file
	if !strings.HasSuffix(filename, "_test.go") {
		_, filename, ln, _ = runtime.Caller(7)
	}

	bf, err := ioutil.ReadFile(filename)

	if err != nil {
		err = fmt.Errorf("failed to open %s", filename)
		return
	}

	lines := strings.Split(string(bf), "\n")[ln-2 : ln+2]

	fl = failingLineData{
		prev:     withSoftTabs(lines[0]),
		content:  withSoftTabs(lines[1]),
		next:     withSoftTabs(lines[2]),
		filename: filename,
		number:   ln,
	}
	return
}

// withSoftTabs returns string after replacing any tabs to soft tabs.
func withSoftTabs(text string) (r string) {
	r = strings.Replace(text, "\t", "  ", -1)
	return
}

// withLeftPadding returns string after replacing new lines with left
// adjusted new lines, with desired padding.
func withLeftPadding(text string, padding int) (r string) {
	pad := "\n"
	for i := 0; i < padding; i++ {
		pad = strings.Join([]string{pad, " "}, "")
	}

	r = strings.Replace(text, "\n", pad, -1)
	return
}
