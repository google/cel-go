// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package compiler provides compatibility aliases forwarding to cel.dev/cel-go/tools/compiler.
package compiler

import (
	"cel.dev/cel-go/cel"
	canoncompiler "cel.dev/cel-go/tools/compiler"
)

// Types
type (
	FileFormat              = canoncompiler.FileFormat
	ExpressionType          = canoncompiler.ExpressionType
	PolicyMetadataEnvOption = canoncompiler.PolicyMetadataEnvOption
	Compiler                = canoncompiler.Compiler
	CustomMetadataCompiler  = canoncompiler.CustomMetadataCompiler
	InputExpression         = canoncompiler.InputExpression
	CompiledExpression      = canoncompiler.CompiledExpression
	FileExpression          = canoncompiler.FileExpression
	RawExpression           = canoncompiler.RawExpression
)

// Constants
const (
	Unspecified = canoncompiler.Unspecified
	BinaryProto = canoncompiler.BinaryProto
	TextProto   = canoncompiler.TextProto
	TextYAML    = canoncompiler.TextYAML
	CELString   = canoncompiler.CELString
	CELPolicy   = canoncompiler.CELPolicy

	ExpressionTypeUnspecified = canoncompiler.ExpressionTypeUnspecified
	CompiledExpressionFile    = canoncompiler.CompiledExpressionFile
	PolicyFile                = canoncompiler.PolicyFile
	ExpressionFile            = canoncompiler.ExpressionFile
	RawExpressionString       = canoncompiler.RawExpressionString
)

// Functions

func NewCompiler(opts ...any) (Compiler, error) {
	return canoncompiler.NewCompiler(opts...)
}

func InferFileFormat(path string) FileFormat {
	return canoncompiler.InferFileFormat(path)
}

func EnvironmentFile(path string) cel.EnvOption {
	return canoncompiler.EnvironmentFile(path)
}

func TypeDescriptorSetFile(path string) cel.EnvOption {
	return canoncompiler.TypeDescriptorSetFile(path)
}
