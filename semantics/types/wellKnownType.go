// Copyright 2018 Google LLC
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
