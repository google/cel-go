// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package policy

import (
	"log"
	"os"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	"gopkg.in/yaml.v3"

	proto3pb "github.com/google/cel-go/test/proto3pb"
)

var (
	policyTests = []struct {
		name      string
		envOpts   []cel.EnvOption
		parseOpts []ParserOption
		expr      string
		expr2     string
	}{
		{
			name: "k8s",
			parseOpts: []ParserOption{func(p *Parser) (*Parser, error) {
				p.TagVisitor = k8sTagHandler()
				return p, nil
			}},
			expr: `
	cel.@block([
	  resource.labels.?environment.orValue("prod"),
	  resource.labels.?break_glass.orValue("false") == "true"],
	  !(@index1 || resource.containers.all(c, c.startsWith(@index0 + ".")))
	    ? optional.of("only %s containers are allowed in namespace %s".format([@index0, resource.namespace]))
	    : optional.none())`,
		},
		{
			name: "nested_rule",
			expr: `
	cel.@block([
	  ["us", "uk", "es"],
	  {"us": false, "ru": false, "ir": false}],
	  ((resource.origin in @index1 && !(resource.origin in @index0))
	    ? optional.of({"banned": true}) : optional.none()).or(
	      optional.of((resource.origin in @index0)
	      ? {"banned": false} : {"banned": true})))`,
		},
		{
			name: "nested_rule2",
			expr: `
	cel.@block([
	  ["us", "uk", "es"],
	  {"us": false, "ru": false, "ir": false}],
	  resource.?user.orValue("").startsWith("bad")
	  ? ((resource.origin in @index1 && !(resource.origin in @index0))
	    ? {"banned": "restricted_region"}
	    : {"banned": "bad_actor"})
	  : (!(resource.origin in @index0)
	    ? {"banned": "unconfigured_region"} : {}))`,
		},
		{
			name: "nested_rule3",
			expr: `
	cel.@block([
	  ["us", "uk", "es"],
	  {"us": false, "ru": false, "ir": false}],
	  resource.?user.orValue("").startsWith("bad")
	  ? optional.of((resource.origin in @index1 && !(resource.origin in @index0))
	    ? {"banned": "restricted_region"} : {"banned": "bad_actor"})
		: (!(resource.origin in @index0)
		  ? optional.of({"banned": "unconfigured_region"}) : optional.none()))`,
		},
		{
			name: "context_pb",
			expr: `
	(single_int32 > google.expr.proto3.test.TestAllTypes{single_int64: 10}.single_int64)
	? optional.of("invalid spec, got single_int32=%d, wanted <= 10".format([single_int32]))
	: ((standalone_enum == google.expr.proto3.test.TestAllTypes.NestedEnum.BAR ||
      google.expr.proto3.test.ImportedGlobalEnum.IMPORT_BAR in imported_enums)
	  ? optional.of("invalid spec, neither nested nor imported enums may refer to BAR or IMPORT_BAR")
	  : optional.none())`,
			envOpts: []cel.EnvOption{
				cel.Types(&proto3pb.TestAllTypes{}),
			},
		},
		{
			name: "pb",
			expr: `
	(spec.single_int32 > google.expr.proto3.test.TestAllTypes{single_int64: 10}.single_int64)
	? optional.of("invalid spec, got single_int32=%d, wanted <= 10".format([spec.single_int32]))
	: ((spec.standalone_enum == google.expr.proto3.test.TestAllTypes.NestedEnum.BAR ||
      google.expr.proto3.test.ImportedGlobalEnum.IMPORT_BAR in spec.imported_enums)
	  ? optional.of("invalid spec, neither nested nor imported enums may refer to BAR or IMPORT_BAR")
	  : optional.none())`,
			envOpts: []cel.EnvOption{
				cel.Types(&proto3pb.TestAllTypes{}),
			},
		},
		{
			name: "required_labels",
			expr: `
	cel.@block([
	  spec.labels,
	  @index0.filter(l, !(l in resource.labels)),
	  resource.labels.filter(l, l in @index0 && @index0[l] != resource.labels[l])],
	(@index1.size() > 0)
	  ? optional.of("missing one or more required labels: %s".format([@index1]))
	  : ((@index2.size() > 0)
	    ? optional.of("invalid values provided on one or more labels: %s".format([@index2]))
		: optional.none()))`,
		},
		{
			name: "restricted_destinations",
			expr: `
    cel.@block([
	  locationCode(origin.ip) == spec.origin,
	  has(request.auth.claims.nationality),
	  @index1 && request.auth.claims.nationality == spec.origin,
	  locationCode(destination.ip) in spec.restricted_destinations,
	  resource.labels.location in spec.restricted_destinations,
	  @index3 || @index4],
	  (@index2 && @index5) ? true : ((!@index1 && @index0 && @index5) ? true : false))`,
			envOpts: []cel.EnvOption{
				cel.Function("locationCode",
					cel.Overload("locationCode_string", []*cel.Type{cel.StringType}, cel.StringType,
						cel.UnaryBinding(func(ip ref.Val) ref.Val {
							switch ip.(types.String) {
							case types.String("10.0.0.1"):
								return types.String("us")
							case types.String("10.0.0.2"):
								return types.String("de")
							default:
								return types.String("ir")
							}
						}))),
			},
		},
		{
			name: "limits",
			expr: `
    cel.@block([
	  "hello",
	  "goodbye",
	  "me",
	  "%s, %s",
	  @index3.format([@index1, @index2])],
	  (now.getHours() >= 20)
	  ? ((now.getHours() < 21)
	    ? optional.of(@index4 + "!")
		: ((now.getHours() < 22)
		  ? optional.of(@index4 + "!!")
		  : ((now.getHours() < 24)
		    ? optional.of(@index4 + "!!!")
			: optional.none())))
	  : optional.of(@index3.format([@index0, @index2])))`,
		},
	}

	policyErrorTests = []struct {
		name         string
		err          string
		compilerOpts []CompilerOption
	}{
		{
			name: "errors",
			err: `ERROR: testdata/errors/policy.yaml:19:1: error configuring import: invalid qualified name: punc.Import!, wanted name of the form 'qualified.name'
 |       punc.Import!
 | ^
ERROR: testdata/errors/policy.yaml:20:12: error configuring import: invalid qualified name: bad import, wanted name of the form 'qualified.name'
 |   - name: "bad import"
 | ...........^
ERROR: testdata/errors/policy.yaml:24:19: undeclared reference to 'spec' (in container '')
 |       expression: spec.labels
 | ..................^
ERROR: testdata/errors/policy.yaml:25:7: invalid variable declaration: overlapping identifier for name 'variables.want'
 |     - name: want
 | ......^
ERROR: testdata/errors/policy.yaml:28:50: Syntax error: mismatched input 'resource' expecting ')'
 |       expression: variables.want.filter(l, !(lin resource.labels))
 | .................................................^
ERROR: testdata/errors/policy.yaml:28:66: Syntax error: extraneous input ')' expecting <EOF>
 |       expression: variables.want.filter(l, !(lin resource.labels))
 | .................................................................^
ERROR: testdata/errors/policy.yaml:30:27: Syntax error: mismatched input '2' expecting {'}', ','}
 |       expression: "{1:305 2:569}"
 | ..........................^
ERROR: testdata/errors/policy.yaml:38:75: Syntax error: extraneous input ']' expecting ')'
 |         "missing one or more required labels: %s".format(variables.missing])
 | ..........................................................................^
ERROR: testdata/errors/policy.yaml:41:67: undeclared reference to 'format' (in container '')
 |         "invalid values provided on one or more labels: %s".format([variables.invalid])
 | ..................................................................^
ERROR: testdata/errors/policy.yaml:45:16: incompatible output types: bool not assignable to string
 |       output: "'false'"
 | ...............^`,
		},
		{
			name: "limits",
			err: `ERROR: testdata/limits/policy.yaml:22:14: variable exceeds nested expression limit
 |     - name: "person"
 | .............^`,
			compilerOpts: []CompilerOption{MaxNestedExpressions(2)},
		},
		{
			name: "limits",
			err: `ERROR: testdata/limits/policy.yaml:30:9: rule exceeds nested expression limit
 |         id: "farewells"
 | ........^`,
			compilerOpts: []CompilerOption{MaxNestedExpressions(5)},
		},
		{
			name: "errors_unreachable",
			err: `ERROR: testdata/errors_unreachable/policy.yaml:28:9: rule creates unreachable outputs
 |         match:
 | ........^
ERROR: testdata/errors_unreachable/policy.yaml:36:13: match creates unreachable outputs
 |           - output: |
 | ............^`,
		},
	}
)

