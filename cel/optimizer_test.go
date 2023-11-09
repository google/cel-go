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

package cel_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"

	proto3pb "github.com/google/cel-go/test/proto3pb"
)

func TestStaticOptimizerUpdateExpr(t *testing.T) {
	expr := `has(a.b)`
	inlined := `[x, y].filter(i, i.size() > 0)[0].z`

	opts := []cel.EnvOption{
		cel.Types(&proto3pb.TestAllTypes{}),
		cel.OptionalTypes(),
		cel.EnableMacroCallTracking(),
		cel.Variable("a", cel.MapType(cel.StringType, cel.StringType)),
		cel.Variable("x", cel.MapType(cel.StringType, cel.StringType)),
		cel.Variable("y", cel.MapType(cel.StringType, cel.StringType)),
	}
	e, err := cel.NewEnv(opts...)
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}
	exprAST, iss := e.Compile(expr)
	if iss.Err() != nil {
		t.Fatalf("Compile() failed: %v", iss.Err())
	}

	inlinedAST, iss := e.Compile(inlined)
	if iss.Err() != nil {
		t.Fatalf("Compile() failed: %v", iss.Err())
	}
	opt := cel.NewStaticOptimizer(&testOptimizer{t: t, inlineExpr: inlinedAST.NativeRep()})
	optAST, iss := opt.Optimize(e, exprAST)
	if iss.Err() != nil {
		t.Fatalf("Optimize() generated an invalid AST: %v", iss.Err())
	}
	optString, err := cel.AstToString(optAST)
	if err != nil {
		t.Fatalf("cel.AstToString() failed: %v", err)
	}
	expected := `has([x, y].filter(i, i.size() > 0)[0].z)`
	if expected != optString {
		t.Errorf("inlined got %q, wanted %q", optString, expected)
	}
}

type testOptimizer struct {
	t          *testing.T
	inlineExpr *ast.AST
}

func (opt *testOptimizer) Optimize(ctx *cel.OptimizerContext, a *ast.AST) *ast.AST {
	opt.t.Helper()
	copy, info := ctx.CopyAST(opt.inlineExpr)
	infoMacroKeys := getMacroKeys(info.MacroCalls())
	for id, call := range info.MacroCalls() {
		a.SourceInfo().SetMacroCall(id, call)
	}
	origID := a.Expr().ID()
	exprID := origID + 100
	presenceTest, hasMacro := ctx.NewHasMacro(exprID, copy)
	macroKeys := getMacroKeys(a.SourceInfo().MacroCalls())
	if len(macroKeys) != 2 {
		opt.t.Errorf("Got %v macro calls, wanted 2", macroKeys)
	}
	ctx.UpdateExpr(a.Expr(), presenceTest)
	macroKeys = getMacroKeys(a.SourceInfo().MacroCalls())
	if _, found := a.SourceInfo().GetMacroCall(origID); found {
		opt.t.Errorf("Got %v macro calls, wanted 1", macroKeys)
	}

	a.SourceInfo().SetMacroCall(exprID, hasMacro)
	macroKeys = getMacroKeys(a.SourceInfo().MacroCalls())
	if !reflect.DeepEqual(macroKeys, append(infoMacroKeys, int(exprID))) {
		opt.t.Errorf("Got %v macro calls, wanted 2", macroKeys)
	}
	return a
}

func getMacroKeys(macroCalls map[int64]ast.Expr) []int {
	keys := []int{}
	for k := range macroCalls {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return keys
}
