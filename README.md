# mongo [![GoDoc](https://godoc.org/github.com/ddspog/bdd?status.svg)](https://godoc.org/github.com/ddspog/bdd) [![Go Report Card](https://goreportcard.com/badge/github.com/ddspog/bdd)](https://goreportcard.com/report/github.com/ddspog/bdd) [![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/) [![Travis CI](https://travis-ci.org/ddspog/bdd.svg?branch=master)](https://travis-ci.org/ddspog/bdd)

by [ddspog](http://github.com/ddspog)

Package **bdd** enables creation of behaviour driven tests with sentences.

## License

You are free to copy, modify and distribute **bdd** package with attribution under the terms of the MIT license. See the [LICENSE](https://github.com/ddspog/bdd/blob/master/LICENSE) file for details.

## Installation

Install **bdd** package with:

```shell
go get github.com/ddspog/bdd
```

## How to use

 Package bdd enables creation of behaviour driven tests with sentences.

This is made through the use of bdd.Sentences(), it will return options
of sentences to return. That will be Given(), Golden() and All(). Those
methods will return the functions needed to make the bdd tests, using
this package, the user can name those function as it desired.

Use bdd.Sentences().Given() when making simple tests, declaring all
cases to be tested on it, like:

```go
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

Use bdd.Sentences().All() when making simple bdd tests, but with lots
of declared test cases for the same type of tests, like:

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

Finally, use bdd.Sentences().Golden() to create tests using golden
files, stored into testdata folder at the same folder of tests, like:

```go
func Test_Golden_Case(t *testing.T) {
    given := bdd.Sentences().Golden()

    input, gold := &struct {
        A int `json:"a"`
        B int `json:"b"`
    }{}, &struct {
        Sum int `json:"sum"`
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

For those tests it's important to have a file GoldenCase.json inside
package testdata folder. The file should contain a structure like:

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

## Contribution

This package has some objectives from now:

* Improve tests on each package and sub-package.
* Eliminate unnecessary code.
* Improve asserts.

Any interest in help is much appreciated.