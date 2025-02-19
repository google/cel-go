// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ext

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

type clauseImpl func(ref.Val) (string, error)

type appendingFormatter struct {
	buf []byte
}

type formattedMapEntry struct {
	key string
	val string
}

func (af *appendingFormatter) format(arg ref.Val) error {
	switch arg.Type() {
	case types.BoolType:
		argBool, ok := arg.Value().(bool)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.BoolType)
		}
		af.buf = strconv.AppendBool(af.buf, argBool)
		return nil
	case types.IntType:
		argInt, ok := arg.Value().(int64)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.IntType)
		}
		af.buf = strconv.AppendInt(af.buf, argInt, 10)
		return nil
	case types.UintType:
		argUint, ok := arg.Value().(uint64)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.UintType)
		}
		af.buf = strconv.AppendUint(af.buf, argUint, 10)
		return nil
	case types.DoubleType:
		argDbl, ok := arg.Value().(float64)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.DoubleType)
		}
		if math.IsNaN(argDbl) {
			af.buf = append(af.buf, "NaN"...)
			return nil
		}
		if math.IsInf(argDbl, -1) {
			af.buf = append(af.buf, "-Infinity"...)
			return nil
		}
		if math.IsInf(argDbl, 1) {
			af.buf = append(af.buf, "Infinity"...)
			return nil
		}
		af.buf = strconv.AppendFloat(af.buf, argDbl, 'f', -1, 64)
		return nil
	case types.BytesType:
		argBytes, ok := arg.Value().([]byte)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.BytesType)
		}
		af.buf = append(af.buf, argBytes...)
		return nil
	case types.StringType:
		argStr, ok := arg.Value().(string)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.StringType)
		}
		af.buf = append(af.buf, argStr...)
		return nil
	case types.DurationType:
		argDur, ok := arg.Value().(time.Duration)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.DurationType)
		}
		af.buf = strconv.AppendFloat(af.buf, argDur.Seconds(), 'f', -1, 64)
		af.buf = append(af.buf, "s"...)
		return nil
	case types.TimestampType:
		argTime, ok := arg.Value().(time.Time)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.TimestampType)
		}
		af.buf = argTime.UTC().AppendFormat(af.buf, time.RFC3339Nano)
		return nil
	case types.NullType:
		af.buf = append(af.buf, "null"...)
		return nil
	case types.TypeType:
		argType, ok := arg.Value().(string)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.TypeType)
		}
		af.buf = append(af.buf, argType...)
		return nil
	case types.ListType:
		argList, ok := arg.(traits.Lister)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.ListType)
		}
		argIter := argList.Iterator()
		af.buf = append(af.buf, "["...)
		if argIter.HasNext() == types.True {
			if err := af.format(argIter.Next()); err != nil {
				return err
			}
			for argIter.HasNext() == types.True {
				af.buf = append(af.buf, ", "...)
				if err := af.format(argIter.Next()); err != nil {
					return err
				}
			}
		}
		af.buf = append(af.buf, "]"...)
		return nil
	case types.MapType:
		argMap, ok := arg.(traits.Mapper)
		if !ok {
			return fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.MapType)
		}
		argIter := argMap.Iterator()
		ents := []formattedMapEntry{}
		for argIter.HasNext() == types.True {
			key := argIter.Next()
			val, ok := argMap.Find(key)
			if !ok {
				return fmt.Errorf("key missing from map: '%s'", key)
			}
			keyStr, err := FormatString(key)
			if err != nil {
				return err
			}
			valStr, err := FormatString(val)
			if err != nil {
				return err
			}
			ents = append(ents, formattedMapEntry{keyStr, valStr})
		}
		sort.SliceStable(ents, func(x, y int) bool {
			return ents[x].key < ents[y].key
		})
		af.buf = append(af.buf, "{"...)
		for i, e := range ents {
			if i > 0 {
				af.buf = append(af.buf, ", "...)
			}
			af.buf = append(af.buf, e.key...)
			af.buf = append(af.buf, ": "...)
			af.buf = append(af.buf, e.val...)
		}
		af.buf = append(af.buf, "}"...)
		return nil
	default:
		return stringFormatError(runtimeID, arg.Type().TypeName())
	}
}

// FormatString returns the string representation of a CEL value.
//
// It is used to implement the %s specifier in the (string).format() extension function.
func FormatString(arg ref.Val) (string, error) {
	var fmter appendingFormatter
	if err := fmter.format(arg); err != nil {
		return "", err
	}
	return string(fmter.buf), nil
}

