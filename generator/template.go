package generator

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"golang.org/x/tools/go/types"
)

type Template struct {
	template *template.Template
}

type TemplateData struct {
	*Package
	Fields    []*TemplateField
	Processed map[interface{}]string
}

type TemplateField struct {
	Name   string
	Path   string
	Fields interface{}
}

func (tf *TemplateField) ValidFields() []*Field {
	return tf.Fields.([]*Field)
}

func (t *Template) Execute(wr io.Writer, data *Package) error {
	var buf bytes.Buffer

	td := &TemplateData{data, []*TemplateField{}, map[interface{}]string{}}
	err := t.template.Execute(&buf, td)
	if err != nil {
		return err
	}

	return prettyfy(buf.Bytes(), wr)
}

func (td *TemplateData) GenType(vi interface{}, path string) string {
	v := reflect.ValueOf(vi)
	sv := v
	if v.Kind() == reflect.Ptr {
		sv = v.Elem()
	}
	if sv.FieldByName("Type").Interface().(string) == "struct" {
		if v.MethodByName("ValidFields").IsValid() {
			return td.LinkStruct(path, vi)
		}
		return ""
	} else {
		k := "Field"
		if v.MethodByName("ContainsMap").Call(nil)[0].Interface().(bool) {
			k = "Map"
		}
		return fmt.Sprintf("%v storable.%v", sv.FieldByName("Name"), k)
	}
}

func (td *TemplateData) LinkStruct(path string, vi interface{}) string {
	v := reflect.ValueOf(vi)
	name := v.Elem().FieldByName("Name").Interface().(string)
	schemaName := "schema" + path + name

	if proc, ok := td.Processed[vi]; ok {
		schemaName = proc
		return name + " *" + schemaName
	}
	td.Processed[vi] = schemaName

	td.Fields = append(td.Fields, &TemplateField{
		Name:   schemaName,
		Path:   path + name,
		Fields: v.MethodByName("ValidFields").Call(nil)[0].Interface(),
	})

	return name + " *" + schemaName
}

func (td *TemplateData) GenVar(vi interface{}, done map[interface{}]bool) string {
	if done == nil {
		done = map[interface{}]bool{}
	}

	v := reflect.ValueOf(vi)
	sv := v
	if v.Kind() == reflect.Ptr {
		sv = v.Elem()
	}

	if done[vi] {
		return sv.FieldByName("Name").Interface().(string) + ": nil,"
	}

	if sv.FieldByName("Type").Interface().(string) == "struct" {
		if v.MethodByName("ValidFields").IsValid() {
			return td.StructValue(vi, done)
		}
		return ""
	} else {
		k := "NewField"
		if v.MethodByName("ContainsMap").Call(nil)[0].Interface().(bool) {
			k = "NewMap"
		}

		return fmt.Sprintf(
			`%v: storable.%v("%v", "%v"),`,
			sv.FieldByName("Name"),
			k,
			v.MethodByName("GetPath").Call(nil)[0],
			v.MethodByName("FindableType").Call(nil)[0],
		)
	}
}

func (td *TemplateData) StructValue(vi interface{}, done map[interface{}]bool) string {
	v := reflect.ValueOf(vi)
	name := v.Elem().FieldByName("Name").Interface().(string)

	ifc := v.Interface()
	if done[ifc] {
		return name + ": nil,"
	}
	done[ifc] = true

	ret := name + ": &" + td.Processed[vi] + "{"
	for _, v := range v.MethodByName("ValidFields").Call(nil)[0].Interface().([]*Field) {
		ret += "\n" + td.GenVar(v, done)
	}
	ret += "\n},"

	return ret
}

func (td *TemplateData) CallHooks(whenStr, actionStr string, model *Model) string {
	before := whenStr == "before"
	actions := map[HookAction]bool{}
	switch actionStr {
	case "insert":
		actions[InsertHook] = true
		actions[SaveHook] = true
	case "update":
		actions[UpdateHook] = true
		actions[SaveHook] = true
	}

	return callHooksGenerator{before, actions}.do(model, "")
}

type callHooksGenerator struct {
	before  bool
	actions map[HookAction]bool
}

type callHooksNode struct {
	v        interface{}
	typ      types.Type
	elemTyp  types.Type
	name     string
	hooks    []Hook
	children []*callHooksNode
	loop     *callHooksNode
}

func (g callHooksGenerator) do(model *Model, prefix string) string {
	modelHooks := g.generateTree(g.makeTree(model, "", nil), prefix)
	storeHooks := g.generateStoreHooks(model)
	return modelHooks + storeHooks
}

