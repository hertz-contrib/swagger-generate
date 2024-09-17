/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package generator

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/plugin"
	"github.com/hertz-contrib/swagger-generate/common/consts"
	"github.com/hertz-contrib/swagger-generate/common/tpl"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-http-swagger/args"
)

type ServerGenerator struct {
	OutputDir string
}

func NewServerGenerator(ast *parser.Thrift, args *args.Arguments) (*ServerGenerator, error) {
	defaultOutputDir := consts.DefaultOutputDir

	idlPath := ast.Filename
	if idlPath == "" {
		return nil, errors.New("failed to get Thrift file path")
	}

	outputDir := args.OutputDir
	if outputDir == "" {
		outputDir = defaultOutputDir
	}

	return &ServerGenerator{
		OutputDir: outputDir,
	}, nil
}

func (g *ServerGenerator) Generate() ([]*plugin.Generated, error) {
	filePath := filepath.Join(g.OutputDir, consts.DefaultOutputSwaggerFile)

	tmpl, err := template.New("server").Delims("{{", "}}").Parse(consts.CodeGenerationCommentThriftHttp + "\n" + tpl.ServerTemplateHttp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, g)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %v", err)
	}

	return []*plugin.Generated{{
		Content: buf.String(),
		Name:    &filePath,
	}}, nil
}
