// Copyright 2019 Google LLC
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

package interpreter

import (
	"errors"
	"fmt"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// Resolver provides methods for finding values by name and resolving qualified attributes from
// them.
type Resolver interface {
	// ResolveName returns a value for the given input string if it exists.
	ResolveName(string) (interface{}, bool)

	// ResolveQualifiers returns the value at the qualified field path or error if the one or more
	// of the qualifier paths could not be found.
	ResolveQualifiers(Activation, interface{}, []Qualifier) (interface{}, error)
}

// NewResolver returns a default Resolver which is cabable of resolving types by simple names, and
// can resolve qualifiers on CEL values using the supported qualifier types: bool, int, string,
// and uint.
func NewResolver(a ref.TypeAdapter, p ref.TypeProvider) Resolver {
	return &resolver{adapter: a, provider: p}
}

type resolver struct {
	adapter  ref.TypeAdapter
	provider ref.TypeProvider
}

// ResolveName returns a type identifier if one exists.
func (res *resolver) ResolveName(name string) (interface{}, bool) {
	return res.provider.FindIdent(name)
}

// ResolveQualifiers resolves static and dynamic qualifiers on the input object.
//
// Resolution of qualifiers on Go simple and aggregate types does not require marshalling of
// intermediate results to CEL ref.Val instances; however, proto message types and Go structs
// will be marshalled to CEL ref.Val's which can result in slower resolution time. Custom
// Resolvers may be used to improve performance.
func (res *resolver) ResolveQualifiers(vars Activation,
	obj interface{},
	quals []Qualifier) (interface{}, error) {
	var s string
	var i int64
	var u uint64
	var b bool
	var cVal ref.Val
	for idx := 0; idx < len(quals); idx++ {
		isMap := false
		isKey := false
		isIndex := false
		// First, determine the value and type of the qualifier. Most qualifiers have both a
		// simple Go value and a CEL value in order to reduce allocations required for conversion
		// between the two.
		switch qual := quals[idx].(type) {
		case *fieldQualifier:
			// A field qualifier is capable of obtaining the value of the field from the object
			// where it is defined.
			if rv, ok := obj.(ref.Val); ok {
				obj = rv.Value()
			}
			v, err := qual.FieldType.GetFrom(obj)
			if err != nil {
				return nil, err
			}
			obj = v
			continue
		case *stringQualifier:
			s = qual.Value
			cVal = qual.CelValue
			goto QualString
		case *intQualifier:
			i = qual.Value
			cVal = qual.CelValue
			goto QualInt
		case *uintQualifier:
			u = qual.Value
			cVal = qual.CelValue
			goto QualUint
		case *boolQualifier:
			b = qual.Value
			cVal = qual.CelValue
			goto QualBool
		case Attribute:
			v, err := qual.Resolve(vars, res)
			if err != nil {
				return nil, err
			}
			switch q := v.(type) {
			case types.Unknown:
				return q, nil
			case types.String:
				s = string(q)
				cVal = q
				goto QualString
			case types.Int:
				i = int64(q)
				cVal = q
				goto QualInt
			case types.Uint:
				u = uint64(q)
				cVal = q
				goto QualUint
			case types.Bool:
				b = q == types.True
				cVal = q
				goto QualBool
			case string:
				s = q
				cVal = types.String(s)
				goto QualString
			case int:
				i = int64(q)
				cVal = types.Int(i)
				goto QualInt
			case int32:
				i = int64(q)
				cVal = types.Int(i)
				goto QualInt
			case int64:
				i = q
				cVal = types.Int(i)
				goto QualInt
			case uint:
				u = uint64(q)
				cVal = types.Uint(u)
				goto QualUint
			case uint32:
				u = uint64(q)
				cVal = types.Uint(u)
				goto QualUint
			case uint64:
				u = q
				cVal = types.Uint(u)
				goto QualUint
			case bool:
				b = q
				cVal = types.Bool(b)
				goto QualBool
			default:
				return nil, fmt.Errorf("unsupported attribute type: %T", q)
			}
		}

	QualString:
		switch o := obj.(type) {
		case map[string]interface{}:
			isMap = true
			obj, isKey = o[s]
		case map[string]string:
			isMap = true
			obj, isKey = o[s]
		case map[string]int:
			isMap = true
			obj, isKey = o[s]
		case map[string]int32:
			isMap = true
			obj, isKey = o[s]
		case map[string]int64:
			isMap = true
			obj, isKey = o[s]
		case map[string]uint:
			isMap = true
			obj, isKey = o[s]
		case map[string]uint32:
			isMap = true
			obj, isKey = o[s]
		case map[string]uint64:
			isMap = true
			obj, isKey = o[s]
		case map[string]float32:
			isMap = true
			obj, isKey = o[s]
		case map[string]float64:
			isMap = true
			obj, isKey = o[s]
		case map[string]bool:
			isMap = true
			obj, isKey = o[s]
		case types.Unknown:
			return o, nil
		default:
			elem, err := res.refResolve(cVal, obj)
			if err != nil {
				return nil, err
			}
			if types.IsUnknown(elem) {
				return res.fmtUnknown(elem, quals[idx]), nil
			}
			obj = elem
			isMap = true
			isKey = true
		}
		if isMap && !isKey {
			return nil, fmt.Errorf("no such key: %v", s)
		}
		continue

	QualInt:
		switch o := obj.(type) {
		case map[int]interface{}:
			isMap = true
			obj, isKey = o[int(i)]
		case map[int32]interface{}:
			isMap = true
			obj, isKey = o[int32(i)]
		case map[int64]interface{}:
			isMap = true
			obj, isKey = o[i]
		case []interface{}:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []string:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []int:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []int32:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []int64:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []uint:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []uint32:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []uint64:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []float32:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []float64:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case []bool:
			isIndex = i >= 0 && i < int64(len(o))
			if isIndex {
				obj = o[i]
			}
		case types.Unknown:
			return o, nil
		default:
			elem, err := res.refResolve(cVal, obj)
			if err != nil {
				return nil, err
			}
			if types.IsUnknown(elem) {
				return res.fmtUnknown(elem, quals[idx]), nil
			}
			isMap = true
			isKey = true
			obj = elem
		}
		if isMap && !isKey {
			return nil, fmt.Errorf("no such key: %v", i)
		}
		if !isMap && !isIndex {
			return nil, fmt.Errorf("index out of bounds: %v", i)
		}
		continue

	QualUint:
		switch o := obj.(type) {
		case map[uint]interface{}:
			isMap = true
			obj, isKey = o[uint(u)]
		case map[uint32]interface{}:
			isMap = true
			obj, isKey = o[uint32(u)]
		case map[uint64]interface{}:
			isMap = true
			obj, isKey = o[u]
		case types.Unknown:
			return o, nil
		default:
			elem, err := res.refResolve(cVal, obj)
			if err != nil {
				return nil, err
			}
			if types.IsUnknown(elem) {
				return res.fmtUnknown(elem, quals[idx]), nil
			}
			isMap = true
			isKey = true
			obj = elem
		}
		if isMap && !isKey {
			return nil, fmt.Errorf("no such key: %v", i)
		}
		if !isMap && !isIndex {
			return nil, fmt.Errorf("index out of bounds: %v", i)
		}
		continue

	QualBool:
		switch o := obj.(type) {
		case map[bool]interface{}:
			obj, isKey = o[b]
		case map[bool]string:
			obj, isKey = o[b]
		case map[bool]int:
			obj, isKey = o[b]
		case map[bool]int32:
			obj, isKey = o[b]
		case map[bool]int64:
			obj, isKey = o[b]
		case map[bool]uint:
			obj, isKey = o[b]
		case map[bool]uint32:
			obj, isKey = o[b]
		case map[bool]uint64:
			obj, isKey = o[b]
		case map[bool]float32:
			obj, isKey = o[b]
		case map[bool]float64:
			obj, isKey = o[b]
		case types.Unknown:
			return o, nil
		default:
			elem, err := res.refResolve(cVal, obj)
			if err != nil {
				return nil, err
			}
			if types.IsUnknown(elem) {
				return res.fmtUnknown(elem, quals[idx]), nil
			}
			isMap = true
			isKey = true
			obj = elem
		}
		if !isKey {
			return nil, fmt.Errorf("no such key: %v", b)
		}
		continue
	}
	return obj, nil
}

func (res *resolver) fmtUnknown(elem ref.Val, qual Qualifier) ref.Val {
	unk := elem.(types.Unknown)
	if len(unk) == 0 {
		return types.Unknown{qual.ID()}
	}
	return unk
}

func (res *resolver) refResolve(idx ref.Val, obj interface{}) (ref.Val, error) {
	celVal := res.adapter.NativeToValue(obj)
	mapper, isMapper := celVal.(traits.Mapper)
	if isMapper {
		elem, found := mapper.Find(idx)
		if !found {
			return nil, fmt.Errorf("no such key: %v", idx)
		}
		if types.IsError(elem) {
			return nil, elem.Value().(error)
		}
		return elem, nil
	}
	indexer, isIndexer := celVal.(traits.Indexer)
	if isIndexer {
		elem := indexer.Get(idx)
		if types.IsError(elem) {
			return nil, elem.Value().(error)
		}
		return elem, nil
	}
	return nil, errors.New("no such overload")
}