func k8sTagHandler() TagVisitor {
	return k8sAdmissionTagHandler{TagVisitor: DefaultTagVisitor()}
}

type k8sAdmissionTagHandler struct {
	TagVisitor
}

func (k8sAdmissionTagHandler) PolicyTag(ctx ParserContext, id int64, tagName string, node *yaml.Node, policy *Policy) {
	switch tagName {
	case "kind":
		policy.SetMetadata("kind", ctx.NewString(node).Value)
	case "metadata":
		m := k8sMetadata{}
		if err := node.Decode(&m); err != nil {
			ctx.ReportErrorAtID(id, "invalid yaml metadata node: %v, error: %w", node, err)
			return
		}
	case "spec":
		spec := ctx.ParseRule(ctx, policy, node)
		policy.SetRule(spec)
	default:
		ctx.ReportErrorAtID(id, "unsupported policy tag: %s", tagName)
	}
}

func (k8sAdmissionTagHandler) RuleTag(ctx ParserContext, id int64, tagName string, node *yaml.Node, policy *Policy, r *Rule) {
	switch tagName {
	case "failurePolicy":
		policy.SetMetadata(tagName, ctx.NewString(node).Value)
	case "matchConstraints":
		m := k8sMatchConstraints{}
		if err := node.Decode(&m); err != nil {
			ctx.ReportErrorAtID(id, "invalid yaml matchConstraints node: %v, error: %w", node, err)
			return
		}
	case "validations":
		id := ctx.CollectMetadata(node)
		if node.LongTag() != "tag:yaml.org,2002:seq" {
			ctx.ReportErrorAtID(id, "invalid 'validations' type, expected list got: %s", node.LongTag())
			return
		}
		for _, val := range node.Content {
			r.AddMatch(ctx.ParseMatch(ctx, policy, val))
		}
	default:
		ctx.ReportErrorAtID(id, "unsupported rule tag: %s", tagName)
	}
}

