package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tyba/mongogen/generator"
)

func main() {
	filename := flag.String("input", "", "input file name")
	flag.Parse()
	if *filename == "" {
		flag.Usage()
		return
	}

	gen := generator.NewGenerator(*filename)
	err := gen.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}
