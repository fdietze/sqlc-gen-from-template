# sqlc-gen-from-template

[sqlc](https://sqlc.dev/) plugin to generate type-safe code for SQL queries using a template.

Related project: [cornerman/scala-db-codegen](https://github.com/cornerman/scala-db-codegen)

## Installation

I recommend using [devbox](https://www.jetpack.io/devbox) and installing this plugin into your project using a flake: 

```bash
# replace <rev> with the latest commit from this repo
devbox add sqlc github:fdietze/sqlc-gen-from-template/<rev>
```

Alternatively, you can build a binary yourself and put it into your `PATH`:

```bash
go build
```

## Usage

Example usage for [Scala](https://www.scala-lang.org/) with the [magnum](https://github.com/AugustNagro/magnum) database library:

`sqlc.yml`
```yml
version: "2"
plugins:
- name: sqlc-gen-from-template
  process:
    cmd: sqlc-gen-from-template # https://github.com/fdietze/sqlc-gen-from-template
sql:
  - engine: "sqlite"
    queries: "queries.sql"
    schema: "schema.sql"
    codegen:
    - out: backend/src/queries
      plugin: sqlc-gen-from-template
      options:
        query_template: "query_template.go.tmpl"
        query_file_extension: "scala"
        # optional formatter command to format generated code
        formatter_cmd: ".devbox/nix/profile/default/bin/scalafmt --stdin"
```

`schema.sql`
```sql
create table post(
  id integer primary key autoincrement -- rowid
  , parent_id integer
) strict;

```

`queries.sql`
```sql
-- name: getReplyIds :many
select id
from post
where parent_id = ?;
```

`query_template.go.tmpl`
```tmpl
{{- /* 
https://pkg.go.dev/text/template
https://github.com/sqlc-dev/sqlc/blob/main/protos/plugin/codegen.proto
https://github.com/AugustNagro/magnum?tab=readme-ov-file
*/ -}}

{{- define "ScalaType" -}}
{{- $scalaType := .Type.Name -}}
{{- if eq .Type.Name "integer"}}{{ $scalaType = "Long" }}
{{- else if eq .Type.Name "text"}}{{ $scalaType = "String" }}
{{- end -}}
{{- $scalaType }}
{{- end -}}

package backend.queries

import com.augustnagro.magnum
import com.augustnagro.magnum.*

{{range .Comments}}// {{.}}
{{end}}

{{$rowType := printf "Row_%s" .Name -}}
{{- if or (eq .Cmd ":many") (eq .Cmd ":one") }}
  {{- if gt (len .Columns) 1 -}}
    case class {{ $rowType }}({{- range .Columns}}
    {{.Name}}:
    {{- if not .NotNull }}Option[{{end}}
    {{- template "ScalaType" .}}
    {{- if not .NotNull }}]{{end}},
    {{- end}}
)
  {{- else -}}


    type {{ $rowType }} = 
    {{- if not (index .Columns 0).NotNull }}Option[{{end}}
    {{- template "ScalaType" (index .Columns 0) }}
    {{- if not (index .Columns 0).NotNull }}]{{end}}
  {{- end}}

{{end}}


{{- $returnType := "__DEFAULT__" -}}
{{- if eq .Cmd ":exec" }}
  {{- $returnType = "Unit" -}}
{{- else if eq .Cmd ":many" }}
  {{- $returnType = printf "Vector[%s]" $rowType -}}
{{- else if eq .Cmd ":one" }}
  {{- $returnType = $rowType -}}
{{- else -}}
  {{- $returnType = "__UNKNOWN_QUERY_ANNOTATION__" -}}
{{- end -}}


def {{.Name}}({{range .Params}}
  {{.Column.Name}}:{{template "ScalaType" .Column}},
{{- end}}
)(using con: DbCon): {{ $returnType }} = {
  Frag("""
  {{ .Text }}
  """, params = IArray({{range .Params}}
  {{.Column.Name}},
  {{end}}))
  {{- if eq .Cmd ":exec" }}.update.run(){{end}}
  {{- if eq .Cmd ":many" }}.query[{{ $rowType }}].run(){{end}}
  {{- if eq .Cmd ":one" }}.query[{{ $rowType }}].run().head{{end}}
}
```

Running `sqlc generate` generates files for every query:

`getReplyIds.scala`
```scala
package backend.queries

import com.augustnagro.magnum
import com.augustnagro.magnum.*

type Row_getReplyIds = Long

def getReplyIds(
  parent_id: Long
)(using con: DbCon): Vector[Row_getReplyIds] = {
  Frag(
    """
  
select id
from post
where parent_id = ?
  """,
    params = IArray(
      parent_id
    ),
  ).query[Row_getReplyIds].run()
}
```
