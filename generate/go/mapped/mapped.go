// Package mapped is a map[string]string phrasebook generator.
//
// The phrasebook will be available as `const phrasebook = map[string]string{}`.
//
// Like yesql + bindata, dotsql etc. there is no no compile time guarantee your
// SQL will exist in the phrasebook. Using the `const` or `accessors` templates
// is suggested but YMMV.
package mapped

import (
	"bytes"
	"io"

	"github.com/techspaceco/phrasebook/generate"
	"github.com/techspaceco/phrasebook/generate/template"
)

const tmpl = `// Generated with ❤ ; DO NOT EDIT
// generator: github.com/techspaceco/phrasebook
// source:    {{ .Source }}
// checksum:  {{ checksum .Source }}
// timestamp: {{ current_time }}
package {{ .Package }}

// Phrasebook SQL.
const phrasebook = map[string]string{
	{{ range .Exports }}
	{{- range lines .Comment }}
	// {{ . }}
	{{- end }}
	{{ printf "%q" .Name }}: {{ printf "%q" . | multiline }},
	{{ end }}
}
`

func init() {
	generate.Register("map", driver)
	generate.Register("mapped", driver)
}

func driver(io.Reader) (generate.Generator, error) {
	return template.New(bytes.NewBufferString(tmpl))
}
