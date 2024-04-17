package policy

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/ast"
)

type semanticType int

const (
	Unspecified semanticType = iota
	FirstMatch
	LastMatch
	LogicalAnd
	LogicalOr
	Accumulate
)

type policy struct {
	name     policyString
	rule     *rule
	semantic semanticType
	info     *ast.SourceInfo
	source   *Source
}

type rule struct {
	id          *policyString
	description *policyString
	variables   []*variable
	matches     []*match
}

type variable struct {
	name       policyString
	expression policyString
}

type match struct {
	condition policyString
	output    *policyString
	rule      *rule
}

type policyString struct {
	id    int64
	value string
}

func parse(src *Source) (*policy, *cel.Issues) {
	info := ast.NewSourceInfo(src)
	errs := common.NewErrors(src)
	iss := cel.NewIssuesWithSourceInfo(errs, info)
	p := newParser(info, src, iss)
	policy := p.parseYaml(src)
	if iss.Err() != nil {
		return nil, iss
	}
	policy.source = src
	policy.info = p.info
	return policy, nil
}

func (p *parser) parseYaml(src *Source) *policy {
	// Parse yaml representation from the source to an object model.
	var docNode yaml.Node
	err := sourceToYaml(src, &docNode)
	if err != nil {
		p.iss.ReportErrorAtID(0, err.Error())
		return nil
	}
	p.collectMetadata(1, &docNode)
	return p.parse(docNode.Content[0])
}

func sourceToYaml(src *Source, docNode *yaml.Node) error {
	err := yaml.Unmarshal([]byte(src.Content()), docNode)
	if err != nil {
		return err
	}
	if docNode.Kind != yaml.DocumentNode {
		return fmt.Errorf("got yaml node of kind %v, wanted mapping node", docNode.Kind)
	}
	return nil
}

func newParser(info *ast.SourceInfo, src *Source, iss *cel.Issues) *parser {
	return &parser{
		info: info,
		src:  src,
		iss:  iss,
	}
}

type parser struct {
	id   int64
	info *ast.SourceInfo
	src  *Source
	iss  *cel.Issues
}

func (p *parser) nextID() int64 {
	p.id++
	return p.id
}

func (p *parser) parse(node *yaml.Node) *policy {
	pol := &policy{}
	id := p.nextID()
	p.collectMetadata(id, node)
	if p.assertYamlType(id, node, yamlMap) == nil {
		return pol
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		id := p.nextID()
		p.collectMetadata(id, key)
		fieldName := key.Value
		val := node.Content[i+1]
		if val.Style == yaml.FoldedStyle || val.Style == yaml.LiteralStyle {
			val.Line++
			val.Column = key.Column + 2
		}
		switch fieldName {
		case "name":
			pol.name = p.parseString(val)
		case "rule":
			rule := p.parseRule(val)
			pol.rule = rule
		default:
			p.reportErrorAtID(id, "unexpected field name: %s", fieldName)
		}
	}
	return pol
}

func (p *parser) parseRule(node *yaml.Node) *rule {
	r := &rule{}
	id := p.nextID()
	p.collectMetadata(id, node)
	if p.assertYamlType(id, node, yamlMap) == nil {
		return r
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		id := p.nextID()
		p.collectMetadata(id, key)
		fieldName := key.Value
		val := node.Content[i+1]
		if val.Style == yaml.FoldedStyle || val.Style == yaml.LiteralStyle {
			val.Line++
			val.Column = key.Column + 2
		}
		switch fieldName {
		case "id":
			ruleID := p.parseString(val)
			r.id = &ruleID
		case "description":
			desc := p.parseString(val)
			r.description = &desc
		case "variables":
			r.variables = p.parseVariables(val)
		case "match":
			r.matches = p.parseMatches(val)
		default:
			p.reportErrorAtID(id, "unexpected field name: %s", fieldName)
		}
	}
	return r
}

func (p *parser) parseVariables(node *yaml.Node) []*variable {
	vars := []*variable{}
	id := p.nextID()
	p.collectMetadata(id, node)
	if p.assertYamlType(id, node, yamlList) == nil {
		return vars
	}
	for _, val := range node.Content {
		vars = append(vars, p.parseVariable(val))
	}
	return vars
}

func (p *parser) parseVariable(node *yaml.Node) *variable {
	v := &variable{}
	id := p.nextID()
	p.collectMetadata(id, node)
	if p.assertYamlType(id, node, yamlMap) == nil {
		return v
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		id := p.nextID()
		p.collectMetadata(id, key)
		fieldName := key.Value
		val := node.Content[i+1]
		if val.Style == yaml.FoldedStyle || val.Style == yaml.LiteralStyle {
			val.Line++
			val.Column = key.Column + 2
		}
		switch fieldName {
		case "name":
			v.name = p.parseString(val)
		case "expression":
			v.expression = p.parseString(val)
		default:
			p.reportErrorAtID(id, "unexpected field name: %s", fieldName)
		}
	}
	return v
}

