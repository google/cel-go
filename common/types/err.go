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

import (
	"fmt"
	refpb "github.com/google/cel-go/common/types/ref"
	"reflect"
)

// Err type which extends the built-in go error and implements refpb.Value.
type Err struct {
	error
}

var (
	// ErrType singleton.
	ErrType = NewTypeValue("error")
)

func NewErr(format string, args ...interface{}) *Err {
	return &Err{fmt.Errorf(format, args...)}
}

func (e *Err) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	return nil, e.error
}

func (e *Err) ConvertToType(typeVal refpb.Type) refpb.Value {
	// Errors are not convertible to other representations.
	return e
}

func (e *Err) Equal(other refpb.Value) refpb.Value {
	// An error cannot be equal to any other value, so it returns itself.
	return e
}

func (e *Err) String() string {
	return e.error.Error()
}

func (e *Err) Type() refpb.Type {
	return ErrType
}

func (e *Err) Value() interface{} {
	return e.error
}

// IsError returns whether the input element refpb.Type or refpb.Value is equal to
// the ErrType singleton.
func IsError(elem interface{}) bool {
	switch elem.(type) {
	case refpb.Type:
		return elem == ErrType
	case refpb.Value:
		return IsError(elem.(refpb.Value).Type())
	}
	return false
}
