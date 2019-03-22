// Copyright 2019 Google LLC
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

package cel

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/parser"

	descpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// EnvOption is a functional interface for configuring the environment.
type EnvOption func(e *env) (*env, error)

// ClearBuiltIns option removes all standard types, operators, and macros from the environment.
//
// Note: This option must be specified before Declarations and/or Macros if used together.
func ClearBuiltIns() EnvOption {
	return func(e *env) (*env, error) {
		e.declarations = []*exprpb.Decl{}
		e.macros = parser.NoMacros
		e.enableBuiltins = false
		return e, nil
	}
}

// ClearMacros options clears all parser macros.
//
// Clearing macros will ensure CEL expressions can only contain linear evaluation paths, as
// comprehensions such as `all` and `exists` are enabled only via macros.
//
// Note: This option is a no-op when used with ClearBuiltIns, and must be used before Macros
// if used together.
func ClearMacros() EnvOption {
	return func(e *env) (*env, error) {
		e.macros = parser.NoMacros
		return e, nil
	}
}

// CustomTypeAdapter swaps the default ref.TypeAdapter implementation with a custom one.
//
// Note: This option must be specified before the Types and TypeDescs options when used together.
func CustomTypeAdapter(adapter ref.TypeAdapter) EnvOption {
	return func(e *env) (*env, error) {
		e.adapter = adapter
		return e, nil
	}
}

// CustomTypeProvider swaps the default ref.TypeProvider implementation with a custom one.
//
// Note: This option must be specified before the Types and TypeDescs options when used together.
func CustomTypeProvider(provider ref.TypeProvider) EnvOption {
	return func(e *env) (*env, error) {
		e.provider = provider
		return e, nil
	}
}

// Declarations option extends the declaration set configured in the environment.
//
// Note: This option must be specified after ClearBuiltIns if both are used together.
func Declarations(decls ...*exprpb.Decl) EnvOption {
	// TODO: provide an alternative means of specifying declarations that doesn't refer
	// to the underlying proto implementations.
	return func(e *env) (*env, error) {
		e.declarations = append(e.declarations, decls...)
		return e, nil
	}
}

// HomogeneousAggregateLiterals option ensures that list and map literal entry types must agree
// during type-checking.
//
// Note, it is still possible to have heterogeneous aggregates when provided as variables to the
// expression, as well as via conversion of well-known dynamic types, or with unchecked
// expressions.
func HomogeneousAggregateLiterals() EnvOption {
	return func(e *env) (*env, error) {
		e.enableDynamicAggregateLiterals = false
		return e, nil
	}
}

// Macros option extends the macro set configured in the environment.
//
// Note: This option must be specified after ClearBuiltIns and/or ClearMacros if used together.
func Macros(macros ...parser.Macro) EnvOption {
	return func(e *env) (*env, error) {
		e.macros = append(e.macros, macros...)
		return e, nil
	}
}

// Container sets the container for resolving variable names. Defaults to an empty container.
//
// If all references within an expression are relative to a protocol buffer package, then
// specifying a container of `google.type` would make it possible to write expressions such as
// `Expr{expression: 'a < b'}` instead of having to write `google.type.Expr{...}`.
func Container(pkg string) EnvOption {
	return func(e *env) (*env, error) {
		e.pkg = packages.NewPackage(pkg)
		return e, nil
	}
}

// Types adds one or more type declarations to the environment, allowing for construction of
// type-literals whose definitions are included in the common expression built-in set.
//
// The input types may either be instances of `proto.Message` or `ref.Type`. Any other type
// provided to this option will result in an error.
//
// Well-known protobuf types within the `google.protobuf.*` package are included in the standard
// environment by default.
//
// Note: This option must be specified after the CustomTypeProvider option when used together.
func Types(addTypes ...interface{}) EnvOption {
	return func(e *env) (*env, error) {
		reg, isReg := e.provider.(ref.TypeRegistry)
		if !isReg {
			return nil, fmt.Errorf("custom types not supported by provider: %T", e.provider)
		}
		for _, t := range addTypes {
			switch t.(type) {
			case proto.Message:
				err := reg.RegisterMessage(t.(proto.Message))
				if err != nil {
					return nil, err
				}
			case ref.Type:
				err := reg.RegisterType(t.(ref.Type))
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unsupported type: %T", t)
			}
		}
		return e, nil
	}
}

// TypeDescs adds type declarations for one or more protocol buffer
// FileDescriptorProtos or FileDescriptorSets.  Note that types added
// via descriptor will not be able to instantiate messages, and so are
// only useful for Check() operations.
func TypeDescs(descs ...interface{}) EnvOption {
	return func(e *env) (*env, error) {
		reg, isReg := e.provider.(ref.TypeRegistry)
		if !isReg {
			return nil, fmt.Errorf("custom types not supported by provider: %T", e.provider)
		}
		for _, d := range descs {
			switch p := d.(type) {
			case *descpb.FileDescriptorSet:
				for _, fd := range p.File {
					err := reg.RegisterDescriptor(fd)
					if err != nil {
						return nil, err
					}
				}
			case *descpb.FileDescriptorProto:
				err := reg.RegisterDescriptor(p)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unsupported type descriptor: %T", d)
			}
		}
		return e, nil
	}
}

// ProgramOption is a functional interface for configuring evaluation bindings and behaviors.
type ProgramOption func(p *prog) (*prog, error)

// Functions adds function overloads that extend or override the set of CEL built-ins.
func Functions(funcs ...*functions.Overload) ProgramOption {
	return func(p *prog) (*prog, error) {
		if err := p.dispatcher.Add(funcs...); err != nil {
			return nil, err
		}
		return p, nil
	}
}

// Globals sets the global variable values for a given program. These values may be shadowed by
// variables with the same name provided to the Eval() call.
//
// The vars value may either be an `interpreter.Activation` instance or a `map[string]interface{}`.
func Globals(vars interface{}) ProgramOption {
	return func(p *prog) (*prog, error) {
		defaultVars, err :=
			interpreter.NewAdaptingActivation(p.adapter, vars)
		if err != nil {
			return nil, err
		}
		p.defaultVars = defaultVars
		return p, nil
	}
}

// EvalOption indicates an evaluation option that may affect the evaluation behavior or information
// in the output result.
type EvalOption int

const (
	// OptTrackState will cause the runtime to return an immutable EvalState value in the Result.
	OptTrackState EvalOption = 1 << iota

	// OptExhaustiveEval causes the runtime to disable short-circuits and track state.
	OptExhaustiveEval EvalOption = 1<<iota | OptTrackState

	// OptFoldConstants evaluates functions and operators with constants as arguments at program
	// creation time. This flag is useful when the expression will be evaluated repeatedly against
	// a series of different inputs.
	OptFoldConstants EvalOption = 1 << iota
)

// EvalOptions sets one or more evaluation options which may affect the evaluation or Result.
func EvalOptions(opts ...EvalOption) ProgramOption {
	return func(p *prog) (*prog, error) {
		for _, opt := range opts {
			p.evalOpts |= opt
		}
		return p, nil
	}
}
