package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tyba/mongogen"
)

type CmdGenerate struct {
	Input  string `short:"" long:"input" description:"input package directory" default:"."`
	Output string `short:"" long:"output" description:"output file name" default:"base.go"`
}

func (c *CmdGenerate) Execute(args []string) error {
	if !isDirectory(c.Input) {
		return fmt.Errorf("Input path should be a directory %s", c.Input)
	}

	p := mongogen.NewProcessor(c.Input)
	name, models, err := p.Do()
	if err != nil {
		return nil
	}

	gen := mongogen.NewGenerator(filepath.Join(c.Input, c.Output))
	gen.Generate(name, models)

	return nil
}

func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}

	return info.IsDir()
}
