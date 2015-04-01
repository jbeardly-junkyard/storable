package generator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"strconv"

	. "github.com/tyba/mongogen/model"
)

type Processor struct {
	filename string
	fileset  *token.FileSet
	file     *ast.File
	Models   []*Model
}

func NewProcessor(filename string) *Processor {
	return &Processor{
		filename: filename,
		fileset:  token.NewFileSet(),
		Models:   make([]*Model, 0),
	}
}

func (p *Processor) Filename() string {
	return p.filename
}

func (p *Processor) Process() error {
	err := p.parseFile()
	if err == nil {
		p.extractModels()
	}

	return err
}

func (p *Processor) Package() string {
	return p.file.Name.Name
}

func (p *Processor) parseFile() error {
	var err error
	p.file, err = parser.ParseFile(
		p.fileset, p.filename, nil, parser.ParseComments|parser.Trace,
	)

	return err
}

func (p *Processor) extractModels() {
	for _, decl := range p.file.Decls {
		ast.Walk(p, decl)
	}
}

func (p *Processor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.GenDecl:
		if node.Tok == token.TYPE {
			p.processTypeDecl(node)
		}
	}

	return p
}

func (p *Processor) processTypeDecl(td *ast.GenDecl) {
	for _, spec := range td.Specs {
		ts := spec.(*ast.TypeSpec)
		if st, ok := ts.Type.(*ast.StructType); ok {
			name := ts.Name.Name
			p.processStruct(name, st)
		}
	}
}

func (p *Processor) processStruct(name string, st *ast.StructType) {

	sp := newStructProcessor(p)
	model := sp.process(name, st)
	if model != nil {
		p.Models = append(p.Models, model)
	}
}

type structProcessor struct {
	processor  *Processor
	collection string
	Model      *Model
}

func newStructProcessor(p *Processor) *structProcessor {
	return &structProcessor{
		processor:  p,
		collection: "",
		Model:      nil,
	}
}

func (sp *structProcessor) process(name string, st *ast.StructType) *Model {
	fields := make([]Field, 0)
	for _, fieldDef := range st.Fields.List {
		more := sp.processFieldDef(fieldDef)
		fields = append(fields, more...)
	}

	if sp.collection != "" {
		return &Model{
			Name:       name,
			Collection: sp.collection,
			Fields:     fields,
		}
	}

	return nil
}

func (sp *structProcessor) processFieldDef(def *ast.Field) []Field {

	typeStr := sp.processor.nodeString(def.Type)
	structTag := sp.extractTag(def)
	if structTag != nil {
		collection := structTag.Get("collection")
		if collection != "" {
			sp.collection = collection
		}
	}

	fields := sp.makeFields(def.Names, typeStr, structTag)

	return fields
}

func (sp *structProcessor) makeFields(
	names []*ast.Ident,
	typeStr string,
	tag *reflect.StructTag) []Field {

	fields := make([]Field, 0)
	for _, ident := range names {
		field := Field{
			Name: sp.processor.nodeString(ident),
			Type: typeStr,
			Tag:  tag,
		}
		fields = append(fields, field)

	}
	return fields
}

func (sp *structProcessor) extractTag(f *ast.Field) *reflect.StructTag {
	tag := f.Tag
	if tag != nil {
		if tagStr, err := strconv.Unquote(tag.Value); err == nil {
			st := reflect.StructTag(tagStr)
			return &st
		}
	}

	return nil
}

func (p *Processor) nodeString(n ast.Node) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, p.fileset, n)
	return buf.String()
}
