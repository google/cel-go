// Copyright 2022 Google LLC
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

package ext

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

var (
	nativeObjTraitMask = traits.FieldTesterType | traits.IndexerType
)

type nativeObj struct {
	val      any
	valType  *nativeType
	refValue reflect.Value
}

type nativeType struct {
	refType reflect.Type
}

func NativeObject(val any) ref.Val {
	if val == nil {
		return types.NullValue
	}
	refValue := reflect.ValueOf(val)
	return &nativeObj{
		val:      val,
		valType:  &nativeType{refType: refValue.Type()},
		refValue: refValue,
	}
}

func (o *nativeObj) ConvertToNative(typeDesc reflect.Type) (any, error) {
	if o.valType.refType == typeDesc {
		return o.val, nil
	}
	return nil, fmt.Errorf("type conversion error from '%T' to '%v'", o.val, typeDesc)
}

func (o *nativeObj) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	default:
		if o.Type().TypeName() == typeVal.TypeName() {
			return o
		}
	case types.TypeType:
		return o.valType
	}
	return types.NewErr("type conversion error from '%s' to '%s'", o.Type().TypeName(), typeVal)
}

func (o *nativeObj) Equal(other ref.Val) ref.Val {
	otherNtv, ok := other.(*nativeObj)
	return types.Bool(ok && reflect.DeepEqual(o.val, otherNtv.val))
}

func (o *nativeObj) IsZeroValue() bool {
	return o.refValue.IsZero()
}

func (o *nativeObj) Get(field ref.Val) ref.Val {
	fieldName, ok := field.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(field)
	}
	fieldVal := o.refValue.FieldByName(string(fieldName))
	if !fieldVal.IsValid() {
		return types.NewErr("no such field: %s", fieldName)
	}
	// Convert the field value to a ref.Val
	return nil
}

func (o *nativeObj) Type() ref.Type {
	return o.valType
}

func (o *nativeObj) Value() any {
	return o.val
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (t *nativeType) ConvertToNative(typeDesc reflect.Type) (any, error) {
	// TODO: replace the internal type representation with a proto-value.
	return nil, fmt.Errorf("type conversion not supported for 'type'")
}

// ConvertToType implements ref.Val.ConvertToType.
func (t *nativeType) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case types.TypeType:
		return types.TypeType
	case types.StringType:
		return types.String(t.TypeName())
	}
	return types.NewErr("type conversion error from '%s' to '%s'", types.TypeType, typeVal)
}

func (t *nativeType) Equal(other ref.Val) ref.Val {
	otherType, ok := other.(ref.Type)
	return types.Bool(ok && t.TypeName() == otherType.TypeName())
}

func (t *nativeType) HasTrait(trait int) bool {
	return nativeObjTraitMask&trait == trait
}

func (t *nativeType) String() string {
	return t.refType.Name()
}

func (t *nativeType) Type() ref.Type {
	return types.TypeType
}

func (t *nativeType) TypeName() string {
	return t.refType.Name()
}

func (t *nativeType) Value() any {
	return t.refType.Name()
}
