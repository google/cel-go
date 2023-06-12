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

// Package decls contains function and variable declaration structs and helper methods.
package decls

import (
	"fmt"
	"strings"

	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// Kind indicates a CEL type's kind which is used to differentiate quickly between simple and complex types.
type Kind uint

const (
	// DynKind represents a dynamic type. This kind only exists at type-check time.
	DynKind Kind = iota

	// AnyKind represents a google.protobuf.Any type. This kind only exists at type-check time.
	AnyKind

	// BoolKind represents a boolean type.
	BoolKind

	// BytesKind represents a bytes type.
	BytesKind

	// DoubleKind represents a double type.
	DoubleKind

	// DurationKind represents a CEL duration type.
	DurationKind

	// IntKind represents an integer type.
	IntKind

	// ListKind represents a list type.
	ListKind

	// MapKind represents a map type.
	MapKind

	// NullTypeKind represents a null type.
	NullTypeKind

	// OpaqueKind represents an abstract type which has no accessible fields.
	OpaqueKind

	// StringKind represents a string type.
	StringKind

	// StructKind represents a structured object with typed fields.
	StructKind

	// TimestampKind represents a a CEL time type.
	TimestampKind

	// TypeKind represents the CEL type.
	TypeKind

	// TypeParamKind represents a parameterized type whose type name will be resolved at type-check time, if possible.
	TypeParamKind

	// UintKind represents a uint type.
	UintKind
)

var (
	// AnyType represents the google.protobuf.Any type.
	AnyType = &Type{
		Kind:        AnyKind,
		runtimeType: types.NewTypeValue("google.protobuf.Any"),
	}
	// BoolType represents the bool type.
	BoolType = &Type{
		Kind:        BoolKind,
		runtimeType: types.BoolType,
	}
	// BytesType represents the bytes type.
	BytesType = &Type{
		Kind:        BytesKind,
		runtimeType: types.BytesType,
	}
	// DoubleType represents the double type.
	DoubleType = &Type{
		Kind:        DoubleKind,
		runtimeType: types.DoubleType,
	}
	// DurationType represents the CEL duration type.
	DurationType = &Type{
		Kind:        DurationKind,
		runtimeType: types.DurationType,
	}
	// DynType represents a dynamic CEL type whose type will be determined at runtime from context.
	DynType = &Type{
		Kind:        DynKind,
		runtimeType: types.NewTypeValue("dyn"),
	}
	// IntType represents the int type.
	IntType = &Type{
		Kind:        IntKind,
		runtimeType: types.IntType,
	}
	// NullType represents the type of a null value.
	NullType = &Type{
		Kind:        NullTypeKind,
		runtimeType: types.NullType,
	}
	// StringType represents the string type.
	StringType = &Type{
		Kind:        StringKind,
		runtimeType: types.StringType,
	}
	// TimestampType represents the time type.
	TimestampType = &Type{
		Kind:        TimestampKind,
		runtimeType: types.TimestampType,
	}
	// TypeType represents a CEL type
	TypeType = &Type{
		Kind:        TypeKind,
		runtimeType: types.TypeType,
	}
	// UintType represents a uint type.
	UintType = &Type{
		Kind:        UintKind,
		runtimeType: types.UintType,
	}
)

// Type holds a reference to a runtime type with an optional type-checked set of type parameters.
type Type struct {
	// Kind indicates general category of the type.
	Kind Kind

	// Parameters holds the optional type-checked set of type Parameters that are used during static analysis.
	Parameters []*Type

	// runtimeType is the runtime type of the declaration.
	runtimeType ref.Type

	// isAssignableType function determines whether one type is assignable to this type.
	// A nil value for the isAssignableType function falls back to equality of kind, runtimeType, and parameters.
	isAssignableType func(other *Type) bool

	// isAssignableRuntimeType function determines whether the runtime type (with erasure) is assignable to this type.
	// A nil value for the isAssignableRuntimeType function falls back to the equality of the type or type name.
	isAssignableRuntimeType func(other ref.Val) bool
}

// IsType indicates whether two types have the same kind, type name, and parameters.
func (t *Type) IsType(other *Type) bool {
	if t.Kind != other.Kind || len(t.Parameters) != len(other.Parameters) {
		return false
	}
	if t.Kind != TypeParamKind && t.RuntimeTypeName() != other.RuntimeTypeName() {
		return false
	}
	for i, p := range t.Parameters {
		if !p.IsType(other.Parameters[i]) {
			return false
		}
	}
	return true
}

// IsAssignableType determines whether the current type is type-check assignable from the input fromType.
func (t *Type) IsAssignableType(fromType *Type) bool {
	if t.isAssignableType != nil {
		return t.isAssignableType(fromType)
	}
	return t.defaultIsAssignableType(fromType)
}

// IsAssignableRuntimeType determines whether the current type is runtime assignable from the input runtimeType.
//
// At runtime, parameterized types are erased and so a function which type-checks to support a map(string, string)
// will have a runtime assignable type of a map.
func (t *Type) IsAssignableRuntimeType(val ref.Val) bool {
	if t.isAssignableRuntimeType != nil {
		return t.isAssignableRuntimeType(val)
	}
	return t.defaultIsAssignableRuntimeType(val)
}

// RuntimeTypeName indicates the type-erased type name associated with the type.
func (t *Type) RuntimeTypeName() string {
	return t.runtimeType.TypeName()
}

// String returns a human-readable definition of the type name.
func (t *Type) String() string {
	if len(t.Parameters) == 0 {
		return t.runtimeType.TypeName()
	}
	params := make([]string, len(t.Parameters))
	for i, p := range t.Parameters {
		params[i] = p.String()
	}
	return fmt.Sprintf("%s(%s)", t.runtimeType.TypeName(), strings.Join(params, ", "))
}

// isDyn indicates whether the type is dynamic in any way.
func (t *Type) isDyn() bool {
	return t.Kind == DynKind || t.Kind == AnyKind || t.Kind == TypeParamKind
}

// defaultIsAssignableType provides the standard definition of what it means for one type to be assignable to another
// where any of the following may return a true result:
// - The from types are the same instance
// - The target type is dynamic
// - The fromType has the same kind and type name as the target type, and all parameters of the target type
//
//	are IsAssignableType() from the parameters of the fromType.
func (t *Type) defaultIsAssignableType(fromType *Type) bool {
	if t == fromType || t.isDyn() {
		return true
	}
	if t.Kind != fromType.Kind ||
		t.runtimeType.TypeName() != fromType.runtimeType.TypeName() ||
		len(t.Parameters) != len(fromType.Parameters) {
		return false
	}
	for i, tp := range t.Parameters {
		fp := fromType.Parameters[i]
		if !tp.IsAssignableType(fp) {
			return false
		}
	}
	return true
}

// defaultIsAssignableRuntimeType inspects the type and in the case of list and map elements, the key and element types
// to determine whether a ref.Val is assignable to the declared type for a function signature.
func (t *Type) defaultIsAssignableRuntimeType(val ref.Val) bool {
	valType := val.Type()
	if !(t.runtimeType == valType || t.isDyn() || t.runtimeType.TypeName() == valType.TypeName()) {
		return false
	}
	switch t.runtimeType {
	case types.ListType:
		elemType := t.Parameters[0]
		l := val.(traits.Lister)
		if l.Size() == types.IntZero {
			return true
		}
		it := l.Iterator()
		for it.HasNext() == types.True {
			elemVal := it.Next()
			return elemType.IsAssignableRuntimeType(elemVal)
		}
	case types.MapType:
		keyType := t.Parameters[0]
		elemType := t.Parameters[1]
		m := val.(traits.Mapper)
		if m.Size() == types.IntZero {
			return true
		}
		it := m.Iterator()
		for it.HasNext() == types.True {
			keyVal := it.Next()
			elemVal := m.Get(keyVal)
			return keyType.IsAssignableRuntimeType(keyVal) && elemType.IsAssignableRuntimeType(elemVal)
		}
	}
	return true
}

// ListType creates an instances of a list type value with the provided element type.
func ListType(elemType *Type) *Type {
	return &Type{
		Kind:        ListKind,
		runtimeType: types.ListType,
		Parameters:  []*Type{elemType},
	}
}

// MapType creates an instance of a map type value with the provided key and value types.
func MapType(keyType, valueType *Type) *Type {
	return &Type{
		Kind:        MapKind,
		runtimeType: types.MapType,
		Parameters:  []*Type{keyType, valueType},
	}
}

// NullableType creates an instance of a nullable type with the provided wrapped type.
//
// Note: only primitive types are supported as wrapped types.
func NullableType(wrapped *Type) *Type {
	return &Type{
		Kind:        wrapped.Kind,
		runtimeType: wrapped.runtimeType,
		Parameters:  wrapped.Parameters,
		isAssignableType: func(other *Type) bool {
			return NullType.IsAssignableType(other) || wrapped.IsAssignableType(other)
		},
		isAssignableRuntimeType: func(other ref.Val) bool {
			return NullType.IsAssignableRuntimeType(other) || wrapped.IsAssignableRuntimeType(other)
		},
	}
}

// OptionalType creates an abstract parameterized type instance corresponding to CEL's notion of optional.
func OptionalType(param *Type) *Type {
	return OpaqueType("optional", param)
}

// OpaqueType creates an abstract parameterized type with a given name.
func OpaqueType(name string, params ...*Type) *Type {
	return &Type{
		Kind:        OpaqueKind,
		runtimeType: types.NewTypeValue(name),
		Parameters:  params,
	}
}

// ObjectType creates a type references to an externally defined type, e.g. a protobuf message type.
func ObjectType(typeName string) *Type {
	// Function sanitizes object types on the fly
	if wkt, found := checkedWellKnowns[typeName]; found {
		return wkt
	}
	return &Type{
		Kind:        StructKind,
		runtimeType: types.NewObjectTypeValue(typeName),
	}
}

// TypeParamType creates a parameterized type instance.
func TypeParamType(paramName string) *Type {
	return &Type{
		Kind:        TypeParamKind,
		runtimeType: types.NewTypeValue(paramName),
	}
}

// TypeTypeWithParam creates a type with a type parameter.
// Used for type-checking purposes, but equivalent to TypeType otherwise.
func TypeTypeWithParam(param *Type) *Type {
	return &Type{
		Kind:        TypeKind,
		Parameters:  []*Type{param},
		runtimeType: types.TypeType,
	}
}

// NewFunction creates a new function declaration with a set of function options to configure overloads
// and function definitions (implementations).
//
// Functions are checked for name collisions and singleton redefinition.
func NewFunction(name string, opts ...FunctionOpt) (*FunctionDecl, error) {
	fn := &FunctionDecl{
		Name:      name,
		Overloads: map[string]*OverloadDecl{},
	}
	var err error
	for _, opt := range opts {
		fn, err = opt(fn)
		if err != nil {
			return nil, err
		}
	}
	if len(fn.Overloads) == 0 {
		return nil, fmt.Errorf("function %s must have at least one overload", name)
	}
	return fn, nil
}

// FunctionDecl defines a function name, overload set, and optionally a singleton definition for all
// overload instances.
type FunctionDecl struct {
	// Name of the function in human-readable terms, e.g. 'contains' of 'math.least'
	Name string

	// Overloads associated with the function name.
	Overloads map[string]*OverloadDecl

	// Singleton implementation of the function for all overloads.
	//
	// If this option is set, an error will occur if any overloads specify a per-overload implementation
	// or if another function with the same name attempts to redefine the singleton.
	Singleton *functions.Overload

	// disableTypeGuards is a performance optimization to disable detailed runtime type checks which could
	// add overhead on common operations. Setting this option true leaves error checks and argument checks
	// intact.
	disableTypeGuards bool

	// declarationDisabled indicates that the binding should be provided on the runtime, but the method should
	// not be exposed as a declaration available for use.
	declarationDisabled bool
}

// IsDeclarationDisabled indicates that the function implementation should be added to the dispatcher, but the
// declaration should not be exposed for use in expressions.
func (f *FunctionDecl) IsDeclarationDisabled() bool {
	return f.declarationDisabled
}

// Merge combines an existing function declaration with another.
//
// If a function is extended, by say adding new overloads to an existing function, then it is merged with the
// prior definition of the function at which point its overloads must not collide with pre-existing overloads
// and its bindings (singleton, or per-overload) must not conflict with previous definitions either.
func (f *FunctionDecl) Merge(other *FunctionDecl) (*FunctionDecl, error) {
	if f == other {
		return f, nil
	}
	if f.Name != other.Name {
		return nil, fmt.Errorf("cannot merge unrelated functions. %s and %s", f.Name, other.Name)
	}
	merged := &FunctionDecl{
		Name:      f.Name,
		Overloads: make(map[string]*OverloadDecl, len(f.Overloads)),
		Singleton: f.Singleton,
	}
	for oID, o := range f.Overloads {
		merged.Overloads[oID] = o
	}
	for _, o := range other.Overloads {
		err := merged.AddOverload(o)
		if err != nil {
			return nil, fmt.Errorf("function declaration merge failed: %v", err)
		}
	}
	if other.Singleton != nil {
		if merged.Singleton != nil {
			return nil, fmt.Errorf("function already has singleton binding: %s", f.Name)
		}
		merged.Singleton = other.Singleton
	}
	return merged, nil
}

// AddOverload ensures that the new overload does not collide with an existing overload signature;
// however, if the function signatures are identical, the implementation may be rewritten as its
// difficult to compare functions by object identity.
func (f *FunctionDecl) AddOverload(overload *OverloadDecl) error {
	for oID, o := range f.Overloads {
		if o.ID != overload.ID && o.SignatureOverlaps(overload) {
			return fmt.Errorf("overload signature collision in function %s: %s collides with %s", f.Name, o.ID, overload.ID)
		}
		if o.ID == overload.ID {
			if o.SignatureEquals(overload) && o.NonStrict == overload.NonStrict {
				// Allow redefinition of an overload implementation so long as the signatures match.
				f.Overloads[oID] = overload
				return nil
			}
			return fmt.Errorf("overload redefinition in function. %s: %s has multiple definitions", f.Name, o.ID)
		}
	}
	f.Overloads[overload.ID] = overload
	return nil
}

// Bindings produces a set of function bindings, if any are defined.
func (f *FunctionDecl) Bindings() ([]*functions.Overload, error) {
	overloads := []*functions.Overload{}
	nonStrict := false
	for _, o := range f.Overloads {
		if o.hasBinding() {
			overload := &functions.Overload{
				Operator:     o.ID,
				Unary:        o.guardedUnaryOp(f.Name, f.disableTypeGuards),
				Binary:       o.guardedBinaryOp(f.Name, f.disableTypeGuards),
				Function:     o.guardedFunctionOp(f.Name, f.disableTypeGuards),
				OperandTrait: o.OperandTrait,
				NonStrict:    o.NonStrict,
			}
			overloads = append(overloads, overload)
			nonStrict = nonStrict || o.NonStrict
		}
	}
	if f.Singleton != nil {
		if len(overloads) != 0 {
			return nil, fmt.Errorf("singleton function incompatible with specialized overloads: %s", f.Name)
		}
		overloads = []*functions.Overload{
			{
				Operator:     f.Name,
				Unary:        f.Singleton.Unary,
				Binary:       f.Singleton.Binary,
				Function:     f.Singleton.Function,
				OperandTrait: f.Singleton.OperandTrait,
			},
		}
		// fall-through to return single overload case.
	}
	if len(overloads) == 0 {
		return overloads, nil
	}
	// Single overload. Replicate an entry for it using the function name as well.
	if len(overloads) == 1 {
		if overloads[0].Operator == f.Name {
			return overloads, nil
		}
		return append(overloads, &functions.Overload{
			Operator:     f.Name,
			Unary:        overloads[0].Unary,
			Binary:       overloads[0].Binary,
			Function:     overloads[0].Function,
			NonStrict:    overloads[0].NonStrict,
			OperandTrait: overloads[0].OperandTrait,
		}), nil
	}
	// All of the defined overloads are wrapped into a top-level function which
	// performs dynamic dispatch to the proper overload based on the argument types.
	bindings := append([]*functions.Overload{}, overloads...)
	funcDispatch := func(args ...ref.Val) ref.Val {
		for _, o := range f.Overloads {
			// During dynamic dispatch over multiple functions, signature agreement checks
			// are preserved in order to assist with the function resolution step.
			switch len(args) {
			case 1:
				if o.UnaryOp != nil && o.matchesRuntimeSignature( /* disableTypeGuards=*/ false, args...) {
					return o.UnaryOp(args[0])
				}
			case 2:
				if o.BinaryOp != nil && o.matchesRuntimeSignature( /* disableTypeGuards=*/ false, args...) {
					return o.BinaryOp(args[0], args[1])
				}
			}
			if o.FunctionOp != nil && o.matchesRuntimeSignature( /* disableTypeGuards=*/ false, args...) {
				return o.FunctionOp(args...)
			}
			// eventually this will fall through to the noSuchOverload below.
		}
		return MaybeNoSuchOverload(f.Name, args...)
	}
	function := &functions.Overload{
		Operator:  f.Name,
		Function:  funcDispatch,
		NonStrict: nonStrict,
	}
	return append(bindings, function), nil
}

// MaybeNoSuchOverload determines whether to propagate an error if one is provided as an argument, or
// to return an unknown set, or to produce a new error for a missing function signature.
func MaybeNoSuchOverload(funcName string, args ...ref.Val) ref.Val {
	argTypes := make([]string, len(args))
	var unk types.Unknown
	for i, arg := range args {
		if types.IsError(arg) {
			return arg
		}
		if types.IsUnknown(arg) {
			unk = append(unk, arg.(types.Unknown)...)
		}
		argTypes[i] = arg.Type().TypeName()
	}
	if len(unk) != 0 {
		return unk
	}
	signature := strings.Join(argTypes, ", ")
	return types.NewErr("no such overload: %s(%s)", funcName, signature)
}

// FunctionOpt defines a functional option for mutating a function declaration.
type FunctionOpt func(*FunctionDecl) (*FunctionDecl, error)

// DisableTypeGuards disables automatically generated function invocation guards on direct overload calls.
// Type guards remain on during dynamic dispatch for parsed-only expressions.
func DisableTypeGuards(value bool) FunctionOpt {
	return func(fn *FunctionDecl) (*FunctionDecl, error) {
		fn.disableTypeGuards = value
		return fn, nil
	}
}

// DisableDeclaration indicates that the function declaration should be disabled, but the runtime function
// binding should be provided. Marking a function as runtime-only is a safe way to manage deprecations
// of function declarations while still preserving the runtime behavior for previously compiled expressions.
func DisableDeclaration(value bool) FunctionOpt {
	return func(fn *FunctionDecl) (*FunctionDecl, error) {
		fn.declarationDisabled = value
		return fn, nil
	}
}

// SingletonUnaryBinding creates a singleton function definition to be used for all function overloads.
//
// Note, this approach works well if operand is expected to have a specific trait which it implements,
// e.g. traits.ContainerType. Otherwise, prefer per-overload function bindings.
func SingletonUnaryBinding(fn functions.UnaryOp, traits ...int) FunctionOpt {
	trait := 0
	for _, t := range traits {
		trait = trait | t
	}
	return func(f *FunctionDecl) (*FunctionDecl, error) {
		if f.Singleton != nil {
			return nil, fmt.Errorf("function already has a singleton binding: %s", f.Name)
		}
		f.Singleton = &functions.Overload{
			Operator:     f.Name,
			Unary:        fn,
			OperandTrait: trait,
		}
		return f, nil
	}
}

// SingletonBinaryBinding creates a singleton function definition to be used with all function overloads.
//
// Note, this approach works well if operand is expected to have a specific trait which it implements,
// e.g. traits.ContainerType. Otherwise, prefer per-overload function bindings.
func SingletonBinaryBinding(fn functions.BinaryOp, traits ...int) FunctionOpt {
	trait := 0
	for _, t := range traits {
		trait = trait | t
	}
	return func(f *FunctionDecl) (*FunctionDecl, error) {
		if f.Singleton != nil {
			return nil, fmt.Errorf("function already has a singleton binding: %s", f.Name)
		}
		f.Singleton = &functions.Overload{
			Operator:     f.Name,
			Binary:       fn,
			OperandTrait: trait,
		}
		return f, nil
	}
}

// SingletonFunctionBinding creates a singleton function definition to be used with all function overloads.
//
// Note, this approach works well if operand is expected to have a specific trait which it implements,
// e.g. traits.ContainerType. Otherwise, prefer per-overload function bindings.
func SingletonFunctionBinding(fn functions.FunctionOp, traits ...int) FunctionOpt {
	trait := 0
	for _, t := range traits {
		trait = trait | t
	}
	return func(f *FunctionDecl) (*FunctionDecl, error) {
		if f.Singleton != nil {
			return nil, fmt.Errorf("function already has a singleton binding: %s", f.Name)
		}
		f.Singleton = &functions.Overload{
			Operator:     f.Name,
			Function:     fn,
			OperandTrait: trait,
		}
		return f, nil
	}
}

// Overload defines a new global overload with an overload id, argument types, and result type. Through the
// use of OverloadOpt options, the overload may also be configured with a binding, an operand trait, and to
// be non-strict.
//
// Note: function bindings should be commonly configured with Overload instances whereas operand traits and
// strict-ness should be rare occurrences.
func Overload(overloadID string, args []*Type, resultType *Type, opts ...OverloadOpt) FunctionOpt {
	return newOverload(overloadID, false, args, resultType, opts...)
}

// MemberOverload defines a new receiver-style overload (or member function) with an overload id, argument types,
// and result type. Through the use of OverloadOpt options, the overload may also be configured with a binding,
// an operand trait, and to be non-strict.
//
// Note: function bindings should be commonly configured with Overload instances whereas operand traits and
// strict-ness should be rare occurrences.
func MemberOverload(overloadID string, args []*Type, resultType *Type, opts ...OverloadOpt) FunctionOpt {
	return newOverload(overloadID, true, args, resultType, opts...)
}

func newOverload(overloadID string, memberFunction bool, args []*Type, resultType *Type, opts ...OverloadOpt) FunctionOpt {
	return func(f *FunctionDecl) (*FunctionDecl, error) {
		overload, err := newOverloadInternal(overloadID, memberFunction, args, resultType, opts...)
		if err != nil {
			return nil, err
		}
		err = f.AddOverload(overload)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}

func newOverloadInternal(overloadID string, memberFunction bool, args []*Type, resultType *Type, opts ...OverloadOpt) (*OverloadDecl, error) {
	overload := &OverloadDecl{
		ID:               overloadID,
		ArgTypes:         args,
		ResultType:       resultType,
		IsMemberFunction: memberFunction,
	}
	var err error
	for _, opt := range opts {
		overload, err = opt(overload)
		if err != nil {
			return nil, err
		}
	}
	return overload, nil
}

// OverloadDecl contains the definition of a single overload id with a specific signature, and an optional
// implementation.
type OverloadDecl struct {
	// ID mirrors the overload signature and provides a unique id which may be referenced within the type-checker
	// and interpreter to optimize performance.
	//
	// The ID format is usually one of two styles:
	// global: <functionName>_<argType>_<argTypeN>
	// member: <memberType>_<functionName>_<argType>_<argTypeN>
	ID string

	// ArgTypes contains the set of argument types expected by the overload.
	//
	// For member functions ArgTypes[0] represents the member operand type.
	ArgTypes []*Type

	// ResultType indicates the output type from calling the function.
	ResultType *Type

	// IsMemberFunction indicates whether the overload is a member function
	IsMemberFunction bool

	// Function implementation options. Optional, but encouraged.
	// UnaryOp is a function binding that takes a single argument.
	UnaryOp functions.UnaryOp
	// BinaryOp is a function binding that takes two arguments.
	BinaryOp functions.BinaryOp
	// FunctionOp is a catch-all for zero-arity and three-plus arity functions.
	FunctionOp functions.FunctionOp

	// NonStrict indicates that the function will accept error and unknown arguments as inputs.
	NonStrict bool

	// OperandTrait indicates whether the member argument should have a specific type-trait.
	//
	// This is useful for creating overloads which operate on a type-interface rather than a concrete type.
	OperandTrait int
}

// SignatureEquals determines whether the incoming overload declaration signature is equal to the current signature.
//
// Result type, operand trait, and strict-ness are not considered as part of signature equality.
func (o *OverloadDecl) SignatureEquals(other *OverloadDecl) bool {
	if o.ID != other.ID || o.IsMemberFunction != other.IsMemberFunction || len(o.ArgTypes) != len(other.ArgTypes) {
		return false
	}
	for i, at := range o.ArgTypes {
		oat := other.ArgTypes[i]
		if !at.IsType(oat) {
			return false
		}
	}
	return o.ResultType.IsType(other.ResultType)
}

// SignatureOverlaps indicates whether two functions have non-equal, but overloapping function signatures.
//
// For example, list(dyn) collides with list(string) since the 'dyn' type can contain a 'string' type.
func (o *OverloadDecl) SignatureOverlaps(other *OverloadDecl) bool {
	if o.IsMemberFunction != other.IsMemberFunction || len(o.ArgTypes) != len(other.ArgTypes) {
		return false
	}
	argsOverlap := true
	for i, argType := range o.ArgTypes {
		otherArgType := other.ArgTypes[i]
		argsOverlap = argsOverlap &&
			(argType.IsAssignableType(otherArgType) ||
				otherArgType.IsAssignableType(argType))
	}
	return argsOverlap
}

// hasBinding indicates whether the overload already has a definition.
func (o *OverloadDecl) hasBinding() bool {
	return o.UnaryOp != nil || o.BinaryOp != nil || o.FunctionOp != nil
}

// guardedUnaryOp creates an invocation guard around the provided unary operator, if one is defined.
func (o *OverloadDecl) guardedUnaryOp(funcName string, disableTypeGuards bool) functions.UnaryOp {
	if o.UnaryOp == nil {
		return nil
	}
	return func(arg ref.Val) ref.Val {
		if !o.matchesRuntimeUnarySignature(disableTypeGuards, arg) {
			return MaybeNoSuchOverload(funcName, arg)
		}
		return o.UnaryOp(arg)
	}
}

// guardedBinaryOp creates an invocation guard around the provided binary operator, if one is defined.
func (o *OverloadDecl) guardedBinaryOp(funcName string, disableTypeGuards bool) functions.BinaryOp {
	if o.BinaryOp == nil {
		return nil
	}
	return func(arg1, arg2 ref.Val) ref.Val {
		if !o.matchesRuntimeBinarySignature(disableTypeGuards, arg1, arg2) {
			return MaybeNoSuchOverload(funcName, arg1, arg2)
		}
		return o.BinaryOp(arg1, arg2)
	}
}

// guardedFunctionOp creates an invocation guard around the provided variadic function binding, if one is provided.
func (o *OverloadDecl) guardedFunctionOp(funcName string, disableTypeGuards bool) functions.FunctionOp {
	if o.FunctionOp == nil {
		return nil
	}
	return func(args ...ref.Val) ref.Val {
		if !o.matchesRuntimeSignature(disableTypeGuards, args...) {
			return MaybeNoSuchOverload(funcName, args...)
		}
		return o.FunctionOp(args...)
	}
}

// matchesRuntimeUnarySignature indicates whether the argument type is runtime assiganble to the overload's expected argument.
func (o *OverloadDecl) matchesRuntimeUnarySignature(disableTypeGuards bool, arg ref.Val) bool {
	return matchRuntimeArgType(o.NonStrict, disableTypeGuards, o.ArgTypes[0], arg) &&
		matchOperandTrait(o.OperandTrait, arg)
}

// matchesRuntimeBinarySignature indicates whether the argument types are runtime assiganble to the overload's expected arguments.
func (o *OverloadDecl) matchesRuntimeBinarySignature(disableTypeGuards bool, arg1, arg2 ref.Val) bool {
	return matchRuntimeArgType(o.NonStrict, disableTypeGuards, o.ArgTypes[0], arg1) &&
		matchRuntimeArgType(o.NonStrict, disableTypeGuards, o.ArgTypes[1], arg2) &&
		matchOperandTrait(o.OperandTrait, arg1)
}

// matchesRuntimeSignature indicates whether the argument types are runtime assiganble to the overload's expected arguments.
func (o *OverloadDecl) matchesRuntimeSignature(disableTypeGuards bool, args ...ref.Val) bool {
	if len(args) != len(o.ArgTypes) {
		return false
	}
	if len(args) == 0 {
		return true
	}
	for i, arg := range args {
		if !matchRuntimeArgType(o.NonStrict, disableTypeGuards, o.ArgTypes[i], arg) {
			return false
		}
	}
	return matchOperandTrait(o.OperandTrait, args[0])
}

func matchRuntimeArgType(nonStrict, disableTypeGuards bool, argType *Type, arg ref.Val) bool {
	if nonStrict && (disableTypeGuards || types.IsUnknownOrError(arg)) {
		return true
	}
	if types.IsUnknownOrError(arg) {
		return false
	}
	return disableTypeGuards || argType.IsAssignableRuntimeType(arg)
}

func matchOperandTrait(trait int, arg ref.Val) bool {
	return trait == 0 || arg.Type().HasTrait(trait) || types.IsUnknownOrError(arg)
}

// OverloadOpt is a functional option for configuring a function overload.
type OverloadOpt func(*OverloadDecl) (*OverloadDecl, error)

// UnaryBinding provides the implementation of a unary overload. The provided function is protected by a runtime
// type-guard which ensures runtime type agreement between the overload signature and runtime argument types.
func UnaryBinding(binding functions.UnaryOp) OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		if o.hasBinding() {
			return nil, fmt.Errorf("overload already has a binding: %s", o.ID)
		}
		if len(o.ArgTypes) != 1 {
			return nil, fmt.Errorf("unary function bound to non-unary overload: %s", o.ID)
		}
		o.UnaryOp = binding
		return o, nil
	}
}

