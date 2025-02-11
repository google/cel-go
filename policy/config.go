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

	"google.golang.org/protobuf/proto"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/env"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
)

// NewConfig returns a YAML serializable policy environment.
func NewConfig(e *env.Config) *Config {
	return &Config{Config: e}
}

// Config represents a YAML serializable CEL environment configuration.
type Config struct {
	*env.Config
}

// AsEnvOptions converts the Config value to a collection of cel environment options.
func (c *Config) AsEnvOptions(provider types.Provider) ([]cel.EnvOption, error) {
	envOpts := []cel.EnvOption{}
	// Configure the standard lib subset.
	if c.StdLib != nil {
		if c.StdLib.Disabled {
			envOpts = append(envOpts, func(e *cel.Env) (*cel.Env, error) {
				if !e.HasLibrary("cel.lib.std") {
					return e, nil
				}
				return cel.NewCustomEnv()
			})
		} else {
			envOpts = append(envOpts, func(e *cel.Env) (*cel.Env, error) {
				return cel.NewCustomEnv(cel.StdLib(cel.StdLibSubset(c.StdLib)))
			})
		}
	}

	// Configure the container
	if c.Container != "" {
		envOpts = append(envOpts, cel.Container(c.Container))
	}

	// Configure abbreviations
	for _, imp := range c.Imports {
		envOpts = append(envOpts, cel.Abbrevs(imp.Name))
	}

	// Configure the context variable declaration
	if c.ContextVariable != nil {
		if len(c.Variables) > 0 {
			return nil, errors.New("either the context_variable or the variables may be set, but not both")
		}
		typeName := c.ContextVariable.TypeName
		if typeName == "" {
			return nil, errors.New("invalid context variable, must set type name field")
		}
		if _, found := provider.FindStructType(typeName); !found {
			return nil, fmt.Errorf("could not find context proto type name: %s", typeName)
		}
		// Attempt to instantiate the proto in order to reflect to its descriptor
		msg := provider.NewValue(typeName, map[string]ref.Val{})
		pbMsg, ok := msg.Value().(proto.Message)
		if !ok {
			return nil, fmt.Errorf("type name was not a protobuf: %T", msg.Value())
		}
		envOpts = append(envOpts, cel.DeclareContextProto(pbMsg.ProtoReflect().Descriptor()))
	}

	if len(c.Variables) != 0 {
		vars := make([]*decls.VariableDecl, 0, len(c.Variables))
		for _, v := range c.Variables {
			vDef, err := v.AsCELVariable(provider)
			if err != nil {
				return nil, err
			}
			vars = append(vars, vDef)
		}
		envOpts = append(envOpts, cel.VariableDecls(vars...))
	}
	if len(c.Functions) != 0 {
		funcs := make([]*decls.FunctionDecl, 0, len(c.Functions))
		for _, f := range c.Functions {
			fnDef, err := f.AsCELFunction(provider)
			if err != nil {
				return nil, err
			}
			funcs = append(funcs, fnDef)
		}
		envOpts = append(envOpts, cel.FunctionDecls(funcs...))
	}
	for _, e := range c.Extensions {
		opt, err := extensionEnvOption(e)
		if err != nil {
			return nil, err
		}
		envOpts = append(envOpts, opt)
	}
	return envOpts, nil
}

// extensionEnvOption converts an ExtensionConfig value to a CEL environment option.
func extensionEnvOption(ec *env.Extension) (cel.EnvOption, error) {
	fac, found := extFactories[ec.Name]
	if !found {
		return nil, fmt.Errorf("unrecognized extension: %s", ec.Name)
	}
	// If the version is 'latest', set the version value to the max uint.
	ver, err := ec.GetVersion()
	if err != nil {
		return nil, err
	}
	return fac(ver), nil
}

// extensionFactory accepts a version and produces a CEL environment associated with the versioned extension.
type extensionFactory func(uint32) cel.EnvOption

var extFactories = map[string]extensionFactory{
	"bindings": func(version uint32) cel.EnvOption {
		return ext.Bindings(ext.BindingsVersion(version))
	},
	"encoders": func(version uint32) cel.EnvOption {
		return ext.Encoders(ext.EncodersVersion(version))
	},
	"lists": func(version uint32) cel.EnvOption {
		return ext.Lists(ext.ListsVersion(version))
	},
	"math": func(version uint32) cel.EnvOption {
		return ext.Math(ext.MathVersion(version))
	},
	"optional": func(version uint32) cel.EnvOption {
		return cel.OptionalTypes(cel.OptionalTypesVersion(version))
	},
	"protos": func(version uint32) cel.EnvOption {
		return ext.Protos(ext.ProtosVersion(version))
	},
	"sets": func(version uint32) cel.EnvOption {
		return ext.Sets(ext.SetsVersion(version))
	},
	"strings": func(version uint32) cel.EnvOption {
		return ext.Strings(ext.StringsVersion(version))
	},
	"two-var-comprehensions": func(version uint32) cel.EnvOption {
		return ext.TwoVarComprehensions(ext.TwoVarComprehensionsVersion(version))
	},
}
