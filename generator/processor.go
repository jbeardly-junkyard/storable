package generator

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

const BaseDocument = "github.com/tyba/storable.Document"

type Processor struct {
	Path   string
	Ignore map[string]bool
}

func NewProcessor(path string, ignore []string) *Processor {
	i := make(map[string]bool, 0)
	for _, file := range ignore {
		i[file] = true
	}

	return &Processor{Path: path, Ignore: i}
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
		if _, ok := p.Ignore[filename]; ok {
			continue
		}

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
		str := p.tryGetStruct(s.Lookup(name).Type())
		if str == nil {
			continue
		}

		if m := p.processStruct(name, str); m != nil {
			r = append(r, m)
		}
	}

	return r
}

func (p *Processor) tryGetStruct(typ types.Type) *types.Struct {
	switch t := typ.(type) {
	case *types.Named:
		return p.tryGetStruct(t.Underlying())
	case *types.Pointer:
		return p.tryGetStruct(t.Elem())
	case *types.Slice:
		return p.tryGetStruct(t.Elem())
	case *types.Struct:
		return t
	}

	return nil
}

func (p *Processor) processStruct(name string, s *types.Struct) *Model {
	m := NewModel(name)

	var base int
	if base, m.Fields = p.getFields(s); base == -1 {
		return nil
	}

	p.procesBaseField(m, m.Fields[base])
	return m
}

func (p *Processor) getFields(s *types.Struct) (base int, fields []*Field) {
	c := s.NumFields()

	base = -1
	fields = make([]*Field, c)

	for i := 0; i < c; i++ {
		f := s.Field(i)
		t := reflect.StructTag(s.Tag(i))

		if f.Type().String() == BaseDocument {
			base = i
		}

		field := NewField(f.Name(), f.Type().String(), t)

		str := p.tryGetStruct(f.Type())
		if f.Type().String() != BaseDocument && str != nil {
			_, subfields := p.getFields(str)
			field.SetFields(subfields)
			field.Type = "struct"
		}

		fields[i] = field
	}

	return
}

func (p *Processor) procesBaseField(m *Model, f *Field) {
	m.Collection = f.Tag.Get("collection")
}

func joinDirectory(directory string, files []string) []string {
	r := make([]string, len(files))
	for i, file := range files {
		r[i] = filepath.Join(directory, file)
	}

	return r
}
