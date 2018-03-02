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

var Null Type = &nullType{}

type nullType struct {
}

var _ Type = &nullType{}

func (n *nullType) Kind() TypeKind {
	return KindNull
}

func (n *nullType) Equals(t Type) bool {
	_, ok := t.(*nullType)
	return ok
}

func (n *nullType) String() string {
	return "null"
}
