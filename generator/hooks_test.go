package generator_test

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/types"
	. "gopkg.in/check.v1"

	"github.com/tyba/storable/generator"
)

type HooksSuite struct{}

var _ = Suite(&HooksSuite{})

func (s *HooksSuite) TestFindHooks(c *C) {
	genPkg := s.fixture(c)

	m := genPkg.Models[0]
	c.Assert(len(m.Hooks), Equals, 2)
	c.Assert(m.Hooks[0], Equals, generator.Hook{
		Before: true,
		Action: generator.InsertHook,
	})
	c.Assert(m.Hooks[1], Equals, generator.Hook{
		Before: true,
		Action: generator.SaveHook,
	})

	c.Assert(len(m.Fields[2].Hooks), Equals, 2)
	c.Assert(m.Fields[2].Hooks[0], Equals, generator.Hook{
		Before: false,
		Action: generator.InsertHook,
	})
	c.Assert(m.Fields[2].Hooks[1], Equals, generator.Hook{
		Before: false,
		Action: generator.SaveHook,
	})

	c.Assert(len(m.Fields[3].Hooks), Equals, 1)
	c.Assert(m.Fields[3].Hooks[0], Equals, generator.Hook{
		Before: true,
		Action: generator.SaveHook,
	})

	c.Assert(len(m.Fields[4].Hooks), Equals, 1)
	c.Assert(m.Fields[4].Hooks[0], Equals, generator.Hook{
		Before: true,
		Action: generator.SaveHook,
	})
}

func (s *HooksSuite) TestGenerateHooks(c *C) {
	genPkg := s.fixture(c)
	td := &generator.TemplateData{genPkg, nil, map[interface{}]string{}}

	got := td.CallHooks("before", "insert", genPkg.Models[0])
	bs, _ := format.Source([]byte("func main() {\n" + got + "\n}"))
	fmt.Println(string(bs))
}

func (s *HooksSuite) fixture(c *C) *generator.Package {
	fixtureSrc := `
	package fixture

	import 	"github.com/tyba/storable"

	type Recur struct {
		storable.Document
		Foo string
		R *Other
		Things map[string][]*Thing
		MoreThings []Thing
	}

	func (r Recur) BeforeInsert() error {
		return nil
	}

	func (r *Recur) BeforeSave() error {
		return nil
	}

	type Other struct {
		R2 *Recur
	}

	func (r Other) AfterInsert() error {
		return nil
	}

	func (r *Other) AfterSave() error {
		return nil
	}

	func (r *Other) BeforeSave() { // Bad signature.
		return nil
	}

	type Thing int

	func (t Thing) BeforeSave() error {
		return nil
	}
	`

	fset := &token.FileSet{}
	astFile, _ := parser.ParseFile(fset, "fixture.go", fixtureSrc, 0)
	cfg := &types.Config{}
	p, _ := cfg.Check("github.com/tcard/navpatch/navpatch", fset, []*ast.File{astFile}, nil)

	prc := generator.NewProcessor("fixture", nil)
	prc.TypesPkg = p
	genPkg, err := prc.ProcessTypesPkg()

	c.Assert(err, IsNil)

	return genPkg
}