// BinaryBinding provides the implementation of a binary overload. The provided function is protected by a runtime
// type-guard which ensures runtime type agreement between the overload signature and runtime argument types.
func BinaryBinding(binding functions.BinaryOp) OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		if o.hasBinding() {
			return nil, fmt.Errorf("overload already has a binding: %s", o.ID)
		}
		if len(o.ArgTypes) != 2 {
			return nil, fmt.Errorf("binary function bound to non-binary overload: %s", o.ID)
		}
		o.BinaryOp = binding
		return o, nil
	}
}

// FunctionBinding provides the implementation of a variadic overload. The provided function is protected by a runtime
// type-guard which ensures runtime type agreement between the overload signature and runtime argument types.
func FunctionBinding(binding functions.FunctionOp) OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		if o.hasBinding() {
			return nil, fmt.Errorf("overload already has a binding: %s", o.ID)
		}
		o.FunctionOp = binding
		return o, nil
	}
}

// OverloadIsNonStrict enables the function to be called with error and unknown argument values.
//
// Note: do not use this option unless absoluately necessary as it should be an uncommon feature.
func OverloadIsNonStrict() OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		o.NonStrict = true
		return o, nil
	}
}

// OverloadOperandTrait configures a set of traits which the first argument to the overload must implement in order to be
// successfully invoked.
func OverloadOperandTrait(trait int) OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		o.OperandTrait = trait
		return o, nil
	}
}

var (
	checkedWellKnowns = map[string]*Type{
		// Wrapper types.
		"google.protobuf.BoolValue":   NullableType(BoolType),
		"google.protobuf.BytesValue":  NullableType(BytesType),
		"google.protobuf.DoubleValue": NullableType(DoubleType),
		"google.protobuf.FloatValue":  NullableType(DoubleType),
		"google.protobuf.Int64Value":  NullableType(IntType),
		"google.protobuf.Int32Value":  NullableType(IntType),
		"google.protobuf.UInt64Value": NullableType(UintType),
		"google.protobuf.UInt32Value": NullableType(UintType),
		"google.protobuf.StringValue": NullableType(StringType),
		// Well-known types.
		"google.protobuf.Any":       AnyType,
		"google.protobuf.Duration":  DurationType,
		"google.protobuf.Timestamp": TimestampType,
		// Json types.
		"google.protobuf.ListValue": ListType(DynType),
		"google.protobuf.NullValue": NullType,
		"google.protobuf.Struct":    MapType(StringType, DynType),
		"google.protobuf.Value":     DynType,
	}
)
