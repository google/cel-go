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

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/pb"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var (
	nativeObjTraitMask = traits.FieldTesterType | traits.IndexerType
)

func NativeTypes(refTypes ...any) cel.EnvOption {
	return func(env *cel.Env) (*cel.Env, error) {
		tp, err := newNativeTypeProvider(env.TypeAdapter(), env.TypeProvider(), refTypes...)
		if err != nil {
			return nil, err
		}
		env, err = cel.CustomTypeAdapter(tp)(env)
		if err != nil {
			return nil, err
		}
		return cel.CustomTypeProvider(tp)(env)
	}
}

func newNativeTypeProvider(adapter ref.TypeAdapter, provider ref.TypeProvider, refTypes ...any) (*nativeTypeProvider, error) {
	nativeTypes := make(map[string]*nativeType, len(refTypes))
	for _, refType := range refTypes {
		switch rt := refType.(type) {
		case reflect.Type:
			if rt.Kind() != reflect.Struct {
				return nil, fmt.Errorf("unsupported reflect.Type %v, must be reflect.Struct", rt)
			}
			t := &nativeType{refType: rt}
			nativeTypes[t.TypeName()] = t
		case reflect.Value:
			rt = reflect.Indirect(rt)
			if rt.Kind() != reflect.Struct {
				return nil, fmt.Errorf("unsupported reflect.Type %v, must be reflect.Struct", rt)
			}
			t := &nativeType{refType: rt.Type()}
			nativeTypes[t.TypeName()] = t
		default:
			return nil, fmt.Errorf("unsupported native type: %v (%T) must be reflect.Type or reflect.Value", rt, rt)
		}
	}
	return &nativeTypeProvider{
		nativeTypes:  nativeTypes,
		baseAdapter:  adapter,
		baseProvider: provider,
	}, nil
}

type nativeTypeProvider struct {
	nativeTypes  map[string]*nativeType
	baseAdapter  ref.TypeAdapter
	baseProvider ref.TypeProvider
}

func (tp *nativeTypeProvider) EnumValue(enumName string) ref.Val {
	return tp.baseProvider.EnumValue(enumName)
}

func (tp *nativeTypeProvider) FindIdent(typeName string) (ref.Val, bool) {
	if t, found := tp.nativeTypes[typeName]; found {
		return t, true
	}
	return tp.baseProvider.FindIdent(typeName)
}

func (tp *nativeTypeProvider) FindType(typeName string) (*exprpb.Type, bool) {
	if _, found := tp.nativeTypes[typeName]; found {
		return decls.NewObjectType(typeName), true
	}
	return tp.baseProvider.FindType(typeName)
}

func (tp *nativeTypeProvider) FindFieldType(typeName, fieldName string) (*ref.FieldType, bool) {
	t, found := tp.nativeTypes[typeName]
	if !found {
		return tp.baseProvider.FindFieldType(typeName, fieldName)
	}
	refField, found := t.refType.FieldByName(fieldName)
	if !found {
		return nil, false
	}
	exprType, ok := convertToExprType(refField.Type)
	if !ok {
		return nil, false
	}
	return &ref.FieldType{
		Type: exprType,
		IsSet: func(obj any) bool {
			refVal := reflect.ValueOf(obj)
			refField := refVal.FieldByName(fieldName)
			return !refField.IsZero()
		},
		GetFrom: func(obj any) (any, error) {
			refVal := reflect.ValueOf(obj)
			refField := refVal.FieldByName(fieldName)
			return refField.Interface(), nil
		},
	}, true
}

func (tp *nativeTypeProvider) NewValue(typeName string, fields map[string]ref.Val) ref.Val {
	t, found := tp.nativeTypes[typeName]
	if !found {
		return tp.baseProvider.NewValue(typeName, fields)
	}
	refVal := reflect.New(t.refType)
	for fieldName, val := range fields {
		_, found := t.refType.FieldByName(fieldName)
		if !found {
			types.NewErr("no such field: %s", fieldName)
		}
		refField := refVal.FieldByName(fieldName)
		if !refField.CanSet() {
			return types.NewErr("cannot set field: %s on %s", fieldName, t.refType.Name())
		}
		fieldVal, err := val.ConvertToNative(refField.Type())
		if err != nil {
			return types.NewErr(err.Error())
		}
		refFieldVal := reflect.ValueOf(fieldVal)
		refField.Set(refFieldVal)
	}
	return tp.NativeToValue(refVal.Interface())
}

