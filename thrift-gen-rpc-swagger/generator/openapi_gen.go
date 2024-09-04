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
	"github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger/args"
	openapi "github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger/thrift"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger/utils"
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
		logs.Errorf("Error getting document option: %s", err)
		return nil
	}
	if extDocument != nil {
		err := utils.MergeStructs(d, extDocument)
		if err != nil {
			logs.Errorf("Error merging document option: %s", err)
			return nil
		}
	}

	g.addPathsToDocument(d, g.ast.Services)

	for len(g.requiredSchemas) > 0 {
		count := len(g.requiredSchemas)
		g.addSchemasForStructsToDocument(d, g.ast.GetStructLikes())
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
			servers = utils.AppendUnique(servers, path.Value.Post.Servers[0].URL)
			allServers = utils.AppendUnique(allServers, path.Value.Post.Servers[0].URL)
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
			{URL: DefaultServerURL},
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

	bytes, err := d.YAMLValue("Generated with thrift-gen-rpc-swagger\n" + infoURL)
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
			path := "/" + f.GetName()

			var inputDesc *thrift_reflection.StructDescriptor
			if len(f.Arguments) >= 1 {
				if len(f.Arguments) > 1 {
					logs.Warnf("function '%s' has more than one argument, but only the first can be used in hertz now", f.GetName())
				}
				inputDesc = g.fileDesc.GetStructDescriptor(f.GetArguments()[0].GetType().GetName())
			}
			outputDesc := g.fileDesc.GetStructDescriptor(f.GetFunctionType().GetName())
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

			op, path2 := g.buildOperation(d, comment, operationID, s.GetName(), path, host, inputDesc, outputDesc)
			methodDesc := g.fileDesc.GetMethodDescriptor(s.GetName(), f.GetName())
			newOp := &openapi.Operation{}
			err := utils.ParseMethodOption(methodDesc, OpenapiOperation, &newOp)
			if err != nil {
				logs.Errorf("Error parsing method option: %s", err)
			}
			err = utils.MergeStructs(op, newOp)
			if err != nil {
				logs.Errorf("Error merging method option: %s", err)
			}
			g.addOperationToDocument(d, op, path2)
		}
		if annotationsCount > 0 {
			comment := g.filterCommentString(s.ReservedComments)
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
) (*openapi.Operation, string) {
	// Parameters array to hold all parameter objects
	var parameters []*openapi.ParameterOrReference

	fieldSchema := &openapi.SchemaOrReference{
		Schema: &openapi.Schema{
			Type: SchemaObjectType,
		},
	}
	parameter := &openapi.Parameter{
		Name:        ParameterNameTTHeader,
		In:          ParameterInQuery,
		Description: ParameterDescription,
		Required:    false,
		Schema:      fieldSchema,
	}
	parameters = append(parameters, &openapi.ParameterOrReference{
		Parameter: parameter,
	})

	var RequestBody *openapi.RequestBodyOrReference
	bodySchema := g.getSchemaByOption(inputDesc)

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
		if !strings.HasPrefix(host, URLDefaultPrefixHTTP) && !strings.HasPrefix(host, URLDefaultPrefixHTTPS) {
			host = URLDefaultPrefixHTTP + host
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

	// Get api.body and api.raw_body option schema
	bodySchema := g.getSchemaByOption(desc)

	var additionalProperties []*openapi.NamedMediaType

	if len(bodySchema.Properties.AdditionalProperties) > 0 {
		refSchema := &openapi.NamedSchemaOrReference{
			Name:  desc.GetName(),
			Value: &openapi.SchemaOrReference{Schema: bodySchema},
		}
		ref := ComponentSchemaPrefix + desc.GetName()
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

	content := &openapi.MediaTypes{
		AdditionalProperties: additionalProperties,
	}

	return StatusOK, headers, content
}

func (g *OpenAPIGenerator) getSchemaByOption(inputDesc *thrift_reflection.StructDescriptor) *openapi.Schema {
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
		extName := field.GetName()

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
			err = utils.MergeStructs(fieldSchema.Schema, newFieldSchema)
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
		Type:       SchemaObjectType,
		Properties: definitionProperties,
	}

	if extSchema != nil {
		err := utils.MergeStructs(schema, extSchema)
		if err != nil {
			logs.Errorf("Error merging struct option: %s", err)
		}
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
		var comment string
		if match[1] != "" {
			// One-line comment
			comment = strings.TrimSpace(match[1])
		} else if match[2] != "" {
			// Multiline comment
			multiLineComment := match[2]
			lines := strings.Split(multiLineComment, "\n")
			for i, line := range lines {
				lines[i] = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*"))
			}
			comment = strings.Join(lines, "\n")
		}
		comments = append(comments, comment)
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
				err = utils.MergeStructs(fieldSchema.Schema, newFieldSchema)
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
			err = utils.MergeStructs(schema, extSchema)
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
	if utils.Contains(g.generatedSchemas, schema.Name) {
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
	if !utils.Contains(g.requiredSchemas, schemaName) {
		g.requiredSchemas = append(g.requiredSchemas, schemaName)
	}
	return ComponentSchemaPrefix + schemaName
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
	ApiBaseDomain    = "api.base_domain"
	ApiBaseURL       = "api.baseurl"
	OpenapiOperation = "openapi.operation"
	OpenapiProperty  = "openapi.property"
	OpenapiSchema    = "openapi.schema"
	OpenapiDocument  = "openapi.document"
)

const (
	OpenAPIVersion        = "3.0.3"
	DefaultOutputFile     = "openapi.yaml"
	DefaultServerURL      = "http://127.0.0.1:8080"
	DefaultInfoTitle      = "API generated by thrift-gen-rpc-swagger"
	DefaultInfoDesc       = "API description"
	DefaultInfoVersion    = "0.0.1"
	infoURL               = "https://github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger"
	URLDefaultPrefixHTTP  = "http://"
	URLDefaultPrefixHTTPS = "https://"

	DefaultResponseDesc = "Successful response"
	StatusOK            = "200"

	ContentTypeJSON       = "application/json"
	SchemaObjectType      = "object"
	ComponentSchemaPrefix = "#/components/schemas/"

	ParameterInQuery      = "query"
	ParameterNameTTHeader = "ttheader"
	ParameterDescription  = "metainfo for request"

	DocumentOptionServiceType = "service"
	DocumentOptionStructType  = "struct"
)
