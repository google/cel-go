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

package interpreter

import (
	"testing"

	"github.com/google/cel-go/common/types"
)

func TestGetterSetter(t *testing.T) {
	evalState := newEvalState(2)
	if val, found := evalState.Value(1); !found || val != nil {
		t.Error("Unexpected value found", val)
	}
	evalState.values[1] = types.String("hello")
	if greeting, found := evalState.Value(1); !found || greeting != types.String("hello") {
		t.Error("Unexpected value found", greeting)
	}
}
