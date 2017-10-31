package ast

import "celgo/common"

type BaseExpression struct {
	id       int64
	location common.Location
}

func (e *BaseExpression) Location() common.Location {
	return e.location
}

func (e *BaseExpression) Id() int64 {
	return e.id
}