type stringFormatter struct{}

func (c *stringFormatter) String(arg ref.Val) (string, error) {
	return FormatString(arg)
}

func (c *stringFormatter) Decimal(arg ref.Val) (string, error) {
	switch arg.Type() {
	case types.IntType:
		argInt, ok := arg.Value().(int64)
		if !ok {
			return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.IntType)
		}
		return strconv.FormatInt(argInt, 10), nil
	case types.UintType:
		argUint, ok := arg.Value().(uint64)
		if !ok {
			return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.UintType)
		}
		return strconv.FormatUint(argUint, 10), nil
	case types.DoubleType:
		argDbl, ok := arg.Value().(float64)
		if !ok {
			return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.DoubleType)
		}
		if math.IsNaN(argDbl) {
			return "NaN", nil
		}
		if math.IsInf(argDbl, -1) {
			return "-Infinity", nil
		}
		if math.IsInf(argDbl, 1) {
			return "Infinity", nil
		}
		return strconv.FormatFloat(argDbl, 'f', -1, 64), nil
	default:
		return "", decimalFormatError(runtimeID, arg.Type().TypeName())
	}
}

func (c *stringFormatter) Fixed(precision int) func(ref.Val) (string, error) {
	return func(arg ref.Val) (string, error) {
		fmtStr := fmt.Sprintf("%%.%df", precision)
		switch arg.Type() {
		case types.IntType:
			argInt, ok := arg.Value().(int64)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.IntType)
			}
			return fmt.Sprintf(fmtStr, argInt), nil
		case types.UintType:
			argUint, ok := arg.Value().(uint64)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.UintType)
			}
			return fmt.Sprintf(fmtStr, argUint), nil
		case types.DoubleType:
			argDbl, ok := arg.Value().(float64)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.DoubleType)
			}
			if math.IsNaN(argDbl) {
				return "NaN", nil
			}
			if math.IsInf(argDbl, -1) {
				return "-Infinity", nil
			}
			if math.IsInf(argDbl, 1) {
				return "Infinity", nil
			}
			return fmt.Sprintf(fmtStr, argDbl), nil
		default:
			return "", fixedPointFormatError(runtimeID, arg.Type().TypeName())
		}
	}
}

func (c *stringFormatter) Scientific(precision int) func(ref.Val) (string, error) {
	return func(arg ref.Val) (string, error) {
		fmtStr := fmt.Sprintf("%%1.%de", precision)
		switch arg.Type() {
		case types.IntType:
			argInt, ok := arg.Value().(int64)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.IntType)
			}
			return fmt.Sprintf(fmtStr, argInt), nil
		case types.UintType:
			argUint, ok := arg.Value().(uint64)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.UintType)
			}
			return fmt.Sprintf(fmtStr, argUint), nil
		case types.DoubleType:
			argDbl, ok := arg.Value().(float64)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.DoubleType)
			}
			if math.IsNaN(argDbl) {
				return "NaN", nil
			}
			if math.IsInf(argDbl, -1) {
				return "-Infinity", nil
			}
			if math.IsInf(argDbl, 1) {
				return "Infinity", nil
			}
			return fmt.Sprintf(fmtStr, argDbl), nil
		default:
			return "", scientificFormatError(runtimeID, arg.Type().TypeName())
		}
	}
}

func (c *stringFormatter) Binary(arg ref.Val) (string, error) {
	switch arg.Type() {
	case types.BoolType:
		argBool, ok := arg.Value().(bool)
		if !ok {
			return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.BoolType)
		}
		if argBool {
			return "1", nil
		}
		return "0", nil
	case types.IntType:
		argInt, ok := arg.Value().(int64)
		if !ok {
			return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.IntType)
		}
		return strconv.FormatInt(argInt, 2), nil
	case types.UintType:
		argUint, ok := arg.Value().(uint64)
		if !ok {
			return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.UintType)
		}
		return strconv.FormatUint(argUint, 2), nil
	default:
		return "", binaryFormatError(runtimeID, arg.Type().TypeName())
	}
}

