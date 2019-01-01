package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shanna/phrasebook"
	"github.com/shanna/phrasebook/generate"
	_ "github.com/shanna/phrasebook/generate/go/constants"
	_ "github.com/shanna/phrasebook/generate/go/functions"
)

func main() {
	// -package by default it will be the folder name.
	// -o filename.go by default it will add .go
	pkg := "test"
	_ = pkg
	inFile := filepath.Join("_testdata", "phrasebook.sql")
	outFile := inFile + ".go"

	dir, err := os.Getwd() // Default, override with -package.
	_ = dir
	exitOnError(err)

	in, err := os.Open(inFile)
	exitOnError(err)
	defer in.Close()

	exports, err := phrasebook.Parse(in)
	exitOnError(err)

	generator, err := generate.New("const")
	exitOnError(err)

	out, err := os.Create(outFile)
	exitOnError(err)
	defer out.Close()

	err = generator.Generate(exports, out)
	exitOnError(err)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}
}
