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
	"fmt"
	"net/netip"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// Network returns a cel.EnvOption to configure extended functions for network
// address parsing, inspection, and CIDR range manipulation.
func Network() cel.EnvOption {
	return func(e *cel.Env) (*cel.Env, error) {
		// Install the library (Types and Functions)
		e, err := cel.Lib(&networkLib{})(e)
		if err != nil {
			return nil, err
		}

		// Install the Adapter (Wrapping the existing one)
		adapter := &networkAdapter{Adapter: e.CELTypeAdapter()}
		return cel.CustomTypeAdapter(adapter)(e)
	}
}

const (
	// Function names matching Kubernetes implementation
	isIPFunc             = "isIP"
	isCIDRFunc           = "isCIDR"
	ipFunc               = "ip"
	cidrFunc             = "cidr"
	familyFunc           = "family"
	isCanonicalFunc      = "ip.isCanonical"
	isLoopbackFunc       = "isLoopback"
	isGlobalUnicastFunc  = "isGlobalUnicast"
	isUnspecifiedFunc    = "isUnspecified"
	isLinkLocalMcastFunc = "isLinkLocalMulticast"
	isLinkLocalUcastFunc = "isLinkLocalUnicast"
	containsIPFunc       = "containsIP"
	containsCIDRFunc     = "containsCIDR"
	maskedFunc           = "masked"
	prefixLengthFunc     = "prefixLength"
)

var (
	// Definitions for the Opaque Types
	IPType   = cel.OpaqueType("net.IP")
	CIDRType = cel.OpaqueType("net.CIDR")
)

type networkLib struct{}

func (*networkLib) LibraryName() string {
	return "cel.lib.ext.network"
}

func (*networkLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// 1. Register Types
		cel.Types(
			IPType,
			CIDRType,
		),

		// 2. Register Functions
		cel.Function(isIPFunc,
			cel.Overload("is_ip", []*cel.Type{cel.StringType}, cel.BoolType,
				cel.UnaryBinding(netIsIP)),
		),
		cel.Function(isCIDRFunc,
			cel.Overload("is_cidr", []*cel.Type{cel.StringType}, cel.BoolType,
				cel.UnaryBinding(netIsCIDR)),
		),
		cel.Function(ipFunc,
			cel.Overload("ip", []*cel.Type{cel.StringType}, IPType,
				cel.UnaryBinding(netIPString)),
			cel.MemberOverload("cidr_ip", []*cel.Type{CIDRType}, IPType,
				cel.UnaryBinding(netCIDRIP)),
		),
		cel.Function(cidrFunc,
			cel.Overload("cidr", []*cel.Type{cel.StringType}, CIDRType,
				cel.UnaryBinding(netCIDRString)),
		),
		cel.Function(familyFunc,
			cel.MemberOverload("ip_family", []*cel.Type{IPType}, cel.IntType,
				cel.UnaryBinding(netIPFamily)),
		),
		cel.Function(isCanonicalFunc,
			cel.Overload("ip_is_canonical", []*cel.Type{cel.StringType}, cel.BoolType,
				cel.UnaryBinding(netIPIsCanonical)),
		),
		cel.Function(isLoopbackFunc,
			cel.MemberOverload("ip_is_loopback", []*cel.Type{IPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsLoopback)),
		),
		cel.Function(isGlobalUnicastFunc,
			cel.MemberOverload("ip_is_global_unicast", []*cel.Type{IPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsGlobalUnicast)),
		),
		cel.Function(isUnspecifiedFunc,
			cel.MemberOverload("ip_is_unspecified", []*cel.Type{IPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsUnspecified)),
		),
		cel.Function(isLinkLocalMcastFunc,
			cel.MemberOverload("ip_is_link_local_multicast", []*cel.Type{IPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsLinkLocalMulticast)),
		),
		cel.Function(isLinkLocalUcastFunc,
			cel.MemberOverload("ip_is_link_local_unicast", []*cel.Type{IPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsLinkLocalUnicast)),
		),
		cel.Function(containsIPFunc,
			cel.MemberOverload("cidr_contains_ip", []*cel.Type{CIDRType, IPType}, cel.BoolType,
				cel.BinaryBinding(netCIDRContainsIP)),
			cel.MemberOverload("cidr_contains_ip_string", []*cel.Type{CIDRType, cel.StringType}, cel.BoolType,
				cel.BinaryBinding(netCIDRContainsIPString)),
		),
		cel.Function(containsCIDRFunc,
			cel.MemberOverload("cidr_contains_cidr", []*cel.Type{CIDRType, CIDRType}, cel.BoolType,
				cel.BinaryBinding(netCIDRContainsCIDR)),
			cel.MemberOverload("cidr_contains_cidr_string", []*cel.Type{CIDRType, cel.StringType}, cel.BoolType,
				cel.BinaryBinding(netCIDRContainsCIDRString)),
		),
		cel.Function(maskedFunc,
			cel.MemberOverload("cidr_masked", []*cel.Type{CIDRType}, CIDRType,
				cel.UnaryBinding(netCIDRMasked)),
		),
		cel.Function(prefixLengthFunc,
			cel.MemberOverload("cidr_prefix_length", []*cel.Type{CIDRType}, cel.IntType,
				cel.UnaryBinding(netCIDRPrefixLength)),
		),
	}
}

