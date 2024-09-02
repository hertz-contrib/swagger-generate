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

package generator

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/plugin"
	"github.com/cloudwego/thriftgo/thrift_reflection"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-http-swagger/args"
	openapi "github.com/hertz-contrib/swagger-generate/thrift-gen-http-swagger/thrift"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-http-swagger/utils"
)

type OpenAPIGenerator struct {
	fileDesc          *thrift_reflection.FileDescriptor
	ast               *parser.Thrift
	generatedSchemas  []string
	requiredSchemas   []string
	commentPattern    *regexp.Regexp
	linterRulePattern *regexp.Regexp
}

// NewOpenAPIGenerator creates a new generator for a protoc plugin invocation.
func NewOpenAPIGenerator(ast *parser.Thrift) *OpenAPIGenerator {
	_, fileDesc := thrift_reflection.RegisterAST(ast)
	return &OpenAPIGenerator{
		fileDesc:          fileDesc,
		ast:               ast,
		generatedSchemas:  make([]string, 0),
		commentPattern:    regexp.MustCompile(`//\s*(.*)|/\*([\s\S]*?)\*/`),
		linterRulePattern: regexp.MustCompile(`\(-- .* --\)`),
	}
}

func (g *OpenAPIGenerator) BuildDocument(arguments *args.Arguments) []*plugin.Generated {
	d := &openapi.Document{}

	version := OpenAPIVersion
	d.Openapi = version
	d.Info = &openapi.Info{
		Title:       DefaultInfoTitle,
		Description: DefaultInfoDesc,
		Version:     DefaultInfoVersion,
	}
	d.Paths = &openapi.Paths{}
	d.Components = &openapi.Components{
		Schemas: &openapi.SchemasOrReferences{
			AdditionalProperties: []*openapi.NamedSchemaOrReference{},
		},
	}

	var extDocument *openapi.Document
	err := g.getDocumentOption(&extDocument)
	if err != nil {
		fmt.Printf("Error getting document option: %s\n", err)
		return nil
	}
	if extDocument != nil {
		utils.MergeStructs(d, extDocument)
	}

	g.addPathsToDocument(d, g.ast.Services)

	for len(g.requiredSchemas) > 0 {
		count := len(g.requiredSchemas)
		g.addSchemasForStructsToDocument(d, g.ast.GetStructLikes())
		g.requiredSchemas = g.requiredSchemas[count:len(g.requiredSchemas)]
	}

	if len(d.Tags) == 1 {
		if d.Info.Title == "" && d.Tags[0].Name != "" {
			d.Info.Title = d.Tags[0].Name + " API"
		}
		if d.Info.Description == "" {
			d.Info.Description = d.Tags[0].Description
		}
		d.Tags[0].Description = ""
	}

	var allServers []string

	// If paths methods has servers, but they're all the same, then move servers to path level
	for _, path := range d.Paths.Path {
		var servers []string
		// Only 1 server will ever be set, per method, by the generator
		if path.Value.Get != nil && len(path.Value.Get.Servers) == 1 {
			servers = utils.AppendUnique(servers, path.Value.Get.Servers[0].URL)
			allServers = utils.AppendUnique(allServers, path.Value.Get.Servers[0].URL)
		}
		if path.Value.Post != nil && len(path.Value.Post.Servers) == 1 {
			servers = utils.AppendUnique(servers, path.Value.Post.Servers[0].URL)
			allServers = utils.AppendUnique(allServers, path.Value.Post.Servers[0].URL)
		}
		if path.Value.Put != nil && len(path.Value.Put.Servers) == 1 {
			servers = utils.AppendUnique(servers, path.Value.Put.Servers[0].URL)
			allServers = utils.AppendUnique(allServers, path.Value.Put.Servers[0].URL)
		}
		if path.Value.Delete != nil && len(path.Value.Delete.Servers) == 1 {
			servers = utils.AppendUnique(servers, path.Value.Delete.Servers[0].URL)
			allServers = utils.AppendUnique(allServers, path.Value.Delete.Servers[0].URL)
		}
		if path.Value.Patch != nil && len(path.Value.Patch.Servers) == 1 {
			servers = utils.AppendUnique(servers, path.Value.Patch.Servers[0].URL)
			allServers = utils.AppendUnique(allServers, path.Value.Patch.Servers[0].URL)
		}

		if len(servers) == 1 {
			path.Value.Servers = []*openapi.Server{{URL: servers[0]}}

			if path.Value.Get != nil {
				path.Value.Get.Servers = nil
			}
			if path.Value.Post != nil {
				path.Value.Post.Servers = nil
			}
			if path.Value.Put != nil {
				path.Value.Put.Servers = nil
			}
			if path.Value.Delete != nil {
				path.Value.Delete.Servers = nil
			}
			if path.Value.Patch != nil {
				path.Value.Patch.Servers = nil
			}
		}
	}

	// Set all servers on API level
	if len(allServers) > 0 {
		d.Servers = []*openapi.Server{}
		for _, server := range allServers {
			d.Servers = append(d.Servers, &openapi.Server{URL: server})
		}
	}

	// If there is only 1 server, we can safely remove all path level servers
	if len(allServers) == 1 {
		for _, path := range d.Paths.Path {
			path.Value.Servers = nil
		}
	}

	{
		pairs := d.Tags
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Name < pairs[j].Name
		})
		d.Tags = pairs
	}

	{
		pairs := d.Paths.Path
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Name < pairs[j].Name
		})
		d.Paths.Path = pairs
	}

	{
		pairs := d.Components.Schemas.AdditionalProperties
		sort.Slice(pairs, func(i, j int) bool {
			return pairs[i].Name < pairs[j].Name
		})
		d.Components.Schemas.AdditionalProperties = pairs
	}

	bytes, err := d.YAMLValue("Generated with thrift-gen-http-swagger\n" + infoURL)
	if err != nil {
		fmt.Printf("Error converting to yaml: %s\n", err)
		return nil
	}
	filePath := filepath.Clean(arguments.OutputDir)
	filePath = filepath.Join(filePath, DefaultOutputFile)
	var ret []*plugin.Generated
	ret = append(ret, &plugin.Generated{
		Content: string(bytes),
		Name:    &filePath,
	})

	return ret
}

