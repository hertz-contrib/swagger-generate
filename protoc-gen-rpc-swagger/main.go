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
 *
 * Copyright 2020 Google LLC. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file may have been modified by CloudWeGo authors. All CloudWeGo
 * Modifications are Copyright 2024 CloudWeGo Authors.
 */

package main

import (
	"flag"
	"path/filepath"
	"strings"

	"github.com/hertz-contrib/swagger-generate/protoc-gen-rpc-swagger/generator"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var flags flag.FlagSet

const DefaultOutputFile = "openapi.yaml"

func main() {
	conf := generator.Configuration{
		Version:        flags.String("version", "3.0.3", "version number text, e.g. 1.2.3"),
		Title:          flags.String("title", "", "name of the API"),
		Description:    flags.String("description", "", "description of the API"),
		Naming:         flags.String("naming", "json", `naming convention. Use "proto" for passing names directly from the proto files`),
		FQSchemaNaming: flags.Bool("fq_schema_naming", false, `schema naming convention. If "true", generates fully-qualified schema names by prefixing them with the proto message package name`),
		EnumType:       flags.String("enum_type", "integer", `type for enum serialization. Use "string" for string-based serialization`),
		OutputMode:     flags.String("output_mode", "merged", `output generation mode. By default, a single openapi.yaml is generated at the out folder. Use "source_relative' to generate a separate '[inputfile].openapi.yaml' next to each '[inputfile].proto'.`),
	}

	serverConf := generator.ServerConfiguration{
		HertzAddr: flags.String("hertz_addr", "127.0.0.1:8080", "hertz server address"),
		KitexAddr: flags.String("kitex_addr", "127.0.0.1:8888", "kitex server address"),
	}

	opts := protogen.Options{
		ParamFunc: flags.Set,
	}

	opts.Run(func(plugin *protogen.Plugin) error {
		// Enable "optional" keyword in front of type (e.g. optional string label = 1;)
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		if *conf.OutputMode == "source_relative" {
			for _, file := range plugin.Files {
				if !file.Generate {
					continue
				}
				outfileName := strings.TrimSuffix(file.Desc.Path(), filepath.Ext(file.Desc.Path())) + "." + DefaultOutputFile
				outputFile := plugin.NewGeneratedFile(outfileName, "")
				gen := generator.NewOpenAPIGenerator(plugin, conf, []*protogen.File{file})
				if err := gen.Run(outputFile); err != nil {
					return err
				}
			}
		} else {
			outputFile := plugin.NewGeneratedFile(DefaultOutputFile, "")
			gen := generator.NewOpenAPIGenerator(plugin, conf, plugin.Files)
			if err := gen.Run(outputFile); err != nil {
				return err
			}
		}
		serverConf.OutputMode = conf.OutputMode
		outputFile := plugin.NewGeneratedFile("swagger.go", "")
		gen, err := generator.NewServerGenerator(serverConf, plugin.Files)
		if err != nil {
			return err
		}
		if err = gen.Generate(outputFile); err != nil {
			return err
		}
		return nil
	})
}
