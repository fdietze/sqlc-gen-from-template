package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"github.com/sqlc-dev/plugin-sdk-go/codegen"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

func main() {
	codegen.Run(generate)
}

type Options struct {
	QueryTemplate      string `json:"query_template" yaml:"query_template"`
	QueryFileExtension string `json:"query_file_extension" yaml:"query_file_extension"`
	Out                string `json:"out" yaml:"out"`
}

func parseOpts(req *plugin.GenerateRequest) (*Options, error) {
	var options Options
	if len(req.PluginOptions) == 0 {
		return &options, nil
	}
	if err := json.Unmarshal(req.PluginOptions, &options); err != nil {
		return nil, fmt.Errorf("unmarshalling plugin options: %w", err)
	}

	return &options, nil
}

func generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	// fmt.Println(req)
	options, _ := parseOpts(req)
	templateFileName := options.QueryTemplate

	tmpl, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Fatalf("Error parsing template file: %v", err)
	}

	resp := plugin.GenerateResponse{}

	for _, query := range req.Queries {
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, query)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}

		resp.Files = append(resp.Files, &plugin.File{
			Name:     fmt.Sprintf("%s.%s", query.Name, options.QueryFileExtension),
			Contents: buf.Bytes(),
		})
	}

	return &resp, nil
}