func (g *OpenAPIGenerator) getDocumentOption(obj interface{}) error {
	serviceOrStruct, name := g.getDocumentAnnotationInWhichServiceOrStruct()
	if serviceOrStruct == DocumentOptionServiceType {
		serviceDesc := g.fileDesc.GetServiceDescriptor(name)
		err := utils.ParseServiceOption(serviceDesc, OpenapiDocument, obj)
		if err != nil {
			return err
		}
	} else if serviceOrStruct == DocumentOptionStructType {
		structDesc := g.fileDesc.GetStructDescriptor(name)
		err := utils.ParseStructOption(structDesc, OpenapiDocument, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *OpenAPIGenerator) addPathsToDocument(d *openapi.Document, services []*parser.Service) {
	for _, s := range services {
		annotationsCount := 0
		for _, f := range s.Functions {
			comment := g.filterCommentString(f.ReservedComments)
			operationID := s.GetName() + "_" + f.GetName()
			rs := utils.GetAnnotations(f.Annotations, HttpMethodAnnotations)
			if len(rs) == 0 {
				continue
			}

			var inputDesc *thrift_reflection.StructDescriptor
			if len(f.Arguments) >= 1 {
				if len(f.Arguments) > 1 {
					logs.Warnf("function '%s' has more than one argument, but only the first can be used in hertz now", f.GetName())
				}
				inputDesc = g.fileDesc.GetStructDescriptor(f.GetArguments()[0].GetType().GetName())
			}
			outputDesc := g.fileDesc.GetStructDescriptor(f.GetFunctionType().GetName())
			for methodName, path := range rs {
				if methodName != "" {
					annotationsCount++
					var host string
					hostOrNil := utils.GetAnnotation(f.Annotations, ApiBaseURL)

					if len(hostOrNil) != 0 {
						host = utils.GetAnnotation(f.Annotations, ApiBaseURL)[0]
					}

					if host == "" {
						hostOrNil = utils.GetAnnotation(s.Annotations, ApiBaseDomain)
						if len(hostOrNil) != 0 {
							host = utils.GetAnnotation(s.Annotations, ApiBaseDomain)[0]
						}
					}

					op, path2 := g.buildOperation(d, methodName, comment, operationID, s.GetName(), path[0], host, inputDesc, outputDesc)
					methodDesc := g.fileDesc.GetMethodDescriptor(s.GetName(), f.GetName())
					newOp := &openapi.Operation{}
					err := utils.ParseMethodOption(methodDesc, OpenapiOperation, &newOp)
					if err != nil {
						logs.Errorf("Error parsing method option: %s", err)
					}
					utils.MergeStructs(op, newOp)
					g.addOperationToDocument(d, op, path2, methodName)
				}
			}
		}
		if annotationsCount > 0 {
			comment := g.filterCommentString(s.ReservedComments)
			d.Tags = append(d.Tags, &openapi.Tag{Name: s.GetName(), Description: comment})
		}
	}
}

func (g *OpenAPIGenerator) buildOperation(
	d *openapi.Document,
	methodName string,
	description string,
	operationID string,
	tagName string,
	path string,
	host string,
	inputDesc *thrift_reflection.StructDescriptor,
	outputDesc *thrift_reflection.StructDescriptor,
) (*openapi.Operation, string) {
	// Parameters array to hold all parameter objects
	var parameters []*openapi.ParameterOrReference

	for _, v := range inputDesc.GetFields() {
		var paramName, paramIn, paramDesc string
		var fieldSchema *openapi.SchemaOrReference
		required := false

		extOrNil := v.Annotations[ApiQuery]
		if len(extOrNil) > 0 {
			if ext := v.Annotations[ApiQuery][0]; ext != "" {
				paramIn = ParameterInQuery
				paramName = ext
				paramDesc = g.filterCommentString(v.Comments)
				fieldSchema = g.schemaOrReferenceForField(v.Type)
				extPropertyOrNil := v.Annotations[OpenapiProperty]
				if len(extPropertyOrNil) > 0 {
					newFieldSchema := &openapi.Schema{}
					err := utils.ParseFieldOption(v, OpenapiProperty, &newFieldSchema)
					if err != nil {
						logs.Errorf("Error parsing field option: %s", err)
					}
					utils.MergeStructs(fieldSchema.Schema, newFieldSchema)
				}
			}
		}
		extOrNil = v.Annotations[ApiPath]
		if len(extOrNil) > 0 {
			if ext := v.Annotations[ApiPath][0]; ext != "" {
				paramIn = ParameterInPath
				paramName = ext
				paramDesc = g.filterCommentString(v.Comments)
				fieldSchema = g.schemaOrReferenceForField(v.Type)
				extPropertyOrNil := v.Annotations[OpenapiProperty]
				if len(extPropertyOrNil) > 0 {
					newFieldSchema := &openapi.Schema{}
					err := utils.ParseFieldOption(v, OpenapiProperty, &newFieldSchema)
					if err != nil {
						logs.Errorf("Error parsing field option: %s", err)
					}
					utils.MergeStructs(fieldSchema.Schema, newFieldSchema)
				}
				required = true
			}
		}
		extOrNil = v.Annotations[ApiCookie]
		if len(extOrNil) > 0 {
			if ext := v.Annotations[ApiCookie][0]; ext != "" {
				paramIn = ParameterInCookie
				paramName = ext
				paramDesc = g.filterCommentString(v.Comments)
				fieldSchema = g.schemaOrReferenceForField(v.Type)
				extPropertyOrNil := v.Annotations[OpenapiProperty]
				if len(extPropertyOrNil) > 0 {
					newFieldSchema := &openapi.Schema{}
					err := utils.ParseFieldOption(v, OpenapiProperty, &newFieldSchema)
					if err != nil {
						logs.Errorf("Error parsing field option: %s", err)
					}
					utils.MergeStructs(fieldSchema.Schema, newFieldSchema)
				}
			}
		}
		extOrNil = v.Annotations[ApiHeader]
		if len(extOrNil) > 0 {
			if ext := v.Annotations[ApiHeader][0]; ext != "" {
				paramIn = ParameterInHeader
				paramName = ext
				paramDesc = g.filterCommentString(v.Comments)
				fieldSchema = g.schemaOrReferenceForField(v.Type)
				extPropertyOrNil := v.Annotations[OpenapiProperty]
				if len(extPropertyOrNil) > 0 {
					newFieldSchema := &openapi.Schema{}
					err := utils.ParseFieldOption(v, OpenapiProperty, &newFieldSchema)
					if err != nil {
						logs.Errorf("Error parsing field option: %s", err)
					}
					utils.MergeStructs(fieldSchema.Schema, newFieldSchema)
				}
			}
		}

		parameter := &openapi.Parameter{
			Name:        paramName,
			In:          paramIn,
			Description: paramDesc,
			Required:    required,
			Schema:      fieldSchema,
		}

		var extParameter *openapi.Parameter
		err := utils.ParseFieldOption(v, OpenapiParameter, &extParameter)
		if err != nil {
			logs.Errorf("Error parsing field option: %s", err)
		}
		utils.MergeStructs(parameter, extParameter)

		// Append the parameter to the parameters array if it was set
		if paramName != "" && paramIn != "" {
			parameters = append(parameters, &openapi.ParameterOrReference{
				Parameter: parameter,
			})
		}
	}

	var RequestBody *openapi.RequestBodyOrReference
	if methodName != "GET" && methodName != "HEAD" && methodName != "DELETE" {
		bodySchema := g.getSchemaByOption(inputDesc, ApiBody)
		formSchema := g.getSchemaByOption(inputDesc, ApiForm)
		rawBodySchema := g.getSchemaByOption(inputDesc, ApiRawBody)

		var additionalProperties []*openapi.NamedMediaType
		if len(bodySchema.Properties.AdditionalProperties) > 0 {
			additionalProperties = append(additionalProperties, &openapi.NamedMediaType{
				Name: ContentTypeJSON,
				Value: &openapi.MediaType{
					Schema: &openapi.SchemaOrReference{
						Schema: bodySchema,
					},
				},
			})
		}

		if len(formSchema.Properties.AdditionalProperties) > 0 {
			additionalProperties = append(additionalProperties, &openapi.NamedMediaType{
				Name: ContentTypeFormMultipart,
				Value: &openapi.MediaType{
					Schema: &openapi.SchemaOrReference{
						Schema: formSchema,
					},
				},
			})

			additionalProperties = append(additionalProperties, &openapi.NamedMediaType{
				Name: ContentTypeFormURLEncoded,
				Value: &openapi.MediaType{
					Schema: &openapi.SchemaOrReference{
						Schema: formSchema,
					},
				},
			})
		}

		if len(rawBodySchema.Properties.AdditionalProperties) > 0 {
			additionalProperties = append(additionalProperties, &openapi.NamedMediaType{
				Name: ContentTypeRawBody,
				Value: &openapi.MediaType{
					Schema: &openapi.SchemaOrReference{
						Schema: rawBodySchema,
					},
				},
			})
		}

		if len(additionalProperties) > 0 {
			RequestBody = &openapi.RequestBodyOrReference{
				RequestBody: &openapi.RequestBody{
					Description: g.filterCommentString(inputDesc.Comments),
					Content: &openapi.MediaTypes{
						AdditionalProperties: additionalProperties,
					},
				},
			}
		}

	}

	name, header, content := g.getResponseForStruct(d, outputDesc)
	desc := g.filterCommentString(outputDesc.Comments)

	if desc == "" {
		desc = DefaultResponseDesc
	}

	var headerOrEmpty *openapi.HeadersOrReferences

	if len(header.AdditionalProperties) != 0 {
		headerOrEmpty = header
	}

	var contentOrEmpty *openapi.MediaTypes

	if len(content.AdditionalProperties) != 0 {
		contentOrEmpty = content
	}

	var responses *openapi.Responses
	if headerOrEmpty != nil || contentOrEmpty != nil {
		responses = &openapi.Responses{
			ResponseOrReference: []*openapi.NamedResponseOrReference{
				{
					Name: name,
					Value: &openapi.ResponseOrReference{
						Response: &openapi.Response{
							Description: desc,
							Headers:     headerOrEmpty,
							Content:     contentOrEmpty,
						},
					},
				},
			},
		}
	}

	re := regexp.MustCompile(`:(\w+)`)
	path = re.ReplaceAllString(path, `{$1}`)

	op := &openapi.Operation{
		Tags:        []string{tagName},
		Description: description,
		OperationID: operationID,
		Parameters:  parameters,
		Responses:   responses,
		RequestBody: RequestBody,
	}

	if host != "" {
		if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
			host = "http://" + host
		}
		op.Servers = append(op.Servers, &openapi.Server{URL: host})
	}

	return op, path
}

func (g *OpenAPIGenerator) getDocumentAnnotationInWhichServiceOrStruct() (string, string) {
	var ret string
	for _, s := range g.ast.Services {
		v := s.Annotations.Get(OpenapiDocument)
		if len(v) > 0 {
			ret = s.GetName()
			return DocumentOptionServiceType, ret
		}
	}
	for _, s := range g.ast.Structs {
		v := s.Annotations.Get(OpenapiDocument)
		if len(v) > 0 {
			ret = s.GetName()
			return DocumentOptionStructType, ret
		}
	}
	return "", ret
}

func (g *OpenAPIGenerator) getResponseForStruct(d *openapi.Document, desc *thrift_reflection.StructDescriptor) (string, *openapi.HeadersOrReferences, *openapi.MediaTypes) {
	headers := &openapi.HeadersOrReferences{AdditionalProperties: []*openapi.NamedHeaderOrReference{}}

	for _, field := range desc.Fields {
		if len(field.Annotations[ApiHeader]) < 1 {
			continue
		}
		if ext := field.Annotations[ApiHeader][0]; ext != "" {
			headerName := ext
			header := &openapi.Header{
				Description: g.filterCommentString(field.Comments),
				Schema:      g.schemaOrReferenceForField(field.Type),
			}
			headers.AdditionalProperties = append(headers.AdditionalProperties, &openapi.NamedHeaderOrReference{
				Name: headerName,
				Value: &openapi.HeaderOrReference{
					Header: header,
				},
			})
		}
	}

	// Get api.body and api.raw_body option schema
	bodySchema := g.getSchemaByOption(desc, ApiBody)
	rawBodySchema := g.getSchemaByOption(desc, ApiRawBody)
	var additionalProperties []*openapi.NamedMediaType

	if len(bodySchema.Properties.AdditionalProperties) > 0 {
		refSchema := &openapi.NamedSchemaOrReference{
			Name:  desc.GetName() + "Body",
			Value: &openapi.SchemaOrReference{Schema: bodySchema},
		}
		ref := "#/components/schemas/" + desc.GetName() + "Body"
		g.addSchemaToDocument(d, refSchema)
		additionalProperties = append(additionalProperties, &openapi.NamedMediaType{
			Name: ContentTypeJSON,
			Value: &openapi.MediaType{
				Schema: &openapi.SchemaOrReference{
					Reference: &openapi.Reference{Xref: ref},
				},
			},
		})
	}

	if len(rawBodySchema.Properties.AdditionalProperties) > 0 {
		refSchema := &openapi.NamedSchemaOrReference{
			Name:  desc.GetName() + "RawBody",
			Value: &openapi.SchemaOrReference{Schema: rawBodySchema},
		}
		ref := "#/components/schemas/" + desc.GetName() + "RawBody"
		g.addSchemaToDocument(d, refSchema)
		additionalProperties = append(additionalProperties, &openapi.NamedMediaType{
			Name: ContentTypeRawBody,
			Value: &openapi.MediaType{
				Schema: &openapi.SchemaOrReference{
					Reference: &openapi.Reference{Xref: ref},
				},
			},
		})
	}

	content := &openapi.MediaTypes{
		AdditionalProperties: additionalProperties,
	}

	return StatusOK, headers, content
}

func (g *OpenAPIGenerator) getSchemaByOption(inputDesc *thrift_reflection.StructDescriptor, option string) *openapi.Schema {
	definitionProperties := &openapi.Properties{
		AdditionalProperties: make([]*openapi.NamedSchemaOrReference, 0),
	}

	var allRequired []string
	var extSchema *openapi.Schema
	err := utils.ParseStructOption(inputDesc, OpenapiSchema, &extSchema)
	if err != nil {
		logs.Errorf("Error parsing struct option: %s", err)
	}
	if extSchema != nil {
		if extSchema.Required != nil {
			allRequired = extSchema.Required
		}
	}

	var required []string
	for _, field := range inputDesc.GetFields() {
		if field.Annotations[option] != nil {
			extName := field.GetName()
			if field.Annotations[option] != nil && field.Annotations[option][0] != "" {
				extName = field.Annotations[option][0]
			}

			if utils.Contains(allRequired, extName) {
				required = append(required, extName)
			}

			// Get the field description from the comments.
			description := g.filterCommentString(field.Comments)
			fieldSchema := g.schemaOrReferenceForField(field.Type)
			if fieldSchema == nil {
				continue
			}

			if fieldSchema.IsSetSchema() {
				fieldSchema.Schema.Description = description
				newFieldSchema := &openapi.Schema{}
				err := utils.ParseFieldOption(field, OpenapiProperty, &newFieldSchema)
				if err != nil {
					logs.Errorf("Error parsing field option: %s", err)
				}
				utils.MergeStructs(fieldSchema.Schema, newFieldSchema)
			}

			definitionProperties.AdditionalProperties = append(
				definitionProperties.AdditionalProperties,
				&openapi.NamedSchemaOrReference{
					Name:  extName,
					Value: fieldSchema,
				},
			)
		}
	}

	schema := &openapi.Schema{
		Type:       SchemaObjectType,
		Properties: definitionProperties,
	}

	if extSchema != nil {
		utils.MergeStructs(schema, extSchema)
	}

	schema.Required = required
	return schema
}

func (g *OpenAPIGenerator) getStructLikeByName(name string) *parser.StructLike {
	for _, s := range g.ast.GetStructLikes() {
		if s.GetName() == name {
			return s
		}
	}
	return nil
}

// filterCommentString removes linter rules from comments.
func (g *OpenAPIGenerator) filterCommentString(str string) string {
	var comments []string
	matches := g.commentPattern.FindAllStringSubmatch(str, -1)

	for _, match := range matches {
		if match[1] != "" {
			// Handle one-line comments
			comments = append(comments, strings.TrimSpace(match[1]))
		} else if match[2] != "" {
			// Handle multiline comments
			multiLineComment := match[2]
			lines := strings.Split(multiLineComment, "\n")

			// Find the minimum indentation level (excluding empty lines)
			minIndent := -1
			for _, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if trimmedLine != "" {
					lineIndent := len(line) - len(strings.TrimLeft(line, " "))
					if minIndent == -1 || lineIndent < minIndent {
						minIndent = lineIndent
					}
				}
			}

			// Remove the minimum indentation and any leading '*' from each line
			for i, line := range lines {
				if minIndent > 0 && len(line) >= minIndent {
					line = line[minIndent:]
				}
				lines[i] = strings.TrimPrefix(line, "*")
			}

			// Remove leading and trailing empty lines from the comment block
			comments = append(comments, strings.TrimSpace(strings.Join(lines, "\n")))
		}
	}

	return strings.Join(comments, "\n")
}

