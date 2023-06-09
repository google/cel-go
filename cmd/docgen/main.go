package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/google/cel-go/checker"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var tmpl = *template.Must(template.New("standard_definitions").Parse(
	`These descriptions are automatically generated based on the cel-go implementation.

<table style="width=100%" border="1">
	<col width="15%">
	<col width="40%">
	<col width="45%">
	<tr>
		<th>Symbol</th>
		<th>Type</th>
		<th>Description</th>
	</tr>
	{{- range $k,$func := . -}}
	{{- range $i, $ol := $func.Overloads}}
	<tr>
		{{- if not $i}}
		<th rowspan="{{len $func.Overloads}}">
			{{ $func.Symbol }}
		</th>
		{{- end}}
		<td>
			{{ $ol.Type }}
		</td>
		<td>
			{{ $ol.Description }}
		</td>
	</tr>
	{{- end}}
	{{- end}}
</table>
`))

type SortableDecls []*expr.Decl

func (s SortableDecls) Len() int { return len(s) }

func (s SortableDecls) Less(i, j int) bool {
	switch s[i].DeclKind.(type) {
	case *expr.Decl_Ident:
		if _, ok := s[j].DeclKind.(*expr.Decl_Function); ok {
			return true
		}
	case *expr.Decl_Function:
		if _, ok := s[j].DeclKind.(*expr.Decl_Ident); ok {
			return false
		}
	}
	return s[i].Name < s[j].Name
}

func (s SortableDecls) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type Function struct {
	Symbol    string
	Overloads []*Overload
}

type Overload struct {
	Type        string
	Description string
}

func TypeToText(exprType *expr.Type) string {
	switch kind := exprType.TypeKind.(type) {
	case *expr.Type_Primitive:
		switch kind.Primitive {
		case expr.Type_INT64:
			return "int"
		case expr.Type_UINT64:
			return "uint"
		default:
			return strings.ToLower(kind.Primitive.String())
		}
	case *expr.Type_WellKnown:
		switch expr.Type_WellKnownType(kind.WellKnown.Number()) {
		case expr.Type_DURATION:
			return "google.protobuf.Duration"
		case expr.Type_TIMESTAMP:
			return "google.protobuf.Timestamp"
		}
	case *expr.Type_MessageType:
		return kind.MessageType
	case *expr.Type_Null:
		return "null"
	case *expr.Type_Type:
		if t := TypeToText(kind.Type); t != "type()" {
			return fmt.Sprintf("type(%s)", TypeToText(kind.Type))
		}
		return "type(dyn)"
	case *expr.Type_TypeParam:
		return kind.TypeParam
	case *expr.Type_MapType_:
		return fmt.Sprintf("map(%s, %s)", TypeToText(kind.MapType.KeyType), TypeToText(kind.MapType.ValueType))
	case *expr.Type_ListType_:
		return fmt.Sprintf("list(%s)", TypeToText(kind.ListType.ElemType))
	case *expr.Type_Error:
		return "error"
	case *expr.Type_Dyn:
		return "dyn"
	}
	return ""
}

func main() {
	list := SortableDecls(checker.StandardDeclarations())
	sort.Sort(list)
	functions := map[string]*Function{}
	for _, decl := range list {
		if strings.HasPrefix(decl.Name, "@") {
			continue
		}
		switch t := decl.DeclKind.(type) {
		case *expr.Decl_Ident:
			typeDesc := TypeToText(t.Ident.Type)
			doc := "type denotation"
			if fn, ok := functions[decl.Name]; ok {
				fn.Overloads = append(fn.Overloads, &Overload{
					Type:        typeDesc,
					Description: doc,
				})
			} else {
				functions[decl.Name] = &Function{
					Symbol: decl.Name,
					Overloads: []*Overload{{
						Type:        typeDesc,
						Description: doc,
					}},
				}
			}
		case *expr.Decl_Function:
			for _, ol := range t.Function.Overloads {
				var in []string
				for _, p := range ol.Params {
					in = append(in, TypeToText(p))
				}
				typeDesc := ""
				if ol.IsInstanceFunction {
					typeDesc = in[0] + "."
					in = in[1:]
				}
				typeDesc += "(" + strings.Join(in, ", ") + ") -> " + TypeToText(ol.ResultType)
				doc := "TODO"
				if ol.Doc != "" {
					doc = ol.Doc
				}
				if fn, ok := functions[decl.Name]; ok {
					fn.Overloads = append(fn.Overloads, &Overload{
						Type:        typeDesc,
						Description: doc,
					})
				} else {
					functions[decl.Name] = &Function{
						Symbol: decl.Name,
						Overloads: []*Overload{{
							Type:        typeDesc,
							Description: doc,
						}},
					}
				}
			}
		}
	}
	buffer := bytes.NewBufferString("")
	if err := tmpl.Execute(buffer, functions); err != nil {
		panic(err)
	}
	fmt.Println(buffer.String())
}
