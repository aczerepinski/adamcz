*** title ***
Unit Testing in Go

*** date ***
12/3/17

*** tags ***
go

*** description ***
An overview of test writing using the Go standard library

*** body ***
In this blog post I will present an overview of unit testing, using only the Go standard library. This is intended for new Go developers who have a basic familiarity with the language but haven't yet ramped up on writing tests.

When I first started writing Go, I immediately reached for testing packages that offer reduced boilerplate and sleek abstractions. Coming out of Ruby and Elixir, all of the error checking and curly braces felt like problems to solve. In the year since then, I've softened a lot on that, and prefer to embrace Go for what it is - a very explicit language that requires a few extra lines of boilerplate now and then.

The cost of reaching for abstractions - test helpers or otherwise - too early is that you don't become intimately familiar with the layers beneath them. Go's testing package is great, and I would highly recommend writing dependency-free tests at least in your early days. If you move on to other testing libraries later, you'll better understand what they're bringing to the table, including their limitations.

### The Basics

Let's say we have a package called `taco` that, among other things, can return the best taqueria in a given zip code. Our file structure for the package might look like this:

```go
/taco
  taco.go
  taco_test.go
```

In the `taco_test.go` file a simple test might look like this:

```go
package taco

func TestBestTaqueriaByZip(t *testing.T) {
  zip := "02446"
  expected := "Dorado Tacos"
  actual, _ := BestTaqueriaByZip(zip)
  if expected != actual {
    t.Errorf("expected best taqueria in %s to be %s, got %s",
      zip, expected, actual)
  }
}
```

Notice that the test function (by convention) is the name of the function it is testing, prepended by "Test". Its input is a pointer to a `T` from the testing package. It's well worth your time to read [the docs on T](https://golang.org/pkg/testing/#T) and its methods, as they are your bread and butter for writing tests. In the above test, I use `Errorf` to mark the test failed and log some info about what went wrong.

Here are some examples of how you'd run tests from the command line:

```bash
# run a single test
go test -run TestBestTaqueriaByZip

# run all tests from the taco package
go test ./taco/... 

# run tests in all packages, excluding the vendor directory (Go 1.8)
go test $(go list ./... | grep -v /vendor/)

# run tests in all packages, excluding the vendor directory (Go 1.9)
go test ./...
```

### Test Tables

What if we want to test some more conditions? A naive attempt would be to copy/paste and modify what we had:

```go
func TestBestTaqueriaByZip(t *testing.T) {
  zip := "02446"
  expected := "Dorado Tacos"
  actual, _ := BestTaqueriaByZip(zip)
  if expected != actual {
    t.Errorf("expected best taqueria in %s to be %s, got %s",
      zip, expected, actual)
  }

  zip = "01702"
  expected = "El Maya"
  actual, _ = BestTaqueriaByZip(zip)
  if expected != actual {
    t.Errorf("expected best taqueria in %s to be %s, got %s",
      zip, expected, actual)
  }
}
```

For two test cases this sort of repetition may be acceptable, but any more and you'll go crazy. The idiomatic solution, found all over the Go language's own source code, is to write a test table:

```go
func TestBestTaqueriaByZip(t *testing.T) {
  tacoTests := []struct{
    zip string
    taqueria string
    err error
  }{
    {"02446", "Dorado Tacos", nil},
    {"01702", "El Maya", nil},
  }

  for _, test := range tacoTests {
    actual, err := BestTaqueriaByZip(test.zip)
    if test.taqueria != actual {
      t.Errorf("expected best taqueria in %s to be %s, got %s",
        test.zip, test.taqueria, actual)
    }
    if test.err != err {
      t.Errorf("%s: expected error to be %v, got %v",
        test.taqueria, test.err, err)
    }
  }
}
```

First we define tacoTests as a slice of anonymous structs, which contain the values we'll need for an individual test case. With this setup out of the way, adding an additional test case later will require only a single line of code. For instance we could test an edge case by inserting just this line: `{"GIMME TACOS!", "", ErrInvalidZip}`.