func (g *OpenAPIGenerator) addSchemasForStructsToDocument(d *openapi.Document, structs []*parser.StructLike) {
	// Handle nested structs
	for _, s := range structs {
		var sls []*parser.StructLike
		for _, f := range s.GetFields() {
			if f.GetType().GetCategory().IsStruct() {
				sls = append(sls, g.getStructLikeByName(f.GetType().GetName()))
			}
		}
		g.addSchemasForStructsToDocument(d, sls)

		schemaName := s.GetName()
		// Only generate this if we need it and haven't already generated it.
		if !utils.Contains(g.requiredSchemas, schemaName) ||
			utils.Contains(g.generatedSchemas, schemaName) {
			continue
		}

		structDesc := g.fileDesc.GetStructDescriptor(s.GetName())

		// Get the description from the comments.
		messageDescription := g.filterCommentString(structDesc.Comments)

		// Build an array holding the fields of the message.
		definitionProperties := &openapi.Properties{
			AdditionalProperties: make([]*openapi.NamedSchemaOrReference, 0),
		}

		for _, field := range structDesc.Fields {
			// Get the field description from the comments.
			description := g.filterCommentString(field.Comments)
			fieldSchema := g.schemaOrReferenceForField(field.Type)
			if fieldSchema == nil {
				continue
			}

			if fieldSchema.IsSetSchema() {
				fieldSchema.Schema.Description = description
				newFieldSchema := &openapi.Schema{}
				err := utils.ParseFieldOption(field, OpenapiProperty, &newFieldSchema)
				if err != nil {
					logs.Errorf("Error parsing field option: %s", err)
				}
				utils.MergeStructs(fieldSchema.Schema, newFieldSchema)
			}

			extName := field.GetName()
			options := []string{ApiHeader, ApiBody, ApiForm, ApiRawBody}
			for _, option := range options {
				if field.Annotations[option] != nil && field.Annotations[option][0] != "" {
					extName = field.Annotations[option][0]
				}
			}

			definitionProperties.AdditionalProperties = append(
				definitionProperties.AdditionalProperties,
				&openapi.NamedSchemaOrReference{
					Name:  extName,
					Value: fieldSchema,
				},
			)
		}

		schema := &openapi.Schema{
			Type:        SchemaObjectType,
			Description: messageDescription,
			Properties:  definitionProperties,
		}

		var extSchema *openapi.Schema
		err := utils.ParseStructOption(structDesc, OpenapiSchema, &extSchema)
		if err != nil {
			logs.Errorf("Error parsing struct option: %s", err)
		}
		if extSchema != nil {
			utils.MergeStructs(schema, extSchema)
		}

		// Add the schema to the components.schema list.
		g.addSchemaToDocument(d, &openapi.NamedSchemaOrReference{
			Name: schemaName,
			Value: &openapi.SchemaOrReference{
				Schema: schema,
			},
		})
	}
}

