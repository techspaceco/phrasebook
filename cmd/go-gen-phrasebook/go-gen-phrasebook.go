package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/techspaceco/phrasebook"
	"github.com/techspaceco/phrasebook/generate"
	_ "github.com/techspaceco/phrasebook/generate/go/constants"
	_ "github.com/techspaceco/phrasebook/generate/go/functions"
)

func main() {
	// TODO(shane): Add a -parser to support sources other than SQL.
	// TODO(shane): Add a -package but default to the CWD folder name.
	// TODO(shane): Add a -template allowing a text/template source file to be supplied.
	// TODO(shane): Add a -output but allow nil leaving it up to the generator based on the filename.
	pkg := "test"
	_ = pkg // TODO:

	inFile := filepath.Join("_testdata", "phrasebook.sql") // TODO: Duh.
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
