package io

import (
	"fmt"
	"reflect"

	"celgo/ast"
	"celgo/common"
	api "celgo/proto"
)

func Serialize(e ast.Expression) *api.ParsedExpr {
	s := serializer{
		nextId:    1,
		locations: make(map[int64]common.Location),
	}

	expr := s.transform(e)
	sourceInfo := s.buildSourceInfo()
	return &api.ParsedExpr{
		Expr:       expr,
		SourceInfo: sourceInfo,
	}
}

type serializer struct {
	nextId    int64
	locations map[int64]common.Location
}

func (s *serializer) transform(e ast.Expression) *api.Expr {
	if e == nil {
		return nil
	}

	switch e.(type) {
	case *ast.IdentExpression:
		return s.transformIdent(e.(*ast.IdentExpression))
	case *ast.CallExpression:
		return s.transformCall(e.(*ast.CallExpression))
	case *ast.SelectExpression:
		return s.transformSelect(e.(*ast.SelectExpression))
	case *ast.CreateMessageExpression:
		return s.transformCreateMessage(e.(*ast.CreateMessageExpression))
	case *ast.CreateStructExpression:
		return s.transformCreateStruct(e.(*ast.CreateStructExpression))
	case *ast.CreateListExpression:
		return s.transformCreateList(e.(*ast.CreateListExpression))
	case *ast.ComprehensionExpression:
		return s.transformComprehension(e.(*ast.ComprehensionExpression))
	case *ast.Int64Constant:
		return s.transformInt64Constant(e.(*ast.Int64Constant))
	case *ast.Uint64Constant:
		return s.transformUint64Constant(e.(*ast.Uint64Constant))
	case *ast.DoubleConstant:
		return s.transformDoubleConstant(e.(*ast.DoubleConstant))
	case *ast.StringConstant:
		return s.transformStringConstant(e.(*ast.StringConstant))
	case *ast.BytesConstant:
		return s.transformBytesConstant(e.(*ast.BytesConstant))
	case *ast.BoolConstant:
		return s.transformBoolConstant(e.(*ast.BoolConstant))
	case *ast.NullConstant:
		return s.transformNullConstant(e.(*ast.NullConstant))
	}

	panic(fmt.Sprintf("Unknown expression type during serialization: %v", reflect.TypeOf(e)))
	return nil
}

func (s *serializer) transformIdent(e *ast.IdentExpression) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_IdentExpr{
			IdentExpr: &api.Expr_Ident{
				Name: e.Name,
			},
		},
	}
}

func (s *serializer) transformCall(e *ast.CallExpression) *api.Expr {
	id := s.track(e)

	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_CallExpr{
			CallExpr: &api.Expr_Call{
				Target:   s.transform(e.Target),
				Function: e.Function,
				Args:     s.transformSlice(e.Args),
			},
		},
	}
}

func (s *serializer) transformSelect(e *ast.SelectExpression) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_SelectExpr{
			SelectExpr: &api.Expr_Select{
				Operand:  s.transform(e.Target),
				Field:    e.Field,
				TestOnly: e.TestOnly,
			},
		},
	}
}

func (s *serializer) transformCreateMessage(e *ast.CreateMessageExpression) *api.Expr {
	entries := make([]*api.Expr_CreateStruct_Entry, len(e.Fields))

	id := s.track(e)

	for i, f := range e.Fields {
		e := &api.Expr_CreateStruct_Entry{
			Id: s.track(f),
			KeyKind: &api.Expr_CreateStruct_Entry_FieldKey{
				FieldKey: f.Name,
			},
			Value: s.transform(f.Initializer),
		}

		entries[i] = e
	}

	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_StructExpr{
			StructExpr: &api.Expr_CreateStruct{
				MessageName: e.MessageName,
				Entries:     entries,
			},
		},
	}
}