func (c *stringFormatter) Hex(useUpper bool) func(ref.Val) (string, error) {
	return func(arg ref.Val) (string, error) {
		var fmtStr string
		if useUpper {
			fmtStr = "%X"
		} else {
			fmtStr = "%x"
		}
		switch arg.Type() {
		case types.IntType:
			argInt, ok := arg.Value().(int64)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.IntType)
			}
			return fmt.Sprintf(fmtStr, argInt), nil
		case types.UintType:
			argUint, ok := arg.Value().(uint64)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.UintType)
			}
			return fmt.Sprintf(fmtStr, argUint), nil
		case types.StringType:
			argStr, ok := arg.Value().(string)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.StringType)
			}
			return fmt.Sprintf(fmtStr, argStr), nil
		case types.BytesType:
			argBytes, ok := arg.Value().([]byte)
			if !ok {
				return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.BytesType)
			}
			return fmt.Sprintf(fmtStr, argBytes), nil
		default:
			return "", hexFormatError(runtimeID, arg.Type().TypeName())
		}
	}
}

func (c *stringFormatter) Octal(arg ref.Val) (string, error) {
	switch arg.Type() {
	case types.IntType:
		argInt, ok := arg.Value().(int64)
		if !ok {
			return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.IntType)
		}
		return strconv.FormatInt(argInt, 8), nil
	case types.UintType:
		argUint, ok := arg.Value().(uint64)
		if !ok {
			return "", fmt.Errorf("type conversion error from '%s' to '%s'", arg.Type(), types.UintType)
		}
		return strconv.FormatUint(argUint, 8), nil
	default:
		return "", octalFormatError(runtimeID, arg.Type().TypeName())
	}
}

// stringFormatValidator implements the cel.ASTValidator interface allowing for static validation
// of string.format calls.
type stringFormatValidator struct{}

// Name returns the name of the validator.
func (stringFormatValidator) Name() string {
	return "cel.validator.string_format"
}

// Configure implements the ASTValidatorConfigurer interface and augments the list of functions to skip
// during homogeneous aggregate literal type-checks.
func (stringFormatValidator) Configure(config cel.MutableValidatorConfig) error {
	functions := config.GetOrDefault(cel.HomogeneousAggregateLiteralExemptFunctions, []string{}).([]string)
	functions = append(functions, "format")
	return config.Set(cel.HomogeneousAggregateLiteralExemptFunctions, functions)
}

// Validate parses all literal format strings and type checks the format clause against the argument
// at the corresponding ordinal within the list literal argument to the function, if one is specified.
func (stringFormatValidator) Validate(env *cel.Env, _ cel.ValidatorConfig, a *ast.AST, iss *cel.Issues) {
	root := ast.NavigateAST(a)
	formatCallExprs := ast.MatchDescendants(root, matchConstantFormatStringWithListLiteralArgs(a))
	for _, e := range formatCallExprs {
		call := e.AsCall()
		formatStr := call.Target().AsLiteral().Value().(string)
		args := call.Args()[0].AsList().Elements()
		formatCheck := &stringFormatChecker{
			args: args,
			ast:  a,
		}
		// use a placeholder locale, since locale doesn't affect syntax
		_, err := parseFormatString(formatStr, formatCheck, formatCheck)
		if err != nil {
			iss.ReportErrorAtID(getErrorExprID(e.ID(), err), "%v", err)
			continue
		}
		seenArgs := formatCheck.argsRequested
		if len(args) > seenArgs {
			iss.ReportErrorAtID(e.ID(),
				"too many arguments supplied to string.format (expected %d, got %d)", seenArgs, len(args))
		}
	}
}

// getErrorExprID determines which list literal argument triggered a type-disagreement for the
// purposes of more accurate error message reports.
func getErrorExprID(id int64, err error) int64 {
	fmtErr, ok := err.(formatError)
	if ok {
		return fmtErr.id
	}
	wrapped := errors.Unwrap(err)
	if wrapped != nil {
		return getErrorExprID(id, wrapped)
	}
	return id
}