func (g callHooksGenerator) makeTree(v interface{}, name string, stack []*callHooksNode) *callHooksNode {
	obj := reflect.Indirect(reflect.ValueOf(v))
	var typ types.Type
	switch v := obj.FieldByName("CheckedNode").Interface().(type) {
	case *types.Var:
		typ = v.Type()
	case types.Type:
		typ = v
	default:
		panic("unexpected")
	}
	var elemTyp types.Type
	switch v := typ.(type) {
	case *types.Pointer:
		elemTyp = v.Elem()
	case *types.Slice:
		elemTyp = v.Elem()
	case *types.Map:
		elemTyp = v.Elem()
	default:
		elemTyp = v
	}

	ret := &callHooksNode{
		v:       v,
		typ:     typ,
		elemTyp: elemTyp,
		name:    name,
		hooks:   obj.FieldByName("Hooks").Interface().([]Hook),
	}

	for _, elem := range stack {
		if elem.elemTyp == elemTyp {
			ret.loop = elem
			return ret
		}
	}

	stack = append(stack, ret)

	fields := obj.FieldByName("Fields").Interface().([]*Field)
	for _, f := range fields {
		ret.children = append(ret.children, g.makeTree(f, f.Name, stack))
	}

	return ret
}

func (g callHooksGenerator) generateTree(n *callHooksNode, prefix string) string {
	if n.loop != nil {
		return fmt.Sprintf("// Loop: %v.%v %v\n", prefix, n.name, n.typ)
	}

	ret := ""

	wrap := func(s string) string { return s }
	varName := func(s string) string { return s }
	typ := n.typ
out:
	for i := 0; ; i++ {
		idx := fmt.Sprintf("k%d", i)
		prevWrap, prevVarName := wrap, varName
		switch vtyp := typ.(type) {
		case *types.Pointer:
			wrap = func(s string) string { return prevWrap("if doc" + prevVarName(prefix) + " != nil {\n" + s + "\n}\n") }
			typ = vtyp.Elem()
		case *types.Slice, *types.Map:
			wrap = func(s string) string {
				return prevWrap(fmt.Sprintf("for %v, _ := range doc%v {\n%v\n}\n", idx, prevVarName(prefix), s))
			}
			varName = func(s string) string { return prevVarName(s) + "[" + idx + "]" }
			typ = vtyp.(interface {
				Elem() types.Type
			}).Elem()
		default:
			break out
		}
	}

	for _, hook := range n.hooks {
		if hook.Before == g.before && g.actions[hook.Action] {
			ret += g.generateCall(varName(prefix), hook.MethodName())
		}
	}

	for _, cn := range n.children {
		ret += g.generateTree(cn, varName(prefix)+"."+cn.name)
	}

	if len(ret) > 0 {
		ret = wrap(ret)
	}

	return ret
}

func (g callHooksGenerator) generateCall(sel string, method string) string {
	return fmt.Sprintf(
		`if err := doc%s.%s(); err != nil {
		return storable.HookError{
			Hook: "%[2]s",
			Field: "%[1]s",
			Cause: err,
		}
	}
	`, sel, method)
}

func (g callHooksGenerator) generateStoreHooks(model *Model) string {
	ret := ""
	for _, hook := range model.StoreHooks {
		if hook.Before == g.before && g.actions[hook.Action] {
			ret += fmt.Sprintf(
				`if err := doc.%s(s); err != nil {
			return storable.HookError{
				Hook: "%[1]s",
				Field: ".",
				Cause: err,
			}
		}
		`, hook.MethodName())
		}
	}
	return ret
}

func prettyfy(input []byte, wr io.Writer) error {
	output, err := format.Source(input)
	if err != nil {
		printDocumentWithNumbers(string(input))
		return err
	}

	_, err = wr.Write(output)
	return err
}

func printDocumentWithNumbers(code string) {
	for i, line := range strings.Split(code, "\n") {
		fmt.Printf("%.3d %s\n", i+1, line)
	}
}

func loadTemplateText(filename string) string {
	filename = filepath.Join(build.Default.GOPATH, "src/github.com/tyba/storable/generator", filename)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(f); err != nil {
		panic(err)
	}

	return strings.Replace(buf.String(), "\\\n", " ", -1)
}

func makeTemplate(name string, filename string) *template.Template {
	text := loadTemplateText(filename)
	return template.Must(template.New(name).Parse(text))
}

func addTemplate(base *template.Template, name string, filename string) *template.Template {
	text := loadTemplateText(filename)
	return template.Must(base.New(name).Parse(text))
}

var base *template.Template = makeTemplate("base", "templates/base.tgo")
var schema *template.Template = addTemplate(base, "schema", "templates/schema.tgo")
var model *template.Template = addTemplate(base, "model", "templates/model.tgo")
var query *template.Template = addTemplate(model, "query", "templates/query.tgo")
var resultset *template.Template = addTemplate(model, "resultset", "templates/resultset.tgo")

var Base *Template = &Template{template: base}