func (k8sAdmissionTagHandler) MatchTag(ctx ParserContext, id int64, tagName string, node *yaml.Node, policy *Policy, m *Match) {
	if m.Output().Value == "" {
		m.SetOutput(ValueString{Value: "'invalid admission request'"})
	}
	switch tagName {
	case "expression":
		// The K8s expression to validate must return false in order to generate a violation message.
		condition := ctx.NewString(node)
		condition.Value = "!(" + condition.Value + ")"
		m.SetCondition(condition)
	case "messageExpression":
		m.SetOutput(ctx.NewString(node))
	}
}

type k8sMetadata struct {
	Name string `yaml:"name"`
}

type k8sMatchConstraints struct {
	ResourceRules []k8sResourceRule `yaml:"resourceRules"`
}

type k8sResourceRule struct {
	APIGroups   []string `yaml:"apiGroups"`
	APIVersions []string `yaml:"apiVersions"`
	Operations  []string `yaml:"operations"`
}

func readPolicy(t testing.TB, fileName string) *Source {
	t.Helper()
	policyBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) failed: %v", fileName, err)
	}
	return ByteSource(policyBytes, fileName)
}

func readPolicyConfig(t testing.TB, fileName string) *Config {
	t.Helper()
	testCaseBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) failed: %v", fileName, err)
	}
	config := &Config{}
	err = yaml.Unmarshal(testCaseBytes, config)
	if err != nil {
		log.Fatalf("yaml.Unmarshal(%s) error: %v", fileName, err)
	}
	return config
}

func readTestSuite(t testing.TB, fileName string) *TestSuite {
	t.Helper()
	testCaseBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) failed: %v", fileName, err)
	}
	suite := &TestSuite{}
	err = yaml.Unmarshal(testCaseBytes, suite)
	if err != nil {
		log.Fatalf("yaml.Unmarshal(%s) error: %v", fileName, err)
	}
	return suite
}
