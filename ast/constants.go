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

package ast

import "github.com/google/cel-go/common"

type Constant interface {
}

type Int64Constant struct {
	BaseExpression
	Constant

	Value int64
}

func (e *Int64Constant) String() string {
	return ToDebugString(e)
}

func (e *Int64Constant) writeDebugString(w *debugWriter) {
	w.appendFormat("%d", e.Value)
	w.adorn(e)
}

type Uint64Constant struct {
	BaseExpression
	Constant

	Value uint64
}

func (e *Uint64Constant) String() string {
	return ToDebugString(e)
}

func (e *Uint64Constant) writeDebugString(w *debugWriter) {
	w.appendFormat("%du", e.Value)
	w.adorn(e)
}

type DoubleConstant struct {
	BaseExpression
	Constant

	Value float64
}

func (e *DoubleConstant) String() string {
	return ToDebugString(e)
}

func (e *DoubleConstant) writeDebugString(w *debugWriter) {
	w.appendFormat("%v", e.Value)
	w.adorn(e)
}

type StringConstant struct {
	BaseExpression
	Constant

	Value string
}

func (e *StringConstant) String() string {
	return ToDebugString(e)
}

func (e *StringConstant) writeDebugString(w *debugWriter) {
	// TODO: escape
	w.append(`"`)
	w.append(e.Value)
	w.append(`"`)
	w.adorn(e)
}

type BytesConstant struct {
	BaseExpression
	Constant

	Value []byte
}

func (e *BytesConstant) String() string {
	return ToDebugString(e)
}

func (e *BytesConstant) writeDebugString(w *debugWriter) {
	w.append(`b"`)
	w.append(string(e.Value))
	w.append(`"`)
	w.adorn(e)
}

type BoolConstant struct {
	BaseExpression
	Constant

	Value bool
}

func (e *BoolConstant) String() string {
	return ToDebugString(e)
}

func (e *BoolConstant) writeDebugString(w *debugWriter) {
	s := "false"
	if e.Value {
		s = "true"
	}
	w.append(s)
	w.adorn(e)
}

type NullConstant struct {
	BaseExpression
	Constant

	Value *NullConstant
}

func (e *NullConstant) String() string {
	return ToDebugString(e)
}

func (e *NullConstant) writeDebugString(w *debugWriter) {
	w.append("null")
	w.adorn(e)
}

func NewInt64Constant(id int64, l common.Location, value int64) *Int64Constant {
	return &Int64Constant{
		BaseExpression: BaseExpression{id: id, location: l},
		Value:          value,
	}
}

func NewUint64Constant(id int64, l common.Location, value uint64) *Uint64Constant {
	return &Uint64Constant{
		BaseExpression: BaseExpression{id: id, location: l},
		Value:          value,
	}
}

func NewDoubleConstant(id int64, l common.Location, value float64) *DoubleConstant {
	return &DoubleConstant{
		BaseExpression: BaseExpression{id: id, location: l},
		Value:          value,
	}
}

func NewStringConstant(id int64, l common.Location, value string) *StringConstant {
	return &StringConstant{
		BaseExpression: BaseExpression{id: id, location: l},
		Value:          value,
	}
}

func NewBytesConstant(id int64, l common.Location, value []byte) *BytesConstant {
	return &BytesConstant{
		BaseExpression: BaseExpression{id: id, location: l},
		Value:          value,
	}
}

func NewBoolConstant(id int64, l common.Location, value bool) *BoolConstant {
	return &BoolConstant{
		BaseExpression: BaseExpression{id: id, location: l},
		Value:          value,
	}
}

func NewNullConstant(id int64, l common.Location) *NullConstant {
	r := &NullConstant{
		BaseExpression: BaseExpression{id: id, location: l},
	}
	r.Value = r
	return r
}
