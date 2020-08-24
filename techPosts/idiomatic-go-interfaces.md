*** title ***
Idiomatic Go Interfaces

*** date ***
1/15/19

*** tags ***
go

*** description ***
An exploration of Go Interfaces

*** body ***
In this blog post, I'm going to discuss interfaces; Go's core abstraction for polymorphism and code reuse. I'm going to presume experience with programming in general, and at least a passing familiarity with Go's syntax and semantics.

If you haven't seen Go interfaces before, you'll see several examples here. In addition, it's worth noting that Go interfaces are implemented *implicitly*. There isn't an "implements" keyword like in Java, and the implementation can happen in any package.

I can tell you from experience that it's possible to write production Go code for many months without developing a sense of what shape interfaces should take, when to use (or skip) them, or how to leverage the existing interfaces in the Go standard library.

To that end, I'll be sharing five standard library interfaces with two primary goals in mind. First, to provide examples of how to interact with Go's existing interfaces, and second, to draw attention to the conventions of idiomatic interfaces (e.g. small, focused, and well-named). So let's get started!


## STRINGER
```go
type Stringer interface {
    String() string
}

```

The `fmt` package's `Stringer` is the prototypical Go interface; it requires just one function, and its name is that function plus the "er" suffix. Any struct with a `String()` method is an implementation of the Stringer interface. One of the most common ways to take advantage of Stringer is when logging and debugging. For example, let's say you have a billing app that defines monthly and annual billing cycles as enums:

```go
type BillingCycle int

const (
    Monthly BillingCycle = iota
    Annual
)
```

If you wanted to log which billing cycle corresponds to a particular customer order, out of the box, you'd get simply `0` or `1`... not very useful, especially if you weren't the developer who defined these constants. It's great to follow enum definitions with a string representation:

```go
func (i BillingCycle) String() string {
    if i == 0 {
        return "Monthly"
    }
    return "Annual"
}
```

Now that BillingCycles meet the Stringer interface, you can pass one to any function that expects to receive a Stringer. For instance, you can pass a BillingCycle to `fmt.Sprintf` and get human readable output:

```go
fmt.Sprintf("Current Billing Cycle: %s", order.BillingCycle)
// returns "Current Billing Cycle: Monthly"
```

With only two values in play, hand writing that String function is the way to go. If you had a longer list of enums to account for (for example enums representing the 12 months of the year), it might make more sense to let the `generate` tool write this function for you.

To take advantage of code generation, leave a comment like this where your String() definition would otherwise go...`//go:generate stringer -type=BillingCycle` ...and run `go generate` from the terminal.

Another example Stringer implementation might be to represent a User as `fmt.Sprintf("%s %s", u.FirstName, u.LastName)`.

## ERROR

If you've written more than 10 lines of Go, you're certainly familiar with functions that return errors. But you may not know that the built-in error type is actually an interface.

```go
type error interface {
    Error() string
}
```

As with any other interface, you can implement this yourself by creating a struct with an Error method that returns a string. A number of third party libraries have done exactly this, such as Dave Cheney's [errors](https://github.com/pkg/errors) package that encourages you to wrap errors with additional context.

If you maintain a pizza ordering app, you might find value in defining a custom error implementation like this:

```go
type IngredientUnavailable struct {
    Ingredient string
    Reason string
}

func (e IngredientUnavailable) Error() string {
    return fmt.Sprintf("%s is unavailable due to %s", e.Ingredient, e.Reason)
}
```

This struct can then be used as the `error` return value of any function, and can be passed directly to functions in the fmt package where it will be correctly formatted.

## SORT.INTERFACE
```go
type Interface interface {
    Len() int
    Less(i, j int) bool
    Swap(i, j int)
}

```
The `sort` package's "Interface" interface can be implemented on any collection type (typically a slice), and is particularly useful when you have a collection that is sorted by multiple criteria throughout your codebase.

For example, if your new startup is Tinder for Accordions, you will soon find yourself implementing sort.Interface many times over as you sort accordions by weight, color, number of keys, and so on.

```go
type ByWeight []Accordion

func (a ByWeight) Len() int { return len(a) }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight < a[j].Weight }
func (a ByWeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

```
With this interface fulfilled, you can sort slices like this... `sort.Sort(ByWeight(accordions))` ...in a way that is very succinct and easy to read. Under the hood, Sort will use different sorting algorithms depending on the size that your `Len()` function returns.


## READER & WRITER

Finally, let's kill two birds with one stone by heading to the `io` package, where over a dozen small, clearly-named interfaces are defined. We'll focus on `Reader` and `Writer` - meaning things that you can read bytes from or write bytes to.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```
In many cases, you won't need to write your own `Reader`s and `Writer`s because the implementations that ship with Go already cover most use cases. Most user-defined implementations piggyback on those from the standard library and add a bit of functionality.

A Reader example that you will find in most web applications is an http `Response.Body`. A common example of a Writer is a `json.Encoder`. An then you have examples like `os.File` that work as both a Reader and Writer (since you can both read from and write to files on your computer). In fact, File is a ReadWriter:

```go
type ReadWriter interface {
    Reader
    Writer
}
```

So what do you gain from knowing that there are Readers and Writers at play in Go applications? You can use `io.Copy` to copy contents from any Reader to any Writer (i.e. from a file to standard out, from an http response to an md5 `hash.Hash`, etc). You can combine the output of many Writers into an `io.MultiWriter`. You can write a custom logger, such as [this one by Brad Fitzpatrick](https://play.golang.org/p/PM0Fx-o6Drz) that indents based on the call stack's current depth.

There are also a number of open source logging packages (such as [logrus](https://github.com/sirupsen/logrus)) that are `io.Writer`s that add functionality and play nice with the standard library (i.e. as the `ErrorLog` set on `http.Server{}`).

## TAKEAWAYS

Here are a few of the things that I'm hoping you'll take away from this post:

1. Go interfaces should be small and focused. An interface with 1 or 2 required methods is much more likely to be implemented multiple times than an interface with 10 methods.

2. Don't write single implementation interfaces. That's just noise and indirection in your codebase.

3. Interfaces are often named after the verb that they do plus "er" (within reason - we don't call it `errorer`!).

4. Interfaces should generally be defined in that package that consumes them. Writers are implemented all over the standard library, but the definition resides in `io` where the output is actually used. 

5. Interfaces are wonderful when used as the input to a function, but are often a headache if used as the return value.

6. Small interfaces can be composed together, as seen in the ReadWriter example (or to take it one step further, in ReadWriteCloser).

And that's all I've got for now. If you find any errors in this post or have any other feedback, feel free to reach out to me on Twitter. Thanks! -Adam