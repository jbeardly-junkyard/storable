package generator

import (
	"bytes"
	"fmt"
	"go/build"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Template struct {
	template *template.Template
}

func (t *Template) Execute(wr io.Writer, data interface{}) error {
	var buf bytes.Buffer

	err := t.template.Execute(&buf, data)
	if err != nil {
		return err
	}

	return prettyfy(buf.Bytes(), wr)
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
