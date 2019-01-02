package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/techspaceco/phrasebook"
	"github.com/techspaceco/phrasebook/generate"
	_ "github.com/techspaceco/phrasebook/generate/go/constants"
	_ "github.com/techspaceco/phrasebook/generate/go/mapped"
)

func main() {
	pkg := flag.String("package", "", "Output class/package name.")
	out := flag.String("output", "", "Output pathname. default: stdout.")
	tmpl := flag.String("template", "constants", "Output template. default: constants")
	flag.Parse()

	in := flag.Arg(0)

	// TODO(shane): Move most of this into the generator so it can be run and tested as a lib.

	input, err := os.Open(in)
	exitOnError(err)
	defer input.Close()

	if *pkg == "" {
		*pkg = path.Base(path.Dir(in))
	}

	output := os.Stdout
	if *out != "" {
		err := os.MkdirAll(path.Dir(*out), os.ModePerm)
		exitOnError(err)

		output, err := os.Create(*out)
		exitOnError(err)
		defer output.Close()
	}

	var template io.Reader
	if *tmpl != "" && !generate.HasDriver(*tmpl) {
		template, err = os.Open(*tmpl)
		exitOnError(err)
		*tmpl = "template"
	}

	generator, err := generate.New(*tmpl, template)
	exitOnError(err)

	exports, err := phrasebook.Parse(input)
	exitOnError(err)

	file := &generate.File{
		Source:  path.Base(in),
		Package: *pkg,
		Exports: exports,
	}

	err = generator.Generate(file, output)
	exitOnError(err)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n\n", err)
		flag.PrintDefaults()
		os.Exit(1)
	}
}
