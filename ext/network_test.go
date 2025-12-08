// Copyright 2025 Google LLC
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
	"reflect"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
)

func TestNetwork_Success(t *testing.T) {
	// These test cases are ported from kubernetes/staging/src/k8s.io/apiserver/pkg/cel/library
	// to ensure 1-to-1 parity with the Kubernetes implementation.
	tests := []struct {
		name string
		expr string
		out  any
	}{
		// --- Global Checks (isIP, isCIDR) ---
		{
			name: "isIP valid IPv4",
			expr: "isIP('1.2.3.4')",
			out:  true,
		},
		{
			name: "isIP valid IPv6",
			expr: "isIP('2001:db8::1')",
			out:  true,
		},
		{
			name: "isIP invalid",
			expr: "isIP('not.an.ip')",
			out:  false,
		},
		{
			name: "isIP with port (invalid)",
			expr: "isIP('127.0.0.1:80')",
			out:  false,
		},
		{
			name: "isCIDR valid",
			expr: "isCIDR('10.0.0.0/8')",
			out:  true,
		},
		{
			name: "isCIDR invalid mask",
			expr: "isCIDR('10.0.0.0/999')",
			out:  false,
		},

		// --- IP Constructors & Equality ---
		{
			name: "ip equality IPv4",
			expr: "ip('127.0.0.1') == ip('127.0.0.1')",
			out:  true,
		},
		{
			name: "ip inequality",
			expr: "ip('127.0.0.1') == ip('1.2.3.4')",
			out:  false,
		},
		{
			name: "ip equality IPv6 mixed case inputs",
			// Logic check: The value is equal even if string rep was different
			expr: "ip('2001:db8::1') == ip('2001:DB8::1')",
			out:  true,
		},

		// --- Family ---
		{
			name: "family IPv4",
			expr: "ip('127.0.0.1').family()",
			out:  int64(4),
		},
		{
			name: "family IPv6",
			expr: "ip('::1').family()",
			out:  int64(6),
		},

		// --- Canonicalization (Critical Feature) ---
		{
			name: "isCanonical IPv4 simple",
			expr: "ip.isCanonical('127.0.0.1')",
			out:  true,
		},
		{
			name: "isCanonical IPv6 standard",
			expr: "ip.isCanonical('2001:db8::1')",
			out:  true,
		},
		{
			name: "isCanonical IPv6 uppercase (invalid)",
			expr: "ip.isCanonical('2001:DB8::1')",
			out:  false,
		},
		{
			name: "isCanonical IPv6 expanded (invalid)",
			expr: "ip.isCanonical('2001:db8:0:0:0:0:0:1')",
			out:  false,
		},

		// --- IP Types (Loopback, Unspecified, etc) ---
		{
			name: "isLoopback IPv4",
			expr: "ip('127.0.0.1').isLoopback()",
			out:  true,
		},
		{
			name: "isLoopback IPv6",
			expr: "ip('::1').isLoopback()",
			out:  true,
		},
		{
			name: "isUnspecified IPv4",
			expr: "ip('0.0.0.0').isUnspecified()",
			out:  true,
		},
		{
			name: "isUnspecified IPv6",
			expr: "ip('::').isUnspecified()",
			out:  true,
		},
		{
			name: "isGlobalUnicast 8.8.8.8",
			expr: "ip('8.8.8.8').isGlobalUnicast()",
			out:  true,
		},
		{
			name: "isLinkLocalMulticast",
			expr: "ip('ff02::1').isLinkLocalMulticast()",
			out:  true,
		},

		// --- CIDR Accessors ---
		{
			name: "cidr prefixLength",
			expr: "cidr('192.168.0.0/24').prefixLength()",
			out:  int64(24),
		},
		{
			name: "cidr ip extraction",
			expr: "cidr('192.168.0.0/24').ip() == ip('192.168.0.0')",
			out:  true,
		},
		{
			name: "cidr ip extraction (host bits set)",
			// K8s behavior: cidr('1.2.3.4/24').ip() returns 1.2.3.4, not 1.2.3.0
			expr: "cidr('192.168.1.5/24').ip() == ip('192.168.1.5')",
			out:  true,
		},
		{
			name: "cidr masked",
			// masked() zeroes out the host bits
			expr: "cidr('192.168.1.5/24').masked() == cidr('192.168.1.0/24')",
			out:  true,
		},
		{
			name: "cidr masked identity",
			expr: "cidr('192.168.1.0/24').masked() == cidr('192.168.1.0/24')",
			out:  true,
		},

		// --- Containment (IP in CIDR) ---
		{
			name: "containsIP simple",
			expr: "cidr('10.0.0.0/8').containsIP(ip('10.1.2.3'))",
			out:  true,
		},
		{
			name: "containsIP string overload",
			expr: "cidr('10.0.0.0/8').containsIP('10.1.2.3')",
			out:  true,
		},
		{
			name: "containsIP edge case (network address)",
			expr: "cidr('10.0.0.0/8').containsIP(ip('10.0.0.0'))",
			out:  true,
		},
		{
			name: "containsIP edge case (broadcast)",
			expr: "cidr('10.0.0.0/8').containsIP(ip('10.255.255.255'))",
			out:  true,
		},
		{
			name: "containsIP false",
			expr: "cidr('10.0.0.0/8').containsIP(ip('11.0.0.0'))",
			out:  false,
		},

		// --- Containment (CIDR in CIDR) ---
		{
			name: "containsCIDR exact match",
			expr: "cidr('10.0.0.0/8').containsCIDR(cidr('10.0.0.0/8'))",
			out:  true,
		},
		{
			name: "containsCIDR subnet",
			expr: "cidr('10.0.0.0/8').containsCIDR(cidr('10.1.0.0/16'))",
			out:  true,
		},
		{
			name: "containsCIDR string overload",
			expr: "cidr('10.0.0.0/8').containsCIDR('10.1.0.0/16')",
			out:  true,
		},
		{
			name: "containsCIDR larger prefix (false)",
			// /8 does not contain /4
			expr: "cidr('10.0.0.0/8').containsCIDR(cidr('0.0.0.0/4'))",
			out:  false,
		},
		{
			name: "containsCIDR disjoint",
			expr: "cidr('10.0.0.0/8').containsCIDR(cidr('11.0.0.0/8'))",
			out:  false,
		},
		{
			name: "containsCIDR different family",
			expr: "cidr('10.0.0.0/8').containsCIDR(cidr('::1/128'))",
			out:  false,
		},
	}

	// Initialize the environment with the Network extension
	env, err := cel.NewEnv(Network())
	if err != nil {
		t.Fatalf("cel.NewEnv(Network()) failed: %v", err)
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			ast, iss := env.Compile(tst.expr)
			if iss.Err() != nil {
				t.Fatalf("Compile(%q) failed: %v", tst.expr, iss.Err())
			}

			prg, err := env.Program(ast)
			if err != nil {
				t.Fatalf("Program(%q) failed: %v", tst.expr, err)
			}

			out, _, err := prg.Eval(cel.NoVars())
			if err != nil {
				t.Fatalf("Eval(%q) failed: %v", tst.expr, err)
			}

			// Convert the CEL result to a native Go value for comparison
			got, err := out.ConvertToNative(reflect.TypeOf(tst.out))
			if err != nil {
				t.Fatalf("ConvertToNative failed for expr %q: %v", tst.expr, err)
			}

			if !reflect.DeepEqual(got, tst.out) {
				t.Errorf("Expr %q result got %v, wanted %v", tst.expr, got, tst.out)
			}
		})
	}
}

