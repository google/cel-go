// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package policy

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"google.golang.org/protobuf/proto"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
)

// Config represents a YAML serializable CEL environment configuration.
type Config struct {
	Name        string             `yaml:"name"`
	Description string             `yaml:"description"`
	Container   string             `yaml:"container"`
	Extensions  []*ExtensionConfig `yaml:"extensions"`
	Variables   []*VariableDecl    `yaml:"variables"`
	Functions   []*FunctionDecl    `yaml:"functions"`
}

// AsEnvOptions converts the Config value to a collection of cel environment options.
func (c *Config) AsEnvOptions(baseEnv *cel.Env) ([]cel.EnvOption, error) {
	envOpts := []cel.EnvOption{}
	if c.Container != "" {
		envOpts = append(envOpts, cel.Container(c.Container))
	}
	for _, e := range c.Extensions {
		opt, err := e.AsEnvOption(baseEnv)
		if err != nil {
			return nil, err
		}
		envOpts = append(envOpts, opt)
	}
	for _, v := range c.Variables {
		opt, err := v.AsEnvOption(baseEnv)
		if err != nil {
			return nil, err
		}
		envOpts = append(envOpts, opt)
	}
	for _, f := range c.Functions {
		opt, err := f.AsEnvOption(baseEnv)
		if err != nil {
			return nil, err
		}
		envOpts = append(envOpts, opt)
	}
	return envOpts, nil
}

// ExtensionFactory accepts a version number and produces a CEL environment associated with the versioned
// extension.
type ExtensionFactory func(uint32) cel.EnvOption

// ExtensionResolver provides a way to lookup ExtensionFactory instances by extension name.
type ExtensionResolver interface {
	// ResolveExtension returns an ExtensionFactory bound to the given name, if one exists.
	ResolveExtension(name string) (ExtensionFactory, bool)
}

// ExtensionConfig represents a YAML serializable definition of a versioned extension library reference.
type ExtensionConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	ExtensionResolver
}

// AsEnvOption converts an ExtensionConfig value to a CEL environment option.
func (ec *ExtensionConfig) AsEnvOption(baseEnv *cel.Env) (cel.EnvOption, error) {
	fac, found := extFactories[ec.Name]
	if !found && ec.ExtensionResolver != nil {
		fac, found = ec.ResolveExtension(ec.Name)
	}
	if !found {
		return nil, fmt.Errorf("unrecognized extension: %s", ec.Name)
	}
	// If the version is 'latest', set the version value to the max uint.
	if ec.Version == "latest" {
		return fac(math.MaxUint32), nil
	}
	if ec.Version == "" {
		return fac(0), nil
	}
	ver, err := strconv.ParseUint(ec.Version, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing uint version: %w", err)
	}
	return fac(uint32(ver)), nil
}

// VariableDecl represents a YAML serializable CEL variable declaration.
type VariableDecl struct {
	Name         string    `yaml:"name"`
	Type         *TypeDecl `yaml:"type"`
	ContextProto string    `yaml:"context_proto"`
}

// AsEnvOption converts a VariableDecl type to a CEL environment option.
//
// Note, variable definitions with differing type definitions will result in an error during
// the compile step.
func (vd *VariableDecl) AsEnvOption(baseEnv *cel.Env) (cel.EnvOption, error) {
	if vd.Name != "" {
		t, err := vd.Type.AsCELType(baseEnv)
		if err != nil {
			return nil, fmt.Errorf("invalid variable type for '%s': %w", vd.Name, err)
		}
		return cel.Variable(vd.Name, t), nil
	}
	if vd.ContextProto != "" {
		if _, found := baseEnv.CELTypeProvider().FindStructType(vd.ContextProto); !found {
			return nil, fmt.Errorf("could not find context proto type name: %s", vd.ContextProto)
		}
		// Attempt to instantiate the proto in order to reflect to its descriptor
		msg := baseEnv.CELTypeProvider().NewValue(vd.ContextProto, map[string]ref.Val{})
		pbMsg, ok := msg.Value().(proto.Message)
		if !ok {
			return nil, fmt.Errorf("type name was not a protobuf: %T", msg.Value())
		}
		return cel.DeclareContextProto(pbMsg.ProtoReflect().Descriptor()), nil
	}
	return nil, errors.New("invalid variable, must set 'name' or 'context_proto' field")
}

// TypeDecl represents a YAML serializable CEL type reference.
type TypeDecl struct {
	TypeName    string      `yaml:"type_name"`
	Params      []*TypeDecl `yaml:"params"`
	IsTypeParam bool        `yaml:"is_type_param"`
}

