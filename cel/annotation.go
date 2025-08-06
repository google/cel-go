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

package cel

import (
	"fmt"

	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/decls"
)

// EnableAnnotations registers a function used internally to annotate expressions.
//
// Annotations are represented as a list associated with an expression.
// This design enables annotations from optimized subtrees to be aggregated and preserved when constructing new nodes during optimization passes,
// ensuring that annotations from different parts of the AST are retained without name conflicts.
func EnableAnnotations() EnvOption {
	t := types.NewTypeParamType("T")
	// cel.@annotation(T, map(string, dyn)) -> T
	return Function("cel.@annotation",
		decls.Overload("cel_annotation",
			[]*types.Type{t, types.NewListType(types.NewMapType(types.StringType, types.DynType))},
			t,
		),
	)
}


// Annotation represents the structure of an annotation, used to associate metadata with a AST node.
type Annotation struct {
	Name   string
	Value  any
}

// NewAnnotation creates a new annotation with the given name and value.
//
// Note: Evaluable CEL expressions as annotations are not supported yet.
func NewAnnotation(name string, value any) *Annotation {
	return &Annotation{
		Name:   name,
		Value:  value,
	}
}

// AnnotationFactory is an interface for generating annotations for a given expression.
type AnnotationFactory interface {
	// GenerateAnnotation generates an annotation for the given expression.
	//
	// The function receives the current expression and the original AST.
	// It returns an Annotation instance if successful, or an error if generation fails.
	GenerateAnnotation(expr ast.Expr, a *ast.AST) (*Annotation, error)
}

// AnnotationOptimizer is an optimizer that applies the AST with annotations.
type AnnotationOptimizer struct {
	env       *Env
	factories []AnnotationFactory
}

// AnnotationOptimizerOption is used to configure the AnnotationOptimizer.
type AnnotationOptimizerOption func(*AnnotationOptimizer) (*AnnotationOptimizer, error)

// AnnotationFactories sets the annotation factories for the AnnotationOptimizer.
func AnnotationFactories(factories []AnnotationFactory) AnnotationOptimizerOption {
	return func(opt *AnnotationOptimizer) (*AnnotationOptimizer, error) {
		opt.factories = append(opt.factories, factories...)
		return opt, nil
	}
}

// AnnotationEnv sets the environment for the AnnotationOptimizer.
func AnnotationEnv(env *Env) AnnotationOptimizerOption {
	return func(opt *AnnotationOptimizer) (*AnnotationOptimizer, error) {
		opt.env = env
		return opt, nil
	}
}

// NewAnnotationOptimizer returns an AnnotationOptimizer instance with the given factories applied.
func NewAnnotationOptimizer(opts ...AnnotationOptimizerOption) (*AnnotationOptimizer, error) {
	annotator := &AnnotationOptimizer{}
	var err error
	for _, opt := range opts {
		annotator, err = opt(annotator)
		if err != nil {
			return nil, err
		}
	}
	return annotator, nil
}

// Optimize applies the annotations recursively to the input AST.
func (opt *AnnotationOptimizer) Optimize(ctx *OptimizerContext, a *ast.AST) *ast.AST {
	node := a.Expr()
	newRoot, err := opt.annotateExpr(ctx, node, a)
	if err != nil {
		ctx.ReportErrorAtID(node.ID(), "error annotating expression: %v", err)
		return a
	}
	return ctx.NewAST(newRoot)
}

// generateAnnotations generates annotations for the given expression using the configured factories.
func (opt *AnnotationOptimizer) generateAnnotations(ctx *OptimizerContext, expr ast.Expr, originalAST *ast.AST) []*Annotation {
	var annotations []*Annotation
	for _, factory := range opt.factories {
		annotation, err := factory.GenerateAnnotation(expr, originalAST)
		if err != nil {
			ctx.ReportErrorAtID(expr.ID(), "error generating annotation: %v", err)
			// Ignore the error and continue with the next annotation factory.
			continue
		}
		if annotation != nil {
			annotations = append(annotations, annotation)
		}
	}
	return annotations
}

func (opt *AnnotationOptimizer) annotateExpr(ctx *OptimizerContext, currentExpr ast.Expr, originalAST *ast.AST) (ast.Expr, error) {
	// Convert to a NavigableExpr to get access to parent and children.
	navigableExpr := ast.NavigateExpr(originalAST, currentExpr)

	// Recursively annotate children first.
	var annotatedChildren []ast.Expr
	for _, child := range navigableExpr.Children() {
		annotatedChild, err := opt.annotateExpr(ctx, child, originalAST)
		if err != nil {
			return nil, err
		}
		annotatedChildren = append(annotatedChildren, annotatedChild)
	}

	// Reconstruct the current expression with the annotated children.
	reconstructedExpr, err := reconstructExpr(currentExpr, ctx, annotatedChildren)
	if err != nil {
		return nil, err
	}

	// Generate annotations for the reconstructed expression.
	annotations := opt.generateAnnotations(ctx, reconstructedExpr, originalAST)

	// If no annotations are generated, return the reconstructed expression.
	if len(annotations) == 0 {
		return reconstructedExpr, nil
	}

	var annotationExprs []ast.Expr
	for _, ann := range annotations {
		annotationExprs = append(annotationExprs, createAnnotationExpr(ctx, ann))
	}
	// Create the annotations node.
	annotationsNode := createAnnotationNode(ctx, annotationExprs)

	// Create the final annotated node.
	finalAnnotatedNode := ctx.NewCall("cel.@annotation", reconstructedExpr, annotationsNode)
	return finalAnnotatedNode, nil
}

