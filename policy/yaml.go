// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package policy

import (
	"go.yaml.in/yaml/v3"
)

// YAMLHelper provides helper methods for working with YAML nodes.
type YAMLHelper struct{}

// IsList returns true if the YAML node is a list.
func (YAMLHelper) IsList(node *yaml.Node) bool {
	return isTypeTag(node, yamlList)
}

// RangeList iterates over the list and calls the provided function for each element. If the
// function returns false, the iteration will stop.
func (yc YAMLHelper) RangeList(node *yaml.Node, fn func(int, *yaml.Node) bool) {
	if !yc.IsList(node) {
		return
	}
	for i, val := range node.Content {
		if !fn(i, val) {
			break
		}
	}
}

// IsMap returns true if the YAML node is a map.
func (YAMLHelper) IsMap(node *yaml.Node) bool {
	return isTypeTag(node, yamlMap) && len(node.Content)%2 == 0
}

// RangeMap iterates over the map and calls the provided function for each key-value pair. If the
// function returns false, the iteration will stop.
func (yc YAMLHelper) RangeMap(node *yaml.Node, fn func(key, val *yaml.Node) bool) {
	if !yc.IsMap(node) {
		return
	}
	for i := 0; i < len(node.Content); i += 2 {
		key, val := normalizeEntry(node.Content, i)
		if !fn(key, val) {
			break
		}
	}
}

// IsString returns true if the YAML node is a string.
func (YAMLHelper) IsString(node *yaml.Node) bool {
	return isTypeTag(node, yamlString)
}

// IsBool returns true if the YAML node is a boolean.
func (YAMLHelper) IsBool(node *yaml.Node) bool {
	return isTypeTag(node, yamlBool)
}

// IsNull returns true if the YAML node is a null.
func (YAMLHelper) IsNull(node *yaml.Node) bool {
	return isTypeTag(node, yamlNull)
}

// IsNumber returns true if the YAML node is a number.
func (YAMLHelper) IsNumber(node *yaml.Node) bool {
	return isTypeTag(node, yamlInt) || isTypeTag(node, yamlDouble)
}

// IsInteger returns true if the YAML node is an integer.
func (YAMLHelper) IsInteger(node *yaml.Node) bool {
	return isTypeTag(node, yamlInt)
}

// IsDouble returns true if the YAML node is a double.
func (YAMLHelper) IsDouble(node *yaml.Node) bool {
	return isTypeTag(node, yamlDouble)
}

// IsTimestamp returns true if the YAML node is a timestamp.
func (YAMLHelper) IsTimestamp(node *yaml.Node) bool {
	return isTypeTag(node, yamlTimestamp)
}

func isTypeTag(node *yaml.Node, wantTag yamlNodeType) bool {
	if yt, found := yamlTypes[node.LongTag()]; found {
		return yt == wantTag
	}
	return false
}
