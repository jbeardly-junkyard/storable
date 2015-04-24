package mongogen

import (
	"io"
	"os"
)

type Generator struct {
	filename string
}

func NewGenerator(filename string) *Generator {
	return &Generator{filename}
}

func (g *Generator) Generate(pkgName string, m []*Model) error {
	return g.writeFile(pkgName, m)
}

func (g *Generator) writeFile(pkgName string, m []*Model) error {
	file, err := os.Create(g.filename)
	if err != nil {
		return err
	}

	defer file.Close()
	return g.runTemplates(pkgName, m, file)
}

func (g *Generator) runTemplates(name string, m []*Model, wr io.Writer) error {
	data := TemplateData{
		Package: name,
		Models:  m,
	}

	return Base.Execute(wr, data)
}
