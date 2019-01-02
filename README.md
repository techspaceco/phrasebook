# phrasebook

Data phrasebook pattern and code generation for SQL.

SQL is the right way to interact with relational databases! By adding simple
annotations to your SQL phrasebook creates clear separation between SQL and Go
and allows for advanced template driven code generation. Use `go:generate` to
export constants, functions, documentation, tests and more from SQL.

## Usage

Start with SQL, say `repository.sql` and add an `-- export` and `-- end`
annotations to any SQL you'd like to export.
```sql
-- List staff with a birthday for @date so you can buy them a beer.
-- export: SelectStaffBirthdaysSQL
select u.*
from users u
where
  role_in(u.id, 'staff')
  and u.birthday = coalesce(@date, current_date)::date;
-- end
```

Create a `repository.go` package and use `const` to generated SQL string constants.
```go
package repository

// go:generate go-gen-phrasebook -package repository -template=const repository.sql

func ListStaffBirthdays() ([]*User, error) {
  rows, err := db.Query(SelectStaffBirthdaysSQL, sql.Named("date", time.Now()))
  // ...
}
```

The `go:generate` const template generates `repository.sql.go` with string constants.
```go
// DO NOT EDIT!
// Generated with ❤ by github.com/techspaceco/phrasebook
// repository.sql
// 2019-01-01 22:14:08.212693 +0000 GMT m=+0.002107042
package repository

// List staff with a birthday for @date so you can buy them a beer.
const SelectStaffBirthdaysSQL = "select u.*
from users u
where
  role_in(u.id, 'staff')
  and u.birthday = coalesce(@date, current_date)::date;"
```

## Advanced

See the generated documentation and examples.

https://godoc.org/github.com/techspaceco/phrasebook

## Pattern

This phrasebook pattern was a lot more common in the early 2000's before popular ORM libraries
started to popularize limiting relational database interaction with SQL abstractions.

https://www.perl.com/pub/2002/10/22/phrasebook.html

## Implementations

This library differs from other Go SQL implementations in the following areas:
* Optional and template driven code generation.
* Optional JSON meta-data retained in export tree allows flexible templates.
* Explicit start/end annotations allow for un-exported comments and source code.
* Leading export comments are matched and retained in generated code.
* The stdlib is the only dependency.

### Go
* https://github.com/nleof/goyesql
* https://github.com/gchaincl/dotsql

### Clojure
* https://www.hugsql.org/
* https://github.com/krisajenkins/yesql

### Perl
* https://metacpan.org/pod/Data::Phrasebook::SQL

## Known Issues

Regexp driven so debugging missing exports due to format errors is a known
weakness of the current parser implementation. Luckily the annotation format
is very simple.

## License

MIT