package generator_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/types"
	. "gopkg.in/check.v1"

	"github.com/tyba/storable/generator"
)

func Test(t *testing.T) { TestingT(t) }

type ProcessorSuite struct{}

var _ = Suite(&ProcessorSuite{})

func (s *ProcessorSuite) TestRecursiveStruct(c *C) {
	fixtureSrc := `
	package fixture

	import 	"github.com/tyba/storable"

	type Recur struct {
		storable.Document
		Foo string
		R *Recur
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
	c.Assert(genPkg.Models[0].Fields[2].Fields[2], Equals, genPkg.Models[0].Fields[2], Commentf("direct type recursivity not handled correctly."))
}

func (s *ProcessorSuite) TestDeepRecursiveStruct(c *C) {
	fixtureSrc := `
	package fixture

	import 	"github.com/tyba/storable"

	type Recur struct {
		storable.Document
		Foo string
		R *Other
	}

	type Other struct {
		R *Recur
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
	c.Assert(genPkg.Models[0].Fields[2].Fields[0].Fields[2], Equals, genPkg.Models[0].Fields[2], Commentf("direct type recursivity not handled correctly."))
}
