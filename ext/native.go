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
	"strings"
	"time"

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
			t, err := newNativeType(rt)
			if err != nil {
				return nil, err
			}
			nativeTypes[t.TypeName()] = t
		case reflect.Value:
			t, err := newNativeType(rt.Type())
			if err != nil {
				return nil, err
			}
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
		return decls.NewTypeType(decls.NewObjectType(typeName)), true
	}
	return tp.baseProvider.FindType(typeName)
}

func (tp *nativeTypeProvider) FindFieldType(typeName, fieldName string) (*ref.FieldType, bool) {
	t, found := tp.nativeTypes[typeName]
	if !found {
		return tp.baseProvider.FindFieldType(typeName, fieldName)
	}
	refField, isDefined := t.HasField(fieldName)
	if !found || !isDefined {
		return nil, false
	}
	exprType, ok := convertToExprType(refField.Type)
	if !ok {
		return nil, false
	}
	return &ref.FieldType{
		Type: exprType,
		IsSet: func(obj any) bool {
			refVal := reflect.Indirect(reflect.ValueOf(obj))
			refField := refVal.FieldByName(fieldName)
			return !refField.IsZero()
		},
		GetFrom: func(obj any) (any, error) {
			refVal := reflect.Indirect(reflect.ValueOf(obj))
			refField := refVal.FieldByName(fieldName)
			return getFieldValue(tp, refField), nil
		},
	}, true
}

func (tp *nativeTypeProvider) NewValue(typeName string, fields map[string]ref.Val) ref.Val {
	t, found := tp.nativeTypes[typeName]
	if !found {
		return tp.baseProvider.NewValue(typeName, fields)
	}
	refPtr := reflect.New(t.refType)
	refVal := refPtr.Elem()
	for fieldName, val := range fields {
		refFieldDef, isDefined := t.HasField(fieldName)
		if !isDefined {
			return types.NewErr("no such field: %s", fieldName)
		}
		fieldVal, err := val.ConvertToNative(refFieldDef.Type)
		if err != nil {
			return types.NewErr(err.Error())
		}
		refField := refVal.FieldByIndex(refFieldDef.Index)
		refFieldVal := reflect.ValueOf(fieldVal)
		refField.Set(refFieldVal)
	}
	return tp.NativeToValue(refPtr.Interface())
}

func (tp *nativeTypeProvider) NativeToValue(val any) ref.Val {
	if val == nil {
		return types.NullValue
	}
	rawVal := reflect.ValueOf(val)
	refVal := rawVal
	if refVal.Kind() == reflect.Ptr {
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
		switch val := val.(type) {
		case ref.Val:
			return val
		case proto.Message, *pb.Map, protoreflect.List, protoreflect.Message, protoreflect.Value,
			time.Time:
			return tp.baseAdapter.NativeToValue(val)
		default:
			return newNativeObject(tp, val, rawVal)
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
		if refType == durationType {
			return decls.Duration, true
		}
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
		if refType == timestampType {
			return decls.Timestamp, true
		}
		return decls.NewObjectType(
			fmt.Sprintf("%s.%s", simplePkgAlias(refType.PkgPath()), refType.Name()),
		), true
	case reflect.Pointer:
		if refType.Implements(pbMsgInterfaceType) {
			pbMsg := reflect.New(refType.Elem()).Interface().(protoreflect.ProtoMessage)
			return decls.NewObjectType(string(pbMsg.ProtoReflect().Descriptor().FullName())), true
		}
		return convertToExprType(refType.Elem())
	}
	return nil, false
}

func newNativeObject(adapter ref.TypeAdapter, val any, refValue reflect.Value) ref.Val {
	valType, err := newNativeType(refValue.Type())
	if err != nil {
		return types.NewErr(err.Error())
	}
	return &nativeObj{
		TypeAdapter: adapter,
		val:         val,
		valType:     valType,
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
	if o.refValue.Type() == typeDesc {
		return o.val, nil
	}
	if o.refValue.Kind() == reflect.Pointer && o.refValue.Type().Elem() == typeDesc {
		return o.refValue.Elem().Interface(), nil
	}
	if typeDesc.Kind() == reflect.Pointer && o.refValue.Type() == typeDesc.Elem() {
		ptr := reflect.New(typeDesc.Elem())
		ptr.Elem().Set(o.refValue)
		return ptr.Interface(), nil
	}
	return nil, fmt.Errorf("type conversion error from '%v' to '%v'", o.Type(), typeDesc)
}

func (o *nativeObj) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case types.TypeType:
		return o.valType
	default:
		if typeVal.TypeName() == o.valType.typeName {
			return o
		}
	}
	return types.NewErr("type conversion error from '%s' to '%s'", o.Type(), typeVal)
}

func (o *nativeObj) Equal(other ref.Val) ref.Val {
	otherNtv, ok := other.(*nativeObj)
	if !ok {
		return types.False
	}
	val := o.val
	otherVal := otherNtv.val
	refVal := o.refValue
	otherRefVal := otherNtv.refValue
	if refVal.Kind() != otherRefVal.Kind() {
		if refVal.Kind() == reflect.Pointer {
			val = refVal.Elem().Interface()
		} else if otherRefVal.Kind() == reflect.Pointer {
			otherVal = otherRefVal.Elem().Interface()
		}
	}
	return types.Bool(reflect.DeepEqual(val, otherVal))
}