func (p *parser) parseMatches(node *yaml.Node) []*match {
	matches := []*match{}
	id := p.nextID()
	p.collectMetadata(id, node)
	if p.assertYamlType(id, node, yamlList) == nil {
		return matches
	}
	for _, val := range node.Content {
		matches = append(matches, p.parseMatch(val))
	}
	return matches
}

func (p *parser) parseMatch(node *yaml.Node) *match {
	m := &match{}
	id := p.nextID()
	p.collectMetadata(id, node)
	if p.assertYamlType(id, node, yamlMap) == nil {
		return m
	}
	m.condition = policyString{id: id, value: "true"}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		id := p.nextID()
		p.collectMetadata(id, key)
		fieldName := key.Value
		val := node.Content[i+1]
		if val.Style == yaml.FoldedStyle || val.Style == yaml.LiteralStyle {
			val.Line++
			val.Column = key.Column + 2
		}
		switch fieldName {
		case "condition":
			m.condition = p.parseString(val)
		case "output":
			outputExpr := p.parseString(val)
			m.output = &outputExpr
		case "rule":
			m.rule = p.parseRule(val)
		default:
			p.reportErrorAtID(id, "unexpected field name: %s", fieldName)
		}
	}
	return m
}

func (p *parser) parseString(node *yaml.Node) policyString {
	id := p.nextID()
	p.collectMetadata(id, node)
	nodeType := p.assertYamlType(id, node, yamlString, yamlText)
	if nodeType == nil {
		return policyString{id: id, value: "*error*"}
	}
	if *nodeType == yamlText {
		return policyString{id: id, value: node.Value}
	}
	if node.Style == yaml.FoldedStyle || node.Style == yaml.LiteralStyle {
		col := node.Column
		line := node.Line
		txt, found := p.src.Snippet(line)
		indent := ""
		for len(indent) < col-1 {
			indent += " "
		}
		var raw strings.Builder
		for found && strings.HasPrefix(txt, indent) {
			line++
			raw.WriteString(txt)
			txt, found = p.src.Snippet(line)
			if found && strings.HasPrefix(txt, indent) {
				raw.WriteString("\n")
			}
		}
		offset := p.info.OffsetRanges()[p.id]
		offsetStart := offset.Start - (int32(node.Column) - 1)
		p.info.SetOffsetRange(p.id, ast.OffsetRange{Start: offsetStart, Stop: offsetStart})
		return policyString{id: id, value: raw.String()}
	}
	return policyString{id: id, value: node.Value}
}

func (p *parser) collectMetadata(id int64, node *yaml.Node) {
	line := node.Line
	col := int32(node.Column)
	switch node.Style {
	case yaml.DoubleQuotedStyle, yaml.SingleQuotedStyle:
		col++
	}
	offsetStart := int32(0)
	if line > 1 {
		offsetStart = p.info.LineOffsets()[line-2]
	}
	p.info.SetOffsetRange(id, ast.OffsetRange{Start: offsetStart + col - 1, Stop: offsetStart + col - 1})
}

func (p *parser) assertYamlType(id int64, node *yaml.Node, nodeTypes ...yamlNodeType) *yamlNodeType {
	nt, found := yamlTypes[node.LongTag()]
	if !found {
		p.reportErrorAtID(id, "unsupported map key type: %v", node.LongTag())
		return nil
	}
	for _, nodeType := range nodeTypes {
		if nt == nodeType {
			return &nt
		}
	}
	p.reportErrorAtID(id, "got yaml node type %v, wanted type(s) %v", node.LongTag(), nodeTypes)
	return nil
}

func (p *parser) reportErrorAtID(id int64, format string, args ...interface{}) {
	p.iss.ReportErrorAtID(id, format, args...)
}

type yamlNodeType int

const (
	yamlText yamlNodeType = iota + 1
	yamlBool
	yamlNull
	yamlString
	yamlInt
	yamlDouble
	yamlList
	yamlMap
	yamlTimestamp
)

var (
	// yamlTypes map of the long tag names supported by the Go YAML v3 library.
	yamlTypes = map[string]yamlNodeType{
		"!txt":                        yamlText,
		"tag:yaml.org,2002:bool":      yamlBool,
		"tag:yaml.org,2002:null":      yamlNull,
		"tag:yaml.org,2002:str":       yamlString,
		"tag:yaml.org,2002:int":       yamlInt,
		"tag:yaml.org,2002:float":     yamlDouble,
		"tag:yaml.org,2002:seq":       yamlList,
		"tag:yaml.org,2002:map":       yamlMap,
		"tag:yaml.org,2002:timestamp": yamlTimestamp,
	}
)