// AsCELType converts the TypeDecl value to a cel.Type value using the input base environment.
//
// All extension types referenced by name within the `TypeDecl.TypeName` field must be configured
// within the base CEL environment argument.
func (td *TypeDecl) AsCELType(baseEnv *cel.Env) (*cel.Type, error) {
	var err error
	switch td.TypeName {
	case "dyn":
		return cel.DynType, nil
	case "map":
		if len(td.Params) == 2 {
			kt, err := td.Params[0].AsCELType(baseEnv)
			if err != nil {
				return nil, err
			}
			vt, err := td.Params[1].AsCELType(baseEnv)
			if err != nil {
				return nil, err
			}
			return cel.MapType(kt, vt), nil
		}
		return nil, fmt.Errorf("map type has unexpected param count: %d", len(td.Params))
	case "list":
		if len(td.Params) == 1 {
			et, err := td.Params[0].AsCELType(baseEnv)
			if err != nil {
				return nil, err
			}
			return cel.ListType(et), nil
		}
		return nil, fmt.Errorf("list type has unexpected param count: %d", len(td.Params))
	default:
		if td.IsTypeParam {
			return cel.TypeParamType(td.TypeName), nil
		}
		if msgType, found := baseEnv.CELTypeProvider().FindStructType(td.TypeName); found {
			// First parameter is the type name.
			return msgType.Parameters()[0], nil
		}
		t, found := baseEnv.CELTypeProvider().FindIdent(td.TypeName)
		if !found {
			return nil, fmt.Errorf("undefined type name: %v", td.TypeName)
		}
		_, ok := t.(*cel.Type)
		if ok && len(td.Params) == 0 {
			return t.(*cel.Type), nil
		}
		params := make([]*cel.Type, len(td.Params))
		for i, p := range td.Params {
			params[i], err = p.AsCELType(baseEnv)
			if err != nil {
				return nil, err
			}
		}
		return cel.OpaqueType(td.TypeName, params...), nil
	}
}

// FunctionDecl represents a YAML serializable declaration of a CEL function.
type FunctionDecl struct {
	Name      string          `yaml:"name"`
	Overloads []*OverloadDecl `yaml:"overloads"`
}

// AsEnvOption converts a FunctionDecl value into a cel.EnvOption using the input environment.
func (fd *FunctionDecl) AsEnvOption(baseEnv *cel.Env) (cel.EnvOption, error) {
	overloads := make([]cel.FunctionOpt, len(fd.Overloads))
	var err error
	for i, o := range fd.Overloads {
		overloads[i], err = o.AsFunctionOption(baseEnv)
		if err != nil {
			return nil, err
		}
	}
	return cel.Function(fd.Name, overloads...), nil
}

// OverloadDecl represents a YAML serializable declaration of a CEL function overload.
type OverloadDecl struct {
	OverloadID string      `yaml:"id"`
	Target     *TypeDecl   `yaml:"target"`
	Args       []*TypeDecl `yaml:"args"`
	Return     *TypeDecl   `yaml:"return"`
}

// AsFunctionOption converts an OverloadDecl value into a cel.FunctionOpt using the input environment.
func (od *OverloadDecl) AsFunctionOption(baseEnv *cel.Env) (cel.FunctionOpt, error) {
	args := make([]*cel.Type, len(od.Args))
	var err error
	for i, a := range od.Args {
		args[i], err = a.AsCELType(baseEnv)
		if err != nil {
			return nil, err
		}
	}

	if od.Return == nil {
		return nil, fmt.Errorf("missing return type on overload: %v", od.OverloadID)
	}
	result, err := od.Return.AsCELType(baseEnv)
	if err != nil {
		return nil, err
	}
	if od.Target != nil {
		t, err := od.Target.AsCELType(baseEnv)
		if err != nil {
			return nil, err
		}
		args = append([]*cel.Type{t}, args...)
		return cel.MemberOverload(od.OverloadID, args, result), nil
	}
	return cel.Overload(od.OverloadID, args, result), nil
}

var extFactories = map[string]ExtensionFactory{
	"bindings": func(version uint32) cel.EnvOption {
		return ext.Bindings()
	},
	"encoders": func(version uint32) cel.EnvOption {
		return ext.Encoders()
	},
	"lists": func(version uint32) cel.EnvOption {
		return ext.Lists()
	},
	"math": func(version uint32) cel.EnvOption {
		return ext.Math()
	},
	"optional": func(version uint32) cel.EnvOption {
		return cel.OptionalTypes(cel.OptionalTypesVersion(version))
	},
	"protos": func(version uint32) cel.EnvOption {
		return ext.Protos()
	},
	"sets": func(version uint32) cel.EnvOption {
		return ext.Sets()
	},
	"strings": func(version uint32) cel.EnvOption {
		return ext.Strings(ext.StringsVersion(version))
	},
}
