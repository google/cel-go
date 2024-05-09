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
	unspecified semanticType = iota
	firstMatch
)

type Policy struct {
	name     ValueString
	rule     *Rule
	semantic semanticType
	info     *ast.SourceInfo
	source   *Source
}

func (p *Policy) Source() *Source {
	return p.source
}

func (p *Policy) SourceInfo() *ast.SourceInfo {
	return p.info
}

func (p *Policy) Name() ValueString {
	return p.name
}

func (p *Policy) Rule() *Rule {
	return p.rule
}

func (p *Policy) SetName(name ValueString) {
	p.name = name
}

func (p *Policy) SetRule(r *Rule) {
	p.rule = r
}

type Rule struct {
	id          *ValueString
	description *ValueString
	variables   []*Variable
	matches     []*Match
}

func (r *Rule) ID() ValueString {
	if r.id != nil {
		return *r.id
	}
	return ValueString{}
}

func (r *Rule) Description() ValueString {
	if r.description != nil {
		return *r.description
	}
	return ValueString{}
}

func (r *Rule) Matches() []*Match {
	return r.matches[:]
}

func (r *Rule) Variables() []*Variable {
	return r.variables[:]
}

func (r *Rule) SetID(id ValueString) {
	r.id = &id
}

func (r *Rule) SetDescription(desc ValueString) {
	r.description = &desc
}

func (r *Rule) AddMatch(m *Match) {
	r.matches = append(r.matches, m)
}

func (r *Rule) AddVariable(v *Variable) {
	r.variables = append(r.variables, v)
}

type Variable struct {
	name       ValueString
	expression ValueString
}

func (v *Variable) Name() ValueString {
	return v.name
}

func (v *Variable) Expression() ValueString {
	return v.expression
}

func (v *Variable) SetName(name ValueString) {
	v.name = name
}

func (v *Variable) SetExpression(e ValueString) {
	v.expression = e
}

type Match struct {
	condition ValueString
	output    *ValueString
	rule      *Rule
}

func (m *Match) Condition() ValueString {
	return m.condition
}

func (m *Match) HasOutput() bool {
	return m.output != nil
}

func (m *Match) Output() ValueString {
	if m.HasOutput() {
		return *m.output
	}
	return ValueString{}
}

func (m *Match) HasRule() bool {
	return m.rule != nil
}

func (m *Match) Rule() *Rule {
	return m.rule
}

func (m *Match) SetCondition(c ValueString) {
	m.condition = c
}

func (m *Match) SetOutput(o ValueString) {
	m.output = &o
}

func (m *Match) SetRule(r *Rule) {
	m.rule = r
}

type ValueString struct {
	ID    int64
	Value string
}

type ParserContext interface {
	NextID() int64
	CollectMetadata(*yaml.Node) int64
	ReportErrorAtID(id int64, msg string, args ...any)

	NewPolicy(*yaml.Node) (*Policy, int64)
	NewRule(*yaml.Node) (*Rule, int64)
	NewVariable(*yaml.Node) (*Variable, int64)
	NewMatch(*yaml.Node) (*Match, int64)
	NewString(*yaml.Node) ValueString
}

type TagVisitor interface {
	PolicyTag(ctx ParserContext, id int64, fieldName string, val *yaml.Node, p *Policy)
	RuleTag(ctx ParserContext, id int64, fieldName string, val *yaml.Node, r *Rule)
	MatchTag(ctx ParserContext, id int64, fieldName string, val *yaml.Node, m *Match)
	VariableTag(ctx ParserContext, id int64, fieldName string, val *yaml.Node, v *Variable)
}

type defaultTagVisitor struct{}

func (defaultTagVisitor) PolicyTag(ctx ParserContext, id int64, fieldName string, node *yaml.Node, p *Policy) {
	ctx.ReportErrorAtID(id, "unsupported policy tag: %s", fieldName)
}

func (defaultTagVisitor) RuleTag(ctx ParserContext, id int64, fieldName string, node *yaml.Node, r *Rule) {
	ctx.ReportErrorAtID(id, "unsupported rule tag: %s", fieldName)
}

