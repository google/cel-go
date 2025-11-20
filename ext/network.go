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
	"net"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// Network returns a cel.EnvOption to configure extended functions for network
// address parsing, inspection, and CIDR range manipulation.
//
// Note: This library defines global functions `ip`, `cidr`, `isIP`, and
// `isCIDR`. If you are currently using variables named `ip` or `cidr`, these
// functions will likely work as intended, however there is a chance for
// collision.
//
// The library closely mirrors the behavior of the Kubernetes CEL network
// libraries, treating IP addresses and CIDR ranges as opaque types with
// specific member functions.
//
// # IP Addresses
//
// The `ip` function converts a string to an IP address (IPv4 or IPv6). If the
// string is not a valid IP, an error is returned. The `isIP` function checks
// if a string is a valid IP address without throwing an error.
//
//	ip(string) -> ip
//	isIP(string) -> bool
//
// Examples:
//
//	ip('127.0.0.1')
//	ip('::1')
//	isIP('1.2.3.4') // true
//	isIP('invalid') // false
//
// # CIDR Ranges
//
// The `cidr` function converts a string to a Classless Inter-Domain Routing
// (CIDR) range. If the string is not valid, an error is returned. The `isCIDR`
// function checks if a string is a valid CIDR notation.
//
//	cidr(string) -> cidr
//	isCIDR(string) -> bool
//
// Examples:
//
//	cidr('192.168.0.0/24')
//	cidr('::1/128')
//	isCIDR('10.0.0.0/8') // true
//
// # IP Member Functions
//
// IP objects support various inspection methods.
//
//	<ip>.family() -> int
//	<ip>.isCanonical() -> bool
//	<ip>.isLoopback() -> bool
//	<ip>.isGlobalUnicast() -> bool
//	<ip>.isLinkLocalMulticast() -> bool
//	<ip>.isLinkLocalUnicast() -> bool
//	<ip>.isUnspecified() -> bool
//
// Note on Canonicalization: `isCanonical()` returns true if the input string
// used to construct the IP matches the RFC 5952 canonical string representation
// of that address.
//
// Examples:
//
//	ip('127.0.0.1').family() == 4
//	ip('::1').family() == 6
//	ip('127.0.0.1').isLoopback() == true
//	ip('2001:db8::1').isCanonical() == true  // RFC 5952 format
//	ip('2001:DB8::1').isCanonical() == false // Uppercase is not canonical
//	ip('2001:db8:0:0:0:0:0:1').isCanonical() == false // Expanded is not canonical
//
// # CIDR Member Functions
//
// CIDR objects support containment checks and property extraction.
//
//	<cidr>.containsIP(ip|string) -> bool
//	<cidr>.containsCIDR(cidr|string) -> bool
//	<cidr>.ip() -> ip
//	<cidr>.masked() -> cidr
//	<cidr>.prefixLength() -> int
//
// Examples:
//
//	cidr('10.0.0.0/8').containsIP(ip('10.0.0.1')) == true
//	cidr('10.0.0.0/8').containsIP('10.0.0.1') == true
//	cidr('10.0.0.0/8').containsCIDR('10.1.0.0/16') == true
//	cidr('192.168.1.5/24').ip() == ip('192.168.1.5')
//	cidr('192.168.1.5/24').masked() == cidr('192.168.1.0/24')
//	cidr('192.168.1.0/24').prefixLength() == 24
func Network() cel.EnvOption {
	return cel.Lib(&networkLib{})
}

const (
	// Function names
	isIPFunc             = "isIP"
	isCIDRFunc           = "isCIDR"
	ipFunc               = "ip"
	cidrFunc             = "cidr"
	familyFunc           = "family"
	isCanonicalFunc      = "isCanonical"
	isLoopbackFunc       = "isLoopback"
	isGlobalUnicastFunc  = "isGlobalUnicast"
	isUnspecifiedFunc    = "isUnspecified"
	isLinkLocalMcastFunc = "isLinkLocalMulticast"
	isLinkLocalUcastFunc = "isLinkLocalUnicast"
	containsIPFunc       = "containsIP"
	containsCIDRFunc     = "containsCIDR"
	maskedFunc           = "masked"
	prefixLengthFunc     = "prefixLength"
	ipFromCIDRFunc       = "ip"
)

