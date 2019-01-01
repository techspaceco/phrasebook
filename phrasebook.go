// Package phrasebook implements the pattern for SQL queries.
//
// TODO(shane): Add a -debug mode for use in development like a lot of the bindata packages do.
package phrasebook

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// commentRE matches leading comments before an export.
var commentRE = regexp.MustCompile(`(?m)(?:^|\n)(?P<comment>(?:(?:--|#[\t ])\s*[^\n]*\n)+)(?:--|#[\t ])\s*export[: ]`)

// commentTrimRE is a subexpression matches that can be used to extract comment text without the delimiter.
var commentTrimRE = regexp.MustCompile(`(?m)(?:^|\n)(?:--|#[\t ])\s*([^\n]*\n)`)

// ExportRE expression matches a query to be exported.
//
// -- List staff with a birthday for @date so you can buy them a beer.
// -- export: ListStaffBirthday, {"json":"optional","meta":"data"}
// select
//   u.*
// from users u
// where
//   role_in(u.id, 'staff')
//   and u.birthday = coalesce(@date, current_date)::date;
// -- end
var ExportRE = regexp.MustCompile(`(?m)(?:^|\n)(?:(?:--|#[\t ])\s*[^\n]*\n)*(?:--|#[\t ])\s*export[: ]\s*(?P<name>[\p{L}][\p{L}\p{N}_]*)(?:\s*[, ]\s*(?P<meta>[^\n]+))?(?P<query>(?:\n.+)+)\n(?:--|#[\t ])\s*end`) // Hello from no (?x) land :P

// Exports phrasebook collection.
type Exports map[string]*Export

// Export SQL query.
type Export struct {
	Name    string
	Comment string
	Query   string
	Meta    json.RawMessage
}

func (e Export) String() string {
	return e.Query
}

// TODO: Advance to the next comment where possible or it will fill the buffer
// if it never matches anything which isn't ideal. Not that I expect many > 65kb
// phrasebooks without a single match.
func split(data []byte, eof bool) (advance int, token []byte, err error) {
	match := ExportRE.FindIndex(data)
	switch {
	case match == nil:
		return 0, nil, nil
	case eof:
		return len(data), data[:len(data)], nil
	default:
		return match[1] - 1, data[match[0]:match[1]], nil
	}
}

// Parse an sql phrasebook.
// -- List staff with a birthday for @date so you can buy them a beer.
// -- export: ListStaffBirthday, {"json":"optional","meta":"data"}
// select
//   u.*
// from users u
// where
//   role_in(u.id, 'staff')
//   and u.birthday = coalesce(@date, current_date)::date;
// -- end
func Parse(sql io.Reader) (Exports, error) {
	exports := make(Exports, 0)

	scan := bufio.NewScanner(sql)
	scan.Split(split)
	for scan.Scan() {
		if len(scan.Bytes()) == 0 {
			continue
		}

		export := &Export{}
		if match := commentRE.FindStringSubmatch(scan.Text()); match != nil {
			export.Comment = commentTrimRE.ReplaceAllString(match[1], "$1")
		}
		match := ExportRE.FindSubmatch(scan.Bytes())
		export.Name = strings.TrimSpace(string(match[1]))
		export.Query = strings.TrimSpace(string(match[3]))
		json.Unmarshal(match[2], &export.Meta)

		if _, ok := exports[export.Name]; ok {
			return nil, fmt.Errorf("export %q appears twice in source", export.Name)
		}
		exports[export.Name] = export
	}

	return exports, nil
}

func MustParse(sql io.Reader) Exports {
	exports, err := Parse(sql)
	if err != nil {
		panic(err)
	}
	return exports
}

func MustParseFile(file string) Exports {
	sql, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	return MustParse(sql)
}
