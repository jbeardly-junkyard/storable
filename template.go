package mongogen

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"text/template"
)

type TemplateData struct {
	Package string
	Models  []*Model
}

type Template struct {
	template *template.Template
}

func (t *Template) Execute(wr io.Writer, data interface{}) error {
	var buf bytes.Buffer
	err := t.template.Execute(&buf, data)
	if err != nil {
		return err
	}

	return prettyfy(buf.String(), wr)
}

func prettyfy(src string, wr io.Writer) error {
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "irrelevant", src, parser.ParseComments)
	if err != nil {
		return err
	}

	return printer.Fprint(wr, fs, file)
}

func loadTemplateText(filename string) string {
	text, err := Asset(filename)
	if err != nil {
		panic(err)
	}

	return string(text)
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
var model *template.Template = addTemplate(base, "model", "templates/model.tgo")
var query *template.Template = addTemplate(model, "query", "templates/query.tgo")
var resultset *template.Template = addTemplate(model, "resultset", "templates/resultset.tgo")

var Base *Template = &Template{template: base}