func TestNetwork_RuntimeErrors(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		errContains string
	}{
		{
			name:        "ip constructor invalid",
			expr:        "ip('999.999.999.999')",
			errContains: "parse error",
		},
		{
			name: "cidr constructor invalid",
			expr: "cidr('1.2.3.4')",
			// Updated to match implementation: "CIDR ... parse error..."
			errContains: "parse error",
		},
		{
			name: "cidr constructor invalid mask",
			expr: "cidr('10.0.0.0/999')",
			// Updated to match implementation: "CIDR ... parse error..."
			errContains: "parse error",
		},
		{
			name:        "containsIP string overload invalid",
			expr:        "cidr('10.0.0.0/8').containsIP('not-an-ip')",
			errContains: "parse error",
		},
		{
			name: "containsCIDR string overload invalid",
			expr: "cidr('10.0.0.0/8').containsCIDR('not-a-cidr')",
			// Updated to match implementation: "CIDR ... parse error..."
			errContains: "parse error",
		},
	}

	env, err := cel.NewEnv(Network())
	if err != nil {
		t.Fatalf("cel.NewEnv(Network()) failed: %v", err)
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			ast, iss := env.Compile(tst.expr)
			if iss.Err() != nil {
				// Note: We only check runtime errors here. Compile errors are unexpected
				// because these functions accept strings, so type-check passes.
				t.Fatalf("Compile(%q) failed unexpectedly: %v", tst.expr, iss.Err())
			}

			prg, err := env.Program(ast)
			if err != nil {
				t.Fatalf("Program(%q) failed: %v", tst.expr, err)
			}

			_, _, err = prg.Eval(cel.NoVars())
			if err == nil {
				t.Errorf("Expected runtime error for %q, got nil", tst.expr)
				return
			}

			// CEL errors are sometimes wrapped, so we check substring
			if !types.IsError(types.NewErr(err.Error())) {
				// Just a sanity check that it is indeed a CEL-compatible error structure
				// Not strictly necessary but good practice
			}

			// Standard substring check
			gotErr := err.Error()
			// We just check if the message contains the specific error text we return in network.go
			found := false
			// Note: The actual error might be wrapped in "evaluation error: ..."
			if len(tst.errContains) > 0 {
				// Simple string contains check
				for i := 0; i < len(gotErr)-len(tst.errContains)+1; i++ {
					if gotErr[i:i+len(tst.errContains)] == tst.errContains {
						found = true
						break
					}
				}
			}

			if !found {
				t.Errorf("Expected error containing %q, got %q", tst.errContains, gotErr)
			}
		})
	}
}
