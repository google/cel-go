// Copyright 2025 Google LLC
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

// Package celtest provides functions for testing CEL policies and expressions.
package celtest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/types"
)

// reportCoverage reports the coverage information for the provided programs.
//   - For the node coverage, the coverage is reported as a percentage of the number of nodes which
//     were evaluated during the test execution and hence are present in the program Coverage report.
//   - For the branch coverage, every node which has a boolean return type is considered as a branch.
//     The number of branches which were evaluated during the test execution and hence are present in
//     the program Coverage report are reported as the branch coverage percentage.
func reportCoverage(t *testing.T, programs []Program) {
	t.Helper()
	for _, p := range programs {
		exprString, err := cel.AstToString(p.Ast)
		if err != nil {
			t.Logf("Error converting AST to string for a program: %v", err)
			continue
		}
		rootNavigableExpr := ast.NavigateAST(p.Ast.NativeRep())
		// Initialize coverage metrics
		cr := &coverageReport{
			nodes:                  0,
			coveredNodes:           0,
			branches:               0,
			coveredBooleanOutcomes: 0,
			unencounteredNodes:     []string{},
			unencounteredBranches:  []string{},
		}
		traverseAndCalculateCoverage(t, rootNavigableExpr, p, true, "", cr)
		printCoverageReport(t, exprString, cr)
	}
}

type coverageReport struct {
	nodes                  int64
	coveredNodes           int64
	branches               int64
	coveredBooleanOutcomes int64
	unencounteredNodes     []string
	unencounteredBranches  []string
}

func traverseAndCalculateCoverage(t *testing.T, expr ast.NavigableExpr, p Program, logUnencountered bool,
	preceedingTabs string, cr *coverageReport) {
	t.Helper()
	if expr == nil || len(p.CoverageStats) == 0 {
		return
	}
	nodeID := expr.ID()
	cr.nodes++
	interestingBoolNode := expr.Type() == types.BoolType && expr.AsLiteral() == nil && expr.AsCall().FunctionName() != "cel.@block"
	// Check for Node Coverage
	if _, isCovered := p.CoverageStats[nodeID]; isCovered {
		cr.coveredNodes++
	} else if logUnencountered {
		if interestingBoolNode {
			exprText, _ := cel.ExprToString(expr, p.Ast.NativeRep().SourceInfo())
			cr.unencounteredNodes = append(cr.unencounteredNodes,
				fmt.Sprintf("\nExpression ID %d ('%s')", nodeID, exprText))
		}
		logUnencountered = false
	}
	// Check for Branch Coverage if the node is a boolean type
	if interestingBoolNode {
		cr.branches += 2
		exprText, _ := cel.ExprToString(expr, p.Ast.NativeRep().SourceInfo())
		if info, found := p.CoverageStats[nodeID]; !found {
			if logUnencountered {
				cr.unencounteredBranches = append(cr.unencounteredBranches,
					"\n"+preceedingTabs+fmt.Sprintf("Expression ID %d ('%s'): Never evaluated (neither true nor false)", nodeID, exprText))
				preceedingTabs = preceedingTabs + "\t\t"
			}
		} else {
			if _, ok := info[types.True]; ok {
				cr.coveredBooleanOutcomes++
			} else if logUnencountered {
				cr.unencounteredBranches = append(cr.unencounteredBranches,
					"\n"+preceedingTabs+fmt.Sprintf("Expression ID %d ('%s'): Never evaluated to 'true'", nodeID, exprText))
				preceedingTabs = preceedingTabs + "\t\t"

			}
			if _, ok := info[types.False]; ok {
				cr.coveredBooleanOutcomes++
			} else if logUnencountered {
				cr.unencounteredBranches = append(cr.unencounteredBranches,
					"\n"+preceedingTabs+fmt.Sprintf("Expression ID %d ('%s'): Never evaluated to 'false'", nodeID, exprText))
				preceedingTabs = preceedingTabs + "\t\t"
			}
		}
	}
	for _, child := range expr.Children() {
		traverseAndCalculateCoverage(t, child.(ast.NavigableExpr), p, logUnencountered, preceedingTabs, cr)
	}
}

func printCoverageReport(t *testing.T, exprString string, cr *coverageReport) {
	t.Helper()
	t.Logf("--- Start Coverage Report ---\nExpression: %s", exprString)
	if cr.nodes == 0 {
		t.Logf("No coverage stats found")
		return
	}
	// Log Node Coverage results
	nodeCoverage := float64(cr.coveredNodes) / float64(cr.nodes) * 100.0
	t.Logf("AST Node Coverage: %.2f%% (%d out of %d nodes covered)", nodeCoverage, cr.coveredNodes, cr.nodes)
	if len(cr.unencounteredNodes) > 0 {
		t.Logf("Interesting Unencountered Nodes:\n%s", strings.Join(cr.unencounteredNodes, "\n"))
	}
	// Log Branch Coverage results
	branchCoverage := 0.0
	if cr.branches > 0 {
		branchCoverage = float64(cr.coveredBooleanOutcomes) / float64(cr.branches) * 100.0
	}
	t.Logf("AST Branch Coverage: %.2f%% (%d out of %d branch outcomes covered)", branchCoverage,
		cr.coveredBooleanOutcomes, cr.branches)
	if len(cr.unencounteredBranches) > 0 {
		t.Logf("Interesting Unencountered Branch Paths:\n%s", strings.Join(cr.unencounteredBranches, "\n"))
	}
	t.Logf("--- End Coverage Report ---\n")
}
