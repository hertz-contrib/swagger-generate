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
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/plugin"
	"github.com/cloudwego/thriftgo/thrift_reflection"
	"github.com/hertz-contrib/swagger-generate/common/consts"
	common "github.com/hertz-contrib/swagger-generate/common/utils"
	openapi "github.com/hertz-contrib/swagger-generate/idl/thrift"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger/args"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger/utils"
)

type OpenAPIGenerator struct {
	fileDesc         *thrift_reflection.FileDescriptor
	ast              *parser.Thrift
	generatedSchemas []string
	requiredSchemas  []string
	requiredTypeDesc []*thrift_reflection.StructDescriptor
}

// NewOpenAPIGenerator creates a new generator for a protoc plugin invocation.
func NewOpenAPIGenerator(ast *parser.Thrift) *OpenAPIGenerator {
	_, fileDesc := thrift_reflection.RegisterAST(ast)
	return &OpenAPIGenerator{
		fileDesc:         fileDesc,
		ast:              ast,
		generatedSchemas: make([]string, 0),
	}
}

func (g *OpenAPIGenerator) BuildDocument(arguments *args.Arguments) []*plugin.Generated {
	d := &openapi.Document{}

	version := consts.OpenAPIVersion
	d.Openapi = version
	d.Info = &openapi.Info{
		Title:       consts.DefaultInfoTitle + consts.PluginNameThriftRpcSwagger,
		Description: consts.DefaultInfoDesc,
		Version:     consts.DefaultInfoVersion,
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
		logs.Errorf("Error getting document option: %s", err)
		return nil
	}
	if extDocument != nil {
		err := common.MergeStructs(d, extDocument)
		if err != nil {
			logs.Errorf("Error merging document option: %s", err)
			return nil
		}
	}

	g.addPathsToDocument(d, g.fileDesc.GetServices())

	for len(g.requiredSchemas) > 0 {
		count := len(g.requiredSchemas)
		g.addSchemasForStructsToDocument(d, g.requiredTypeDesc)
		g.requiredSchemas = g.requiredSchemas[count:len(g.requiredSchemas)]
	}

	// If there is only 1 service, then use it's title for the
	// document, if the document is missing it.
	if len(d.Tags) == 1 {
		if d.Info.Title == "" && d.Tags[0].Name != "" {
			d.Info.Title = d.Tags[0].Name + " API"
		}
		if d.Info.Description == "" {
			d.Info.Description = d.Tags[0].Description
		}
	}

	var allServers []string

	// If paths methods has servers, but they're all the same, then move servers to path level
	for _, path := range d.Paths.Path {
		var servers []string
		// Only 1 server will ever be set, per method, by the generator
		if path.Value.Post != nil && len(path.Value.Post.Servers) == 1 {
			servers = common.AppendUnique(servers, path.Value.Post.Servers[0].URL)
			allServers = common.AppendUnique(allServers, path.Value.Post.Servers[0].URL)
		}

		if len(servers) == 1 {
			path.Value.Servers = []*openapi.Server{{URL: servers[0]}}

			if path.Value.Post != nil {
				path.Value.Post.Servers = nil
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

	// If there are no servers, add a default one
	if len(allServers) == 0 {
		d.Servers = []*openapi.Server{
			{URL: consts.DefaultServerURL},
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

	bytes, err := d.YAMLValue("Generated with " + consts.PluginNameThriftRpcSwagger + "\n" + consts.InfoURL + consts.PluginNameThriftRpcSwagger)
	if err != nil {
		logs.Errorf("Error converting to yaml: %s", err)
		return nil
	}
	outputDir := arguments.OutputDir
	if outputDir == "" {
		outputDir = consts.DefaultOutputDir
	}
	filePath := filepath.Join(outputDir, consts.DefaultOutputYamlFile)
	var ret []*plugin.Generated
	ret = append(ret, &plugin.Generated{
		Content: string(bytes),
		Name:    &filePath,
	})

	return ret
}

func (g *OpenAPIGenerator) getDocumentOption(obj interface{}) error {
	serviceOrStruct, name := g.getDocumentAnnotationInWhichServiceOrStruct()
	if serviceOrStruct == consts.DocumentOptionServiceType {
		serviceDesc := g.fileDesc.GetServiceDescriptor(name)
		err := utils.ParseServiceOption(serviceDesc, consts.OpenapiDocument, obj)
		if err != nil {
			return err
		}
	} else if serviceOrStruct == consts.DocumentOptionStructType {
		structDesc := g.fileDesc.GetStructDescriptor(name)
		err := utils.ParseStructOption(structDesc, consts.OpenapiDocument, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *OpenAPIGenerator) addPathsToDocument(d *openapi.Document, services []*thrift_reflection.ServiceDescriptor) {
	var err error
	for _, s := range services {
		annotationsCount := 0
		for _, m := range s.GetMethods() {
			var inputDesc, outputDesc, throwDesc *thrift_reflection.StructDescriptor

			if len(m.Args) > 0 {
				if len(m.Args) > 1 {
					logs.Warnf("function '%s' has more than one argument, but only the first can be used in plugin now", m.GetName())
				}
				// TODO: support more argument types
				if m.Args[0].GetType().IsStruct() {
					inputDesc, err = m.Args[0].GetType().GetStructDescriptor()
					if err != nil {
						logs.Errorf("Error getting arguments descriptor: %s", err)
					}
				} else {
					logs.Errorf("now only support struct type for input, but got %s", m.Args[0].GetType().GetName())
				}
			}

			// TODO: support more response types
			if m.Response.IsStruct() {
				outputDesc, err = m.Response.GetStructDescriptor()
				if err != nil {
					logs.Errorf("Error getting response descriptor: %s", err)
				}
			} else if m.Response.Name != "void" {
				logs.Errorf("now only support struct type for output, but got %s", m.Response.Name)
			}

			if len(m.ThrowExceptions) > 0 {
				throwDesc, err = m.ThrowExceptions[0].GetType().GetExceptionDescriptor()
				if err != nil {
					logs.Errorf("Error getting exception descriptor: %s", err)
				}
			}
			var host string

			if urls, ok := m.Annotations[consts.ApiBaseURL]; ok && len(urls) > 0 {
				host = urls[0]
			} else if domains, ok := s.Annotations[consts.ApiBaseDomain]; ok && len(domains) > 0 {
				host = domains[0]
			}

			annotationsCount++
			operationID := s.GetName() + "_" + m.GetName()
			path := "/" + m.GetName()
			comment := g.filterCommentString(m.Comments)

			op, path2 := g.buildOperation(d, comment, operationID, s.GetName(), path, host, inputDesc, outputDesc, throwDesc)

			newOp := &openapi.Operation{}
			err = utils.ParseMethodOption(m, consts.OpenapiOperation, &newOp)
			if err != nil {
				logs.Errorf("Error parsing method option: %s", err)
			}
			err = common.MergeStructs(op, newOp)
			if err != nil {
				logs.Errorf("Error merging method option: %s", err)
			}

			g.addOperationToDocument(d, op, path2)
		}
		if annotationsCount > 0 {
			comment := g.filterCommentString(s.Comments)
			d.Tags = append(d.Tags, &openapi.Tag{Name: s.GetName(), Description: comment})
		}
	}
}

func (g *OpenAPIGenerator) buildOperation(
	d *openapi.Document,
	description string,
	operationID string,
	tagName string,
	path string,
	host string,
	inputDesc *thrift_reflection.StructDescriptor,
	outputDesc *thrift_reflection.StructDescriptor,
	throwDesc *thrift_reflection.StructDescriptor,
) (*openapi.Operation, string) {
	// Parameters array to hold all parameter objects
	var parameters []*openapi.ParameterOrReference

	fieldSchema := &openapi.SchemaOrReference{
		Schema: &openapi.Schema{
			Type: consts.SchemaObjectType,
		},
	}
	parameter := &openapi.Parameter{
		Name:        consts.ParameterNameTTHeader,
		In:          consts.ParameterInQuery,
		Description: consts.ParameterDescription,
		Required:    false,
		Schema:      fieldSchema,
	}
	parameters = append(parameters, &openapi.ParameterOrReference{
		Parameter: parameter,
	})

	var RequestBody *openapi.RequestBodyOrReference

	if inputDesc != nil {
		bodySchema := g.getSchemaByOption(inputDesc)

		var additionalProperties []*openapi.NamedMediaType
		if len(bodySchema.Properties.AdditionalProperties) > 0 {
			refSchema := &openapi.NamedSchemaOrReference{
				Name:  inputDesc.GetName(),
				Value: &openapi.SchemaOrReference{Schema: bodySchema},
			}

			ref := consts.ComponentSchemaPrefix + inputDesc.GetName()

			g.addSchemaToDocument(d, refSchema)

			additionalProperties = append(additionalProperties, &openapi.NamedMediaType{
				Name: consts.ContentTypeJSON,
				Value: &openapi.MediaType{
					Schema: &openapi.SchemaOrReference{
						Reference: &openapi.Reference{Xref: ref},
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

	var (
		desc                    string
		contentOrEmpty          *openapi.MediaTypes
		exceptionDesc           string
		exceptionContentOrEmpty *openapi.MediaTypes
		responses               *openapi.Responses
	)

	if outputDesc != nil {
		name, content := g.getResponseForStruct(d, outputDesc)
		desc = g.filterCommentString(outputDesc.Comments)

		if desc == "" {
			desc = consts.DefaultResponseDesc
		}

		if len(content.AdditionalProperties) != 0 {
			contentOrEmpty = content
		}

		if contentOrEmpty != nil || exceptionContentOrEmpty != nil {
			responses = &openapi.Responses{
				ResponseOrReference: []*openapi.NamedResponseOrReference{
					{
						Name: name,
						Value: &openapi.ResponseOrReference{
							Response: &openapi.Response{
								Description: desc,
								Content:     contentOrEmpty,
							},
						},
					},
				},
			}
		}
	}

	if throwDesc != nil {
		exceptionName, exceptionContent := g.getExceptionForStruct(d, throwDesc)
		exceptionDesc = g.filterCommentString(throwDesc.Comments)

		if exceptionDesc == "" {
			exceptionDesc = consts.DefaultExceptionDesc
		}

		if len(exceptionContent.AdditionalProperties) != 0 {
			exceptionContentOrEmpty = exceptionContent
		}

		if responses == nil {
			responses = &openapi.Responses{
				ResponseOrReference: []*openapi.NamedResponseOrReference{},
			}
		}

		if contentOrEmpty != nil || exceptionContentOrEmpty != nil {
			responses = &openapi.Responses{
				ResponseOrReference: append(responses.ResponseOrReference, &openapi.NamedResponseOrReference{
					Name: exceptionName,
					Value: &openapi.ResponseOrReference{
						Response: &openapi.Response{
							Description: exceptionDesc,
							Content:     exceptionContentOrEmpty,
						},
					},
				}),
			}
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
		if !strings.HasPrefix(host, consts.URLDefaultPrefixHTTP) && !strings.HasPrefix(host, consts.URLDefaultPrefixHTTPS) {
			host = consts.URLDefaultPrefixHTTP + host
		}
		op.Servers = append(op.Servers, &openapi.Server{URL: host})
	}

	return op, path
}

func (g *OpenAPIGenerator) getDocumentAnnotationInWhichServiceOrStruct() (string, string) {
	var ret string
	for _, s := range g.ast.Services {
		v := s.Annotations.Get(consts.OpenapiDocument)
		if len(v) > 0 {
			ret = s.GetName()
			return consts.DocumentOptionServiceType, ret
		}
	}
	for _, s := range g.ast.Structs {
		v := s.Annotations.Get(consts.OpenapiDocument)
		if len(v) > 0 {
			ret = s.GetName()
			return consts.DocumentOptionStructType, ret
		}
	}
	return "", ret
}

func (g *OpenAPIGenerator) getResponseForStruct(d *openapi.Document, desc *thrift_reflection.StructDescriptor) (string, *openapi.MediaTypes) {
	bodySchema := g.getSchemaByOption(desc)

	var additionalProperties []*openapi.NamedMediaType

	if len(bodySchema.Properties.AdditionalProperties) > 0 {
		refSchema := &openapi.NamedSchemaOrReference{
			Name:  desc.GetName(),
			Value: &openapi.SchemaOrReference{Schema: bodySchema},
		}
		ref := consts.ComponentSchemaPrefix + desc.GetName()
		g.addSchemaToDocument(d, refSchema)
		additionalProperties = append(additionalProperties, &openapi.NamedMediaType{
			Name: consts.ContentTypeJSON,
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

	return consts.StatusOK, content
}

func (g *OpenAPIGenerator) getExceptionForStruct(d *openapi.Document, desc *thrift_reflection.StructDescriptor) (string, *openapi.MediaTypes) {
	bodySchema := g.getSchemaByOption(desc)

	var additionalProperties []*openapi.NamedMediaType

	if len(bodySchema.Properties.AdditionalProperties) > 0 {
		refSchema := &openapi.NamedSchemaOrReference{
			Name:  desc.GetName(),
			Value: &openapi.SchemaOrReference{Schema: bodySchema},
		}
		ref := consts.ComponentSchemaPrefix + desc.GetName()
		g.addSchemaToDocument(d, refSchema)
		additionalProperties = append(additionalProperties, &openapi.NamedMediaType{
			Name: consts.ContentTypeJSON,
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

	return consts.StatusBadRequest, content
}

func (g *OpenAPIGenerator) getSchemaByOption(inputDesc *thrift_reflection.StructDescriptor) *openapi.Schema {
	definitionProperties := &openapi.Properties{
		AdditionalProperties: make([]*openapi.NamedSchemaOrReference, 0),
	}

	var allRequired []string
	var extSchema *openapi.Schema
	err := utils.ParseStructOption(inputDesc, consts.OpenapiSchema, &extSchema)
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
		extName := field.GetName()

		if common.Contains(allRequired, extName) {
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
			err := utils.ParseFieldOption(field, consts.OpenapiProperty, &newFieldSchema)
			if err != nil {
				logs.Errorf("Error parsing field option: %s", err)
			}
			err = common.MergeStructs(fieldSchema.Schema, newFieldSchema)
			if err != nil {
				logs.Errorf("Error merging field option: %s", err)
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
		Type:       consts.SchemaObjectType,
		Properties: definitionProperties,
	}

	if extSchema != nil {
		err := common.MergeStructs(schema, extSchema)
		if err != nil {
			logs.Errorf("Error merging struct option: %s", err)
		}
	}

	schema.Required = required
	return schema
}

// filterCommentString removes linter rules from comments.
func (g *OpenAPIGenerator) filterCommentString(str string) string {
	var comments []string
	matches := regexp.MustCompile(consts.CommentPatternRegexp).FindAllStringSubmatch(str, -1)

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

func (g *OpenAPIGenerator) addSchemasForStructsToDocument(d *openapi.Document, structs []*thrift_reflection.StructDescriptor) {
	for _, s := range structs {
		var sls []*thrift_reflection.StructDescriptor
		for _, f := range s.GetFields() {
			fieldType := f.GetType()
			if fieldType == nil {
				logs.Errorf("Warning: field type is nil for field: %s\n", f.GetName())
				continue
			}
			if fieldType.IsStruct() {
				structDesc, _ := fieldType.GetStructDescriptor()
				sls = append(sls, structDesc)
			}
		}
		if len(sls) > 0 {
			g.addSchemasForStructsToDocument(d, sls)
		}

		schemaName := s.GetName()

		// Only generate this if we need it and haven't already generated it.
		if !common.Contains(g.requiredSchemas, schemaName) ||
			common.Contains(g.generatedSchemas, schemaName) {
			continue
		}

		// Get the description from the comments.
		messageDescription := g.filterCommentString(s.Comments)

		// Build an array holding the fields of the message.
		definitionProperties := &openapi.Properties{
			AdditionalProperties: make([]*openapi.NamedSchemaOrReference, 0),
		}

		for _, field := range s.Fields {
			// Get the field description from the comments.
			description := g.filterCommentString(field.Comments)
			fieldSchema := g.schemaOrReferenceForField(field.Type)
			if fieldSchema == nil {
				continue
			}

			if fieldSchema.IsSetSchema() {
				fieldSchema.Schema.Description = description
				newFieldSchema := &openapi.Schema{}
				err := utils.ParseFieldOption(field, consts.OpenapiProperty, &newFieldSchema)
				if err != nil {
					logs.Errorf("Error parsing field option: %s", err)
				}
				err = common.MergeStructs(fieldSchema.Schema, newFieldSchema)
				if err != nil {
					logs.Errorf("Error merging field option: %s", err)
				}
			}

			fName := field.GetName()

			definitionProperties.AdditionalProperties = append(
				definitionProperties.AdditionalProperties,
				&openapi.NamedSchemaOrReference{
					Name:  fName,
					Value: fieldSchema,
				},
			)
		}

		schema := &openapi.Schema{
			Type:        consts.SchemaObjectType,
			Description: messageDescription,
			Properties:  definitionProperties,
		}

		var extSchema *openapi.Schema
		err := utils.ParseStructOption(s, consts.OpenapiSchema, &extSchema)
		if err != nil {
			logs.Errorf("Error parsing struct option: %s", err)
		}
		if extSchema != nil {
			err = common.MergeStructs(schema, extSchema)
			if err != nil {
				logs.Errorf("Error merging struct option: %s", err)
			}
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
	if common.Contains(g.generatedSchemas, schema.Name) {
		return
	}
	g.generatedSchemas = append(g.generatedSchemas, schema.Name)
	d.Components.Schemas.AdditionalProperties = append(d.Components.Schemas.AdditionalProperties, schema)
}

func (g *OpenAPIGenerator) addOperationToDocument(d *openapi.Document, op *openapi.Operation, path string) {
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

	selectedPathItem.Value.Post = op
}

func (g *OpenAPIGenerator) schemaReferenceForMessage(message *thrift_reflection.StructDescriptor) string {
	schemaName := message.GetName()
	if !common.Contains(g.requiredSchemas, schemaName) {
		g.requiredSchemas = append(g.requiredSchemas, schemaName)
		g.requiredTypeDesc = append(g.requiredTypeDesc, message)
	}
	return consts.ComponentSchemaPrefix + schemaName
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
				Type: consts.SchemaObjectType,
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
	case fieldType.IsTypedef():
		typedefDesc, err := fieldType.GetTypedefDescriptor()
		if err != nil {
			logs.Errorf("Error getting typedef descriptor: %s", err)
			return nil
		}
		kindSchema = g.schemaOrReferenceForField(typedefDesc.Type)

	case fieldType.IsEnum():
		enumDesc, err := fieldType.GetEnumDescriptor()
		if err != nil {
			logs.Errorf("Error getting enum descriptor: %s", err)
			return nil
		}
		kindSchema = &openapi.SchemaOrReference{Schema: &openapi.Schema{}}
		kindSchema.Schema.Type = "string"
		kindSchema.Schema.Format = "enum"
		kindSchema.Schema.Enum = make([]*openapi.Any, 0, len(enumDesc.GetValues()))
		for _, v := range enumDesc.GetValues() {
			kindSchema.Schema.Enum = append(kindSchema.Schema.Enum, &openapi.Any{Yaml: v.GetName()})
		}

	case fieldType.IsUnion():
		unionDesc, err := fieldType.GetUnionDescriptor()
		if err != nil {
			logs.Errorf("Error getting union descriptor: %s", err)
			return nil
		}
		kindSchema = &openapi.SchemaOrReference{Schema: &openapi.Schema{}}
		kindSchema.Schema.OneOf = make([]*openapi.SchemaOrReference, 0, len(unionDesc.GetFields()))
		for _, f := range unionDesc.GetFields() {
			fieldSchema := g.schemaOrReferenceForField(f.Type)
			kindSchema.Schema.OneOf = append(kindSchema.Schema.OneOf, fieldSchema)
		}

	case fieldType.IsException():
		logs.Errorf("Error: exception type not supported: %s for field", fieldType.GetName())

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