func (defaultTagVisitor) MatchTag(ctx ParserContext, id int64, fieldName string, node *yaml.Node, m *Match) {
	ctx.ReportErrorAtID(id, "unsupported match tag: %s", fieldName)
}

func (defaultTagVisitor) VariableTag(ctx ParserContext, id int64, fieldName string, node *yaml.Node, v *Variable) {
	ctx.ReportErrorAtID(id, "unsupported match tag: %s", fieldName)
}

type Parser struct {
	tagVisitor TagVisitor
}

func NewParser() *Parser {
	return &Parser{tagVisitor: defaultTagVisitor{}}
}

// Parse generates an internal parsed policy representation from a YAML input file.
// The internal representation ensures that CEL expressions are tracked relative to
// where they occur within the file, thus making error messages relative to the whole
// file rather than the individual expression.
func (parser *Parser) Parse(src *Source) (*Policy, *cel.Issues) {
	info := ast.NewSourceInfo(src)
	errs := common.NewErrors(src)
	iss := cel.NewIssuesWithSourceInfo(errs, info)
	p := newParserImpl(parser.tagVisitor, info, src, iss)
	policy := p.parseYaml(src)
	if iss.Err() != nil {
		return nil, iss
	}
	return policy, nil
}

func (p *parserImpl) parseYaml(src *Source) *Policy {
	// Parse yaml representation from the source to an object model.
	var docNode yaml.Node
	err := sourceToYaml(src, &docNode)
	if err != nil {
		p.iss.ReportErrorAtID(0, err.Error())
		return nil
	}
	return p.parsePolicy(p, docNode.Content[0])
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

func newParserImpl(visitor TagVisitor, info *ast.SourceInfo, src *Source, iss *cel.Issues) *parserImpl {
	return &parserImpl{
		visitor: visitor,
		info:    info,
		src:     src,
		iss:     iss,
	}
}

type parserImpl struct {
	id      int64
	visitor TagVisitor
	info    *ast.SourceInfo
	src     *Source
	iss     *cel.Issues
}

func (p *parserImpl) NextID() int64 {
	p.id++
	return p.id
}

func (p *parserImpl) NewPolicy(node *yaml.Node) (*Policy, int64) {
	policy := &Policy{}
	policy.source = p.src
	policy.info = p.info
	policy.semantic = firstMatch
	id := p.CollectMetadata(node)
	return policy, id
}

func (p *parserImpl) NewRule(node *yaml.Node) (*Rule, int64) {
	r := &Rule{}
	id := p.CollectMetadata(node)
	return r, id
}

func (p *parserImpl) NewVariable(node *yaml.Node) (*Variable, int64) {
	v := &Variable{}
	id := p.CollectMetadata(node)
	return v, id
}

func (p *parserImpl) NewMatch(node *yaml.Node) (*Match, int64) {
	m := &Match{}
	id := p.CollectMetadata(node)
	return m, id
}

func (p *parserImpl) NewString(node *yaml.Node) ValueString {
	id := p.CollectMetadata(node)
	nodeType := p.assertYamlType(id, node, yamlString, yamlText)
	if nodeType == nil {
		return ValueString{ID: id, Value: "*error*"}
	}
	if *nodeType == yamlText {
		return ValueString{ID: id, Value: node.Value}
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
		return ValueString{ID: id, Value: raw.String()}
	}
	return ValueString{ID: id, Value: node.Value}
}

func (p *parserImpl) CollectMetadata(node *yaml.Node) int64 {
	id := p.NextID()
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
	return id
}

func (p *parserImpl) parsePolicy(ctx ParserContext, node *yaml.Node) *Policy {
	ctx.CollectMetadata(node)
	policy, id := ctx.NewPolicy(node)
	if p.assertYamlType(id, node, yamlMap) == nil {
		return policy
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		p.CollectMetadata(key)
		fieldName := key.Value
		val := node.Content[i+1]
		switch fieldName {
		case "name":
			policy.SetName(ctx.NewString(val))
		case "rule":
			policy.SetRule(p.parseRule(ctx, val))
		default:
			id := ctx.CollectMetadata(val)
			p.visitor.PolicyTag(ctx, id, fieldName, val, policy)
		}
	}
	return policy
}

func (p *parserImpl) parseRule(ctx ParserContext, node *yaml.Node) *Rule {
	r, id := ctx.NewRule(node)
	if p.assertYamlType(id, node, yamlMap) == nil {
		return r
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		ctx.CollectMetadata(key)
		fieldName := key.Value
		val := node.Content[i+1]
		if val.Style == yaml.FoldedStyle || val.Style == yaml.LiteralStyle {
			val.Line++
			val.Column = key.Column + 2
		}
		switch fieldName {
		case "id":
			r.SetID(ctx.NewString(val))
		case "description":
			r.SetDescription(ctx.NewString(val))
		case "variables":
			p.parseVariables(ctx, r, val)
		case "match":
			p.parseMatches(ctx, r, val)
		default:
			id := ctx.CollectMetadata(val)
			p.visitor.RuleTag(ctx, id, fieldName, val, r)
		}
	}
	return r
}

func (p *parserImpl) parseVariables(ctx ParserContext, r *Rule, node *yaml.Node) {
	id := ctx.CollectMetadata(node)
	if p.assertYamlType(id, node, yamlList) == nil {
		return
	}
	for _, val := range node.Content {
		r.AddVariable(p.parseVariable(ctx, val))
	}
}

func (p *parserImpl) parseVariable(ctx ParserContext, node *yaml.Node) *Variable {
	v, id := ctx.NewVariable(node)
	if p.assertYamlType(id, node, yamlMap) == nil {
		return v
	}
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		ctx.CollectMetadata(key)
		fieldName := key.Value
		val := node.Content[i+1]
		if val.Style == yaml.FoldedStyle || val.Style == yaml.LiteralStyle {
			val.Line++
			val.Column = key.Column + 2
		}
		switch fieldName {
		case "name":
			v.SetName(ctx.NewString(val))
		case "expression":
			v.SetExpression(ctx.NewString(val))
		default:
			id := ctx.CollectMetadata(node)
			p.visitor.VariableTag(ctx, id, fieldName, val, v)
		}
	}
	return v
}