var (
	// Definitions for the Opaque Types
	networkIPType   = types.NewTypeValue("network.IP", traits.ReceiverType)
	networkCIDRType = types.NewTypeValue("network.CIDR", traits.ReceiverType)
)

type networkLib struct{}

func (*networkLib) LibraryName() string {
	return "cel.lib.ext.network"
}

func (*networkLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		// 1. Register the types
		cel.Types(
			networkIPType,
			networkCIDRType,
		),
		// 2. Register the Adapter (Correctly placed here)
		cel.CustomTypeAdapter(&networkAdapter{
			Adapter: types.DefaultTypeAdapter,
		}),
		// 3. Register the Functions
		// Global Checkers
		cel.Function(isIPFunc,
			cel.Overload("isIP_string", []*cel.Type{cel.StringType}, cel.BoolType,
				cel.UnaryBinding(netIsIP)),
		),
		cel.Function(isCIDRFunc,
			cel.Overload("isCIDR_string", []*cel.Type{cel.StringType}, cel.BoolType,
				cel.UnaryBinding(netIsCIDR)),
		),
		// Constructors
		cel.Function(ipFunc,
			cel.Overload("ip_string", []*cel.Type{cel.StringType}, networkIPType,
				cel.UnaryBinding(netIPString)),
		),
		cel.Function(cidrFunc,
			cel.Overload("cidr_string", []*cel.Type{cel.StringType}, networkCIDRType,
				cel.UnaryBinding(netCIDRString)),
		),
		// IP Member Functions
		cel.Function(familyFunc,
			cel.MemberOverload("ip_family", []*cel.Type{networkIPType}, cel.IntType,
				cel.UnaryBinding(netIPFamily)),
		),
		cel.Function(isCanonicalFunc,
			cel.MemberOverload("ip_isCanonical", []*cel.Type{networkIPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsCanonical)),
		),
		cel.Function(isLoopbackFunc,
			cel.MemberOverload("ip_isLoopback", []*cel.Type{networkIPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsLoopback)),
		),
		cel.Function(isGlobalUnicastFunc,
			cel.MemberOverload("ip_isGlobalUnicast", []*cel.Type{networkIPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsGlobalUnicast)),
		),
		cel.Function(isUnspecifiedFunc,
			cel.MemberOverload("ip_isUnspecified", []*cel.Type{networkIPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsUnspecified)),
		),
		cel.Function(isLinkLocalMcastFunc,
			cel.MemberOverload("ip_isLinkLocalMulticast", []*cel.Type{networkIPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsLinkLocalMulticast)),
		),
		cel.Function(isLinkLocalUcastFunc,
			cel.MemberOverload("ip_isLinkLocalUnicast", []*cel.Type{networkIPType}, cel.BoolType,
				cel.UnaryBinding(netIPIsLinkLocalUnicast)),
		),
		// CIDR Member Functions
		cel.Function(containsIPFunc,
			cel.MemberOverload("cidr_containsIP_ip", []*cel.Type{networkCIDRType, networkIPType}, cel.BoolType,
				cel.BinaryBinding(netCIDRContainsIP)),
			cel.MemberOverload("cidr_containsIP_string", []*cel.Type{networkCIDRType, cel.StringType}, cel.BoolType,
				cel.BinaryBinding(netCIDRContainsIPString)),
		),
		cel.Function(containsCIDRFunc,
			cel.MemberOverload("cidr_containsCIDR_cidr", []*cel.Type{networkCIDRType, networkCIDRType}, cel.BoolType,
				cel.BinaryBinding(netCIDRContainsCIDR)),
			cel.MemberOverload("cidr_containsCIDR_string", []*cel.Type{networkCIDRType, cel.StringType}, cel.BoolType,
				cel.BinaryBinding(netCIDRContainsCIDRString)),
		),
		cel.Function(maskedFunc,
			cel.MemberOverload("cidr_masked", []*cel.Type{networkCIDRType}, networkCIDRType,
				cel.UnaryBinding(netCIDRMasked)),
		),
		cel.Function(prefixLengthFunc,
			cel.MemberOverload("cidr_prefixLength", []*cel.Type{networkCIDRType}, cel.IntType,
				cel.UnaryBinding(netCIDRPrefixLength)),
		),
		cel.Function(ipFromCIDRFunc,
			cel.MemberOverload("cidr_ip", []*cel.Type{networkCIDRType}, networkIPType,
				cel.UnaryBinding(netCIDRIP)),
		),
	}
}