// reconstructExpr reconstructs the given expression with the given annotated children.
//
// The function receives the current expression and a list of annotated children.
// It then reconstructs the expression by using the annotated children at the appropriate
// positions.
func reconstructExpr(currentExpr ast.Expr, ctx *OptimizerContext, annotatedChildren []ast.Expr) (ast.Expr, error) {
	var reconstructedExpr ast.Expr

	switch currentExpr.Kind() {

	case ast.CallKind:
		originalCall := currentExpr.AsCall()
		if originalCall.IsMemberFunction() {
			reconstructedExpr = ctx.NewMemberCall(originalCall.FunctionName(), annotatedChildren[0], annotatedChildren[1:]...)
		} else {
			reconstructedExpr = ctx.NewCall(originalCall.FunctionName(), annotatedChildren...)
		}

	case ast.ListKind:
		originalList := currentExpr.AsList()
		reconstructedExpr = ctx.NewList(annotatedChildren, originalList.OptionalIndices())

	case ast.MapKind:
		originalMap := currentExpr.AsMap()
		var newMapEntries []ast.EntryExpr
		for i := 0; i < len(annotatedChildren); i += 2 {
			originalEntry := originalMap.Entries()[i/2].AsMapEntry()
			newMapEntries = append(newMapEntries, ctx.NewMapEntry(annotatedChildren[i], annotatedChildren[i+1], originalEntry.IsOptional()))
		}
		reconstructedExpr = ctx.NewMap(newMapEntries)

	case ast.StructKind:
		originalStruct := currentExpr.AsStruct()
		var newStructFields []ast.EntryExpr
		for i, originalField := range originalStruct.Fields() {
			originalStructField := originalField.AsStructField()
			newStructFields = append(newStructFields, ctx.NewStructField(originalStructField.Name(), annotatedChildren[i], originalStructField.IsOptional()))
		}
		reconstructedExpr = ctx.NewStruct(originalStruct.TypeName(), newStructFields)

	case ast.SelectKind:
		originalSelect := currentExpr.AsSelect()
		if originalSelect.IsTestOnly() {
			reconstructedExpr = ctx.NewPresenceTest(annotatedChildren[0], originalSelect.FieldName())
		} else {
			reconstructedExpr = ctx.NewSelect(annotatedChildren[0], originalSelect.FieldName())
		}

	case ast.ComprehensionKind:
		originalComprehension := currentExpr.AsComprehension()
		if originalComprehension.HasIterVar2() {
			reconstructedExpr = ctx.NewComprehensionTwoVar(annotatedChildren[0], originalComprehension.IterVar(), originalComprehension.IterVar2(), originalComprehension.AccuVar(), annotatedChildren[1], annotatedChildren[2], annotatedChildren[3], annotatedChildren[4])
		} else {
			reconstructedExpr = ctx.NewComprehension(annotatedChildren[0], originalComprehension.IterVar(), originalComprehension.AccuVar(), annotatedChildren[1], annotatedChildren[2], annotatedChildren[3], annotatedChildren[4])
		}

	case ast.LiteralKind:
		reconstructedExpr = ctx.NewLiteral(currentExpr.AsLiteral())

	case ast.IdentKind:
		reconstructedExpr = ctx.NewIdent(currentExpr.AsIdent())

	case ast.UnspecifiedExprKind:
		reconstructedExpr = ctx.NewUnspecifiedExpr()

	default:
		return nil, fmt.Errorf("unsupported expression kind for annotation: %v", currentExpr.Kind())
	}
	return reconstructedExpr, nil
}

// createAnnotationNode creates a list-typed AST node that holds annotations for a given expression.
// This list node is designed to store a list of annotations directly associated with the current AST node.
func createAnnotationNode(ctx *OptimizerContext, annotations []ast.Expr) ast.Expr {
	return ctx.NewList(annotations, nil)
}

// createAnnotationExpr creates an AST node for a single Annotation struct (a map literal).
func createAnnotationExpr(ctx *OptimizerContext, ann *Annotation) ast.Expr {
	return ctx.NewMap(
		[]ast.EntryExpr{
			ctx.NewMapEntry(ctx.NewLiteral(types.String("name")), ctx.NewLiteral(types.String(ann.Name)), false),
			ctx.NewMapEntry(ctx.NewLiteral(types.String("value")), ctx.NewLiteral(types.DefaultTypeAdapter.NativeToValue(ann.Value)), false),
		})
}