### More Testing.T Methods

Occasionally you'll test something that can be broken into sequential parts, where latter parts don't make sense if a preceding part fails. `t.FailNow()` marks a test failed and stops execution immediately. Similarly, `t.Fatal()` and `t.Fatalf()` can log a message and then immediately stop execution of the test.

I'll share an example of when something like this may be useful:

```go
if err := xml.Unmarshal(responseFromTacoAPI, &giantTacoStruct); err != nil {
  t.Fatalf("unable to unmarshal xml: %v", err)
}
// assertions about contents of giantTacoStruct...
```

If the unmarshaling didn't work, there's no need to continue checking various attributes on the struct it was meant to populate.

Another helpful method is `t.Skip()`, which is frequently used alongside `testing.Short()`.

```go
func TestThatMakesAnActualHTTPRequest(t *testing.T) {
  if testing.Short() {
    t.Skip("skipping test that makes an actual http request")
  }
}
```

Tests have a lot more value when the suite is fast enough that you can run it constantly. When a project reaches a certain size this may no longer be possible, but you can omit the longest running tests with the above code. From the command line the flag is: `go test -short`. 

Lastly, one of the coolest methods available is `t.Parallel()`, which marks a test safe for parallel execution. It isn't very helpful at first, because Go test suites run extremely quickly by default. But as your test suite grows, it may be worth sprinkling some `t.Parallel()`s around. Just be careful with tests that access constrained and/or mutable resources, like a database. The gains are somewhat limited because the lion's share of test time will still be spent building the packages, but execution time can be cut way down. In the case of an API package that makes many network requests, I saw `t.Parallel()` shave over 20 seconds of execution time.

### TestMain

`TestMain` is an option that you can turn to if tests in your package require significant setup and/or teardown. If you choose to include one of these, you'll be responsible for manually calling `m.Run()` and `os.Exit()`. Here's an example:

```go
func TestMain(m *testing.M) {
  // setup code, such as running migrations on a test database
  os.Exit(m.Run())
}
```

One thing to be aware of is that Go packages are likely to be comprised of many files, and `TestMain` could theoretically be sitting in any of them. It is possible to inject globals into the rest of a package's tests from here, but doing so obscures the behavior of tests in other files (from the same package). Whenever possible, I prefer that test setup occur with a normal factory function that can be called explicitly from each test. The latter option is only one line of code per test and it reduces the chance for confusion.

### Code Coverage

You can generate a really great code coverage report using only tools from a standard Go installation.  Running `go test -coverprofile=coverage.out` generates a coverage file that can be opened with `go tool cover -html=coverage.out`. It renders lines of code that have test coverage in green, and line that don't in red. Additionally, some text editors such as [VS Code](https://github.com/Microsoft/vscode-go) can run this for you and display the results right in your editor.


### Benchmarking

So far we've been looking at `TestMyFunction` funcs that take a pointer to `testing.T`. The testing package also provides benchmarking with a similar API, `BenchmarkMyFunction(b *testing.B)`. The general pattern looks like this:

```go
func BenchmarkBestTaqueriaByZip(b *testing.B) {
  for n := 0; n < b.N; n++ {
    bestTaqueriaByZip("02215")
  }
}
```

When you run your benchmark (`go test -bench=.`), you don't need to specify the value of `N`. Go will pick a (usually very high) number of iterations that varies based on the execution time and consistency of the function you're testing. 

If you ever discover a hot path that needs to be optimized, remember that Go offers great benchmarking that can be set up in just a few lines of test code, and executed in ~1 second. Feel free to open your terminal and text editor side by side, and run the benchmark repeatedly as you edit.

And that's all I have for now! I hope you enjoyed this post. If you want to reach out to me with any questions or corrections, feel free to do so [on  Twitter](https://twitter.com/adamczerepinski). Thanks for reading!

-Adam