func (p *parserImpl) parseMatches(ctx ParserContext, r *Rule, node *yaml.Node) {
	id := ctx.CollectMetadata(node)
	if p.assertYamlType(id, node, yamlList) == nil {
		return
	}
	for _, val := range node.Content {
		r.AddMatch(p.parseMatch(ctx, val))
	}
}

func (p *parserImpl) parseMatch(ctx ParserContext, node *yaml.Node) *Match {
	m, id := ctx.NewMatch(node)
	if p.assertYamlType(id, node, yamlMap) == nil {
		return m
	}
	m.SetCondition(ValueString{ID: ctx.NextID(), Value: "true"})
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i]
		id := ctx.CollectMetadata(key)
		fieldName := key.Value
		val := node.Content[i+1]
		if val.Style == yaml.FoldedStyle || val.Style == yaml.LiteralStyle {
			val.Line++
			val.Column = key.Column + 2
		}
		switch fieldName {
		case "condition":
			if m.HasRule() {
				p.ReportErrorAtID(id, "only the rule or the condition may be set")
			}
			m.SetCondition(ctx.NewString(val))
		case "output":
			m.SetOutput(ctx.NewString(val))
		case "rule":
			m.SetRule(p.parseRule(ctx, val))
		default:
			id := ctx.CollectMetadata(node)
			p.visitor.MatchTag(ctx, id, fieldName, val, m)
		}
	}
	return m
}

func (p *parserImpl) assertYamlType(id int64, node *yaml.Node, nodeTypes ...yamlNodeType) *yamlNodeType {
	nt, found := yamlTypes[node.LongTag()]
	if !found {
		p.ReportErrorAtID(id, "unsupported map key type: %v", node.LongTag())
		return nil
	}
	for _, nodeType := range nodeTypes {
		if nt == nodeType {
			return &nt
		}
	}
	p.ReportErrorAtID(id, "got yaml node type %v, wanted type(s) %v", node.LongTag(), nodeTypes)
	return nil
}

func (p *parserImpl) ReportErrorAtID(id int64, format string, args ...interface{}) {
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
