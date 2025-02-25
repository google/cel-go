package types

import (
	"strings"

	"github.com/google/cel-go/common/types/ref"
)

type formattable interface {
	format(*strings.Builder)
}

// Format formats the value as a string. The result is only intended for human consumption and ignores errors.
// Do not depend on the output being stable. It may change at any time.
func Format(val ref.Val) string {
	var sb strings.Builder
	formatTo(&sb, val)
	return sb.String()
}

func formatTo(sb *strings.Builder, val ref.Val) {
	if fmtable, ok := val.(formattable); ok {
		fmtable.format(sb)
		return
	}
}
