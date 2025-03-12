package compiler

import (
	"reflect"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/env"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/policy"
	"gopkg.in/yaml.v3"
)

func TestEnvironmentFileCompareTextprotoAndYAML(t *testing.T) {
	t.Run("compare textproto and yaml environment files", func(t *testing.T) {
		protoConfig, err := parseEnv(t, "proto_config", "testdata/config.textproto")
		if err != nil {
			t.Fatalf("parseProtoEnv() failed: %v", err)
		}
		config, err := parseEnv(t, "yaml_config", "testdata/config.yaml")
		if err != nil {
			t.Fatalf("parseYamlEnv() failed: %v", err)
		}
		if protoConfig.Container != config.Container {
			t.Fatalf("Container got %q, wanted %q", protoConfig.Container, config.Container)
		}
		if !reflect.DeepEqual(protoConfig.Imports, config.Imports) {
			t.Fatalf("Imports got %v, wanted %v", protoConfig.Imports, config.Imports)
		}
		if !reflect.DeepEqual(protoConfig.StdLib, config.StdLib) {
			t.Fatalf("StdLib got %v, wanted %v", protoConfig.StdLib, config.StdLib)
		}
		if len(protoConfig.Extensions) != len(config.Extensions) {
			t.Fatalf("Extensions count got %d, wanted %d", len(protoConfig.Extensions), len(config.Extensions))
		}
		for _, protoConfigExt := range protoConfig.Extensions {
			found := false
			for _, configExt := range config.Extensions {
				if reflect.DeepEqual(protoConfigExt, configExt) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("Extensions got %v, wanted %v", protoConfig.Extensions, config.Extensions)
			}
		}
		if !reflect.DeepEqual(protoConfig.ContextVariable, config.ContextVariable) {
			t.Fatalf("ContextVariable got %v, wanted %v", protoConfig.ContextVariable, config.ContextVariable)
		}
		if len(protoConfig.Variables) != len(config.Variables) {
			t.Fatalf("Variables count got %d, wanted %d", len(protoConfig.Variables), len(config.Variables))
		} else {
			for i, v := range protoConfig.Variables {
				for j, p := range v.TypeDesc.Params {
					if p.TypeName == "google.protobuf.Any" &&
						config.Variables[i].TypeDesc.Params[j].TypeName == "dyn" {
						p.TypeName = "dyn"
					}
				}
				if !reflect.DeepEqual(v, config.Variables[i]) {
					t.Fatalf("Variables[%d] not equal, got %v, wanted %v", i, v, config.Variables[i])
				}
			}
		}
		if len(protoConfig.Functions) != len(config.Functions) {
			t.Fatalf("Functions count got %d, wanted %d", len(protoConfig.Functions), len(config.Functions))
		} else {
			for i, f := range protoConfig.Functions {
				if !reflect.DeepEqual(f, config.Functions[i]) {
					t.Fatalf("Functions[%d] not equal, got %v, wanted %v", i, f, config.Functions[i])
				}
			}
		}
		if len(protoConfig.Features) != len(config.Features) {
			t.Fatalf("Features count got %d, wanted %d", len(protoConfig.Features), len(config.Features))
		} else {
			for i, f := range protoConfig.Features {
				if !reflect.DeepEqual(f, config.Features[i]) {
					t.Fatalf("Features[%d] not equal, got %v, wanted %v", i, f, config.Features[i])
				}
			}
		}
		if len(protoConfig.Validators) != len(config.Validators) {
			t.Fatalf("Validators count got %d, wanted %d", len(protoConfig.Validators), len(config.Validators))
		} else {
			for i, v := range protoConfig.Validators {
				if !reflect.DeepEqual(v, config.Validators[i]) {
					t.Fatalf("Validators[%d] not equal, got %v, wanted %v", i, v, config.Validators[i])
				}
			}
		}
	})
}

