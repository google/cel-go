// Copyright 2021 Google LLC
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

package parser

import "fmt"

type options struct {
	maxRecursionDepth  int
	errorRecoveryLimit int
	macros             map[string]Macro
}

// Option configures the behavior of the parser.
type Option func(*options) error

// MaxRecursionDepth is an option which limits the maximum depth the parser will attempt to
// parse the expression before giving up.
func MaxRecursionDepth(maxRecursionDepth int) Option {
	return func(opts *options) error {
		if maxRecursionDepth < -1 {
			return fmt.Errorf("max recursion depth must be greater than or equal to -1: %d", maxRecursionDepth)
		}
		opts.maxRecursionDepth = maxRecursionDepth
		return nil
	}
}

// ErrorRecoveryLimit is an option which limits the number of attempts the parser will perform
// to recover from an error.
func ErrorRecoveryLimit(errorRecoveryLimit int) Option {
	return func(opts *options) error {
		if errorRecoveryLimit < -1 {
			return fmt.Errorf("error recovery limit must be greater than or equal to -1: %d", errorRecoveryLimit)
		}
		opts.errorRecoveryLimit = errorRecoveryLimit
		return nil
	}
}

// Macros adds the given macros to the parser.
func Macros(macros ...Macro) Option {
	return func(opts *options) error {
		for _, m := range macros {
			if m != nil {
				if opts.macros == nil {
					opts.macros = make(map[string]Macro)
				}
				opts.macros[m.MacroKey()] = m
			}
		}
		return nil
	}
}
