package types

type Type interface {
	Kind() TypeKind
	Equals(other Type) bool
	String() string
}

type TypeKind int

const (
	KindPrimitive TypeKind = iota
	KindError
	KindDynamic
	KindMessage
	KindList
	KindMap
	KindFunction
	KindTypeParameter
	KindWrapper
	KindWellKnown
	KindNull
	KindType
)
