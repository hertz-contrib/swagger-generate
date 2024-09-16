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
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/plugin"
	"github.com/hertz-contrib/swagger-generate/common/consts"
	"github.com/hertz-contrib/swagger-generate/common/tpl"
	"github.com/hertz-contrib/swagger-generate/common/utils"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger/args"
)

type ServerGenerator struct {
	IdlPath   string
	KitexAddr string
	OutputDir string
}

func NewServerGenerator(ast *parser.Thrift, args *args.Arguments) (*ServerGenerator, error) {
	defaultKitexAddr := consts.DefaultKitexAddr
	defaultOutputDir := consts.DefaultOutputDir

	idlPath := ast.Filename
	if idlPath == "" {
		return nil, errors.New("failed to get Thrift file path")
	}

	kitexAddr := args.KitexAddr
	if kitexAddr == "" {
		kitexAddr = defaultKitexAddr
	}

	outputDir := args.OutputDir
	if outputDir == "" {
		outputDir = defaultOutputDir
	}

	if err := validateAddress(kitexAddr); err != nil {
		return nil, err
	}

	return &ServerGenerator{
		IdlPath:   idlPath,
		KitexAddr: kitexAddr,
		OutputDir: outputDir,
	}, nil
}

func (g *ServerGenerator) Generate() ([]*plugin.Generated, error) {
	filePath := filepath.Join(g.OutputDir, consts.DefaultOutputSwaggerFile)

	if utils.FileExists(filePath) {
		updatedContent, err := updateVariables(filePath, g.KitexAddr, g.IdlPath)
		if err != nil {
			return nil, err
		}
		return []*plugin.Generated{{
			Content: updatedContent,
			Name:    &filePath,
		}}, nil
	}

	tmpl, err := template.New("server").Delims("{{", "}}").Parse(consts.CodeGenerationCommentThriftRpc + "\n" + tpl.ServerTemplateRpc)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, g)
	if err != nil {
		return nil, err
	}

	return []*plugin.Generated{{
		Content: buf.String(),
		Name:    &filePath,
	}}, nil
}

func updateVariables(filePath, newKitexAddr, newIdlPath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	kitexAddrPattern := regexp.MustCompile(`kitexAddr\s*=\s*"(.*?)"`)
	idlPathPattern := regexp.MustCompile(`idlFile\s*=\s*"(.*?)"`)

	updatedContent := kitexAddrPattern.ReplaceAllString(string(content), fmt.Sprintf(`kitexAddr = "%s"`, newKitexAddr))
	updatedContent = idlPathPattern.ReplaceAllString(updatedContent, fmt.Sprintf(`idlFile = "%s"`, newIdlPath))

	return updatedContent, nil
}

func validateAddress(addr string) error {
	if !strings.Contains(addr, ":") {
		return errors.New("address must include a port (e.g., '127.0.0.1:8888')")
	}
	return nil
}
