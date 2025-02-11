// Copyright 2025 Google LLC
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

// Package env provides a representation of a CEL environment.
package env

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/types"
)

// NewConfig creates an instance of a YAML serializable CEL environment configuration.
func NewConfig() *Config {
	return &Config{
		Imports:    []*Import{},
		Extensions: []*Extension{},
		Variables:  []*Variable{},
		Functions:  []*Function{},
	}
}

// Config represents a serializable form of the CEL environment configuration.
//
// Note: custom validations, feature flags, and performance tuning parameters are
// not (yet) considered part of the core CEL environment configuration and should
// be managed separately until a common convention for such configuration settings
// can be developed.
type Config struct {
	Name            string           `yaml:"name"`
	Description     string           `yaml:"description,omitempty"`
	Container       string           `yaml:"container,omitempty"`
	Imports         []*Import        `yaml:"imports,omitempty"`
	StdLib          *LibrarySubset   `yaml:"stdlib,omitempty"`
	Extensions      []*Extension     `yaml:"extensions,omitempty"`
	ContextVariable *ContextVariable `yaml:"context_variable,omitempty"`
	Variables       []*Variable      `yaml:"variables,omitempty"`
	Functions       []*Function      `yaml:"functions,omitempty"`
}

// Import represents a type name that will be appreviated by its simple name using
// the cel.Abbrevs() option.
type Import struct {
	Name string `yaml:"name"`
}

// Variable represents a typed variable declaration which will be published via the
// cel.VariableDecls() option.
type Variable struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	// Type represents the type declaration for the variable.
	//
	// Deprecated: use the embedded *TypeDesc fields directly.
	Type *TypeDesc `yaml:"type,omitempty"`

	// TypeDesc is an embedded set of fields allowing for the specification of the Variable type.
	*TypeDesc `yaml:",inline"`
}

// GetType returns the variable type description.
//
// Note, if both the embedded TypeDesc and the field Type are non-nil, the embedded TypeDesc will
// take precedence.
func (vd *Variable) GetType() *TypeDesc {
	if vd == nil {
		return nil
	}
	if vd.TypeDesc != nil {
		return vd.TypeDesc
	}
	if vd.Type != nil {
		return vd.Type
	}
	return nil
}

// AsCELVariable converts the serializable form of the Variable into a CEL environment declaration.
func (vd *Variable) AsCELVariable(tp types.Provider) (*decls.VariableDecl, error) {
	if vd == nil {
		return nil, errors.New("nil Variable cannot be converted to a VariableDecl")
	}
	if vd.Name == "" {
		return nil, errors.New("invalid variable, must declare a name")
	}
	if vd.GetType() != nil {
		t, err := vd.GetType().AsCELType(tp)
		if err != nil {
			return nil, fmt.Errorf("invalid variable type for '%s': %w", vd.Name, err)
		}
		return decls.NewVariable(vd.Name, t), nil
	}
	return nil, fmt.Errorf("invalid variable '%s', no type specified", vd.Name)
}

// ContextVariable represents a structured message whose fields are to be treated as the top-level
// variable identifiers within CEL expressions.
type ContextVariable struct {
	// TypeName represents the fully qualified typename of the context variable.
	// Currently, only protobuf types are supported.
	TypeName string `yaml:"type_name"`
}

// Function represents the serializable format of a function and its overloads.
type Function struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Overloads   []*Overload `yaml:"overloads"`
}

// AsCELFunction converts the serializable form of the Function into CEL environment declaration.
func (fn *Function) AsCELFunction(tp types.Provider) (*decls.FunctionDecl, error) {
	if fn == nil {
		return nil, errors.New("nil Function cannot be converted to a FunctionDecl")
	}
	if fn.Name == "" {
		return nil, errors.New("invalid function, must declare a name")
	}
	if len(fn.Overloads) == 0 {
		return nil, fmt.Errorf("invalid function %s, must declare an overload", fn.Name)
	}
	overloads := make([]decls.FunctionOpt, len(fn.Overloads))
	var err error
	for i, o := range fn.Overloads {
		overloads[i], err = o.AsFunctionOption(tp)
		if err != nil {
			return nil, err
		}
	}
	return decls.NewFunction(fn.Name, overloads...)
}