func (*networkLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

// networkAdapter implements types.Adapter to handle net.IP and *net.IPNet conversion.
type networkAdapter struct {
	types.Adapter
}

func (a *networkAdapter) NativeToValue(value any) ref.Val {
	switch v := value.(type) {
	case net.IP:
		// If passing a raw net.IP, we assume standard string representation
		return IP{IP: v, Str: v.String()}
	case *net.IPNet:
		return CIDR{IPNet: v, Str: v.String()}
	}
	return a.Adapter.NativeToValue(value)
}

// --- Implementation Logic ---

func netIsIP(val ref.Val) ref.Val {
	s, ok := val.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	return types.Bool(net.ParseIP(string(s)) != nil)
}

func netIsCIDR(val ref.Val) ref.Val {
	s, ok := val.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	_, _, err := net.ParseCIDR(string(s))
	return types.Bool(err == nil)
}

func netIPString(val ref.Val) ref.Val {
	s, ok := val.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	str := string(s)
	ip := net.ParseIP(str)
	if ip == nil {
		return types.NewErr("invalid ip address: %s", str)
	}
	return IP{IP: ip, Str: str}
}

func netCIDRString(val ref.Val) ref.Val {
	s, ok := val.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	str := string(s)
	ip, ipNet, err := net.ParseCIDR(str)
	if err != nil {
		return types.NewErr("invalid cidr range: %s", str)
	}
	// Store the specific IP (which might have host bits set) alongside the network
	return CIDR{IPNet: ipNet, ExtraIP: ip, Str: str}
}

func netIPFamily(val ref.Val) ref.Val {
	ip, ok := val.(IP)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	if ip.IP.To4() != nil {
		return types.Int(4)
	}
	return types.Int(6)
}

func netIPIsCanonical(val ref.Val) ref.Val {
	ip, ok := val.(IP)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	return types.Bool(ip.Str == ip.IP.String())
}

func netIPIsLoopback(val ref.Val) ref.Val {
	ip, ok := val.(IP)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	return types.Bool(ip.IP.IsLoopback())
}

func netIPIsGlobalUnicast(val ref.Val) ref.Val {
	ip, ok := val.(IP)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	return types.Bool(ip.IP.IsGlobalUnicast())
}

func netIPIsUnspecified(val ref.Val) ref.Val {
	ip, ok := val.(IP)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	return types.Bool(ip.IP.IsUnspecified())
}

func netIPIsLinkLocalMulticast(val ref.Val) ref.Val {
	ip, ok := val.(IP)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	return types.Bool(ip.IP.IsLinkLocalMulticast())
}

func netIPIsLinkLocalUnicast(val ref.Val) ref.Val {
	ip, ok := val.(IP)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	return types.Bool(ip.IP.IsLinkLocalUnicast())
}

func netCIDRContainsIP(lhs, rhs ref.Val) ref.Val {
	cidr, ok := lhs.(CIDR)
	if !ok {
		return types.MaybeNoSuchOverloadErr(lhs)
	}
	ip, ok := rhs.(IP)
	if !ok {
		return types.MaybeNoSuchOverloadErr(rhs)
	}
	return types.Bool(cidr.IPNet.Contains(ip.IP))
}

func netCIDRContainsIPString(lhs, rhs ref.Val) ref.Val {
	cidr, ok := lhs.(CIDR)
	if !ok {
		return types.MaybeNoSuchOverloadErr(lhs)
	}
	s, ok := rhs.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(rhs)
	}
	ip := net.ParseIP(string(s))
	if ip == nil {
		return types.NewErr("invalid ip address: %s", s)
	}
	return types.Bool(cidr.IPNet.Contains(ip))
}

func netCIDRContainsCIDR(lhs, rhs ref.Val) ref.Val {
	parent, ok := lhs.(CIDR)
	if !ok {
		return types.MaybeNoSuchOverloadErr(lhs)
	}
	child, ok := rhs.(CIDR)
	if !ok {
		return types.MaybeNoSuchOverloadErr(rhs)
	}
	ones1, _ := parent.IPNet.Mask.Size()
	ones2, _ := child.IPNet.Mask.Size()
	return types.Bool(parent.IPNet.Contains(child.IPNet.IP) && ones1 <= ones2)
}

