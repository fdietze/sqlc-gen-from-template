# Languages and Database Libraries

To use a template, copy the respective file into your repo, set it in the `sqlc.yaml` and commit it. Don't hesitate to adjust it to your needs. Every project has different requirements.

If you come up with a template for a new language or database library, please contribute it, even if it is not perfect. It might be a good starting point for someone else.

# Creating a new template from scratch

Look at the [protobuf data structures provided by sqlc](https://github.com/sqlc-dev/sqlc/blob/main/protos/plugin/codegen.proto). They provide a list of queries and schemas. Imagine how you would turn them into your desired code. You probably want to loop over the queries and generate a return type and function for every query.

It would look something like this:

```tmpl
Loop over queries:
{{range .Queries }}


define return type:
struct Row_{{ .Name }} {
  {{range .Columns}}
    {{.Name}}: {{.Type.Name}}
  {{end}}
}

decide if the query is :many, :exec, :one etc.
{{if eq .Cmd ":many"}}

define function:
function {{.Name}} (
  {{range .Params}}
    {{.Column.Name}}:{{.Column.Type.Name}},
  {{- end}}
 conn: &rusqlite::Connection
 ): Row_{{.Name}} {

    let query = "{{ .Text }}"
    ... execute query ...
}
{{end}}



{{if eq .Cmd ":one" }}
...
{{end}}


{{end}}
```
