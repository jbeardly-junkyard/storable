package generator_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/types"

	"github.com/tyba/storable/generator"
)

func TestRecursiveStruct(t *testing.T) {
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
	genPkg, err := prc.ProcessTypesPkg(p)

	if err != nil {
		t.Fatal(err)
	}

	if genPkg.Models[0].Fields[2].Fields[2] != genPkg.Models[0].Fields[2] {
		t.Fatalf("direct type recursivity not handled correctly.")
	}
}

func TestDeepRecursiveStruct(t *testing.T) {
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
	genPkg, err := prc.ProcessTypesPkg(p)

	if err != nil {
		t.Fatal(err)
	}

	if genPkg.Models[0].Fields[2].Fields[0].Fields[2] != genPkg.Models[0].Fields[2] {
		t.Fatalf("indirect type recursivity not handled correctly.")
	}
}