// matchConstantFormatStringWithListLiteralArgs matches all valid expression nodes for string
// format checking.
func matchConstantFormatStringWithListLiteralArgs(a *ast.AST) ast.ExprMatcher {
	return func(e ast.NavigableExpr) bool {
		if e.Kind() != ast.CallKind {
			return false
		}
		call := e.AsCall()
		if !call.IsMemberFunction() || call.FunctionName() != "format" {
			return false
		}
		overloadIDs := a.GetOverloadIDs(e.ID())
		if len(overloadIDs) != 0 {
			found := false
			for _, overload := range overloadIDs {
				if overload == overloads.ExtFormatString {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		formatString := call.Target()
		if formatString.Kind() != ast.LiteralKind || formatString.AsLiteral().Type() != cel.StringType {
			return false
		}
		args := call.Args()
		if len(args) != 1 {
			return false
		}
		formatArgs := args[0]
		return formatArgs.Kind() == ast.ListKind
	}
}

// stringFormatChecker implements the formatStringInterpolater interface
type stringFormatChecker struct {
	args          []ast.Expr
	argsRequested int
	currArgIndex  int64
	ast           *ast.AST
}

func (c *stringFormatChecker) String(arg ref.Val) (string, error) {
	formatArg := c.args[c.currArgIndex]
	valid, badID := c.verifyString(formatArg)
	if !valid {
		return "", stringFormatError(badID, c.typeOf(badID).TypeName())
	}
	return "", nil
}

func (c *stringFormatChecker) Decimal(arg ref.Val) (string, error) {
	id := c.args[c.currArgIndex].ID()
	valid := c.verifyTypeOneOf(id, types.IntType, types.UintType, types.DoubleType)
	if !valid {
		return "", decimalFormatError(id, c.typeOf(id).TypeName())
	}
	return "", nil
}

func (c *stringFormatChecker) Fixed(precision int) func(ref.Val) (string, error) {
	return func(arg ref.Val) (string, error) {
		id := c.args[c.currArgIndex].ID()
		valid := c.verifyTypeOneOf(id, types.IntType, types.UintType, types.DoubleType)
		if !valid {
			return "", fixedPointFormatError(id, c.typeOf(id).TypeName())
		}
		return "", nil
	}
}

func (c *stringFormatChecker) Scientific(precision int) func(ref.Val) (string, error) {
	return func(arg ref.Val) (string, error) {
		id := c.args[c.currArgIndex].ID()
		valid := c.verifyTypeOneOf(id, types.IntType, types.UintType, types.DoubleType)
		if !valid {
			return "", scientificFormatError(id, c.typeOf(id).TypeName())
		}
		return "", nil
	}
}

func (c *stringFormatChecker) Binary(arg ref.Val) (string, error) {
	id := c.args[c.currArgIndex].ID()
	valid := c.verifyTypeOneOf(id, types.BoolType, types.IntType, types.UintType)
	if !valid {
		return "", binaryFormatError(id, c.typeOf(id).TypeName())
	}
	return "", nil
}

func (c *stringFormatChecker) Hex(useUpper bool) func(ref.Val) (string, error) {
	return func(arg ref.Val) (string, error) {
		id := c.args[c.currArgIndex].ID()
		valid := c.verifyTypeOneOf(id, types.IntType, types.UintType, types.StringType, types.BytesType)
		if !valid {
			return "", hexFormatError(id, c.typeOf(id).TypeName())
		}
		return "", nil
	}
}

func (c *stringFormatChecker) Octal(arg ref.Val) (string, error) {
	id := c.args[c.currArgIndex].ID()
	valid := c.verifyTypeOneOf(id, types.IntType, types.UintType)
	if !valid {
		return "", octalFormatError(id, c.typeOf(id).TypeName())
	}
	return "", nil
}

func (c *stringFormatChecker) Arg(index int64) (ref.Val, error) {
	c.argsRequested++
	c.currArgIndex = index
	// return a dummy value - this is immediately passed to back to us
	// through one of the FormatCallback functions, so anything will do
	return types.Int(0), nil
}

func (c *stringFormatChecker) Size() int64 {
	return int64(len(c.args))
}

func (c *stringFormatChecker) typeOf(id int64) *cel.Type {
	return c.ast.GetType(id)
}

func (c *stringFormatChecker) verifyTypeOneOf(id int64, validTypes ...*cel.Type) bool {
	t := c.typeOf(id)
	if t == cel.DynType {
		return true
	}
	for _, vt := range validTypes {
		// Only check runtime type compatibility without delving deeper into parameterized types
		if t.Kind() == vt.Kind() {
			return true
		}
	}
	return false
}

func (c *stringFormatChecker) verifyString(sub ast.Expr) (bool, int64) {
	paramA := cel.TypeParamType("A")
	paramB := cel.TypeParamType("B")
	subVerified := c.verifyTypeOneOf(sub.ID(),
		cel.ListType(paramA), cel.MapType(paramA, paramB),
		cel.IntType, cel.UintType, cel.DoubleType, cel.BoolType, cel.StringType,
		cel.TimestampType, cel.BytesType, cel.DurationType, cel.TypeType, cel.NullType)
	if !subVerified {
		return false, sub.ID()
	}
	switch sub.Kind() {
	case ast.ListKind:
		for _, e := range sub.AsList().Elements() {
			// recursively verify if we're dealing with a list/map
			verified, id := c.verifyString(e)
			if !verified {
				return false, id
			}
		}
		return true, sub.ID()
	case ast.MapKind:
		for _, e := range sub.AsMap().Entries() {
			// recursively verify if we're dealing with a list/map
			entry := e.AsMapEntry()
			verified, id := c.verifyString(entry.Key())
			if !verified {
				return false, id
			}
			verified, id = c.verifyString(entry.Value())
			if !verified {
				return false, id
			}
		}
		return true, sub.ID()
	default:
		return true, sub.ID()
	}
}

// helper routines for reporting common errors during string formatting static validation and
// runtime execution.

func binaryFormatError(id int64, badType string) error {
	return newFormatError(id, "only ints, uints, and bools can be formatted as binary, was given %s", badType)
}

func decimalFormatError(id int64, badType string) error {
	return newFormatError(id, "decimal clause can only be used on ints, uints, and doubles, was given %s", badType)
}

func fixedPointFormatError(id int64, badType string) error {
	return newFormatError(id, "fixed-point clause can only be used on ints, uints, and doubles, was given %s", badType)
}

func hexFormatError(id int64, badType string) error {
	return newFormatError(id, "only ints, uints, bytes, and strings can be formatted as hex, was given %s", badType)
}

func octalFormatError(id int64, badType string) error {
	return newFormatError(id, "octal clause can only be used on ints and uints, was given %s", badType)
}

func scientificFormatError(id int64, badType string) error {
	return newFormatError(id, "scientific clause can only be used on ints, uints, and doubles, was given %s", badType)
}

func stringFormatError(id int64, badType string) error {
	return newFormatError(id, "string clause can only be used on strings, bools, bytes, ints, doubles, maps, lists, types, durations, and timestamps, was given %s", badType)
}

type formatError struct {
	id  int64
	msg string
}

func newFormatError(id int64, msg string, args ...any) error {
	return formatError{
		id:  id,
		msg: fmt.Sprintf(msg, args...),
	}
}

func (e formatError) Error() string {
	return e.msg
}

func (e formatError) Is(target error) bool {
	return e.msg == target.Error()
}

// stringArgList implements the formatListArgs interface.
type stringArgList struct {
	args traits.Lister
}

func (c *stringArgList) Arg(index int64) (ref.Val, error) {
	if index >= c.args.Size().Value().(int64) {
		return nil, fmt.Errorf("index %d out of range", index)
	}
	return c.args.Get(types.Int(index)), nil
}

func (c *stringArgList) Size() int64 {
	return c.args.Size().Value().(int64)
}

// formatStringInterpolator is an interface that allows user-defined behavior
// for formatting clause implementations, as well as argument retrieval.
// Each function is expected to support the appropriate types as laid out in
// the string.format documentation, and to return an error if given an inappropriate type.
type formatStringInterpolator interface {
	// String takes a ref.Val and a string representing the current locale identifier
	// and returns the Val formatted as a string, or an error if one occurred.
	String(ref.Val) (string, error)

	// Decimal takes a ref.Val and a string representing the current locale identifier
	// and returns the Val formatted as a decimal integer, or an error if one occurred.
	Decimal(ref.Val) (string, error)

	// Fixed takes an int pointer representing precision (or nil if none was given) and
	// returns a function operating in a similar manner to String and Decimal, taking a
	// ref.Val and locale and returning the appropriate string. A closure is returned
	// so precision can be set without needing an additional function call/configuration.
	Fixed(int) func(ref.Val) (string, error)

	// Scientific functions identically to Fixed, except the string returned from the closure
	// is expected to be in scientific notation.
	Scientific(int) func(ref.Val) (string, error)

	// Binary takes a ref.Val and a string representing the current locale identifier
	// and returns the Val formatted as a binary integer, or an error if one occurred.
	Binary(ref.Val) (string, error)

	// Hex takes a boolean that, if true, indicates the hex string output by the returned
	// closure should use uppercase letters for A-F.
	Hex(bool) func(ref.Val) (string, error)

	// Octal takes a ref.Val and a string representing the current locale identifier and
	// returns the Val formatted in octal, or an error if one occurred.
	Octal(ref.Val) (string, error)
}

// formatListArgs is an interface that allows user-defined list-like datatypes to be used
// for formatting clause implementations.
type formatListArgs interface {
	// Arg returns the ref.Val at the given index, or an error if one occurred.
	Arg(int64) (ref.Val, error)

	// Size returns the length of the argument list.
	Size() int64
}

// parseFormatString formats a string according to the string.format syntax, taking the clause implementations
// from the provided FormatCallback and the args from the given FormatList.
func parseFormatString(formatStr string, callback formatStringInterpolator, list formatListArgs) (string, error) {
	i := 0
	argIndex := 0
	var builtStr strings.Builder
	for i < len(formatStr) {
		if formatStr[i] == '%' {
			if i+1 < len(formatStr) && formatStr[i+1] == '%' {
				err := builtStr.WriteByte('%')
				if err != nil {
					return "", fmt.Errorf("error writing format string: %w", err)
				}
				i += 2
				continue
			} else {
				argAny, err := list.Arg(int64(argIndex))
				if err != nil {
					return "", err
				}
				if i+1 >= len(formatStr) {
					return "", errors.New("unexpected end of string")
				}
				if int64(argIndex) >= list.Size() {
					return "", fmt.Errorf("index %d out of range", argIndex)
				}
				numRead, val, refErr := parseAndFormatClause(formatStr[i:], argAny, callback, list)
				if refErr != nil {
					return "", refErr
				}
				_, err = builtStr.WriteString(val)
				if err != nil {
					return "", fmt.Errorf("error writing format string: %w", err)
				}
				i += numRead
				argIndex++
			}
		} else {
			err := builtStr.WriteByte(formatStr[i])
			if err != nil {
				return "", fmt.Errorf("error writing format string: %w", err)
			}
			i++
		}
	}
	return builtStr.String(), nil
}

// parseAndFormatClause parses the format clause at the start of the given string with val, and returns
// how many characters were consumed and the substituted string form of val, or an error if one occurred.
func parseAndFormatClause(formatStr string, val ref.Val, callback formatStringInterpolator, list formatListArgs) (int, string, error) {
	i := 1
	read, formatter, err := parseFormattingClause(formatStr[i:], callback)
	i += read
	if err != nil {
		return -1, "", newParseFormatError("could not parse formatting clause", err)
	}

	valStr, err := formatter(val)
	if err != nil {
		return -1, "", newParseFormatError("error during formatting", err)
	}
	return i, valStr, nil
}

func parseFormattingClause(formatStr string, callback formatStringInterpolator) (int, clauseImpl, error) {
	i := 0
	read, precision, err := parsePrecision(formatStr[i:])
	i += read
	if err != nil {
		return -1, nil, fmt.Errorf("error while parsing precision: %w", err)
	}
	r := rune(formatStr[i])
	i++
	switch r {
	case 's':
		return i, callback.String, nil
	case 'd':
		return i, callback.Decimal, nil
	case 'f':
		return i, callback.Fixed(precision), nil
	case 'e':
		return i, callback.Scientific(precision), nil
	case 'b':
		return i, callback.Binary, nil
	case 'x', 'X':
		return i, callback.Hex(unicode.IsUpper(r)), nil
	case 'o':
		return i, callback.Octal, nil
	default:
		return -1, nil, fmt.Errorf("unrecognized formatting clause \"%c\"", r)
	}
}

func parsePrecision(formatStr string) (int, int, error) {
	i := 0
	if formatStr[i] != '.' {
		return i, defaultPrecision, nil
	}
	i++
	var buffer strings.Builder
	for {
		if i >= len(formatStr) {
			return -1, -1, errors.New("could not find end of precision specifier")
		}
		if !isASCIIDigit(rune(formatStr[i])) {
			break
		}
		buffer.WriteByte(formatStr[i])
		i++
	}
	precision, err := strconv.Atoi(buffer.String())
	if err != nil {
		return -1, -1, fmt.Errorf("error while converting precision to integer: %w", err)
	}
	if precision < 0 {
		return -1, -1, fmt.Errorf("negative precision: %d", precision)
	}
	return i, precision, nil
}

func isASCIIDigit(r rune) bool {
	return r <= unicode.MaxASCII && unicode.IsDigit(r)
}

type parseFormatError struct {
	msg     string
	wrapped error
}

func newParseFormatError(msg string, wrapped error) error {
	return parseFormatError{msg: msg, wrapped: wrapped}
}

func (e parseFormatError) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.wrapped.Error())
}

func (e parseFormatError) Is(target error) bool {
	return e.Error() == target.Error()
}

func (e parseFormatError) Unwrap() error {
	return e.wrapped
}

const (
	runtimeID = int64(-1)
)