func (s *serializer) transformCreateStruct(e *ast.CreateStructExpression) *api.Expr {
	entries := make([]*api.Expr_CreateStruct_Entry, len(e.Entries))

	id := s.track(e)

	for i, f := range e.Entries {
		e := &api.Expr_CreateStruct_Entry{
			Id: s.track(f),
			KeyKind: &api.Expr_CreateStruct_Entry_MapKey{
				MapKey: s.transform(f.Key),
			},
			Value: s.transform(f.Value),
		}

		entries[i] = e
	}

	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_StructExpr{
			StructExpr: &api.Expr_CreateStruct{
				Entries: entries,
			},
		},
	}
}

func (s *serializer) transformCreateList(e *ast.CreateListExpression) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_ListExpr{
			ListExpr: &api.Expr_CreateList{
				Elements: s.transformSlice(e.Entries),
			},
		},
	}
}

func (s *serializer) transformComprehension(e *ast.ComprehensionExpression) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_ComprehensionExpr{
			ComprehensionExpr: &api.Expr_Comprehension{
				AccuVar:       e.Accumulator,
				IterVar:       e.Variable,
				IterRange:     s.transform(e.Target),
				AccuInit:      s.transform(e.Init),
				LoopCondition: s.transform(e.LoopCondition),
				LoopStep:      s.transform(e.LoopStep),
				Result:        s.transform(e.Result),
			},
		},
	}
}

func (s *serializer) transformInt64Constant(e *ast.Int64Constant) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_ConstExpr{
			ConstExpr: &api.Constant{
				ConstantKind: &api.Constant_Int64Value{
					Int64Value: e.Value,
				},
			},
		},
	}
}

func (s *serializer) transformUint64Constant(e *ast.Uint64Constant) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_ConstExpr{
			ConstExpr: &api.Constant{
				ConstantKind: &api.Constant_Uint64Value{
					Uint64Value: e.Value,
				},
			},
		},
	}
}

func (s *serializer) transformStringConstant(e *ast.StringConstant) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_ConstExpr{
			ConstExpr: &api.Constant{
				ConstantKind: &api.Constant_StringValue{
					StringValue: e.Value,
				},
			},
		},
	}
}

func (s *serializer) transformDoubleConstant(e *ast.DoubleConstant) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_ConstExpr{
			ConstExpr: &api.Constant{
				ConstantKind: &api.Constant_DoubleValue{
					DoubleValue: e.Value,
				},
			},
		},
	}
}

func (s *serializer) transformBytesConstant(e *ast.BytesConstant) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_ConstExpr{
			ConstExpr: &api.Constant{
				ConstantKind: &api.Constant_BytesValue{
					BytesValue: e.Value,
				},
			},
		},
	}
}

func (s *serializer) transformBoolConstant(e *ast.BoolConstant) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_ConstExpr{
			ConstExpr: &api.Constant{
				ConstantKind: &api.Constant_BoolValue{
					BoolValue: e.Value,
				},
			},
		},
	}
}

func (s *serializer) transformNullConstant(e *ast.NullConstant) *api.Expr {
	id := s.track(e)
	return &api.Expr{
		Id: id,
		ExprKind: &api.Expr_ConstExpr{
			ConstExpr: &api.Constant{
				ConstantKind: &api.Constant_NullValue{},
			},
		},
	}
}

func (s *serializer) transformSlice(exprs []ast.Expression) []*api.Expr {
	if exprs == nil {
		return []*api.Expr{}
	}

	result := make([]*api.Expr, len(exprs))
	for i, e := range exprs {
		result[i] = s.transform(e)
	}

	return result
}

func (s *serializer) track(e ast.Expression) int64 {
	id := s.nextId
	s.nextId++
	s.locations[id] = e.Location()
	return id
}

func (s *serializer) buildSourceInfo() *api.SourceInfo {
	lines := make(map[int64]int32)
	columns := make(map[int64]int32)
	var locationName string

	for id, loc := range s.locations {
		lines[id] = int32(loc.Line())
		columns[id] = int32(loc.Column())
		locationName = loc.Source().Name()
	}

	return &api.SourceInfo{
		Location: locationName,
		Columns:  columns,
		Lines:    lines,
	}
}