// addSchemaToDocument adds the schema to the document if required
func (g *OpenAPIGenerator) addSchemaToDocument(d *openapi.Document, schema *openapi.NamedSchemaOrReference) {
	if utils.Contains(g.generatedSchemas, schema.Name) {
		return
	}
	g.generatedSchemas = append(g.generatedSchemas, schema.Name)
	d.Components.Schemas.AdditionalProperties = append(d.Components.Schemas.AdditionalProperties, schema)
}

func (g *OpenAPIGenerator) addOperationToDocument(d *openapi.Document, op *openapi.Operation, path, methodName string) {
	var selectedPathItem *openapi.NamedPathItem
	for _, namedPathItem := range d.Paths.Path {
		if namedPathItem.Name == path {
			selectedPathItem = namedPathItem
			break
		}
	}
	// If we get here, we need to create a path item.
	if selectedPathItem == nil {
		selectedPathItem = &openapi.NamedPathItem{Name: path, Value: &openapi.PathItem{}}
		d.Paths.Path = append(d.Paths.Path, selectedPathItem)
	}
	// Set the operation on the specified method.
	switch methodName {
	case "GET":
		selectedPathItem.Value.Get = op
	case "POST":
		selectedPathItem.Value.Post = op
	case "PUT":
		selectedPathItem.Value.Put = op
	case "DELETE":
		selectedPathItem.Value.Delete = op
	case "PATCH":
		selectedPathItem.Value.Patch = op
	case "OPTIONS":
		selectedPathItem.Value.Options = op
	case "HEAD":
		selectedPathItem.Value.Head = op
	}
}