func (o *nativeObj) IsZeroValue() bool {
	fmt.Printf("[nativeObj] is non-zero: %v\n", o.refValue)
	return o.refValue.IsZero()
}

// IsSet tests whether a field which is defined is set to a non-default value.
func (o *nativeObj) IsSet(field ref.Val) ref.Val {
	refField, refErr := o.getReflectedField(field)
	if refErr != nil {
		return refErr
	}
	return types.Bool(!refField.IsZero())
}

func (o *nativeObj) Get(field ref.Val) ref.Val {
	refField, refErr := o.getReflectedField(field)
	if refErr != nil {
		return refErr
	}
	return adaptFieldValue(o, refField)
}

func (o *nativeObj) getReflectedField(field ref.Val) (reflect.Value, ref.Val) {
	fieldName, ok := field.(types.String)
	if !ok {
		return reflect.Value{}, types.MaybeNoSuchOverloadErr(field)
	}
	fieldNameStr := string(fieldName)
	refField, isDefined := o.valType.HasField(fieldNameStr)
	if !isDefined {
		return reflect.Value{}, types.NewErr("no such field: %s", fieldName)
	}
	refVal := reflect.Indirect(o.refValue)
	return refVal.FieldByIndex(refField.Index), nil
}

func (o *nativeObj) Type() ref.Type {
	return o.valType
}

func (o *nativeObj) Value() any {
	return o.val
}

func newNativeType(rawType reflect.Type) (*nativeType, error) {
	refType := rawType
	if refType.Kind() == reflect.Pointer {
		refType = refType.Elem()
	}
	if !isValidObjectType(refType) {
		return nil, fmt.Errorf("unsupported reflect.Type %v, must be reflect.Struct", rawType)
	}
	return &nativeType{
		typeName: fmt.Sprintf("%s.%s", simplePkgAlias(refType.PkgPath()), refType.Name()),
		refType:  refType,
	}, nil
}

type nativeType struct {
	typeName string
	refType  reflect.Type
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (t *nativeType) ConvertToNative(typeDesc reflect.Type) (any, error) {
	return nil, fmt.Errorf("type conversion error for type to '%v'", typeDesc)
}

// ConvertToType implements ref.Val.ConvertToType.
func (t *nativeType) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case types.TypeType:
		return types.TypeType
	}
	return types.NewErr("type conversion error from '%s' to '%s'", types.TypeType, typeVal)
}

func (t *nativeType) Equal(other ref.Val) ref.Val {
	otherType, ok := other.(ref.Type)
	return types.Bool(ok && t.TypeName() == otherType.TypeName())
}

func (t *nativeType) HasField(fieldName string) (reflect.StructField, bool) {
	f, found := t.refType.FieldByName(fieldName)
	if !found || !f.IsExported() || !isSupportedType(f.Type) {
		return reflect.StructField{}, false
	}
	return f, true
}

func (t *nativeType) HasTrait(trait int) bool {
	return nativeObjTraitMask&trait == trait
}

func (t *nativeType) String() string {
	return t.typeName
}

func (t *nativeType) Type() ref.Type {
	return types.TypeType
}

func (t *nativeType) TypeName() string {
	return t.typeName
}

func (t *nativeType) Value() any {
	return t.typeName
}

func adaptFieldValue(adapter ref.TypeAdapter, refField reflect.Value) ref.Val {
	return adapter.NativeToValue(getFieldValue(adapter, refField))
}

func getFieldValue(adapter ref.TypeAdapter, refField reflect.Value) any {
	if refField.IsZero() {
		switch refField.Kind() {
		case reflect.Array, reflect.Slice:
			return types.NewDynamicList(adapter, []ref.Val{})
		case reflect.Map:
			return types.NewDynamicMap(adapter, map[ref.Val]ref.Val{})
		case reflect.Struct:
			if refField.Type() == timestampType {
				return types.Timestamp{Time: time.Unix(0, 0)}
			}
			return reflect.New(refField.Type()).Elem().Interface()
		case reflect.Pointer:
			return reflect.New(refField.Type().Elem()).Interface()
		}
	}
	return refField.Interface()
}

func simplePkgAlias(pkgPath string) string {
	paths := strings.Split(pkgPath, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[len(paths)-1]
}

func isValidObjectType(refType reflect.Type) bool {
	return refType.Kind() == reflect.Struct
}

func isSupportedType(refType reflect.Type) bool {
	switch refType.Kind() {
	case reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.Func, reflect.UnsafePointer, reflect.Uintptr:
		return false
	case reflect.Array, reflect.Slice:
		return isSupportedType(refType.Elem())
	case reflect.Map:
		return isSupportedType(refType.Key()) && isSupportedType(refType.Elem())
	}
	return true
}

var (
	pbMsgInterfaceType = reflect.TypeOf((*protoreflect.ProtoMessage)(nil)).Elem()
	timestampType      = reflect.TypeOf(time.Now())
	durationType       = reflect.TypeOf(time.Nanosecond)
)
