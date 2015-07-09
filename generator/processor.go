package generator

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strings"

	_ "golang.org/x/tools/go/gcimporter"
	"golang.org/x/tools/go/types"
)

const BaseDocument = "github.com/tyba/storable.Document"

type Processor struct {
	Path   string
	Ignore map[string]bool

	TypesPkg     *types.Package
	fieldsForStr map[*types.Struct]*[]*Field
}

func NewProcessor(path string, ignore []string) *Processor {
	i := make(map[string]bool, 0)
	for _, file := range ignore {
		i[file] = true
	}

	return &Processor{
		Path:         path,
		Ignore:       i,
		fieldsForStr: map[*types.Struct]*[]*Field{},
	}
}

func (p *Processor) Do() (*Package, error) {
	files, err := p.getSourceFiles()
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

	var base int
	if base, m.Fields = p.getFields(s); base == -1 {
		return nil
	}

	p.procesBaseField(m, m.Fields[base])
	p.findHooks(m)

	return m
}

func (p *Processor) getFields(s *types.Struct) (base int, fields []*Field) {
	base = p.processFields(s)

	for _, fields := range p.fieldsForStr {
		for _, f := range *fields {
			if f.CheckedNode == nil || len(f.Fields) > 0 {
				continue
			}
			if v := p.fieldsForStr[p.tryGetStruct(f.CheckedNode.Type())]; v != nil {
				f.SetFields(*v)
			}
		}
	}

	fields = *p.fieldsForStr[s]

	return
}

// Returns which field index is an embedded storable.Document, or -1 if none.
func (p *Processor) processFields(s *types.Struct) int {
	c := s.NumFields()

	base := -1
	fields := make([]*Field, 0)
	if _, ok := p.fieldsForStr[s]; !ok {
		p.fieldsForStr[s] = &fields
	}

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
			_, ok := p.fieldsForStr[str]
			if !ok {
				p.processFields(str)
			}
		}

		fields = append(fields, field)
	}

	return base
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

func (p *Processor) findHooks(m *Model) {
	modelType := types.NewPointer(p.TypesPkg.Scope().Lookup(m.Name).Type())
	m.Hooks = p.hookMethods(modelType)

	done := map[interface{}]bool{modelType: true}
	for _, f := range m.Fields {
		p.findFieldHooks(f, done)
	}
}

func (p *Processor) findFieldHooks(f *Field, done map[interface{}]bool) {
	if done[f.CheckedNode.Type()] {
		return
	}
	done[f.CheckedNode.Type()] = true

	if f.CheckedNode == nil {
		return
	}

	typ := f.CheckedNode.Type()
out:
	for {
		switch v := typ.(type) {
		case *types.Slice:
			typ = v.Elem()
		case *types.Map:
			typ = v.Elem()
		default:
			break out
		}
	}

	if _, ok := typ.(*types.Pointer); !ok {
		typ = types.NewPointer(typ)
	}

	f.Hooks = p.hookMethods(typ)

	for _, f := range f.Fields {
		p.findFieldHooks(f, done)
	}
}

func (p *Processor) hookMethods(t types.Type) []Hook {
	var hooks []Hook

	ms := types.NewMethodSet(t)

	actions := []HookAction{
		InsertHook,
		UpdateHook,
		SaveHook,
	}

	for _, action := range actions {
		hook := Hook{Before: false, Action: action}
		if p.hasHookMethod(ms, hook.MethodName()) {
			hooks = append(hooks, hook)
		}
		hook.Before = true
		if p.hasHookMethod(ms, hook.MethodName()) {
			hooks = append(hooks, hook)
		}
	}

	return hooks
}

func (p *Processor) hasHookMethod(ms *types.MethodSet, methodName string) bool {
	sel := ms.Lookup(p.TypesPkg, methodName)
	if sel == nil {
		return false
	}
	method, ok := sel.Obj().(*types.Func)
	if !ok {
		return false
	}
	sig, ok := method.Type().(*types.Signature)
	if !ok {
		return false
	}

	if params := sig.Params(); params != nil && params.Len() > 0 {
		return false
	}

	ret := sig.Results()
	if ret == nil || ret.Len() != 1 || !isBuiltinError(ret.At(0).Type()) {
		return false
	}

	return true
}

func isBuiltinError(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	return named.Obj().Name() == "error" && named.Obj().Parent() == types.Universe
}
