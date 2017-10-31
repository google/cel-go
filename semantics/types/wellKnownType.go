package types

import "fmt"

type WellKnownType struct {
	kind wellKnownKind
}

var _ Type = &WellKnownType{}

var Duration = newWellKnown(wellKnownKindDuration)
var Timestamp = newWellKnown(wellKnownKindTimestamp)
var Any = newWellKnown(wellKnownKindAny)

type wellKnownKind int

const (
	wellKnownKindDuration wellKnownKind = iota
	wellKnownKindTimestamp
	wellKnownKindAny
)

func (w *WellKnownType) Kind() TypeKind {
	return KindWellKnown
}

func (w *WellKnownType) Equals(t Type) bool {
	other, ok := t.(*WellKnownType)
	if !ok {
		return false
	}

	return w.kind == other.kind
}

func (w *WellKnownType) String() string {
	switch w.kind {
	case wellKnownKindDuration:
		return "google.protobuf.Duration"
	case wellKnownKindTimestamp:
		return "google.protobuf.Timestamp"
	case wellKnownKindAny:
		return "any"
	default:
		panic(fmt.Sprintf("Unknown well-known kind: %v", w.kind))
	}
}

func newWellKnown(kind wellKnownKind) *WellKnownType {
	return &WellKnownType{
		kind: kind,
	}
}
