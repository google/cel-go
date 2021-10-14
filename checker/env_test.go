// Copyright 2018 Google LLC
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

package checker

import (
	"strings"
	"testing"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/containers"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestOverlappingIdentifier(t *testing.T) {
	env := NewStandardEnv(containers.DefaultContainer, newTestRegistry(t))
	err := env.Add(
		decls.NewVar("int", decls.NewTypeType(nil)))
	if err == nil {
		t.Error("Got nil, wanted error")
	} else if !strings.Contains(err.Error(), "overlapping identifier") {
		t.Errorf("Got %v, wanted overlapping identifier error", err)
	}
}

func TestOverlappingMacro(t *testing.T) {
	env := NewStandardEnv(containers.DefaultContainer, newTestRegistry(t))
	err := env.Add(decls.NewFunction("has",
		decls.NewOverload("has", []*exprpb.Type{decls.String}, decls.Bool)))
	if err == nil {
		t.Error("Got nil, wanted error")
	} else if !strings.Contains(err.Error(), "overlapping macro") {
		t.Errorf("Got %v, wanted overlapping macro error", err)
	}
}

func TestOverlappingOverload(t *testing.T) {
	env := NewStandardEnv(containers.DefaultContainer, newTestRegistry(t))
	paramA := decls.NewTypeParamType("A")
	typeParamAList := []string{"A"}
	err := env.Add(decls.NewFunction(overloads.TypeConvertDyn,
		decls.NewParameterizedOverload(overloads.ToDyn,
			[]*exprpb.Type{paramA}, decls.Dyn,
			typeParamAList)))
	if err == nil {
		t.Error("Got nil, wanted error")
	} else if !strings.Contains(err.Error(), "overlapping overload") {
		t.Errorf("Got %v, wanted overlapping overload error", err)
	}
}

func TestSanitizedOverload(t *testing.T) {
	env := NewStandardEnv(containers.DefaultContainer, newTestRegistry(t))
	err := env.Add(decls.NewFunction(operators.Add,
		decls.NewOverload("timestamp_add_int",
			[]*exprpb.Type{decls.NewObjectType("google.protobuf.Timestamp"), decls.Int},
			decls.NewObjectType("google.protobuf.Timestamp"))))
	if err != nil {
		t.Errorf("env.Add('timestamp_add_int') failed: %v", err)
	}
}

func TestSanitizedInstanceOverload(t *testing.T) {
	env := NewStandardEnv(containers.DefaultContainer, newTestRegistry(t))
	err := env.Add(decls.NewFunction("orDefault",
		decls.NewInstanceOverload("int_ordefault_int",
			[]*exprpb.Type{
				decls.NewObjectType("google.protobuf.Int32Value"),
				decls.NewObjectType("google.protobuf.Int32Value")},
			decls.Int)))
	if err != nil {
		t.Errorf("env.Add('int_ordefault_int') failed: %v", err)
	}
}

func newTestRegistry(t *testing.T) ref.TypeRegistry {
	t.Helper()
	reg, err := types.NewRegistry()
	if err != nil {
		t.Fatalf("types.NewRegistry() failed: %v", err)
	}
	return reg
}
