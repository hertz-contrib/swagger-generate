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

package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hertz-contrib/swagger-generate/common/consts"
	"github.com/hertz-contrib/swagger-generate/swagger2idl/converter"
	"github.com/hertz-contrib/swagger-generate/swagger2idl/generate"
	"github.com/hertz-contrib/swagger-generate/swagger2idl/parser"
	"github.com/urfave/cli/v2"
)

var (
	outputType    string
	outputFile    string
	openapiOption bool
	apiOption     bool
	namingOption  bool
)

func main() {
	app := &cli.App{
		Name:  "swagger2idl",
		Usage: "Convert OpenAPI specs to Protobuf or Thrift IDL",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "type",
				Aliases:     []string{"t"},
				Usage:       "Specify output type: 'proto' or 'thrift'. If not provided, inferred from output file extension.",
				Destination: &outputType,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "Specify output file path. If not provided, defaults to output.proto or output.thrift based on the output type.",
				Destination: &outputFile,
			},
			&cli.BoolFlag{
				Name:        "openapi",
				Aliases:     []string{"oa"},
				Usage:       "Include OpenAPI specific options in the output",
				Destination: &openapiOption,
			},
			&cli.BoolFlag{
				Name:        "api",
				Aliases:     []string{"a"},
				Usage:       "Include API specific options in the output",
				Destination: &apiOption,
			},
			&cli.BoolFlag{
				Name:        "naming",
				Aliases:     []string{"n"},
				Usage:       "use naming conventions for the output IDL file",
				Value:       true,
				Destination: &namingOption,
			},
		},
		Action: func(c *cli.Context) error {
			args := c.Args().Slice()

			if len(args) < 1 {
				log.Fatal("Please provide the path to the OpenAPI file.")
			}

			openapiFile := args[0]

			// Automatically determine output type based on file extension if not provided
			if outputType == "" && outputFile != "" {
				ext := filepath.Ext(outputFile)
				switch ext {
				case ".proto":
					outputType = consts.IDLProto
				case ".thrift":
					outputType = consts.IDLThrift
				default:
					log.Fatalf("Cannot determine output type from file extension: %s. Use --type to specify explicitly.", ext)
				}
			}

			if outputFile == "" {
				if outputType == consts.IDLProto {
					outputFile = consts.DefaultProtoFilename
				} else if outputType == consts.IDLThrift {
					outputFile = consts.DefaultThriftFilename
				} else {
					log.Fatal("Output file must be specified if output type is not provided.")
				}
			}

			spec, err := parser.LoadOpenAPISpec(openapiFile)
			if err != nil {
				log.Fatalf("Failed to load OpenAPI file: %v", err)
			}

			converterOption := &converter.ConvertOption{
				OpenapiOption: openapiOption,
				ApiOption:     apiOption,
				NamingOption:  namingOption,
			}

			var idlContent string
			var file *os.File
			var errFile error

			switch outputType {
			case consts.IDLProto:
				protoConv := converter.NewProtoConverter(spec, converterOption)

				if err = protoConv.Convert(); err != nil {
					log.Fatalf("Error during conversion: %v", err)
				}
				protoEngine := generate.NewProtoGenerate()

				idlContent, err = protoEngine.Generate(protoConv.GetIdl())
				if err != nil {
					log.Fatalf("Error generating proto docs: %v", err)
				}

				file, errFile = os.Create(outputFile)
			case consts.IDLThrift:
				thriftConv := converter.NewThriftConverter(spec, converterOption)

				if err = thriftConv.Convert(); err != nil {
					log.Fatalf("Error during conversion: %v", err)
				}
				thriftEngine := generate.NewThriftGenerate()

				idlContent, err = thriftEngine.Generate(thriftConv.GetIdl())
				if err != nil {
					log.Fatalf("Error generating thrift docs: %v", err)
				}

				file, errFile = os.Create(outputFile)
			default:
				log.Fatalf("Invalid output type: %s", outputType)
			}

			if errFile != nil {
				log.Fatalf("Failed to create file: %v", errFile)
			}
			defer func() {
				if err := file.Close(); err != nil {
					log.Printf("Error closing file: %v", err)
				}
			}()

			if _, err = file.WriteString(idlContent); err != nil {
				log.Fatalf("Error writing to file: %v", err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
