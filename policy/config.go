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
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
)

type Config struct {
	Name        string             `yaml:"name"`
	Description string             `yaml:"description"`
	Container   string             `yaml:"container"`
	Extensions  []*ExtensionConfig `yaml:"extensions"`
	Variables   []*VariableDecl    `yaml:"variables"`
	Functions   []*FunctionDecl    `yaml:"functions"`
}

func (c *Config) AsEnvOptions(baseEnv *cel.Env) []cel.EnvOption {
	envOpts := []cel.EnvOption{}
	if c.Container != "" {
		envOpts = append(envOpts, cel.Container(c.Container))
	}
	for _, e := range c.Extensions {
		envOpts = append(envOpts, e.AsEnvOption(baseEnv))
	}
	for _, v := range c.Variables {
		envOpts = append(envOpts, v.AsEnvOption(baseEnv))
	}
	for _, f := range c.Functions {
		envOpts = append(envOpts, f.AsEnvOption(baseEnv))
	}
	return envOpts
}

type ExtensionConfig struct {
	Name    string `yaml:"name"`
	Version int    `yaml:"version"`
}

func (ec *ExtensionConfig) AsEnvOption(baseEnv *cel.Env) cel.EnvOption {
	fac, found := extFactories[ec.Name]
	if !found {
		panic(fmt.Sprintf("unrecognized extension: %s", ec.Name))
	}
	return fac(uint32(ec.Version))
}

type VariableDecl struct {
	Name string    `yaml:"name"`
	Type *TypeDecl `yaml:"type"`
}

func (vd *VariableDecl) AsEnvOption(baseEnv *cel.Env) cel.EnvOption {
	return cel.Variable(vd.Name, vd.Type.AsCelType(baseEnv))
}

type TypeDecl struct {
	TypeName    string      `yaml:"type_name"`
	Params      []*TypeDecl `yaml:"params"`
	IsTypeParam bool        `yaml:"is_type_param"`
}

func (td *TypeDecl) AsCelType(baseEnv *cel.Env) *cel.Type {
	switch td.TypeName {
	case "dyn":
		return cel.DynType
	case "map":
		if len(td.Params) == 2 {
			return cel.MapType(
				td.Params[0].AsCelType(baseEnv),
				td.Params[1].AsCelType(baseEnv))
		}
		panic(fmt.Sprintf("map type has unexpected type params: %v", td.Params))
	case "list":
		if len(td.Params) == 1 {
			return cel.ListType(td.Params[0].AsCelType(baseEnv))
		}
		panic(fmt.Sprintf("list type has unexpected params: %v", td.Params))
	default:
		if td.IsTypeParam {
			return cel.TypeParamType(td.TypeName)
		}
		if msgType, found := baseEnv.CELTypeProvider().FindStructType(td.TypeName); found {
			return msgType
		}
		t, found := baseEnv.CELTypeProvider().FindIdent(td.TypeName)
		if !found {
			panic(fmt.Sprintf("undefined type name: %v", td))
		}
		_, ok := t.(*cel.Type)
		if ok && len(td.Params) == 0 {
			return t.(*cel.Type)
		}
		params := make([]*cel.Type, len(td.Params))
		for i, p := range td.Params {
			params[i] = p.AsCelType(baseEnv)
		}
		return cel.OpaqueType(td.TypeName, params...)
	}
}

type FunctionDecl struct {
	Name      string          `yaml:"name"`
	Overloads []*OverloadDecl `yaml:"overloads"`
}

func (fd *FunctionDecl) AsEnvOption(baseEnv *cel.Env) cel.EnvOption {
	overloads := make([]cel.FunctionOpt, len(fd.Overloads))
	for i, o := range fd.Overloads {
		overloads[i] = o.AsFunctionOption(baseEnv)
	}
	return cel.Function(fd.Name, overloads...)
}

type OverloadDecl struct {
	OverloadID string      `yaml:"id"`
	Target     *TypeDecl   `yaml:"target"`
	Args       []*TypeDecl `yaml:"args"`
	Return     *TypeDecl   `yaml:"return"`
}

func (od *OverloadDecl) AsFunctionOption(baseEnv *cel.Env) cel.FunctionOpt {
	args := make([]*cel.Type, len(od.Args))
	for i, a := range od.Args {
		args[i] = a.AsCelType(baseEnv)
	}
	result := od.Return.AsCelType(baseEnv)
	if od.Target != nil {
		t := od.Target.AsCelType(baseEnv)
		args = append([]*cel.Type{t}, args...)
		return cel.MemberOverload(od.OverloadID, args, result)
	}
	return cel.Overload(od.OverloadID, args, result)
}

type extVersionFactory func(uint32) cel.EnvOption

var extFactories = map[string]extVersionFactory{
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