// Overload represents the serializable format of a function overload.
type Overload struct {
	ID          string      `yaml:"id"`
	Description string      `yaml:"description"`
	Target      *TypeDesc   `yaml:"target"`
	Args        []*TypeDesc `yaml:"args"`
	Return      *TypeDesc   `yaml:"return"`
}

// AsFunctionOption converts the serializable form of the Overload into a function declaration option.
func (od *Overload) AsFunctionOption(tp types.Provider) (decls.FunctionOpt, error) {
	if od == nil {
		return nil, errors.New("nil Overload cannot be converted to a FunctionOpt")
	}
	args := make([]*types.Type, len(od.Args))
	var err error
	for i, a := range od.Args {
		args[i], err = a.AsCELType(tp)
		if err != nil {
			return nil, err
		}
	}
	if od.Return == nil {
		return nil, fmt.Errorf("missing return type on overload: %v", od.ID)
	}
	result, err := od.Return.AsCELType(tp)
	if err != nil {
		return nil, err
	}
	if od.Target != nil {
		t, err := od.Target.AsCELType(tp)
		if err != nil {
			return nil, err
		}
		args = append([]*types.Type{t}, args...)
		return decls.MemberOverload(od.ID, args, result), nil
	}
	return decls.Overload(od.ID, args, result), nil
}

// Extension represents a named and optionally versioned extension library configured in the environment.
type Extension struct {
	// Name is either the LibraryName() or some short-hand simple identifier which is understood by the config-handler.
	Name string `yaml:"name"`

	// Version may either be an unsigned long value or the string 'latest'. If empty, the value is treated as '0'.
	Version string `yaml:"version,omitempty"`
}

// GetVersion returns the parsed version string, or an error if the version cannot be parsed.
func (e *Extension) GetVersion() (uint32, error) {
	if e == nil {
		return 0, errors.New("nil Extension cannot produce a version")
	}
	if e.Version == "latest" {
		return math.MaxUint32, nil
	}
	if e.Version == "" {
		return uint32(0), nil
	}
	ver, err := strconv.ParseUint(e.Version, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing uint version: %w", err)
	}
	return uint32(ver), nil
}

// LibrarySubset indicates a subset of the macros and function supported by a subsettable library.
type LibrarySubset struct {
	// DisableMacros disables macros for the given library.
	DisableMacros bool `yaml:"disable_macros"`

	// IncludeMacros specifies a set of macro function names to include in the subset.
	IncludeMacros []string `yaml:"include_macros"`

	// ExcludeMacros specifies a set of macro function names to exclude from the subset.
	// Note: if IncludeMacros is non-empty, then ExcludeFunctions is ignored.
	ExcludeMacros []string `yaml:"exclude_macros"`

	// IncludeFunctions specifies a set of functions to include in the subset.
	//
	// Note: the overloads specified in the subset need only specify their ID.
	// Note: if IncludeFunctions is non-empty, then ExcludeFunctions is ignored.
	IncludeFunctions []*Function `yaml:"include_functions"`

	// ExcludeFunctions specifies the set of functions to exclude from the subset.
	//
	// Note: the overloads specified in the subset need only specify their ID.
	ExcludeFunctions []*Function `yaml:"exclude_functions"`
}

