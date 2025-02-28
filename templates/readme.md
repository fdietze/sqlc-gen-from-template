# Languages and Database Libraries

To use a template, copy the respective file into your repo, set it in the `sqlc.yaml` and commit it along with the generated code. Don't hesitate to adjust it to your needs. Every project has different requirements. Commiting the template and the generated code allows you to see diffs once you change the template or add new queries.

If you come up with a template for a new language or database library, please contribute it, even if it is not perfect. It might be a good starting point for someone else.

A test project for the rustql version that could be used as a simple project template is available at [https://github.com/ReenigneCA/rust_sqlc_test](https://github.com/ReenigneCA/rust_sqlc_test)


# Creating a new template from scratch

Look at the [protobuf data structures provided by sqlc](https://github.com/sqlc-dev/sqlc/blob/main/protos/plugin/codegen.proto). They provide a list of queries and schemas. Imagine how you would turn them into your desired code. You probably want to loop over the queries and generate a return type and function for every query.

Additionally, some functions are mapped into the template:
 
- strings.ToLower is mapped to ToLower which allows for case-insensitive comparisons particularly for SQL types;
- strings.Contains is mapped to Contains the Rust example uses this to alter options based on the filename being processed;
- GetPluginOption allows you to obtain values under plugin: options: in the sqlc.yaml so you can have custom flags for your templates, it returns "" if no value is set;
- Please see the rustqlite example to see these being used.

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