func netCIDRContainsCIDRString(lhs, rhs ref.Val) ref.Val {
	parent, ok := lhs.(CIDR)
	if !ok {
		return types.MaybeNoSuchOverloadErr(lhs)
	}
	s, ok := rhs.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(rhs)
	}
	_, childNet, err := net.ParseCIDR(string(s))
	if err != nil {
		return types.NewErr("invalid cidr range: %s", s)
	}
	ones1, _ := parent.IPNet.Mask.Size()
	ones2, _ := childNet.Mask.Size()
	return types.Bool(parent.IPNet.Contains(childNet.IP) && ones1 <= ones2)
}

func netCIDRMasked(val ref.Val) ref.Val {
	cidr, ok := val.(CIDR)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	maskedIP := cidr.IPNet.IP.Mask(cidr.IPNet.Mask)
	newNet := &net.IPNet{IP: maskedIP, Mask: cidr.IPNet.Mask}
	return CIDR{IPNet: newNet, ExtraIP: maskedIP, Str: newNet.String()}
}

func netCIDRPrefixLength(val ref.Val) ref.Val {
	cidr, ok := val.(CIDR)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	ones, _ := cidr.IPNet.Mask.Size()
	return types.Int(ones)
}

func netCIDRIP(val ref.Val) ref.Val {
	cidr, ok := val.(CIDR)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	return IP{IP: cidr.ExtraIP, Str: cidr.ExtraIP.String()}
}

// --- Opaque Type Wrappers ---

// IP is an exported CEL value that wraps net.IP.
// It implements ref.Val.
type IP struct {
	net.IP
	Str string
}

func (i IP) ConvertToNative(typeDesc reflect.Type) (any, error) {
	if typeDesc == reflect.TypeOf(net.IP{}) {
		return i.IP, nil
	}
	if typeDesc.Kind() == reflect.String {
		return i.IP.String(), nil
	}
	return nil, fmt.Errorf("unsupported type conversion to '%v'", typeDesc)
}

func (i IP) ConvertToType(typeValue ref.Type) ref.Val {
	switch typeValue {
	case types.StringType:
		return types.String(i.IP.String())
	case networkIPType:
		return i
	case types.TypeType:
		return networkIPType
	}
	return types.NewErr("type conversion error from '%s' to '%s'", networkIPType, typeValue)
}

func (i IP) Equal(other ref.Val) ref.Val {
	o, ok := other.(IP)
	if !ok {
		return types.ValOrErr(other, "no such overload")
	}
	return types.Bool(i.IP.Equal(o.IP))
}

func (i IP) Type() ref.Type {
	return networkIPType
}

func (i IP) Value() any {
	return i.IP
}

// CIDR is an exported CEL value that wraps *net.IPNet.
// It implements ref.Val.
type CIDR struct {
	*net.IPNet
	ExtraIP net.IP
	Str     string
}

func (c CIDR) ConvertToNative(typeDesc reflect.Type) (any, error) {
	if typeDesc == reflect.TypeOf(&net.IPNet{}) {
		return c.IPNet, nil
	}
	if typeDesc.Kind() == reflect.String {
		return c.IPNet.String(), nil
	}
	return nil, fmt.Errorf("unsupported type conversion to '%v'", typeDesc)
}

func (c CIDR) ConvertToType(typeValue ref.Type) ref.Val {
	switch typeValue {
	case types.StringType:
		return types.String(c.IPNet.String())
	case networkCIDRType:
		return c
	case types.TypeType:
		return networkCIDRType
	}
	return types.NewErr("type conversion error from '%s' to '%s'", networkCIDRType, typeValue)
}

func (c CIDR) Equal(other ref.Val) ref.Val {
	o, ok := other.(CIDR)
	if !ok {
		return types.ValOrErr(other, "no such overload")
	}
	return types.Bool(c.IPNet.IP.Equal(o.IPNet.IP) && c.IPNet.Mask.String() == o.IPNet.Mask.String())
}

func (c CIDR) Type() ref.Type {
	return networkCIDRType
}

func (c CIDR) Value() any {
	return c.IPNet
}
