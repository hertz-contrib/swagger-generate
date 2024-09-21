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

	"github.com/hertz-contrib/swagger-generate/common/consts"
	"github.com/hertz-contrib/swagger-generate/common/tpl"
	"github.com/hertz-contrib/swagger-generate/common/utils"
	"google.golang.org/protobuf/compiler/protogen"
)

type ServerConfiguration struct {
	KitexAddr *string
}

type ServerGenerator struct {
	IdlPath   string
	KitexAddr string
}

func NewServerGenerator(conf ServerConfiguration, inputFiles []*protogen.File) (*ServerGenerator, error) {
	kitexAddr := conf.KitexAddr
	if kitexAddr == nil {
		*kitexAddr = consts.DefaultKitexAddr
	}

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
	// Check if Hertz and Kitex addresses are valid (basic validation)
	if err := validateAddress(*kitexAddr); err != nil {
		return nil, fmt.Errorf("invalid Kitex address: %w", err)
	}

	return &ServerGenerator{
		IdlPath:   idlPath,
		KitexAddr: *kitexAddr,
	}, nil
}

func validateAddress(addr string) error {
	if addr == "" {
		return errors.New("address cannot be empty")
	}
	if !strings.Contains(addr, ":") {
		return errors.New("address must include a port (e.g., '127.0.0.1:8080')")
	}
	return nil
}

func (g *ServerGenerator) Generate(outputFile *protogen.GeneratedFile) error {
	filePath := filepath.Join(filepath.Dir(g.IdlPath), consts.DefaultOutputSwaggerFile)
	if utils.FileExists(filePath) {
		updatedContent, err := updateVariables(filePath, g.KitexAddr, g.IdlPath)
		if err != nil {
			return errors.New("failed to update variables in the existing file")
		}
		if _, err = outputFile.Write([]byte(updatedContent)); err != nil {
			return errors.New("failed to write output file")
		}
	} else {
		tmpl, err := template.New("server").Delims("{{", "}}").Parse(consts.CodeGenerationCommentPbRpc + "\n" + tpl.ServerTemplateRpcPb)
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
	}
	return nil
}

func updateVariables(filePath, newKitexAddr, newIdlPath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	kitexAddrPattern := regexp.MustCompile(`kitexAddr\s*=\s*"(.*?)"`)
	idlPathPattern := regexp.MustCompile(`idlFile\s*=\s*"(.*?)"`)

	updatedContent := kitexAddrPattern.ReplaceAllString(string(content), fmt.Sprintf(`kitexAddr = "%s"`, newKitexAddr))
	updatedContent = idlPathPattern.ReplaceAllString(updatedContent, fmt.Sprintf(`idlFile = "%s"`, newIdlPath))

	return updatedContent, nil
}