func (g *OpenAPIGenerator) schemaReferenceForMessage(message *thrift_reflection.StructDescriptor) string {
	schemaName := message.GetName()
	if !utils.Contains(g.requiredSchemas, schemaName) {
		g.requiredSchemas = append(g.requiredSchemas, schemaName)
	}
	return "#/components/schemas/" + schemaName
}

func (g *OpenAPIGenerator) schemaOrReferenceForField(fieldType *thrift_reflection.TypeDescriptor) *openapi.SchemaOrReference {
	var kindSchema *openapi.SchemaOrReference

	switch {
	case fieldType.IsStruct():
		structDesc, err := fieldType.GetStructDescriptor()
		if err != nil {
			logs.Errorf("Error getting struct descriptor: %s", err)
			return nil
		}
		ref := g.schemaReferenceForMessage(structDesc)
		kindSchema = &openapi.SchemaOrReference{
			Reference: &openapi.Reference{Xref: ref},
		}

	case fieldType.IsMap():
		valueSchema := g.schemaOrReferenceForField(fieldType.GetValueType())
		kindSchema = &openapi.SchemaOrReference{
			Schema: &openapi.Schema{
				Type: SchemaObjectType,
				AdditionalProperties: &openapi.AdditionalPropertiesItem{
					SchemaOrReference: valueSchema,
				},
			},
		}

	case fieldType.IsList():
		itemSchema := g.schemaOrReferenceForField(fieldType.GetValueType())
		kindSchema = &openapi.SchemaOrReference{
			Schema: &openapi.Schema{
				Type: "array",
				Items: &openapi.ItemsItem{
					SchemaOrReference: []*openapi.SchemaOrReference{itemSchema},
				},
			},
		}

	default:
		kindSchema = &openapi.SchemaOrReference{Schema: &openapi.Schema{}}
		switch fieldType.GetName() {
		case "string":
			kindSchema.Schema.Type = "string"
		case "binary":
			kindSchema.Schema.Type = "string"
			kindSchema.Schema.Format = "binary"
		case "bool":
			kindSchema.Schema.Type = "boolean"
		case "byte":
			kindSchema.Schema.Type = "string"
			kindSchema.Schema.Format = "byte"
		case "double":
			kindSchema.Schema.Type = "number"
			kindSchema.Schema.Format = "double"
		case "i8":
			kindSchema.Schema.Type = "integer"
			kindSchema.Schema.Format = "int8"
		case "i16":
			kindSchema.Schema.Type = "integer"
			kindSchema.Schema.Format = "int16"
		case "i32":
			kindSchema.Schema.Type = "integer"
			kindSchema.Schema.Format = "int32"
		case "i64":
			kindSchema.Schema.Type = "integer"
			kindSchema.Schema.Format = "int64"
		}
	}

	return kindSchema
}