func (*networkLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

// networkAdapter adapts netip types while preserving existing adapters.
type networkAdapter struct {
	types.Adapter
}

func (a *networkAdapter) NativeToValue(value any) ref.Val {
	switch v := value.(type) {
	case netip.Addr:
		return IP{Addr: v}
	case netip.Prefix:
		return CIDR{Prefix: v}
	}
	// Delegate to the wrapped adapter (e.g., Protobuf adapter)
	return a.Adapter.NativeToValue(value)
}

// --- Implementation Logic ---

// parseIPAddr parses a string into an IP address.
// We use this function to parse IP addresses in the CEL library
// so that we can share the common logic of rejecting IP addresses
// that contain zones or are IPv4-mapped IPv6 addresses.
func parseIPAddr(raw string) (netip.Addr, error) {
	addr, err := netip.ParseAddr(raw)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("IP Address %q parse error during conversion from string: %v", raw, err)
	}

	if addr.Zone() != "" {
		return netip.Addr{}, fmt.Errorf("IP address %q with zone value is not allowed", raw)
	}

	if addr.Is4In6() {
		return netip.Addr{}, fmt.Errorf("IPv4-mapped IPv6 address %q is not allowed", raw)
	}

	return addr, nil
}

// parseCIDR parses a string into a CIDR/Prefix.
// We use this function to parse CIDRs in the CEL library
// so that we can share the common logic of rejecting CIDRs
// that contain zones or are IPv4-mapped IPv6 addresses.
func parseCIDR(raw string) (netip.Prefix, error) {
	prefix, err := netip.ParsePrefix(raw)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("CIDR %q parse error during conversion from string: %v", raw, err)
	}

	// netip.Prefix.Addr() returns the address part of the prefix
	if prefix.Addr().Zone() != "" {
		return netip.Prefix{}, fmt.Errorf("CIDR %q with zone value is not allowed", raw)
	}

	if prefix.Addr().Is4In6() {
		return netip.Prefix{}, fmt.Errorf("IPv4-mapped IPv6 address %q is not allowed", raw)
	}

	return prefix, nil
}

func netIsIP(val ref.Val) ref.Val {
	s := val.(types.String)
	_, err := parseIPAddr(string(s))
	return types.Bool(err == nil)
}

func netIsCIDR(val ref.Val) ref.Val {
	s := val.(types.String)
	_, err := parseCIDR(string(s))
	return types.Bool(err == nil)
}

func netIPString(val ref.Val) ref.Val {
	s := val.(types.String)
	str := string(s)
	addr, err := parseIPAddr(str)
	if err != nil {
		return types.WrapErr(err)
	}
	return IP{Addr: addr}
}

func netCIDRString(val ref.Val) ref.Val {
	s := val.(types.String)
	str := string(s)
	prefix, err := parseCIDR(str)
	if err != nil {
		return types.WrapErr(err)
	}
	return CIDR{Prefix: prefix}
}

func netIPFamily(val ref.Val) ref.Val {
	ip := val.(IP)
	if ip.Addr.Is4() {
		return types.Int(4)
	}
	return types.Int(6)
}

func netIPIsCanonical(val ref.Val) ref.Val {
	s := val.(types.String)
	str := string(s)
	addr, err := parseIPAddr(str)
	if err != nil {
		return types.WrapErr(err)
	}
	return types.Bool(addr.String() == str)
}

func netIPIsLoopback(val ref.Val) ref.Val {
	ip := val.(IP)
	return types.Bool(ip.Addr.IsLoopback())
}

func netIPIsGlobalUnicast(val ref.Val) ref.Val {
	ip := val.(IP)
	return types.Bool(ip.Addr.IsGlobalUnicast())
}

func netIPIsUnspecified(val ref.Val) ref.Val {
	ip := val.(IP)
	return types.Bool(ip.Addr.IsUnspecified())
}

func netIPIsLinkLocalMulticast(val ref.Val) ref.Val {
	ip := val.(IP)
	return types.Bool(ip.Addr.IsLinkLocalMulticast())
}

