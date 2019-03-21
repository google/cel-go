# Common Expression Language

[![Build Status](https://travis-ci.org/google/cel-go.svg?branch=master)](https://travis-ci.org/google/cel-go) [![Go Report Card](https://goreportcard.com/badge/github.com/google/cel-go)](https://goreportcard.com/report/github.com/google/cel-go)
[![GoDoc](https://godoc.org/github.com/google/cel-go?status.svg)][6]

The Common Expression Language (CEL) is a non-Turing complete language designed
for simplicity, speed, safety, and portability. CEL's C-like [syntax][1] looks
nearly identical to equivalent expressions in C++, Go, Java, and TypeScript.

```java
// Check whether a resource name starts with a group name.
resource.name.startsWith("/groups/" + auth.claims.group)
```

```go
// Determine whether the request is in the permitted time window.
request.time - resource.age < duration("24h")
```

```typescript
// Check whether all resource names in a list match a given filter.
auth.claims.email_verified && resources.all(r, r.startsWith(auth.claims.email))
```

A CEL "program" is a single expression. The examples have been tagged as
`java`, `go`, and `typescript` within the markdown to showcase the commonality
of the syntax.

CEL is ideal for lightweight expression evaluation when a fully sandboxed
scripting language is too resource intensive.

## Overview

Determine the variables and functions you want to provide to CEL. Parse and
check an expression to make sure it's valid. Then evaluate the output AST
against some input. Checking is optional, but strongly encouraged.

### Environment Setup

Let's expose `name` and  `group` variables to CEL using the `cel.Declarations`
environment option:

```go
import(
    "github.com/google/cel-go/cel"
    "github.com/google/cel-go/checker/decls"
)

env, err := cel.NewEnv(
    cel.Declarations(
        decls.NewIdent("name", decls.String, nil),
        decls.NewIdent("group", decls.String, nil))
```

That's it, the environment is ready to be use for parsing and type-checking.
CEL supports all the usual primitive types in addition to lists, maps, as well
as first-class support for JSON and Protocol Buffers.

### Parse and Check

The parsing phase indicates whether the expression is syntactically valid and
expands any macros present within the environment. Parsing and checking is
more computationally expensive than evaluation, and it is recommended that
expressions be parsed and checked ahead of time.

```go
parsed, issues := env.Parse(`name.startsWith("/groups/" + group)`)
if issues != nil && issues.Err() != nil {
    log.Fatalf("parse error: %s", issues.Err())
}
checked, issues := env.Check(parsed)
if issues != nil && issues.Err() != nil {
    log.Fatalf("type-check error: %s", issues.Err())
}
prg, err := env.Program(checked)
if err != nil {
    log.Fatalf("program construction error: %s", err)
}
```

The `cel.Program` generated at the end of parse and check is stateless,
thread-safe, and cachable.

Type-checking in an optional, but strongly encouraged step that can reject some
semantically invalid expressions using static analysis. Additionally, the check
produces metadata which can improve function invocation performance and object
field selection at evaluation-time.

#### Macros

Macros are enabled by default and may be disabled. Macros were introduced to
support optional CEL features that might not be desired in all use cases
without the syntactic burden and complexity such features might desire if
they were part of the core CEL syntax. Macros are expanded at parse time and
their expansions are type-checked at check time.

For example, when macros are enabled it is possible to support bounded
iteration / fold operators. The macros `all`, `exists`, `exists_one`, `filter`,
and `map` are particularly useful for evaluating a single predicate against
list and map values.

```javascript
// Ensure all tweets are less than 140 chars
tweets.all(t, t.size() <= 140)
```

The `has` macro is useful for unifying field presence testing logic across
protobuf types and dynamic (JSON-like) types.

```javascript
// Test whether the field is a non-default value if proto-based, or defined
// in the JSON case.
has(message.field)
```

Both cases traditionally require special syntax at the language level, but
these features are exposed via macros in CEL.

### Evaluate

Now, evaluate for fun and profit. The evaluation is thread-safe and side-effect
free. Many different inputs can be send to the same `cel.Program` and if fields
are present in the input, but not referenced in the expression, they are
ignored.

```go
// The `out` var contains the output of a successful evaluation.
// The `details' var would contain intermediate evalaution state if enabled as
// a cel.ProgramOption. This can be useful for visualizing how the `out` value
// was arrive at.
out, details, err := prg.Eval(map[string]interface{}{
    "name": "/groups/acme.co/documents/secret-stuff",
    "group": "acme.co"}})
fmt.Println(out) // 'true'
```

For more examples of how to use CEL, see [cel_test.go](https://github.com/google/cel-go/tree/master/cel/cel_test.go).

#### Partial State

What if `name` hadn't been supplied? CEL is designed for this case. In
distributed apps it is not uncommon to have edge caches and central services.
If possible, evaluation should happen at the edge, but it isn't always possible
to know the full state required for all values and functions present in the
CEL expression.

To improve the odds of successful evaluation with partial state, CEL uses
commutative logical operators `&&`, `||`. If an error or unknown value (not the
same thing) is encountered on the left-hand side, the right hand side is
evaluated also to determine the outcome. While it is possible to implement
evaluation with partial state without this feature, this method was chosen
because it aligns with the semantics of SQL evaluation and because it's more
robust to evaluation against dynamic data types such as JSON inputs.

In the following truth-table, the symbols `<x>` and `<y>` represent error or
unknown values, with the `?` indicating that the branch is not taken due to
short-circuiting. When the result is `<x, y>` this means that the both args
are possibly relevant to the result.

| Expression          | Result   |
|---------------------|----------|
| `false && ?`        | `false`  |
| `true && false`     | `false`  |
| `<x> && false`      | `false`  |
| `true && true`      | `true`   |
| `true && <x>`       | `<x>`    |
| `<x> && true`       | `<x>`    |
| `<x> && <y>`        | `<x, y>` |
| `true \|\| ?`       | `true`   |
| `false \|\| true`   | `true`   |
| `<x> \|\| true`     | `true`   |
| `false \|\| false`  | `false`  |
| `false \|\| <x>`    | `<x>`    |
| `<x> \|\| false`    | `<x>`    |
| `<x> \|\| <y>`      | `<x, y>` |

In the cases where unknowns are expected, `cel.EvalOptions(cel.OptTrackState)`
should be enabled. The `details` value returned by `Eval()` will contain the
intermediate evaluation values and can be provided to the `interpreter.Prune`
function to generate a residual expression. e.g.:

```cpp
// Residual when `name` omitted:
name.startsWith("/groups/acme.co")
```

This technique can be useful when there are variables that are expensive to
compute unless they are absolutely needed. This functionality will be the
focus of many future improvements, so keep an eye out for more goodness here!

### Errors

Parse and check errors have friendly error messages with pointers to where the
issues occur in source:

```sh
ERROR: <input>:1:40: undefined field 'undefined'
    | TestAllTypes{single_int32: 1, undefined: 2}
    | .......................................^`,
```

Both the parsed and checked expressions contain source position information
about each node that appears in the output AST. This information can be used
to determine error locations at evaluation time as well.

## Performance

CEL evaluates very quickly. When the expression does not change frequently,
or is easily cached, the evaluation speed is the more important factor when
considering an expression language.

The following expression was benchmarked between CEL and two other popular
Go expression language libraries, namely https://github.com/antonmedv/expr
and https://github.com/Knetic/govaluate:

```javascript
name.startsWith("/groups/" + group)  // CEL
startsWith(name, concat("/groups/", group)) // Govaluate
StartsWith(name, Concat("/groups/", group)) // Expr
```

The syntax varies slightly between the examples based on the features and
limitations of the various libraries. It is important to keep in mind that
the test setup and motivation for using one library over another is based
on a number of factors, and that benchmarks aren't necessarily indicative
of production behavior. Thus, the following results are purely illustrative
of CEL's evaluation speed.

```
BenchmarkCEL-8           3000000               357 ns/op
BenchmarkGovaluate-8     3000000               572 ns/op
BenchmarkExpr-8          1000000              1402 ns/op
```

Benchmark setup forthcoming in an upcoming wiki.

## Install

CEL-Go supports `modules` and may be installed using either the `get` or
`build` commands depending on your preference. Since CEL uses semantic
versioning prefer using the new go `modules`:

```sh
go mod init <my-cel-app>
go build ./...
```

```sh
go get -u github.com/google/cel-go/...
```

And of course, there is always the option to build from source directly.

## Common Questions

### Why not JavaScript, Lua, or WASM?

JavaScript and Lua are rich languages that require sandboxing to execute
safely. Sandboxing is costly and factors into the "what will I let users
evaluate?" question heavily when the answer is anything more than O(n)
complexity.

CEL evaluates linearly with respect to the size of the expression and the input
being evaluated when macros are disabled. The only functions beyond the
built-ins that may be invoked are provided by the host environment. While
extension functions may be more complex, this is a choice by the application
embedding CEL.

But, why not WASM? WASM is an excellent choice for certain applications and
is far superior to embedded JavaScript and Lua, but it does not have support
for garbage collection and non-primitive object types require semi-expensive
calls across modules. In most cases CEL will be faster and just as portable
for its intended use case, though for node.js and web-based execution CEL
too may offer a WASM evaluator with direct to WASM compilation.

### Do I need to Parse _and_ Check?

Checking is an optional, but strongly suggested, step in CEL expression
validation. It is sufficient in some cases to simply Parse and rely on the
runtime bindings and error handling to do the right thing.

### Where can I learn more about the language?

* See the [CEL Spec][1] for the specification and conformance test suite.
* Ask for support on the [CEL Go Discuss][2] Google group.

### Where can I learn more about the internals?

* See [GoDoc][6] to learn how to integrate CEL into services written in Go.
* See the [CEL C++][3] toolchain (under development) for information about how
  to integrate CEL evaluation into other environments.

### How can I contribute?

* See [CONTRIBUTING.md](./CONTRIBUTING.md) to get started.
* Use [GitHub Issues][4] to request features or report bugs.

### Some tests don't work with `go test`?

A handful of tests rely on [Bazel][5]. In particular dynamic proto support
at check time and the conformance test driver require Bazel to coordinate
the test inputs:

```sh
bazel test ...
```

## License

Released under the [Apache License](LICENSE).

Disclaimer: This is not an official Google product.

[1]:  https://github.com/google/cel-spec
[2]:  https://groups.google.com/forum/#!forum/cel-go-discuss
[3]:  https://github.com/google/cel-cpp
[4]:  https://github.com/google/cel-go/issues
[5]:  https://bazel.build
[6]:  https://godoc.org/github.com/google/cel-go
