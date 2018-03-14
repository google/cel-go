package interpreter

import (
	"github.com/google/cel-go/common"
)

type Metadata interface {
	CharacterOffset(exprId int64) (int32, bool)
	Location(exprId int64) (common.Location, bool)
	Expressions() []int64
}
