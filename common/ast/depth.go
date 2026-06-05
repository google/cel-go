// Copyright 2026 Google LLC
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

// ExceedsMaxDepth reports whether the given expression nests deeper than
// maxDepth. The traversal itself is bounded: it never recurses past
// maxDepth+1 levels, so it is safe to call on adversarially deep inputs that
// would otherwise blow the Go stack during later checking or planning.
//
// A non-positive maxDepth disables the check and always returns false.
func ExceedsMaxDepth(e Expr, maxDepth int) bool {
	if maxDepth <= 0 {
		return false
	}
	return exceedsMaxDepth(e, 0, maxDepth)
}

func exceedsMaxDepth(e Expr, depth, maxDepth int) bool {
	if e == nil {
		return false
	}
	if depth > maxDepth {
		return true
	}
	switch e.Kind() {
	case CallKind:
		c := e.AsCall()
		if c.IsMemberFunction() {
			if exceedsMaxDepth(c.Target(), depth+1, maxDepth) {
				return true
			}
		}
		for _, arg := range c.Args() {
			if exceedsMaxDepth(arg, depth+1, maxDepth) {
				return true
			}
		}
	case ComprehensionKind:
		c := e.AsComprehension()
		if exceedsMaxDepth(c.IterRange(), depth+1, maxDepth) {
			return true
		}
		if exceedsMaxDepth(c.AccuInit(), depth+1, maxDepth) {
			return true
		}
		if exceedsMaxDepth(c.LoopCondition(), depth+1, maxDepth) {
			return true
		}
		if exceedsMaxDepth(c.LoopStep(), depth+1, maxDepth) {
			return true
		}
		if exceedsMaxDepth(c.Result(), depth+1, maxDepth) {
			return true
		}
	case ListKind:
		for _, elem := range e.AsList().Elements() {
			if exceedsMaxDepth(elem, depth+1, maxDepth) {
				return true
			}
		}
	case MapKind:
		for _, entry := range e.AsMap().Entries() {
			me := entry.AsMapEntry()
			if exceedsMaxDepth(me.Key(), depth+1, maxDepth) {
				return true
			}
			if exceedsMaxDepth(me.Value(), depth+1, maxDepth) {
				return true
			}
		}
	case SelectKind:
		if exceedsMaxDepth(e.AsSelect().Operand(), depth+1, maxDepth) {
			return true
		}
	case StructKind:
		for _, f := range e.AsStruct().Fields() {
			if exceedsMaxDepth(f.AsStructField().Value(), depth+1, maxDepth) {
				return true
			}
		}
	}
	return false
}
