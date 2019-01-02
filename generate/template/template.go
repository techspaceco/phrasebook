package template

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
	text "text/template"
	"time"

	"github.com/techspaceco/phrasebook/generate"
)

var _ generate.Generator = (*Template)(nil)

func init() {
	generate.Register("template", driver)
}

// TODO: Logger interface.
func driver(template io.Reader) (generate.Generator, error) {
	return New(template)
}

// Template generator configuration.
type Template struct {
	tmpl *template.Template
}

// New generator instance.
func New(template io.Reader) (*Template, error) {
	buf, err := ioutil.ReadAll(template)
	if err != nil {
		return nil, nil
	}

	tmpl, err := text.New("").Funcs(text.FuncMap{"lines": lines, "checksum": checksum, "current_time": timestamp}).Parse(string(buf))
	if err != nil {
		return nil, err
	}

	return &Template{tmpl: tmpl}, nil
}

// Generate a const string phrasebook.
func (c *Template) Generate(file *generate.File, w io.Writer) error {
	return c.tmpl.Execute(w, file)
}

// Lines helper function.
func lines(s string) []string {
	return strings.Split(strings.TrimSpace(s), "\n")
}

func checksum(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Timestamp helper function.
func timestamp() string {
	return time.Now().UTC().String()
}