// SubsetFunction produces a function declaration which matches the supported subset, or nil
// if the function is not supported by the LibrarySubset.
//
// For IncludeFunctions, if the function does not specify a set of overloads to include, the
// whole function definition is included. If overloads are set, then a new function which
// includes only the specified overloads is produced.
//
// For ExcludeFunctions, if the function does not specify a set of overloads to exclude, the
// whole function definition is excluded. If overloads are set, then a new function which
// includes only the permitted overloads is produced.
func (lib *LibrarySubset) SubsetFunction(fn *decls.FunctionDecl) (*decls.FunctionDecl, bool) {
	// When lib is null, it should indicate that all values are included in the subset.
	if lib == nil {
		return fn, true
	}
	if len(lib.IncludeFunctions) != 0 {
		for _, include := range lib.IncludeFunctions {
			if include.Name != fn.Name() {
				continue
			}
			if len(include.Overloads) == 0 {
				return fn, true
			}
			overloadIDs := make([]string, len(include.Overloads))
			for i, o := range include.Overloads {
				overloadIDs[i] = o.ID
			}
			return fn.Subset(decls.IncludeOverloads(overloadIDs...)), true
		}
		return nil, false
	}
	if len(lib.ExcludeFunctions) != 0 {
		for _, exclude := range lib.ExcludeFunctions {
			if exclude.Name != fn.Name() {
				continue
			}
			if len(exclude.Overloads) == 0 {
				return nil, false
			}
			overloadIDs := make([]string, len(exclude.Overloads))
			for i, o := range exclude.Overloads {
				overloadIDs[i] = o.ID
			}
			return fn.Subset(decls.ExcludeOverloads(overloadIDs...)), true
		}
		return fn, true
	}
	return fn, true
}

// SubsetMacro indicates whether the macro function should be included in the library subset.
func (lib *LibrarySubset) SubsetMacro(macroFunction string) bool {
	// When lib is null, it should indicate that all values are included in the subset.
	if lib == nil {
		return true
	}
	if lib.DisableMacros {
		return false
	}
	if len(lib.IncludeMacros) != 0 {
		for _, name := range lib.IncludeMacros {
			if name == macroFunction {
				return true
			}
		}
		return false
	}
	if len(lib.ExcludeMacros) != 0 {
		for _, name := range lib.ExcludeMacros {
			if name == macroFunction {
				return false
			}
		}
		return true
	}
	return true
}

// TypeDesc represents the serializable format of a CEL *types.Type value.
type TypeDesc struct {
	TypeName    string      `yaml:"type_name"`
	Params      []*TypeDesc `yaml:"params"`
	IsTypeParam bool        `yaml:"is_type_param"`
}

// AsCELType converts the serializable object to a *types.Type value.
func (td *TypeDesc) AsCELType(tp types.Provider) (*types.Type, error) {
	if td == nil {
		return nil, errors.New("nil TypeDesc cannot be converted to a Type instance")
	}
	if td.TypeName == "" {
		return nil, errors.New("invalid type description, declare a type name")
	}
	var err error
	switch td.TypeName {
	case "dyn":
		return types.DynType, nil
	case "map":
		if len(td.Params) == 2 {
			kt, err := td.Params[0].AsCELType(tp)
			if err != nil {
				return nil, err
			}
			vt, err := td.Params[1].AsCELType(tp)
			if err != nil {
				return nil, err
			}
			return types.NewMapType(kt, vt), nil
		}
		return nil, fmt.Errorf("map type has unexpected param count: %d", len(td.Params))
	case "list":
		if len(td.Params) == 1 {
			et, err := td.Params[0].AsCELType(tp)
			if err != nil {
				return nil, err
			}
			return types.NewListType(et), nil
		}
		return nil, fmt.Errorf("list type has unexpected param count: %d", len(td.Params))
	case "optional_type":
		if len(td.Params) == 1 {
			et, err := td.Params[0].AsCELType(tp)
			if err != nil {
				return nil, err
			}
			return types.NewOptionalType(et), nil
		}
		return nil, fmt.Errorf("optional_type has unexpected param count: %d", len(td.Params))
	default:
		if td.IsTypeParam {
			return types.NewTypeParamType(td.TypeName), nil
		}
		if msgType, found := tp.FindStructType(td.TypeName); found {
			// First parameter is the type name.
			return msgType.Parameters()[0], nil
		}
		t, found := tp.FindIdent(td.TypeName)
		if !found {
			return nil, fmt.Errorf("undefined type name: %v", td.TypeName)
		}
		_, ok := t.(*types.Type)
		if ok && len(td.Params) == 0 {
			return t.(*types.Type), nil
		}
		params := make([]*types.Type, len(td.Params))
		for i, p := range td.Params {
			params[i], err = p.AsCELType(tp)
			if err != nil {
				return nil, err
			}
		}
		return types.NewOpaqueType(td.TypeName, params...), nil
	}
}
