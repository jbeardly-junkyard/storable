package mongogen

import (
	"bytes"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"path/filepath"
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
	filename = filepath.Join(build.Default.GOPATH, "src/github.com/tyba/mongogen", filename)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(f); err != nil {
		panic(err)
	}

	return buf.String()
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
