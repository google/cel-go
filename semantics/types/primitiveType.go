package types

import "fmt"

type PrimitiveType struct {
	kind primitiveKind
}

var _ Type = &PrimitiveType{}

var Int64 = newPrimitive(primitiveKindInt64)
var Uint64 = newPrimitive(primitiveKindUint64)
var String = newPrimitive(primitiveKindString)
var Bytes = newPrimitive(primitiveKindBytes)
var Bool = newPrimitive(primitiveKindBool)
var Double = newPrimitive(primitiveKindDouble)

type primitiveKind int

const (
	primitiveKindString primitiveKind = iota
	primitiveKindInt64
	primitiveKindUint64
	primitiveKindBytes
	primitiveKindDouble
	primitiveKindBool
)

func (p *PrimitiveType) Kind() TypeKind {
	return KindPrimitive
}

func (p *PrimitiveType) Equals(t Type) bool {
	other, ok := t.(*PrimitiveType)
	if !ok {
		return false
	}

	return p.kind == other.kind
}

func (p *PrimitiveType) String() string {
	switch p.kind {
	case primitiveKindString:
		return "string"
	case primitiveKindInt64:
		return "int"
	case primitiveKindUint64:
		return "uint"
	case primitiveKindBytes:
		return "bytes"
	case primitiveKindBool:
		return "bool"
	case primitiveKindDouble:
		return "double"
	default:
		panic(fmt.Sprintf("Unknown primitive kind: %v", p.kind))
	}
}

func newPrimitive(kind primitiveKind) *PrimitiveType {
	return &PrimitiveType{
		kind: kind,
	}
}
