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
	"github.com/google/cel-go/test/proto3pb"

	"gopkg.in/yaml.v3"
)

var (
	policyTests = []struct {
		name      string
		envOpts   []cel.EnvOption
		parseOpts []ParserOption
		expr      string
	}{
		{
			name: "k8s",
			parseOpts: []ParserOption{func(p *Parser) (*Parser, error) {
				p.TagVisitor = k8sTagHandler()
				return p, nil
			}},
			expr: `
    cel.bind(variables.env, resource.labels.?environment.orValue("prod"),
	  cel.bind(variables.break_glass, resource.labels.?break_glass.orValue("false") == "true",
	   !(variables.break_glass ||
		 resource.containers.all(c, c.startsWith(variables.env + ".")))
	   ? optional.of("only %s containers are allowed in namespace %s".format([variables.env, resource.namespace]))
	   : optional.none()))`,
		},
		{
			name: "nested_rule",
			expr: `
	cel.bind(variables.permitted_regions, ["us", "uk", "es"],
	  cel.bind(variables.banned_regions, {"us": false, "ru": false, "ir": false},
	  (resource.origin in variables.banned_regions &&
		!(resource.origin in variables.permitted_regions))
		? optional.of({"banned": true}) : optional.none()).or(
			optional.of((resource.origin in variables.permitted_regions)
			? {"banned": false} : {"banned": true})))`,
		},
		{
			name: "pb",
			expr: `(spec.single_int32 > 10)
			? optional.of("invalid spec, got single_int32=%d, wanted <= 10".format([spec.single_int32]))
			: optional.none()`,
			envOpts: []cel.EnvOption{
				cel.Types(&proto3pb.TestAllTypes{}),
			},
		},
		{
			name: "required_labels",
			expr: `
	cel.bind(variables.want, spec.labels,
		cel.bind(variables.missing, variables.want.filter(l, !(l in resource.labels)),
		cel.bind(variables.invalid,
			resource.labels.filter(l, l in variables.want &&
				variables.want[l] != resource.labels[l]),
				(variables.missing.size() > 0)
				? optional.of("missing one or more required labels: %s".format([variables.missing]))
				: ((variables.invalid.size() > 0)
				? optional.of("invalid values provided on one or more labels: %s".format([variables.invalid])) : optional.none()))))`,
		},
		{
			name: "restricted_destinations",
			expr: `
	cel.bind(variables.matches_origin_ip,
	  locationCode(origin.ip) == spec.origin,
	  cel.bind(variables.has_nationality, has(request.auth.claims.nationality),
	    cel.bind(variables.matches_nationality,
		  variables.has_nationality && request.auth.claims.nationality == spec.origin,
		  cel.bind(variables.matches_dest_ip,
			locationCode(destination.ip) in spec.restricted_destinations,
			cel.bind(variables.matches_dest_label,
			  resource.labels.location in spec.restricted_destinations,
              cel.bind(variables.matches_dest,
				variables.matches_dest_ip || variables.matches_dest_label,
				(variables.matches_nationality && variables.matches_dest)
				? true
				: ((!variables.has_nationality && variables.matches_origin_ip && variables.matches_dest)
		        ? true : false)))))))`,
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
			}},
	}

	policyErrorTests = []struct {
		name string
		err  string
	}{
		{
			name: "errors",
			err: `ERROR: testdata/errors/policy.yaml:19:19: undeclared reference to 'spec' (in container '')
 |       expression: spec.labels
 | ..................^
ERROR: testdata/errors/policy.yaml:21:50: Syntax error: mismatched input 'resource' expecting ')'
 |       expression: variables.want.filter(l, !(lin resource.labels))
 | .................................................^
ERROR: testdata/errors/policy.yaml:21:66: Syntax error: extraneous input ')' expecting <EOF>
 |       expression: variables.want.filter(l, !(lin resource.labels))
 | .................................................................^
ERROR: testdata/errors/policy.yaml:23:27: Syntax error: mismatched input '2' expecting {'}', ','}
 |       expression: "{1:305 2:569}"
 | ..........................^
ERROR: testdata/errors/policy.yaml:31:75: Syntax error: extraneous input ']' expecting ')'
 |         "missing one or more required labels: %s".format(variables.missing])
 | ..........................................................................^
ERROR: testdata/errors/policy.yaml:34:67: undeclared reference to 'format' (in container '')
 |         "invalid values provided on one or more labels: %s".format([variables.invalid])
 | ..................................................................^`,
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
