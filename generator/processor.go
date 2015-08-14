package generator

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	_ "golang.org/x/tools/go/gcimporter"
	"golang.org/x/tools/go/types"
)

const BaseDocument = "github.com/tyba/storable.Document"

type Processor struct {
	Path       string
	Ignore     map[string]bool
	TypesPkg   *types.Package
	SourceCode map[string][]byte
}

func NewProcessor(path string, ignore []string) *Processor {
	i := make(map[string]bool, 0)
	for _, file := range ignore {
		i[file] = true
	}

	return &Processor{
		Path:   path,
		Ignore: i,
	}
}

func (p *Processor) Do() (*Package, error) {
	files, err := p.getSourceFiles()
	if err != nil {
		return nil, err
	}

	p.SourceCode, err = p.readSourceFiles(files)
	if err != nil {
		return nil, err
	}

	p.TypesPkg, _ = p.parseSourceFiles(files)
	return p.ProcessTypesPkg()
}

func (p *Processor) ProcessTypesPkg() (*Package, error) {
	pkg := &Package{Name: p.TypesPkg.Name()}
	p.processPackage(pkg)

	return pkg, nil
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

	return joinDirectory(p.Path, p.removeIngoredFiles(files)), nil
}

func (p *Processor) removeIngoredFiles(filenames []string) []string {
	var output []string
	for _, filename := range filenames {
		if _, ok := p.Ignore[filename]; ok {
			continue
		}

		output = append(output, filename)
	}

	return output
}

func (p *Processor) readSourceFiles(filenames []string) (map[string][]byte, error) {
	source := make(map[string][]byte, 0)
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return source, err
		}

		source[filename] = b
	}

	return source, nil
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

	config := types.Config{FakeImportC: true, Error: func(error) {}}
	info := &types.Info{}

	return config.Check(p.Path, fs, files, info)
}

func (p *Processor) processPackage(pkg *Package) {
	var newFuncs []*types.Func

	s := p.TypesPkg.Scope()
	for _, name := range s.Names() {
		fun := p.tryGetFunction(s.Lookup(name))
		if fun != nil {
			pkg.Functions = append(pkg.Functions, name)
			if strings.HasPrefix(fun.Name(), "new") {
				newFuncs = append(newFuncs, fun)
			}
		}

		str := p.tryGetStruct(s.Lookup(name).Type())
		if str == nil {
			continue
		}

		if m := p.processStruct(name, str); m != nil {
			pkg.Models = append(pkg.Models, m)
			m.CheckedNode = s.Lookup(name).Type().(*types.Named)
			m.Package = p.TypesPkg
		} else {
			pkg.Structs = append(pkg.Structs, name)
		}
	}

	for _, fun := range newFuncs {
		p.tryMatchNewFunc(pkg.Models, fun)
	}
}

func (p *Processor) tryMatchNewFunc(models []*Model, fun *types.Func) {
	modelName := fun.Name()[len("new"):]

	for _, m := range models {
		if m.Name != modelName {
			continue
		}

		sig := fun.Type().(*types.Signature)

		if sig.Recv() != nil {
			continue
		}

		res := sig.Results()
		for i := 0; i < res.Len(); i++ {
			if isTypeOrPtrTo(res.At(i).Type(), m.CheckedNode) {
				m.NewFunc = fun
				return
			}
		}
	}
}

func (p *Processor) tryGetFunction(typ types.Object) *types.Func {
	switch t := typ.(type) {
	case *types.Func:
		return t
	}

	return nil
}

func (p *Processor) tryGetStruct(typ types.Type) *types.Struct {
	switch t := typ.(type) {
	case *types.Named:
		return p.tryGetStruct(t.Underlying())
	case *types.Pointer:
		return p.tryGetStruct(t.Elem())
	case *types.Slice:
		return p.tryGetStruct(t.Elem())
	case *types.Map:
		return p.tryGetStruct(t.Elem())
	case *types.Struct:
		return t
	}

	return nil
}

func (p *Processor) processStruct(name string, s *types.Struct) *Model {
	m := NewModel(name)
	m.Events = p.getEvents(name)

	var base int
	if base, m.Fields = p.getFields(s); base == -1 {
		return nil
	}

	p.procesBaseField(m, m.Fields[base])

	return m
}

func (p *Processor) getFields(s *types.Struct) (base int, fields []*Field) {
	base, fields = p.processFields(s, []*types.Struct{})
	return
}

func (p *Processor) getEvents(name string) []Event {
	events := []Event{}

	all := []Event{BeforeInsert, AfterInsert, BeforeUpdate, AfterUpdate}
	for _, e := range all {
		if p.isEventPresent(name, e) {
			events = append(events, e)
		}
	}

	return events
}

func (p *Processor) isEventPresent(name string, e Event) bool {
	re := regexp.MustCompile(fmt.Sprintf("\\*%sStore\\) %s\\(", name, e))

	for _, code := range p.SourceCode {
		if re.Match(code) {
			return true
		}
	}

	return false
}

// Returns which field index is an embedded storable.Document, or -1 if none.
func (p *Processor) processFields(s *types.Struct, done []*types.Struct) (base int, fields []*Field) {
	c := s.NumFields()

	base = -1
	fields = make([]*Field, 0)

	for i := 0; i < c; i++ {
		f := s.Field(i)
		if !f.Exported() {
			continue
		}

		t := reflect.StructTag(s.Tag(i))
		if f.Type().String() == BaseDocument {
			base = i
		}

		field := NewField(f.Name(), f.Type().Underlying().String(), t)
		field.CheckedNode = f
		str := p.tryGetStruct(f.Type())
		if f.Type().String() != BaseDocument && str != nil {
			field.Type = getStructType(f.Type())

			d := false
			for _, v := range done {
				if v == str {
					d = true
					break
				}
			}
			if !d {
				_, subfs := p.processFields(str, append(done, str))
				field.SetFields(subfs)
			}
		}

		fields = append(fields, field)
	}

	return base, fields
}

func getStructType(t types.Type) string {
	ts := t.String()
	if ts != "time.Time" && ts != "bson.ObjectId" {
		return "struct"
	}

	return ts
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
