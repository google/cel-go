# Common Expression Language

[![Build Status](https://travis-ci.org/google/cel-go.svg?branch=master)](https://travis-ci.org/google/cel-go) [![Go Report Card](https://goreportcard.com/badge/github.com/google/cel-go)](https://goreportcard.com/report/github.com/google/cel-go)
[![GoDoc](https://godoc.org/github.com/google/cel-go?status.svg)][6]

The Common Expression Language (CEL) is a non-Turing complete language designed
for simplicity, speed, safety, and portability. CEL's C-like syntax looks nearly
identical to equivalent expressions in C++, Go, Java, and TypeScript. e.g.

```javascript
// Check whether a resource name starts with a group name.
resource.name.startsWith('/groups/' + auth.claims.group)

// Determine whether the request is in the permitted time window.
request.time - resource.age < duration('24h')

// Check whether all resource names in a list match a given filter.
auth.claims.email_verified && resources.all(r, r.startsWith(auth.claims.email))
```

CEL is ideal for lightweight expression evaluation when a fully sandboxed
scripting language is too resource intensive.

## Overview

Determine the attribute names and functions you want to provide to the CEL
stack. Parse and check an expression to make sure it's valid. Then evaluate
the output AST against some input.

### Environment Setup

Let's expose `resource.name` and  `auth.claims` variables to CEL using the
standard builtins:

```go
import(
    "github.com/google/cel-go/cel"
    "github.com/google/cel-go/checker/decls"
)

env, err := cel.NewEnv(
    cel.Declarations(
        decls.NewIdent("resource.name", decls.String, nil),
        decls.NewIdent("auth.claims",
            decls.NewMapType(decls.String, decls.Dyn), nil)))
```

That's it, the environment is ready to be use for parsing and type-checking.
Note, the `auth.claims` variable is exposed as a map with `string` keys and
`dyn` values. The `dyn` in this case represents JSON input as the
`auth.claims` could originate from a JSON Web Token (JWT).

### Parse and Check

The parsing phase indicates whether the expression is syntactically valid and
expands any macros present within the environment. Parsing and checking is
more computationally expensive than evaluation, and it is recommended that
expressions be parsed and checked ahead of time.

```go
parsed, issues := env.Parse(
    `resource.name.startsWith("/groups/" + auth.claims.group)`)
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

Macros are enabled by default and may be disabled. The comprehension macros
`all`, `exists`, `exists_one`, `filter`, and `map` are particularly useful for evaluating a single predicate against list and map values.

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
    "resource.name": "/groups/acme.co/documents/secret-stuff",
    "auth.claims": map[string]interface{}{
        "sub": "alice@acme.co",
        "group": "acme.co"}})
fmt.Println(out) // 'true'
```

For more examples of how to use CEL, see [cel_test.go](https://github.com/google/cel-go/tree/master/cel/cel_test.go).

#### Partial State

What if the `auth.claims` hadn't been supplied? Well, CEL is designed for this.
When some of the input is unknown, the `out` value will contain a list of the
unknown expressions encountered during evaluation if they were relevant to the
outcome. In cases where an alternative conditional branch can produce a
definitive result, CEL takes this branch and skips right over the unknowns.

In the cases where unknowns are expected, state tracking should be enabled. The
`details` field will contain all of the intermediate evaluation values and can
be provided to the `interpreter.Prune` function to generate a residual
expression. e.g.:

```javascript
// Residual when `resource.name` omitted:
resource.name.startsWith("/groups/acme.co")
```

This technique can be useful when there are attributes that are expensive to
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

CEL evaluates very quickly. When the expression set does not change frequently,
or is easily cached, the evaluation speed is the more important factor when
considering an expression language.

The following expression was benchmarked between CEL and two other popular
Go expression language libraries, namely https://github.com/antonmedv/expr
and https://github.com/Knetic/govaluate:

```javascript
resource.name.startsWith("/groups/" + auth.claims)  // CEL
startsWith(resource_name, concat("/groups/", auth_claims)) // Govaluate
StartsWith(resource.name, Concat("/groups/", auth.claims)) // Expr
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
being evaluated. The only functions beyond the built-ins that may be invoked
are provided by the host environment. While extension functions may be more
complex, this is a choice by the application embedding CEL.

But, why not WASM? WASM is an excellent choice for certain applications and
is far superior to embedded JavaScript and Lua, but it does not have support
for garbage collection and non-primitive object types require semi-expensive
calls across modules. In most cases CEL will be faster and just as portable
for its intended use case, though for node.js and web-based execution CEL
too may offer a WASM evaluator with direct to WASM compilation.

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
