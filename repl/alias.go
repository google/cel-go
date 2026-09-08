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

// Package repl provides compatibility aliases forwarding to cel.dev/cel-go/repl.
package repl

import (
	"cel.dev/cel-go/common/env"
	"cel.dev/cel-go/common/types"
	canonrepl "cel.dev/cel-go/repl"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Types
type (
	Optioner          = canonrepl.Optioner
	EvaluationContext = canonrepl.EvaluationContext
	Evaluator         = canonrepl.Evaluator
	Cmder             = canonrepl.Cmder
)

// Functions

func NewEvaluator() (*Evaluator, error) {
	return canonrepl.NewEvaluator()
}

func Parse(line string) (Cmder, error) {
	return canonrepl.Parse(line)
}

func UnparseExprType(t *exprpb.Type) string {
	return canonrepl.UnparseExprType(t)
}

func UnparseType(t *types.Type) string {
	return canonrepl.UnparseType(t)
}

func ParseType(t string) (*env.TypeDesc, error) {
	return canonrepl.ParseType(t)
}
