package constants

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
package {{ .Package }}
{{ range .Exports }}

{{- range lines .Comment }}
// {{ . }}
{{- else }}
// {{ .Name }} SQL
{{- end }}
//
{{- range lines .Query }}
//   {{ . }}
{{- end }}
const {{ .Name }} = {{ . | printf "%q" }}
{{ end }}
`

func init() {
	generate.Register("const", driver)
	generate.Register("constants", driver)
}

func driver(io.Reader) (generate.Generator, error) {
	return template.New(bytes.NewBufferString(tmpl))
}
