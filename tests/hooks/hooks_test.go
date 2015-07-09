package hooks_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"golang.org/x/tools/go/types"
	. "gopkg.in/check.v1"
	"gopkg.in/mgo.v2"

	"github.com/tyba/storable"
	"github.com/tyba/storable/generator"
	"github.com/tyba/storable/tests/hooks"
)

func Test(t *testing.T) { TestingT(t) }

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
	c.Assert(len(m.StoreHooks), Equals, 1)
	c.Assert(m.StoreHooks[0], Equals, generator.Hook{
		Before: true,
		Action: generator.UpdateHook,
	})

	c.Assert(len(m.Fields[2].Hooks), Equals, 3)
	c.Assert(m.Fields[2].Hooks[0], Equals, generator.Hook{
		Before: false,
		Action: generator.InsertHook,
	})
	c.Assert(m.Fields[2].Hooks[1], Equals, generator.Hook{
		Before: false,
		Action: generator.UpdateHook,
	})
	c.Assert(m.Fields[2].Hooks[2], Equals, generator.Hook{
		Before: false,
		Action: generator.SaveHook,
	})

	c.Assert(len(m.Fields[3].Hooks), Equals, 1)
	c.Assert(m.Fields[3].Hooks[0], Equals, generator.Hook{
		Before: true,
		Action: generator.SaveHook,
	})
}

func (s *HooksSuite) TestGenerateHooks(c *C) {
	conn, _ := mgo.Dial("localhost")
	db := conn.DB("storable-test")
	defer db.DropDatabase()

	store := hooks.NewRecurStore(db)

	doc := store.New()
	doc.Foo = "Bar"
	doc.R = &hooks.Other{Name: "MyOther", R2: doc}
	failer := &hooks.Failer{}
	doc.MyFailer = failer
	doc.Things = map[string][]*hooks.Thing{
		"foo": nil,
		"bar": {nil, &hooks.Thing{123}},
		"qux": {&hooks.Thing{456}, &hooks.Thing{789}},
	}

	err := store.Insert(doc)
	c.Assert(hooks.Log, DeepEquals, []string{
		"Called BeforeInsert on Recur with Foo Bar",
		"Called BeforeSave on *Recur with Foo Bar",
		"Called BeforeInsert on *Failer",
	})
	c.Assert(err, Not(IsNil))
	herr := err.(storable.HookError)
	c.Assert(herr.Hook, Equals, "BeforeInsert")
	c.Assert(herr.Field, Equals, ".MyFailer")
	c.Assert(herr.Cause.Error(), Equals, "I failed, sorry!")

	hooks.Log = nil
	doc.MyFailer = nil

	err = store.Insert(doc)
	c.Assert(hooks.Log[:2], DeepEquals, []string{
		"Called BeforeInsert on Recur with Foo Bar",
		"Called BeforeSave on *Recur with Foo Bar",
	})
	c.Assert(s.sortedStrs(hooks.Log[2:5]), DeepEquals, s.sortedStrs([]string{
		"Called BeforeSave on Thing 123",
		"Called BeforeSave on Thing 456",
		"Called BeforeSave on Thing 789",
	}))
	c.Assert(hooks.Log[5:], DeepEquals, []string{
		"Called AfterInsert on Other with Name MyOther",
		"Called AfterSave on *Other with Name MyOther",
	})
	c.Assert(err, IsNil)

	hooks.Log = nil
	doc.MyAfterFailer = &hooks.AfterFailer{}

	err = store.Update(doc)
	c.Assert(hooks.Log[:1], DeepEquals, []string{
		"Called BeforeSave on *Recur with Foo Bar",
	})
	c.Assert(s.sortedStrs(hooks.Log[1:4]), DeepEquals, s.sortedStrs([]string{
		"Called BeforeSave on Thing 123",
		"Called BeforeSave on Thing 456",
		"Called BeforeSave on Thing 789",
	}))
	c.Assert(hooks.Log[4:], DeepEquals, []string{
		"Called BeforeUpdate(s) on *Recur with Foo Bar",
		"Called AfterUpdate on Other with Name MyOther",
		"Called AfterSave on *Other with Name MyOther",
		"Called AfterSave on *AfterFailer",
	})
	c.Assert(err, Not(IsNil))
	herr = err.(storable.HookError)
	c.Assert(herr.Hook, Equals, "AfterSave")
	c.Assert(herr.Field, Equals, ".MyAfterFailer")
	c.Assert(herr.Cause.Error(), Equals, "I failed too late, sorry!")

	hooks.Log = nil
	doc.MyAfterFailer = nil

	err = store.Update(doc)
	c.Assert(hooks.Log[:1], DeepEquals, []string{
		"Called BeforeSave on *Recur with Foo Bar",
	})
	c.Assert(s.sortedStrs(hooks.Log[1:4]), DeepEquals, s.sortedStrs([]string{
		"Called BeforeSave on Thing 123",
		"Called BeforeSave on Thing 456",
		"Called BeforeSave on Thing 789",
	}))
	c.Assert(hooks.Log[4:], DeepEquals, []string{
		"Called BeforeUpdate(s) on *Recur with Foo Bar",
		"Called AfterUpdate on Other with Name MyOther",
		"Called AfterSave on *Other with Name MyOther",
	})
	c.Assert(err, IsNil)
}

func (s *HooksSuite) fixture(c *C) *generator.Package {
	_, thisFile, _, _ := func() (uintptr, string, int, bool) { return runtime.Caller(0) }()
	f, _ := os.Open(filepath.Dir(thisFile) + "/fixture.go")
	stof, _ := os.Open(filepath.Dir(thisFile) + "/storable.go")

	fset := &token.FileSet{}
	fAST, _ := parser.ParseFile(fset, "fixture.go", f, 0)
	stofAST, _ := parser.ParseFile(fset, "storable.go", stof, 0)
	cfg := &types.Config{}
	p, _ := cfg.Check("github.com/tcard/navpatch/navpatch", fset, []*ast.File{fAST, stofAST}, nil)

	prc := generator.NewProcessor("fixture", nil)
	prc.TypesPkg = p
	genPkg, err := prc.ProcessTypesPkg()

	c.Assert(err, IsNil)

	return genPkg
}

func (s *HooksSuite) sortedStrs(sl []string) []string {
	sort.StringSlice(sl).Sort()
	return sl
}
