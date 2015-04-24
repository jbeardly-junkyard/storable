package mongogen

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"

	_ "golang.org/x/tools/go/gcimporter"
	"golang.org/x/tools/go/types"
)

const BaseDocument = "github.com/maxwellhealth/bongo.DocumentBase"

type Processor struct {
	Path string
}

func NewProcessor(path string) *Processor {
	return &Processor{Path: path}
}

func (p *Processor) Do() (string, []*Model, error) {
	files, err := p.getSourceFiles()
	if err != nil {
		return "", nil, err
	}

	pkg, err := p.parseSourceFiles(files)
	if err != nil {
		return "", nil, err
	}

	return pkg.Name(), p.processPackage(pkg), nil
}

func (p *Processor) getSourceFiles() ([]string, error) {
	pkg, err := build.Default.ImportDir(p.Path, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot process directory %s: %s", p.Path, err)
	}

	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)

	if len(files) == 0 {
		return nil, fmt.Errorf("%s: no buildable Go files", p.Path)
	}

	return joinDirectory(p.Path, files), nil
}

func (p *Processor) parseSourceFiles(filenames []string) (*types.Package, error) {
	var files []*ast.File
	fs := token.NewFileSet()
	for _, filename := range filenames {
		file, err := parser.ParseFile(fs, filename, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parsing package: %s: %s", filename, err)
		}

		files = append(files, file)
	}

	config := types.Config{FakeImportC: true}
	info := &types.Info{}

	return config.Check(p.Path, fs, files, info)
}

func (p *Processor) processPackage(pkg *types.Package) []*Model {
	s := pkg.Scope()
	r := make([]*Model, 0)
	for _, name := range s.Names() {
		obj := s.Lookup(name)
		named, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}

		str, ok := named.Underlying().(*types.Struct)
		if !ok {
			continue
		}

		if m := p.procesStruct(name, str); m != nil {
			r = append(r, m)
		}
	}

	return r
}

func (p *Processor) procesStruct(name string, s *types.Struct) *Model {
	m := NewModel(name)
	isValid := false
	for i := 0; i < s.NumFields(); i++ {
		f := s.Field(i)

		if f.Type().String() == BaseDocument {
			isValid = true
		}

		m.Fields = append(m.Fields, &Field{
			Name: f.Name(),
			Type: f.Type().String(),
			Tag:  reflect.StructTag(s.Tag(i)),
		})

	}

	if isValid {
		return m
	}

	return nil
}

func joinDirectory(directory string, files []string) []string {
	r := make([]string, len(files))
	for i, file := range files {
		r[i] = filepath.Join(directory, file)
	}

	return r
}
