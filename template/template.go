package template

import (
	"bytes"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	gotemplate "text/template"

    . "github.com/tyba/mongogen/model"
)

type TemplateData struct {
	Package string
	Models  []*Model
}

type Template struct {
	template *gotemplate.Template
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

func makeTemplate(name string, filename string) *gotemplate.Template {
	text := loadTemplateText(filename)
	return gotemplate.Must(gotemplate.New(name).Parse(text))
}

func addTemplate(base *gotemplate.Template, name string, filename string) *gotemplate.Template {
	text := loadTemplateText(filename)
	return gotemplate.Must(base.New(name).Parse(text))
}

var base *gotemplate.Template = makeTemplate("base", "template/code/base.tgo")
var model *gotemplate.Template = addTemplate(base, "model", "template/code/model.tgo")
var query *gotemplate.Template = addTemplate(model, "query", "template/code/query.tgo")

var Base *Template = &Template{template: base}
