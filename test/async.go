// Copyright 2020 Google LLC
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

package test

import (
	"context"
	"time"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
)

func FakeRPC(timeout time.Duration) functions.AsyncOp {
	return func(ctx context.Context, vars ref.Resolver, args []ref.Val) ref.Val {
		rpcCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		time.Sleep(20 * time.Millisecond)
		select {
		case <-rpcCtx.Done():
			return types.NewErr(rpcCtx.Err().Error())
		default:
			in := args[0].(types.String)
			return in.Add(types.String(" success!"))
		}
	}
}
