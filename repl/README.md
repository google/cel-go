# CEL REPL

This is a simple tool for experimenting with CEL expressions and learning the
syntax.

## Usage
The REPL (Read Evaluate Print Loop) is implemented as a command line tool.

By default, the the input will be interpreted as a cel expression to evaluate.

Special commands (prefixed with '%') are used to update the evaluation
environment.

An example session:

```
# from a cel-go clone
$ cd ./repl

$ go run .
CEL REPL
%exit or EOF to quit.

cel-repl> %let x = 10
cel-repl> %let y = {'abc': {'def': [1, 2, 3]}}
cel-repl> y.abc.def.filter(el, el < x)
[1 2 3] : list(int)
cel-repl> %delete x
cel-repl> y.abc.def.filter(el, el < x)
Expr failed:
ERROR: <input>:1:27: undeclared reference to 'x' (in container '')
 | y.abc.def.filter(el, el < x)
 | ..........................^
cel-repl> %declare x : int
cel-repl> y.abc.def.filter(el, el < x)
Expr failed:
no such attribute: id: 8, names: [x]
no such attribute: id: 8, names: [x] : list(int)
cel-repl> y.abc.def.filter(el, el < x || el > 0)
[1 2 3] : list(int)
cel-repl> %exit
```

### Commands

#### let
`%let` introduces or update a variable or function  declaration and provide a
definition (as another CEL expression). A type hint is optionally provided to
check that the provided expression has the expected type.

`%let <identifier> (: <type>)? = <expr>`

Example:

`%let y = 42`

For functions, result types are mandatory:

`%let <identifier> (<identifier> : <type>, ...) : <type> -> <expr>`

Example:

`%let oracle(x : int) : bool -> x == 42`

#### declare

`%declare` introduces or updates a variable or function declaration with no
definition.

`%declare <identifier> : <type>`

#### delete
`%delete` deletes a variable declaration

`%delete <identifier>`

#### eval
`%eval` evaluate an expression:

`%eval <expr>` or simply `<expr>`

### Evaluation Model

The evaluator considers the let expressions and declarations in order, with
functions defined before variables. Let expressions may refer to earlier
expressions, but the reverse is not true. To prevent breaking dependant
expressions, updates will fail if removing or changing a let prevents a later
let expression from compiling. Let expressions are compiled when declared and
evaluated before the expression in an `%eval` command.

Functions are implicitly defined before variables: let variables may refer to
functions, but functions cannot refer to let variables.

Using curly-braces to indicate scopes, this looks like: 
```
let sum (x : int, y : int) : int -> x + y 
{
    let x = sum(2, 4)
    {
        let y = sum(x, 30) 
        {
            eval sum(x, y) == 42
            // (x) + (x + 30)
            // (2 + 4) + ((2 + 4) + 30)
        }
    }
}
```

## Installing

To build and install as a standalone binary: 

```
$ git clone git@github.com:google/cel-go.git ./cel-go
$ cd ./cel-go/repl
$ go build .
# e.g. to your $PATH
$ mv ./repl <install location>
```