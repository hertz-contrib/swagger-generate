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

	"github.com/hertz-contrib/swagger-generate/common/consts"
	"github.com/hertz-contrib/swagger-generate/common/tpl"
	"github.com/hertz-contrib/swagger-generate/common/utils"
	"google.golang.org/protobuf/compiler/protogen"
)

type ServerGenerator struct {
	IdlPath string
}

func NewServerGenerator(inputFiles []*protogen.File) (*ServerGenerator, error) {
	var idlPath string
	var genFiles []*protogen.File
	for _, f := range inputFiles {
		if f.Generate {
			genFiles = append(genFiles, f)
		}
	}
	if len(genFiles) > 1 {
		return nil, errors.New("only one .proto file is supported for generation swagger")
	} else if len(genFiles) == 1 {
		idlPath = genFiles[0].Desc.Path()
	} else {
		return nil, errors.New("no .proto files marked for generation")
	}

	return &ServerGenerator{
		IdlPath: idlPath,
	}, nil
}

func (g *ServerGenerator) Generate(outputFile *protogen.GeneratedFile) error {
	filePath := filepath.Join(filepath.Dir(g.IdlPath), consts.DefaultOutputSwaggerFile)
	if utils.FileExists(filePath) {
		return errors.New("swagger.go file already exists")
	}

	tmpl, err := template.New("server").Delims("{{", "}}").Parse(consts.CodeGenerationCommentPbHttp + "\n" + tpl.ServerTemplateHttp)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, g)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if _, err = outputFile.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}
	return nil
}
