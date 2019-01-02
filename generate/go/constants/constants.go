package constants

import (
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/techspaceco/phrasebook"
	"github.com/techspaceco/phrasebook/generate"
)

var _ generate.Generator = (*Constant)(nil)

func init() {
	generate.Register("const", driver)
	generate.Register("constants", driver)
}

func driver() (generate.Generator, error) {
	return New()
}

// Constant generator configuration.
type Constant struct{}

// New generator instance.
func New() (*Constant, error) {
	return &Constant{}, nil
}

// Generate a const string phrasebook.
func (c *Constant) Generate(exports phrasebook.Exports, w io.Writer) error {
	return tmpl.Execute(w, struct {
		Filename  string
		Package   string
		Timestamp time.Time
		Exports   phrasebook.Exports
	}{
		Filename:  "TODO(shane): Filename.",
		Package:   "TODO(shane): Package.",
		Timestamp: time.Now(),
		Exports:   exports,
	})
}

var tmpl = template.Must(template.New("").
	Funcs(template.FuncMap{
		"Lines": func(s string) []string {
			return strings.Split(strings.TrimSpace(s), "\n")
		},
	}).
	Parse(`// Generated with ❤ by github.com/techspaceco/phrasebook; DO NOT EDIT
// {{ .Filename }}
// {{ .Timestamp }}
package {{ .Package }}
{{ range .Exports }}

{{- range Lines .Comment }}
// {{ . }}
{{- end }}
const {{ .Name }} = {{ . | printf "%q" }}
{{ end }}
`))