func netIPIsLinkLocalUnicast(val ref.Val) ref.Val {
	ip := val.(IP)
	return types.Bool(ip.Addr.IsLinkLocalUnicast())
}

func netCIDRContainsIP(lhs, rhs ref.Val) ref.Val {
	cidr := lhs.(CIDR)
	ip := rhs.(IP)
	return types.Bool(cidr.Prefix.Contains(ip.Addr))
}

func netCIDRContainsIPString(lhs, rhs ref.Val) ref.Val {
	cidr := lhs.(CIDR)
	s := rhs.(types.String)
	addr, err := parseIPAddr(string(s))
	if err != nil {
		return types.WrapErr(err)
	}
	return types.Bool(cidr.Prefix.Contains(addr))
}

func netCIDRContainsCIDR(lhs, rhs ref.Val) ref.Val {
	parent := lhs.(CIDR)
	child := rhs.(CIDR)
	// Matches K8s logic: Must overlap and parent must be "larger" (smaller or equal bit count)
	return types.Bool(parent.Prefix.Overlaps(child.Prefix) && parent.Prefix.Bits() <= child.Prefix.Bits())
}

func netCIDRContainsCIDRString(lhs, rhs ref.Val) ref.Val {
	parent := lhs.(CIDR)
	s := rhs.(types.String)
	childPrefix, err := parseCIDR(string(s))
	if err != nil {
		return types.WrapErr(err)
	}
	return types.Bool(parent.Prefix.Overlaps(childPrefix) && parent.Prefix.Bits() <= childPrefix.Bits())
}

func netCIDRMasked(val ref.Val) ref.Val {
	cidr := val.(CIDR)
	return CIDR{Prefix: cidr.Prefix.Masked()}
}

func netCIDRPrefixLength(val ref.Val) ref.Val {
	cidr := val.(CIDR)
	return types.Int(cidr.Prefix.Bits())
}

func netCIDRIP(val ref.Val) ref.Val {
	cidr := val.(CIDR)
	return IP{Addr: cidr.Prefix.Addr()}
}

// --- Opaque Type Wrappers ---

// IP is an exported CEL value that wraps netip.Addr.
type IP struct {
	netip.Addr
}

func (i IP) ConvertToNative(typeDesc reflect.Type) (any, error) {
	// Use reflect.TypeFor to avoid instantiating netip.Addr{}
	if typeDesc == reflect.TypeFor[*netip.Addr]() {
		return i.Addr, nil
	}
	if typeDesc.Kind() == reflect.String {
		return i.Addr.String(), nil
	}
	return nil, fmt.Errorf("unsupported type conversion to '%v'", typeDesc)
}

func (i IP) ConvertToType(typeValue ref.Type) ref.Val {
	switch typeValue {
	case types.StringType:
		return types.String(i.Addr.String())
	case IPType:
		return i
	case types.TypeType:
		return IPType
	}
	return types.NewErr("type conversion error from '%s' to '%s'", IPType, typeValue)
}

func (i IP) Equal(other ref.Val) ref.Val {
	o, ok := other.(IP)
	if !ok {
		return types.ValOrErr(other, "no such overload")
	}
	return types.Bool(i.Addr == o.Addr)
}

func (i IP) Type() ref.Type {
	return IPType
}

func (i IP) Value() any {
	return i.Addr
}

// CIDR is an exported CEL value that wraps netip.Prefix.
type CIDR struct {
	netip.Prefix
}

func (c CIDR) ConvertToNative(typeDesc reflect.Type) (any, error) {
	// Use reflect.TypeFor to avoid instantiating netip.Prefix{}
	if typeDesc == reflect.TypeFor[netip.Prefix]() {
		return c.Prefix, nil
	}
	if typeDesc.Kind() == reflect.String {
		return c.Prefix.String(), nil
	}
	return nil, fmt.Errorf("unsupported type conversion to '%v'", typeDesc)
}

func (c CIDR) ConvertToType(typeValue ref.Type) ref.Val {
	switch typeValue {
	case types.StringType:
		return types.String(c.Prefix.String())
	case CIDRType:
		return c
	case types.TypeType:
		return CIDRType
	}
	return types.NewErr("type conversion error from '%s' to '%s'", CIDRType, typeValue)
}

func (c CIDR) Equal(other ref.Val) ref.Val {
	o, ok := other.(CIDR)
	if !ok {
		return types.ValOrErr(other, "no such overload")
	}
	return types.Bool(c.Prefix == o.Prefix)
}

func (c CIDR) Type() ref.Type {
	return CIDRType
}

func (c CIDR) Value() any {
	return c.Prefix
}