func (tp *nativeTypeProvider) NativeToValue(val any) ref.Val {
	refVal := reflect.ValueOf(val)
	if refVal.Kind() == reflect.Ptr {
		if refVal.IsNil() {
			return types.UnsupportedRefValConversionErr(val)
		}
		refVal = reflect.Indirect(refVal)
	}
	// This isn't quite right if you're also supporting proto,
	// but maybe an acceptable limitation.
	switch refVal.Kind() {
	case reflect.Array, reflect.Slice:
		switch val.(type) {
		case []byte:
			return tp.baseAdapter.NativeToValue(val)
		default:
			return types.NewDynamicList(tp, val)
		}
	case reflect.Map:
		return types.NewDynamicMap(tp, val)
	case reflect.Struct:
		switch val.(type) {
		case proto.Message, *pb.Map, protoreflect.List, protoreflect.Message, protoreflect.Value:
			return tp.baseAdapter.NativeToValue(val)
		default:
			return newNativeObject(tp, val, refVal)
		}
	default:
		return tp.baseAdapter.NativeToValue(val)
	}
}

func convertToExprType(refType reflect.Type) (*exprpb.Type, bool) {
	switch refType.Kind() {
	case reflect.Bool:
		return decls.Bool, true
	case reflect.Float32, reflect.Float64:
		return decls.Double, true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return decls.Int, true
	case reflect.String:
		return decls.String, true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return decls.Uint, true
	case reflect.Array, reflect.Slice:
		refElem := refType.Elem()
		if refElem == reflect.TypeOf(byte(0)) {
			return decls.Bytes, true
		}
		elemType, ok := convertToExprType(refElem)
		if !ok {
			return nil, false
		}
		return decls.NewListType(elemType), true
	case reflect.Map:
		keyType, ok := convertToExprType(refType.Key())
		if !ok {
			return nil, false
		}
		// Ensure the key type is a int, bool, uint, string
		elemType, ok := convertToExprType(refType.Elem())
		if !ok {
			return nil, false
		}
		return decls.NewMapType(keyType, elemType), true
	case reflect.Struct:
		return decls.NewObjectType(fmt.Sprintf("%s.%s", refType.PkgPath(), refType.Name())), true
	case reflect.Pointer:
		return convertToExprType(refType.Elem())
	}
	return nil, false
}

func newNativeObject(adapter ref.TypeAdapter, val any, refValue reflect.Value) ref.Val {
	if val == nil {
		return types.NullValue
	}
	return &nativeObj{
		TypeAdapter: adapter,
		val:         val,
		valType:     &nativeType{refType: refValue.Type()},
		refValue:    refValue,
	}
}

type nativeObj struct {
	ref.TypeAdapter
	val      any
	valType  *nativeType
	refValue reflect.Value
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
	f, isDefined := o.valType.refType.FieldByName(string(fieldName))
	if !isDefined {
		return types.NewErr("no such field: %s", fieldName)
	}
	fv := o.refValue.FieldByIndex(f.Index)
	return o.NativeToValue(fv.Interface())
}

func (o *nativeObj) Type() ref.Type {
	return o.valType
}

func (o *nativeObj) Value() any {
	return o.val
}

type nativeType struct {
	refType reflect.Type
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (t *nativeType) ConvertToNative(typeDesc reflect.Type) (any, error) {
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
	return t.TypeName()
}

func (t *nativeType) Type() ref.Type {
	return types.TypeType
}

func (t *nativeType) TypeName() string {
	return fmt.Sprintf("%s.%s", t.refType.PkgPath(), t.refType.Name())
}

func (t *nativeType) Value() any {
	return t.TypeName()
}