func parseEnv(t *testing.T, name, path string) (*env.Config, error) {
	t.Helper()
	opts := EnvironmentFile(path)
	e, err := cel.NewCustomEnv(opts)
	if err != nil {
		return nil, err
	}
	conf, err := e.ToConfig(name)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func TestFileExpressionCustomPolicyParser(t *testing.T) {
	t.Run("test file expression custom policy parser", func(t *testing.T) {
		envOpt := EnvironmentFile("../../policy/testdata/k8s/config.yaml")
		parserOpt := policy.ParserOption(func(p *policy.Parser) (*policy.Parser, error) {
			p.TagVisitor = policy.K8sTestTagHandler()
			return p, nil
		})
		compilerOpts := []any{envOpt, parserOpt}
		compiler, err := NewCompiler(compilerOpts...)
		if err != nil {
			t.Fatalf("NewCompiler() failed: %v", err)
		}
		policyFile := &FileExpression{
			Path: "../../policy/testdata/k8s/policy.yaml",
		}
		k8sAst, _, err := policyFile.CreateAst(compiler)
		if err != nil {
			t.Fatalf("CreateAst() failed: %v", err)
		}
		if k8sAst == nil {
			t.Fatalf("CreateAst() returned nil ast")
		}
	})
}

func TestFileExpressionPolicyMetadataOptions(t *testing.T) {
	t.Run("test file expression policy metadata options", func(t *testing.T) {
		envOpt := EnvironmentFile("testdata/custom_policy_config.yaml")
		parserOpt := policy.ParserOption(func(p *policy.Parser) (*policy.Parser, error) {
			p.TagVisitor = customTagHandler{TagVisitor: policy.DefaultTagVisitor()}
			return p, nil
		})
		policyMetadataOpt := PolicyMetadataEnvOption(ParsePolicyVariables)
		compilerOpts := []any{envOpt, parserOpt, policyMetadataOpt}
		compiler, err := NewCompiler(compilerOpts...)
		if err != nil {
			t.Fatalf("NewCompiler() failed: %v", err)
		}
		policyFile := &FileExpression{
			Path: "testdata/custom_policy.celpolicy",
		}
		ast, _, err := policyFile.CreateAst(compiler)
		if err != nil {
			t.Fatalf("CreateAst() failed: %v", err)
		}
		if ast == nil {
			t.Fatalf("CreateAst() returned nil ast")
		}
	})
}

func ParsePolicyVariables(metadata map[string]any) cel.EnvOption {
	variables := []*decls.VariableDecl{}
	for n, t := range metadata {
		variables = append(variables, decls.NewVariable(n, parseCustomPolicyVariableType(t.(string))))
	}
	return cel.VariableDecls(variables...)
}

func parseCustomPolicyVariableType(t string) *types.Type {
	switch t {
	case "int":
		return types.IntType
	case "string":
		return types.StringType
	default:
		return types.UnknownType
	}
}

type variableType struct {
	VariableName string `yaml:"variable_name"`
	VariableType string `yaml:"variable_type"`
}

type customTagHandler struct {
	policy.TagVisitor
}

func (customTagHandler) PolicyTag(ctx policy.ParserContext, id int64, tagName string, node *yaml.Node, p *policy.Policy) {
	switch tagName {
	case "variable_types":
		varList := []*variableType{}
		if err := node.Decode(&varList); err != nil {
			ctx.ReportErrorAtID(id, "invalid yaml variable_types node: %v, error: %w", node, err)
			return
		}
		for _, v := range varList {
			p.SetMetadata(v.VariableName, v.VariableType)
		}
	default:
		ctx.ReportErrorAtID(id, "unsupported policy tag: %s", tagName)
	}
}

func TestRawExpressionCreateAst(t *testing.T) {
	t.Run("test raw expression create ast", func(t *testing.T) {
		envOpt := EnvironmentFile("testdata/config.yaml")
		compiler, err := NewCompiler(envOpt)
		if err != nil {
			t.Fatalf("NewCompiler() failed: %v", err)
		}
		rawExpr := &RawExpression{
			Value: "locationCode(destination.ip)==locationCode(origin.ip)",
		}
		ast, _, err := rawExpr.CreateAst(compiler)
		if err != nil {
			t.Fatalf("CreateAst() failed: %v", err)
		}
		if ast == nil {
			t.Fatalf("CreateAst() returned nil ast")
		}
	})
}
