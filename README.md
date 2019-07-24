# bdd [![GoDoc](https://godoc.org/github.com/ddsgok/bdd?status.svg)](https://godoc.org/github.com/ddsgok/bdd) [![Go Report Card](https://goreportcard.com/badge/github.com/ddsgok/bdd)](https://goreportcard.com/report/github.com/ddsgok/bdd) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) [![Travis CI](https://travis-ci.org/ddsgok/bdd.svg?branch=master)](https://travis-ci.org/ddsgok/bdd)

by [eduncan911](https://github.com/eduncan911) and forked by [ddspog](http://github.com/ddspog)

Package **bdd** enables creation of behaviour driven tests with sentences.

## License

You are free to copy, modify and distribute **bdd** package with attribution under the terms of the MIT license. See the [LICENSE](https://github.com/ddsgok/bdd/blob/master/LICENSE) file for details.

## Installation

Install **bdd** package with:

```shell
go get github.com/ddsgok/bdd
```

## How to use

Package bdd enables creation of behaviour driven tests with sentences.

This is made through the use of bdd.Sentences(), it will return options of sentences to return. That will be Given(), Golden() and All(). Those methods will return the functions needed to make the bdd tests, using this package, the user can name those function as it desired.

To start using the package, take Dan North's original BDD definitions, you spec code using the Given/When/Then storyline similar to:

```gherkin
Feature X
    Given a context
    When an event occurs
    Then it should do something
```

You represent these thoughts in code using bdd.Sentences().Given():

```go
package product_test

import (
    "github.com/ddsgok/bdd"
    "testing"
)

func Test_Simple_Case(t *testing.T) {
    given := bdd.Sentences().Given()

    given(t, "a Product p", func(when bdd.When) {
        p := newProduct()

        when("p.SetPrice(12) is called", func(it bdd.It) {
            p.SetPrice(12)

            it("p.GetPrice() should return 12", func(assert bdd.Assert) {
                assert.Equal(12, p.GetPrice())
            })
        })
    })
}
```

Use bdd.Sentences().All() when making simple bdd tests, but with lots of declared test cases for the same type of tests, like:

```go
func Test_Multiple_Case(t *testing.T) {
    given, like, s := bdd.Sentences().All()

    given(t, "a Product p", func(when bdd.When) {
        p := newProduct()

        when("p.SetPrice(%[1]v) is called", func(it bdd.It, args ...interface{}) {
            p.SetPrice(args[0].(int))

            it("p.GetPrice() should return %[1]v", func(assert bdd.Assert) {
                assert.Equal(args[0].(int), p.GetPrice())
            })
        })
    })
}
```

Finally, use bdd.Sentences().Golden() to create tests using golden files, stored into testdata folder at the same folder of tests, like:

```go
func Test_Golden_Case(t *testing.T) {
    given := bdd.Sentences().Golden()

    input, gold := &struct {
        A int `json:"a" yaml:"a"`
        B int `json:"b" yaml:"b"`
    }{}, &struct {
        Sum int `json:"sum" yaml:"sum"`
    }{}

    given(t, "two values a = %[input.a]v and b = %[input.b]v", func(when bdd.When, golden bdd.Golden) {
        golden.Load(input, gold)
        a, b := input.A, input.B

        when("sum := a + b is called", func(it bdd.It){
            sum := a + b

            golden.Update(func() interface{} {
                gold.Sum = sum
                return gold
            })

            it("should have sum equal to %[golden.sum]v", func(assert bdd.Assert) {
                assert.Equal(gold.Sum, sum)
            })
        })
    })
}
```

For those tests it's important to have a file GoldenCase.json inside package testdata folder. The file should contain a structure like:

```json
{
    "two values a = %[input.a]v and b = %[input.b]v": [{
        "input": { "a": 0, "b", 1 },
        "golden": { "sum": 1 }
    }, {
        "input": { "a": 2, "b", 3 },
        "golden": { "sum": 5 }
    }]
}
```

You could also use a file GoldenCase.yml with the structure:

```yaml
"two values a = %[input.a]v and b = %[input.b]v":
- input: {a: 0, b: 1}
    golden: {sum: 1}
- input: {a: 2, b: 3}
    golden: {sum: 5}
```

Please use adequate syntax to avoid problems when parsing these files. My package still doesn't have the best error description yet.

## Golden Files

All test names using this package, will name the feature, which removes '_' and change to spaces. The golden files for each test, must be named with the words in test name joined, and with the first letter uppercase.

For example, in a test called:

```go
func Test_Creation_of_a_Product(t *testing.T) {
    given := bdd.Sentences().Golden()
    // ...
}
```

The feature will be named as "Creation of a Product", and the name of golden file for this test should be "CreationOfAProduct.json" or the ".yml" version.

For all the golden files, they need the following structure:

```yaml
"context of test, named on given":
- { input: {}, golden: {} }
- { input: {}, golden: {} }
# ...
"another context":
- { input: {}, golden: {} }
- { input: {}, golden: {} }
# ...
```

All golden files, should start with a dictionary of context to a list of test cases. These context are named on the given sentence, after t argument.

The list of test cases must, for each value, contain two keys: input and golden, with its values as user desire.

Using `-update` flag will update the golden fields on each test case. Actually that would mean the golden file should have a nice layout for filling. In the future, I'll add a way to automatically ensure to have a file with nice structure.

All tests using this package have colored output.

## Contribution

This package has some objectives from now:

* Improve tests on each package and sub-package.
* Eliminate unnecessary code.
* Improve asserts.

Any interest in help is much appreciated.