const (
	ApiGet           = "api.get"
	ApiPost          = "api.post"
	ApiPut           = "api.put"
	ApiPatch         = "api.patch"
	ApiDelete        = "api.delete"
	ApiOptions       = "api.options"
	ApiHEAD          = "api.head"
	ApiAny           = "api.any"
	ApiQuery         = "api.query"
	ApiForm          = "api.form"
	ApiPath          = "api.path"
	ApiHeader        = "api.header"
	ApiCookie        = "api.cookie"
	ApiBody          = "api.body"
	ApiRawBody       = "api.raw_body"
	ApiBaseDomain    = "api.base_domain"
	ApiBaseURL       = "api.baseurl"
	OpenapiOperation = "openapi.operation"
	OpenapiProperty  = "openapi.property"
	OpenapiSchema    = "openapi.schema"
	OpenapiParameter = "openapi.parameter"
	OpenapiDocument  = "openapi.document"
)

var HttpMethodAnnotations = map[string]string{
	ApiGet:     "GET",
	ApiPost:    "POST",
	ApiPut:     "PUT",
	ApiPatch:   "PATCH",
	ApiDelete:  "DELETE",
	ApiOptions: "OPTIONS",
	ApiHEAD:    "HEAD",
	ApiAny:     "ANY",
}

const (
	OpenAPIVersion     = "3.0.3"
	DefaultOutputFile  = "openapi.yaml"
	DefaultInfoTitle   = "API generated by thrift-gen-http-swagger"
	DefaultInfoDesc    = "API description"
	DefaultInfoVersion = "0.0.1"
	infoURL            = "https://github.com/hertz-contrib/swagger-generate/thrift-gen-http-swagger"

	DefaultResponseDesc = "Successful response"
	StatusOK            = "200"
	SchemaObjectType    = "object"

	ContentTypeJSON           = "application/json"
	ContentTypeFormMultipart  = "multipart/form-data"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	ContentTypeRawBody        = "text/plain"

	ParameterInQuery  = "query"
	ParameterInHeader = "header"
	ParameterInPath   = "path"
	ParameterInCookie = "cookie"

	DocumentOptionServiceType = "service"
	DocumentOptionStructType  = "struct"
)