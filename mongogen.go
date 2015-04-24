package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	_ "golang.org/x/tools/go/gcimporter"
	"golang.org/x/tools/go/types"
)

func main() {
	parser := flags.NewNamedParser("mongogen", flags.Default)
	parser.AddCommand(
		"gen",
		"Generate files for types using mongogen document.",
		"",
		&CmdGenerate{},
	)

	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrCommandRequired {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}

}

type CmdGenerate struct {
	Input  string `short:"" long:"input" description:"input package directory" default:"."`
	Output string `short:"" long:"output" description:"output file name" default:"base.go"`
}

func (c *CmdGenerate) Execute(args []string) error {
	if !isDirectory(c.Input) {
		return fmt.Errorf("Input path should be a directory %s", c.Input)
	}

	files, err := getSourceFiles(c.Input)
	if err != nil {
		return err
	}

	pkg, err := parseSourceFiles(c.Input, files)
	if err != nil {
		return err
	}

	processPackage(pkg)

	return nil
}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}

	return info.IsDir()
}

func getSourceFiles(directory string) ([]string, error) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		return nil, fmt.Errorf("cannot process directory %s: %s", directory, err)
	}

	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)

	if len(files) == 0 {
		return nil, fmt.Errorf("%s: no buildable Go files", directory)
	}

	return joinDirectory(directory, files), nil
}

func joinDirectory(directory string, files []string) []string {
	r := make([]string, len(files))
	for i, file := range files {
		r[i] = filepath.Join(directory, file)
	}

	return r
}

func parseSourceFiles(directory string, filenames []string) (*types.Package, error) {
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

	return config.Check(directory, fs, files, info)
}

func processPackage(pkg *types.Package) {
	s := pkg.Scope()
	for _, name := range s.Names() {
		obj := s.Lookup(name)

		fmt.Printf("%s ---> %T\n", name, obj.Type())
		named, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}

		str, ok := named.Underlying().(*types.Struct)
		if !ok {
			continue
		}

		for i := 0; i < str.NumFields(); i++ {
			f := str.Field(i)
			fmt.Println("\t", f.Name(), f.Type(), str.Tag(i))
		}
	}
}
