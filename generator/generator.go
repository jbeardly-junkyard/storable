package generator

import (
	"io"
	"os"
	"path"

	"mongogen/template"
)

type Generator struct {
	processor *Processor
}

func NewGenerator(filename string) *Generator {
	return &Generator{
		processor: NewProcessor(filename),
	}
}

func (g *Generator) Generate() error {
	err := g.processor.Process()
	if err != nil {
		return err
	}

	return g.writeFile()
}

func (g *Generator) writeFile() error {
	filename := g.destFilename()
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	return g.runTemplates(file)
}

func (g *Generator) runTemplates(wr io.Writer) error {
	data := g.getTemplateData()
	err := template.Base.Execute(wr, data)

	return err
}

func (g *Generator) destFilename() string {
	dir, file := path.Split(g.processor.Filename())
	newPath := dir + "base" + file

	return newPath
}

func (g *Generator) getTemplateData() template.TemplateData {
	return template.TemplateData{
		Package: g.processor.Package(),
		Models:  g.processor.Models,
	}
}
