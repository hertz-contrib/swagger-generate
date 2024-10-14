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

package converter

import (
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hertz-contrib/swagger-generate/common/consts"
	common "github.com/hertz-contrib/swagger-generate/common/utils"
	"github.com/hertz-contrib/swagger-generate/swagger2idl/protobuf"
	"github.com/hertz-contrib/swagger-generate/swagger2idl/utils"
)

// ProtoConverter struct, used to convert OpenAPI specifications into Proto files
type ProtoConverter struct {
	spec            *openapi3.T
	ProtoFile       *protobuf.ProtoFile
	converterOption *ConvertOption
}

// NewProtoConverter creates and initializes a ProtoConverter
func NewProtoConverter(spec *openapi3.T, option *ConvertOption) *ProtoConverter {
	return &ProtoConverter{
		spec: spec,
		ProtoFile: &protobuf.ProtoFile{
			PackageName: utils.GetPackageName(spec),
			Messages:    []*protobuf.ProtoMessage{},
			Services:    []*protobuf.ProtoService{},
			Enums:       []*protobuf.ProtoEnum{},
			Imports:     []string{},
			Options:     []*protobuf.Option{},
		},
		converterOption: option,
	}
}

// Convert converts the OpenAPI specification to a Proto file
func (c *ProtoConverter) Convert() error {
	// Convert the go Option to Proto
	err := c.addExtensionsToProtoOptions()
	if err != nil {
		return fmt.Errorf("error parsing extensions to proto options: %w", err)
	}

	// Convert tags into Proto services
	c.convertTagsToProtoServices()

	// Convert components into Proto messages
	err = c.convertComponentsToProtoMessages()
	if err != nil {
		return fmt.Errorf("error converting components to proto messages: %w", err)
	}

	// Convert paths into Proto services
	err = c.convertPathsToProtoServices()
	if err != nil {
		return fmt.Errorf("error converting paths to proto services: %w", err)
	}

	if c.converterOption.OpenapiOption {
		c.addOptionsToProto()
	}

	return nil
}

func (c *ProtoConverter) GetIdl() interface{} {
	return c.ProtoFile
}

// convertTagsToProtoServices converts OpenAPI tags into Proto services and stores them in the ProtoFile
func (c *ProtoConverter) convertTagsToProtoServices() {
	tags := c.spec.Tags
	for _, tag := range tags {
		serviceName := common.ToPascaleCase(tag.Name)
		service := &protobuf.ProtoService{
			Name:        serviceName,
			Description: tag.Description,
		}
		c.ProtoFile.Services = append(c.ProtoFile.Services, service)
	}
}

// convertComponentsToProtoMessages converts OpenAPI components into Proto messages and stores them in the ProtoFile
func (c *ProtoConverter) convertComponentsToProtoMessages() error {
	components := c.spec.Components
	if components == nil {
		return nil
	}

	if components.Schemas == nil {
		return nil
	}

	for name, schemaRef := range components.Schemas {
		schema := schemaRef

		if c.converterOption.NamingOption {
			name = common.ToPascaleCase(name)
		}

		protoType, err := c.ConvertSchemaToProtoType(schema, name, nil)
		if err != nil {
			return fmt.Errorf("error converting schema %s: %w", name, err)
		}

		switch v := protoType.(type) {
		case *protobuf.ProtoField:
			message := &protobuf.ProtoMessage{
				Name:   name,
				Fields: []*protobuf.ProtoField{v},
			}

			if c.converterOption.OpenapiOption {
				optionStr := common.StructToOption(schema.Value, "    ")

				schemaOption := &protobuf.Option{
					Name:  consts.OpenapiSchema,
					Value: optionStr,
				}
				message.Options = append(message.Options, schemaOption)
				c.AddProtoImport(consts.OpenapiProtoFile)
			}
			c.addMessageToProto(message)
		case *protobuf.ProtoMessage:
			if c.converterOption.OpenapiOption {
				optionStr := common.StructToOption(schema.Value, "    ")

				schemaOption := &protobuf.Option{
					Name:  consts.OpenapiSchema,
					Value: optionStr,
				}
				v.Options = append(v.Options, schemaOption)
				c.AddProtoImport(consts.OpenapiProtoFile)
			}
			c.addMessageToProto(v)
		case *protobuf.ProtoEnum:
			if c.converterOption.OpenapiOption {
				optionStr := common.StructToOption(schema.Value, "    ")

				schemaOption := &protobuf.Option{
					Name:  consts.OpenapiSchema,
					Value: optionStr,
				}
				v.Options = append(v.Options, schemaOption)
				c.AddProtoImport(consts.OpenapiProtoFile)
			}
			c.addEnumToProto(v)
		case *protobuf.ProtoOneOf:
			message := &protobuf.ProtoMessage{
				Name:   name,
				OneOfs: []*protobuf.ProtoOneOf{v},
			}

			if c.converterOption.OpenapiOption {
				optionStr := common.StructToOption(schema.Value, "    ")

				schemaOption := &protobuf.Option{
					Name:  consts.OpenapiSchema,
					Value: optionStr,
				}
				message.Options = append(message.Options, schemaOption)
				c.AddProtoImport(consts.OpenapiProtoFile)
			}
			c.addMessageToProto(message)
		}
	}
	return nil
}

// convertPathsToProtoServices converts OpenAPI path items into Proto services and stores them in the ProtoFile
func (c *ProtoConverter) convertPathsToProtoServices() error {
	paths := c.spec.Paths
	services, err := c.ConvertPathsToProtoServices(paths)
	if err != nil {
		return fmt.Errorf("error converting paths to proto services: %w", err)
	}

	c.ProtoFile.Services = append(c.ProtoFile.Services, services...)
	return nil
}

// ConvertPathsToProtoServices converts OpenAPI path items into Proto services
func (c *ProtoConverter) ConvertPathsToProtoServices(paths *openapi3.Paths) ([]*protobuf.ProtoService, error) {
	var services []*protobuf.ProtoService

	for path, pathItem := range paths.Map() {
		for method, operation := range pathItem.Operations() {
			serviceName := utils.GetServiceName(operation)
			methodName := utils.GetMethodName(operation, path, method)

			if c.converterOption.NamingOption {
				serviceName = common.ToPascaleCase(serviceName)
				methodName = common.ToPascaleCase(methodName)
			}

			inputMessage, err := c.generateRequestMessage(operation, methodName)
			if err != nil {
				return nil, fmt.Errorf("error generating request message for %s: %w", methodName, err)
			}

			outputMessage, err := c.generateResponseMessage(operation, methodName)
			if err != nil {
				return nil, fmt.Errorf("error generating response message for %s: %w", methodName, err)
			}

			service := c.findOrCreateService(serviceName)

			if !c.methodExistsInService(service, methodName) {
				protoMethod := &protobuf.ProtoMethod{
					Name:   methodName,
					Input:  inputMessage,
					Output: outputMessage,
				}

				if c.converterOption.ApiOption {
					if optionName, ok := MethodToOption[method]; ok {
						option := &protobuf.Option{
							Name:  optionName,
							Value: fmt.Sprintf("%q", utils.ConvertPath(path)),
						}
						protoMethod.Options = append(protoMethod.Options, option)
						c.AddProtoImport(consts.ApiProtoFile)
					}
				}

				if c.converterOption.OpenapiOption {
					optionStr := common.StructToOption(operation, "     ")

					schemaOption := &protobuf.Option{
						Name:  consts.OpenapiOperation,
						Value: optionStr,
					}
					protoMethod.Options = append(protoMethod.Options, schemaOption)
					c.AddProtoImport(consts.OpenapiProtoFile)

				}
				service.Methods = append(service.Methods, protoMethod)
			}
		}
	}

	return services, nil
}

// generateRequestMessage generates a request message for an operation
func (c *ProtoConverter) generateRequestMessage(operation *openapi3.Operation, methodName string) (string, error) {
	messageName := utils.GetMessageName(operation, methodName, "Request")

	if c.converterOption.NamingOption {
		messageName = common.ToPascaleCase(messageName)
	}

	message := &protobuf.ProtoMessage{Name: messageName}

	if operation.RequestBody == nil && len(operation.Parameters) == 0 {
		c.AddProtoImport(consts.EmptyProtoFile)
		return consts.EmptyMessage, nil
	}

	if operation.RequestBody != nil {
		if operation.RequestBody.Ref != "" {
			return common.ToPascaleCase(utils.ExtractMessageNameFromRef(operation.RequestBody.Ref)), nil
		}

		if operation.RequestBody.Value != nil && len(operation.RequestBody.Value.Content) > 0 {
			for mediaTypeStr, mediaType := range operation.RequestBody.Value.Content {
				schema := mediaType.Schema
				if schema != nil {
					protoType, err := c.ConvertSchemaToProtoType(schema, common.FormatStr(mediaTypeStr), message)
					if err != nil {
						return "", err
					}

					switch v := protoType.(type) {
					case *protobuf.ProtoField:
						if c.converterOption.ApiOption {
							var optionName string
							if mediaTypeStr == "application/json" {
								optionName = "api.body"
							} else if mediaTypeStr == "application/x-www-form-urlencoded" || mediaTypeStr == "multipart/form-data" {
								optionName = "api.form"
							}
							if optionName != "" {
								v.Options = append(v.Options, &protobuf.Option{
									Name:  optionName,
									Value: fmt.Sprintf("%q", v.Name),
								})
								c.AddProtoImport(consts.ApiProtoFile)
							}
						}
						c.addFieldIfNotExists(&message.Fields, v)
					case *protobuf.ProtoMessage:
						for _, field := range v.Fields {
							if c.converterOption.ApiOption {
								var optionName string
								if mediaTypeStr == "application/json" {
									optionName = "api.body"
								} else if mediaTypeStr == "application/x-www-form-urlencoded" || mediaTypeStr == "multipart/form-data" {
									optionName = "api.form"
								}
								if optionName != "" {
									field.Options = append(field.Options, &protobuf.Option{
										Name:  optionName,
										Value: fmt.Sprintf("%q", field.Name),
									})
									c.AddProtoImport(consts.ApiProtoFile)
								}
							}
							c.addFieldIfNotExists(&message.Fields, field)
						}

						message.Enums = append(message.Enums, v.Enums...)

						message.OneOfs = append(message.OneOfs, v.OneOfs...)

						for _, nestedMessage := range v.Messages {
							c.addMessageIfNotExists(&message.Messages, nestedMessage)
						}
					case *protobuf.ProtoEnum:
						name := mediaTypeStr
						if c.converterOption.NamingOption {
							name = common.ToSnakeCase(name)
						} else {
							name = common.FormatStr(name)
						}
						newField := &protobuf.ProtoField{
							Name: name + "_field",
							Type: v.Name,
						}
						if c.converterOption.ApiOption {
							var optionName string
							if mediaTypeStr == "application/json" {
								optionName = "api.body"
							} else if mediaTypeStr == "application/x-www-form-urlencoded" || mediaTypeStr == "multipart/form-data" {
								optionName = "api.form"
							}
							if optionName != "" {
								newField.Options = append(newField.Options, &protobuf.Option{
									Name:  optionName,
									Value: fmt.Sprintf("%q", v.Name),
								})
								c.AddProtoImport(consts.ApiProtoFile)
							}
						}
						if c.converterOption.OpenapiOption {
							optionStr := common.StructToOption(schema.Value, "     ")

							schemaOption := &protobuf.Option{
								Name:  consts.OpenapiProperty,
								Value: optionStr,
							}
							newField.Options = append(newField.Options, schemaOption)
							c.AddProtoImport(consts.OpenapiProtoFile)
						}
						message.Enums = append(message.Enums, v)
						message.Fields = append(message.Fields, newField)
					case *protobuf.ProtoOneOf:
						message.OneOfs = append(message.OneOfs, v)
					}
				}
			}
		}
	}

	if len(operation.Parameters) > 0 {
		for _, param := range operation.Parameters {
			if param.Value.Schema != nil {
				fieldOrMessage, err := c.ConvertSchemaToProtoType(param.Value.Schema, param.Value.Name, message)
				if err != nil {
					return "", err
				}
				description := param.Value.Description
				switch v := fieldOrMessage.(type) {
				case *protobuf.ProtoField:
					if c.converterOption.ApiOption {
						v.Options = append(v.Options, &protobuf.Option{
							Name:  "api." + param.Value.In,
							Value: fmt.Sprintf("%q", param.Value.Name),
						})
						c.AddProtoImport(consts.ApiProtoFile)
					}
					if c.converterOption.OpenapiOption {
						optionStr := common.StructToOption(param.Value, "     ")

						schemaOption := &protobuf.Option{
							Name:  consts.OpenapiParameter,
							Value: optionStr,
						}
						v.Options = append(v.Options, schemaOption)
						c.AddProtoImport(consts.OpenapiProtoFile)
					}
					v.Description = description
					c.addFieldIfNotExists(&message.Fields, v)
				case *protobuf.ProtoMessage:
					for _, field := range v.Fields {
						if c.converterOption.ApiOption {
							field.Options = append(field.Options, &protobuf.Option{
								Name:  "api." + param.Value.In,
								Value: fmt.Sprintf("%q", param.Value.Name),
							})
							c.AddProtoImport(consts.ApiProtoFile)
						}
						if c.converterOption.OpenapiOption {
							optionStr := common.StructToOption(param.Value, "     ")

							schemaOption := &protobuf.Option{
								Name:  consts.OpenapiParameter,
								Value: optionStr,
							}
							field.Options = append(field.Options, schemaOption)
							c.AddProtoImport(consts.OpenapiProtoFile)
						}
						c.addFieldIfNotExists(&message.Fields, field)
					}
					message.Enums = append(message.Enums, v.Enums...)

					message.OneOfs = append(message.OneOfs, v.OneOfs...)

					for _, nestedMessage := range v.Messages {
						c.addMessageIfNotExists(&message.Messages, nestedMessage)
					}
				case *protobuf.ProtoEnum:
					name := param.Value.Name
					if c.converterOption.NamingOption {
						name = common.ToPascaleCase(name)
					}
					newField := &protobuf.ProtoField{
						Name: name + "_field",
						Type: v.Name,
					}
					if c.converterOption.ApiOption {
						newField.Options = append(newField.Options, &protobuf.Option{
							Name:  "api." + param.Value.In,
							Value: fmt.Sprintf("%q", param.Value.Name),
						})
						c.AddProtoImport(consts.ApiProtoFile)
					}
					if c.converterOption.OpenapiOption {
						optionStr := common.StructToOption(param.Value, "     ")

						schemaOption := &protobuf.Option{
							Name:  consts.OpenapiParameter,
							Value: optionStr,
						}
						newField.Options = append(newField.Options, schemaOption)
						c.AddProtoImport(consts.OpenapiProtoFile)
					}
					message.Enums = append(message.Enums, v)
					message.Fields = append(message.Fields, newField)
				case *protobuf.ProtoOneOf:
					message.OneOfs = append(message.OneOfs, v)
				}
			}
		}
	}

	// if there are no fields or messages, return an empty message
	if len(message.Fields) > 0 || len(message.Messages) > 0 || len(message.Enums) > 0 {
		c.addMessageToProto(message)
		return message.Name, nil
	}

	return "", nil
}

// generateResponseMessage generates a response message for an operation
func (c *ProtoConverter) generateResponseMessage(operation *openapi3.Operation, methodName string) (string, error) {
	if operation.Responses == nil {
		return "", nil
	}

	responses := operation.Responses.Map()
	responseCount := 0
	for _, responseRef := range responses {
		if responseRef.Ref == "" && (responseRef.Value == nil || (len(responseRef.Value.Content) == 0 && len(responseRef.Value.Headers) == 0)) {
			continue
		}
		responseCount++
	}

	if responseCount == 1 {
		for _, responseRef := range responses {
			if responseRef.Ref == "" && (responseRef.Value == nil || (len(responseRef.Value.Content) == 0 && len(responseRef.Value.Headers) == 0)) {
				continue
			}
			return c.processSingleResponse("", responseRef, operation, methodName)
		}
	}

	if responseCount == 0 {
		c.AddProtoImport(consts.EmptyProtoFile)
		return consts.EmptyMessage, nil
	}

	// create a wrapper message for multiple responses
	wrapperMessageName := utils.GetMessageName(operation, methodName, "Response")
	if c.converterOption.NamingOption {
		wrapperMessageName = common.ToPascaleCase(wrapperMessageName)
	}

	wrapperMessage := &protobuf.ProtoMessage{Name: wrapperMessageName}

	emptyFlag := true

	for statusCode, responseRef := range responses {
		if responseRef.Ref == "" && (responseRef.Value == nil || len(responseRef.Value.Content) == 0) {
			break
		}
		emptyFlag = false
		messageName, err := c.processSingleResponse(statusCode, responseRef, operation, methodName)
		if err != nil {
			return "", err
		}

		name := "Response_" + statusCode
		if c.converterOption.NamingOption {
			name = common.ToSnakeCase(name)
		}
		field := &protobuf.ProtoField{
			Name: name,
			Type: messageName,
		}
		wrapperMessage.Fields = append(wrapperMessage.Fields, field)
	}

	if emptyFlag {
		c.AddProtoImport(consts.EmptyProtoFile)
		return consts.EmptyMessage, nil
	}

	c.addMessageToProto(wrapperMessage)

	return wrapperMessageName, nil
}

// processSingleResponse deals with a single response in an operation
func (c *ProtoConverter) processSingleResponse(statusCode string, responseRef *openapi3.ResponseRef, operation *openapi3.Operation, methodName string) (string, error) {
	if responseRef.Ref != "" {
		return common.ToPascaleCase(utils.ExtractMessageNameFromRef(responseRef.Ref)), nil
	}

	response := responseRef.Value
	messageName := utils.GetMessageName(operation, methodName, "Response") + common.ToUpperCase(statusCode)

	if c.converterOption.NamingOption {
		messageName = common.ToPascaleCase(messageName)
	}

	message := &protobuf.ProtoMessage{Name: messageName}

	if len(response.Headers) > 0 {
		for headerName, headerRef := range response.Headers {
			if headerRef != nil {

				fieldOrMessage, err := c.ConvertSchemaToProtoType(headerRef.Value.Schema, headerName, message)
				if err != nil {
					return "", err
				}

				switch v := fieldOrMessage.(type) {
				case *protobuf.ProtoField:
					if c.converterOption.ApiOption {
						option := &protobuf.Option{
							Name:  "api.header",
							Value: fmt.Sprintf("%q", headerName),
						}
						v.Options = append(v.Options, option)
						c.AddProtoImport(consts.ApiProtoFile)
					}
					if c.converterOption.OpenapiOption {
						optionStr := common.StructToOption(headerRef.Value, "     ")

						schemaOption := &protobuf.Option{
							Name:  consts.OpenapiProperty,
							Value: optionStr,
						}
						v.Options = append(v.Options, schemaOption)
						c.AddProtoImport(consts.OpenapiProtoFile)
					}
					c.addFieldIfNotExists(&message.Fields, v)
				case *protobuf.ProtoMessage:
					for _, field := range v.Fields {
						if c.converterOption.ApiOption {
							option := &protobuf.Option{
								Name:  "api.header",
								Value: fmt.Sprintf("%q", field.Name),
							}
							field.Options = append(field.Options, option)
							c.AddProtoImport(consts.ApiProtoFile)
						}
						if c.converterOption.OpenapiOption {
							optionStr := common.StructToOption(headerRef.Value, "     ")

							schemaOption := &protobuf.Option{
								Name:  consts.OpenapiProperty,
								Value: optionStr,
							}
							field.Options = append(field.Options, schemaOption)
							c.AddProtoImport(consts.OpenapiProtoFile)
						}
						c.addFieldIfNotExists(&message.Fields, field)
					}
					message.Enums = append(message.Enums, v.Enums...)

					message.OneOfs = append(message.OneOfs, v.OneOfs...)

					for _, nestedMessage := range v.Messages {
						c.addMessageIfNotExists(&message.Messages, nestedMessage)
					}
				case *protobuf.ProtoEnum:
					name := headerName
					if c.converterOption.NamingOption {
						name = common.ToSnakeCase(name)
					}
					newField := &protobuf.ProtoField{
						Name: name + "_field",
						Type: v.Name,
					}
					if c.converterOption.ApiOption {
						option := &protobuf.Option{
							Name:  "api.header",
							Value: fmt.Sprintf("%q", headerName),
						}
						newField.Options = append(newField.Options, option)
						c.AddProtoImport(consts.ApiProtoFile)
					}
					if c.converterOption.OpenapiOption {
						optionStr := common.StructToOption(headerRef.Value, "     ")

						schemaOption := &protobuf.Option{
							Name:  consts.OpenapiProperty,
							Value: optionStr,
						}
						newField.Options = append(newField.Options, schemaOption)
						c.AddProtoImport(consts.OpenapiProtoFile)
					}
					message.Enums = append(message.Enums, v)
					message.Fields = append(message.Fields, newField)
				case *protobuf.ProtoOneOf:
					message.OneOfs = append(message.OneOfs, v)
				}
			}
		}
	}

	for mediaTypeStr, mediaType := range response.Content {
		schema := mediaType.Schema
		if schema != nil {

			protoType, err := c.ConvertSchemaToProtoType(schema, common.FormatStr(mediaTypeStr), message)
			if err != nil {
				return "", err
			}

			switch v := protoType.(type) {
			case *protobuf.ProtoField:
				if c.converterOption.ApiOption && mediaTypeStr == "application/json" {
					option := &protobuf.Option{
						Name:  "api.body",
						Value: fmt.Sprintf("%q", v.Name),
					}
					v.Options = append(v.Options, option)
					c.AddProtoImport(consts.ApiProtoFile)
				}
				if c.converterOption.OpenapiOption {
					optionStr := common.StructToOption(schema.Value, "     ")

					schemaOption := &protobuf.Option{
						Name:  consts.OpenapiProperty,
						Value: optionStr,
					}
					v.Options = append(v.Options, schemaOption)
					c.AddProtoImport(consts.OpenapiProtoFile)
				}
				c.addFieldIfNotExists(&message.Fields, v)
			case *protobuf.ProtoMessage:
				for _, field := range v.Fields {
					if c.converterOption.ApiOption && mediaTypeStr == "application/json" {
						option := &protobuf.Option{
							Name:  "api.body",
							Value: fmt.Sprintf("%q", field.Name),
						}
						field.Options = append(field.Options, option)
						c.AddProtoImport(consts.ApiProtoFile)
					}
					if c.converterOption.OpenapiOption {
						optionStr := common.StructToOption(schema.Value, "     ")

						schemaOption := &protobuf.Option{
							Name:  consts.OpenapiProperty,
							Value: optionStr,
						}
						field.Options = append(field.Options, schemaOption)
						c.AddProtoImport(consts.OpenapiProtoFile)
					}
					c.addFieldIfNotExists(&message.Fields, field)
				}
				message.Enums = append(message.Enums, v.Enums...)

				message.OneOfs = append(message.OneOfs, v.OneOfs...)

				for _, nestedMessage := range v.Messages {
					c.addMessageIfNotExists(&message.Messages, nestedMessage)
				}
			case *protobuf.ProtoEnum:
				name := mediaTypeStr
				if c.converterOption.NamingOption {
					name = common.ToSnakeCase(mediaTypeStr)
				} else {
					name = common.ToUpperCase(name)
				}
				newField := &protobuf.ProtoField{
					Name: name + "_field",
					Type: v.Name,
				}
				if c.converterOption.ApiOption && mediaTypeStr == "application/json" {
					option := &protobuf.Option{
						Name:  "api.body",
						Value: fmt.Sprintf("%q", v.Name),
					}
					newField.Options = append(newField.Options, option)
					c.AddProtoImport(consts.ApiProtoFile)
				}
				if c.converterOption.OpenapiOption {
					optionStr := common.StructToOption(schema.Value, "     ")

					schemaOption := &protobuf.Option{
						Name:  consts.OpenapiProperty,
						Value: optionStr,
					}
					newField.Options = append(newField.Options, schemaOption)
					c.AddProtoImport(consts.OpenapiProtoFile)
				}
				message.Enums = append(message.Enums, v)
				message.Fields = append(message.Fields, newField)
			case *protobuf.ProtoOneOf:
				message.OneOfs = append(message.OneOfs, v)
			}
		}
	}

	if len(message.Fields) > 0 || len(message.Messages) > 0 || len(message.Enums) > 0 {
		c.addMessageToProto(message)
		return message.Name, nil
	}
	return "", nil
}

// ConvertSchemaToProtoType converts an OpenAPI schema to a Proto field or message
func (c *ProtoConverter) ConvertSchemaToProtoType(
	schemaRef *openapi3.SchemaRef,
	protoName string,
	parentMessage *protobuf.ProtoMessage,
) (interface{}, error) {
	var protoType string
	var result interface{}

	// Handle referenced schema
	if schemaRef.Ref != "" {
		name := c.applySnakeCaseNamingOption(utils.ExtractMessageNameFromRef(schemaRef.Ref))
		return &protobuf.ProtoField{
			Name: name,
			Type: common.ToPascaleCase(utils.ExtractMessageNameFromRef(schemaRef.Ref)),
		}, nil
	}

	// Ensure schema value is valid
	if schemaRef.Value == nil {
		return nil, errors.New("schema type is required")
	}

	schema := schemaRef.Value
	description := schema.Description

	// Handle oneOf, allOf, anyOf even if schema.Type is nil
	if len(schema.OneOf) > 0 {
		protoUnion, err := c.handleOneOf(schema.OneOf, protoName, parentMessage)
		if err != nil {
			return nil, err
		}
		return protoUnion, nil
	} else if len(schema.AllOf) > 0 {
		protoMessage, err := c.handleAllOf(schema.AllOf, protoName, parentMessage)
		if err != nil {
			return nil, err
		}
		return protoMessage, nil
	} else if len(schema.AnyOf) > 0 {
		protoMessage, err := c.handleAnyOf(schema.AnyOf, protoName, parentMessage)
		if err != nil {
			return nil, err
		}
		return protoMessage, nil
	}

	// Process schema type
	switch {
	case schema.Type.Includes("string"):
		if schema.Format == "date" || schema.Format == "date-time" {
			protoType = consts.TimestampMessage
			c.AddProtoImport(consts.TimestampProtoFile)
		} else if len(schema.Enum) != 0 {
			var name string
			if parentMessage == nil {
				name = protoName
			} else {
				name = c.applyPascaleCaseNamingOption(common.ToUpperCase(protoName))
			}
			protoEnum := &protobuf.ProtoEnum{
				Name:        name,
				Description: description,
			}
			for i, enumValue := range schema.Enum {
				protoEnum.Values = append(protoEnum.Values, &protobuf.ProtoEnumValue{
					Index: i,
					Value: enumValue,
				})
			}
			result = protoEnum
		} else {
			protoType = "string"
		}

	case schema.Type.Includes("integer"):
		if len(schema.Enum) != 0 {
			var name string
			if parentMessage == nil {
				name = protoName
			} else {
				name = c.applyPascaleCaseNamingOption(common.ToUpperCase(protoName))
			}
			protoEnum := &protobuf.ProtoEnum{
				Name:        name,
				Description: description,
			}
			for i, enumValue := range schema.Enum {
				protoEnum.Values = append(protoEnum.Values, &protobuf.ProtoEnumValue{
					Index: i,
					Value: enumValue,
				})
			}
			result = protoEnum
		} else if schema.Format == "int32" {
			protoType = "int32"
		} else {
			protoType = "int64"
		}

	case schema.Type.Includes("number"):
		if len(schema.Enum) != 0 {
			var name string
			if parentMessage == nil {
				name = protoName
			} else {
				name = c.applyPascaleCaseNamingOption(common.ToUpperCase(protoName))
			}
			protoEnum := &protobuf.ProtoEnum{
				Name:        name,
				Description: description,
			}
			for i, enumValue := range schema.Enum {
				protoEnum.Values = append(protoEnum.Values, &protobuf.ProtoEnumValue{
					Index: i,
					Value: enumValue,
				})
			}
			result = protoEnum
		} else if schema.Format == "float" {
			protoType = "float"
		} else {
			protoType = "double"
		}

	case schema.Type.Includes("boolean"):
		protoType = "bool"

	case schema.Type.Includes("array"):
		if schema.Items != nil {
			fieldOrMessage, err := c.ConvertSchemaToProtoType(schema.Items, protoName+"Item", parentMessage)
			if err != nil {
				return nil, err
			}

			fieldType := ""
			if field, ok := fieldOrMessage.(*protobuf.ProtoField); ok {
				fieldType = field.Type
			} else if nestedMessage, ok := fieldOrMessage.(*protobuf.ProtoMessage); ok {
				fieldType = nestedMessage.Name
				c.addNestedMessageToParent(parentMessage, nestedMessage)
			} else if enum, ok := fieldOrMessage.(*protobuf.ProtoEnum); ok {
				fieldType = enum.Name
				c.addNestedEnumToParent(parentMessage, enum)
			}

			result = &protobuf.ProtoField{
				Name:        c.applySnakeCaseNamingOption(protoName),
				Type:        fieldType,
				Repeated:    true,
				Description: description,
			}
		}

	case schema.Type.Includes("object"):
		var message *protobuf.ProtoMessage
		if parentMessage == nil {
			message = &protobuf.ProtoMessage{Name: protoName}
		} else {
			message = &protobuf.ProtoMessage{Name: c.applyPascaleCaseNamingOption(common.ToUpperCase(protoName))}
		}
		for propName, propSchema := range schema.Properties {
			protoType, err := c.ConvertSchemaToProtoType(propSchema, propName, message)
			if err != nil {
				return nil, err
			}

			if field, ok := protoType.(*protobuf.ProtoField); ok {
				if c.converterOption.OpenapiOption {
					optionStr := common.StructToOption(propSchema.Value, "     ")

					schemaOption := &protobuf.Option{
						Name:  consts.OpenapiProperty,
						Value: optionStr,
					}
					field.Options = append(field.Options, schemaOption)
					c.AddProtoImport(consts.OpenapiProtoFile)
				}
				message.Fields = append(message.Fields, field)
			} else if nestedMessage, ok := protoType.(*protobuf.ProtoMessage); ok {
				var name string
				if c.converterOption.NamingOption {
					name = common.ToSnakeCase(nestedMessage.Name)
				}
				newField := &protobuf.ProtoField{
					Name: name + "_field",
					Type: nestedMessage.Name,
				}
				if c.converterOption.OpenapiOption {
					optionStr := common.StructToOption(propSchema.Value, "     ")

					schemaOption := &protobuf.Option{
						Name:  consts.OpenapiProperty,
						Value: optionStr,
					}
					newField.Options = append(newField.Options, schemaOption)
					c.AddProtoImport(consts.OpenapiProtoFile)
				}
				c.addNestedMessageToParent(message, nestedMessage)
				message.Fields = append(message.Fields, newField)
			} else if enum, ok := protoType.(*protobuf.ProtoEnum); ok {
				c.addNestedEnumToParent(message, enum)
				message.Fields = append(message.Fields, &protobuf.ProtoField{
					Name: c.applySnakeCaseNamingOption(propName + "_field"),
					Type: enum.Name,
				})
			} else if oneOf, ok := protoType.(*protobuf.ProtoOneOf); ok {
				c.addNestedOneOfToParent(message, oneOf)
			}
		}

		if schema.AdditionalProperties.Schema != nil {
			mapValueType := "string"
			additionalPropMessage, err := c.ConvertSchemaToProtoType(schema.AdditionalProperties.Schema, protoName+"AdditionalProperties", parentMessage)
			if err != nil {
				return nil, err
			}
			if msg, ok := additionalPropMessage.(*protobuf.ProtoMessage); ok {
				mapValueType = msg.Name
			} else if enum, ok := additionalPropMessage.(*protobuf.ProtoEnum); ok {
				mapValueType = enum.Name
			}

			message.Fields = append(message.Fields, &protobuf.ProtoField{
				Name: "additional_properties",
				Type: "map<string, " + mapValueType + ">",
			})
		}

		message.Description = description
		result = message
	}

	// If result is still nil, construct a default ProtoField
	if result == nil {
		result = &protobuf.ProtoField{
			Name:        c.applySnakeCaseNamingOption(protoName),
			Type:        protoType,
			Description: description,
		}
	}

	return result, nil
}

// handleOneOf processes oneOf schemas
func (c *ProtoConverter) handleOneOf(oneOfSchemas []*openapi3.SchemaRef, protoName string, parentMessage *protobuf.ProtoMessage) (*protobuf.ProtoOneOf, error) {
	oneOf := &protobuf.ProtoOneOf{
		Name: c.applyPascaleCaseNamingOption(protoName + "OneOf"),
	}

	for i, schemaRef := range oneOfSchemas {
		fieldName := fmt.Sprintf("%sOption%d", protoName, i+1)
		protoType, err := c.ConvertSchemaToProtoType(schemaRef, fieldName, parentMessage)
		if err != nil {
			return nil, err
		}
		switch v := protoType.(type) {
		case *protobuf.ProtoField:
			oneOf.Fields = append(oneOf.Fields, v)
		case *protobuf.ProtoMessage:
			newField := &protobuf.ProtoField{
				Name: c.applySnakeCaseNamingOption(v.Name + "_field"),
				Type: v.Name,
			}
			c.addNestedMessageToParent(parentMessage, v)
			oneOf.Fields = append(oneOf.Fields, newField)
		case *protobuf.ProtoEnum:
			newField := &protobuf.ProtoField{
				Name: c.applySnakeCaseNamingOption(v.Name + "_field"),
				Type: v.Name,
			}
			c.addNestedEnumToParent(parentMessage, v)
			oneOf.Fields = append(oneOf.Fields, newField)
		case *protobuf.ProtoOneOf:
			c.addNestedOneOfToParent(parentMessage, v)
		}
	}
	return oneOf, nil
}

// handleAllOf processes allOf schemas
func (c *ProtoConverter) handleAllOf(allOfSchemas []*openapi3.SchemaRef, protoName string, parentMessage *protobuf.ProtoMessage) (*protobuf.ProtoMessage, error) {
	allOfMessage := &protobuf.ProtoMessage{
		Name: c.applyPascaleCaseNamingOption(protoName + "AllOf"),
	}

	for i, schemaRef := range allOfSchemas {
		fieldName := fmt.Sprintf("%sPart%d", protoName, i+1)
		protoType, err := c.ConvertSchemaToProtoType(schemaRef, fieldName, parentMessage)
		if err != nil {
			return nil, err
		}

		switch v := protoType.(type) {
		case *protobuf.ProtoField:
			allOfMessage.Fields = append(allOfMessage.Fields, v)
		case *protobuf.ProtoMessage:
			newField := &protobuf.ProtoField{
				Name: c.applySnakeCaseNamingOption(v.Name + "_field"),
				Type: v.Name,
			}
			c.addNestedMessageToParent(allOfMessage, v)
			allOfMessage.Fields = append(allOfMessage.Fields, newField)
		case *protobuf.ProtoEnum:
			newField := &protobuf.ProtoField{
				Name: c.applySnakeCaseNamingOption(v.Name + "_field"),
				Type: v.Name,
			}
			c.addNestedEnumToParent(allOfMessage, v)
			allOfMessage.Fields = append(allOfMessage.Fields, newField)
		case *protobuf.ProtoOneOf:
			c.addNestedOneOfToParent(allOfMessage, v)
		}
	}

	return allOfMessage, nil
}

// handleAnyOf processes anyOf schemas
func (c *ProtoConverter) handleAnyOf(anyOfSchemas []*openapi3.SchemaRef, protoName string, parentMessage *protobuf.ProtoMessage) (*protobuf.ProtoMessage, error) {
	anyOfMessage := &protobuf.ProtoMessage{
		Name: c.applyPascaleCaseNamingOption(protoName + "AnyOf"),
	}

	for i, schemaRef := range anyOfSchemas {
		fieldName := fmt.Sprintf("%sOption%d", protoName, i+1)
		protoType, err := c.ConvertSchemaToProtoType(schemaRef, fieldName, parentMessage)
		if err != nil {
			return nil, err
		}

		switch v := protoType.(type) {
		case *protobuf.ProtoField:
			anyOfMessage.Fields = append(anyOfMessage.Fields, v)
		case *protobuf.ProtoMessage:
			newField := &protobuf.ProtoField{
				Name: c.applySnakeCaseNamingOption(v.Name + "_field"),
				Type: v.Name,
			}
			c.addNestedMessageToParent(anyOfMessage, v)
			anyOfMessage.Fields = append(anyOfMessage.Fields, newField)
		case *protobuf.ProtoEnum:
			newField := &protobuf.ProtoField{
				Name: c.applySnakeCaseNamingOption(v.Name + "_field"),
				Type: v.Name,
			}
			c.addNestedEnumToParent(anyOfMessage, v)
			anyOfMessage.Fields = append(anyOfMessage.Fields, newField)
		case *protobuf.ProtoOneOf:
			c.addNestedOneOfToParent(anyOfMessage, v)
		}
	}

	return anyOfMessage, nil
}

// applyPascaleCaseNamingOption applies naming convention based on the converter's naming option
func (c *ProtoConverter) applyPascaleCaseNamingOption(name string) string {
	if c.converterOption.NamingOption {
		return common.ToPascaleCase(name)
	}
	return name
}

// applySnakeCaseNamingOption applies naming convention based on the converter's naming option
func (c *ProtoConverter) applySnakeCaseNamingOption(name string) string {
	if c.converterOption.NamingOption {
		return common.ToSnakeCase(name)
	}
	return name
}

// addOptionsToProto adds options to the ProtoFile
func (c *ProtoConverter) addOptionsToProto() {
	optionStr := common.StructToOption(c.spec, "")

	schemaOption := &protobuf.Option{
		Name:  consts.OpenapiDocument,
		Value: optionStr,
	}
	c.ProtoFile.Options = append(c.ProtoFile.Options, schemaOption)
	c.AddProtoImport(consts.OpenapiProtoFile)
}

// Add a new method to handle structured extensions
func (c *ProtoConverter) addExtensionsToProtoOptions() error {
	// Check for x-option in spec extensions
	if xOption, ok := c.spec.Extensions["x-options"]; ok {
		if optionMap, ok := xOption.(map[string]interface{}); ok {
			for key, value := range optionMap {
				option := &protobuf.Option{
					Name:  key,
					Value: fmt.Sprintf("%q", value),
				}
				c.ProtoFile.Options = append(c.ProtoFile.Options, option)
			}
		}
	}

	// Check for x-option in spec.info.extensions
	if c.spec.Info != nil {
		if xOption, ok := c.spec.Info.Extensions["x-options"]; ok {
			if optionMap, ok := xOption.(map[string]interface{}); ok {
				for key, value := range optionMap {
					option := &protobuf.Option{
						Name:  key,
						Value: fmt.Sprintf("%q", value),
					}
					c.ProtoFile.Options = append(c.ProtoFile.Options, option)
				}
			}
		}
	}

	return nil
}

// addNestedMessageToParent adds a nested message to a parent message
func (c *ProtoConverter) addNestedMessageToParent(parentMessage, nestedMessage *protobuf.ProtoMessage) {
	if parentMessage != nil && nestedMessage != nil {
		parentMessage.Messages = append(parentMessage.Messages, nestedMessage)
	}
}

// addNestedEnum adds a nested Enum to a parent message
func (c *ProtoConverter) addNestedEnumToParent(parentMessage *protobuf.ProtoMessage, nestedEnum *protobuf.ProtoEnum) {
	if parentMessage != nil && nestedEnum != nil {
		parentMessage.Enums = append(parentMessage.Enums, nestedEnum)
	}
}

// addNestedOneOfToParent adds a nested oneOf to a parent message
func (c *ProtoConverter) addNestedOneOfToParent(parentMessage *protobuf.ProtoMessage, nestedOneOf *protobuf.ProtoOneOf) {
	if parentMessage != nil && nestedOneOf != nil {
		parentMessage.OneOfs = append(parentMessage.OneOfs, nestedOneOf)
	}
}

// mergeProtoMessage merges a ProtoMessage into the ProtoFile
func (c *ProtoConverter) addMessageToProto(message *protobuf.ProtoMessage) error {
	var existingMessage *protobuf.ProtoMessage
	for _, msg := range c.ProtoFile.Messages {
		if msg.Name == message.Name {
			existingMessage = msg
			break
		}
	}

	// merge message
	if existingMessage != nil {
		// merge Fields
		fieldNames := make(map[string]struct{})
		for _, field := range existingMessage.Fields {
			fieldNames[field.Name] = struct{}{}
		}
		for _, newField := range message.Fields {
			if _, exists := fieldNames[newField.Name]; !exists {
				existingMessage.Fields = append(existingMessage.Fields, newField)
			}
		}

		// merge Messages
		messageNames := make(map[string]struct{})
		for _, nestedMsg := range existingMessage.Messages {
			messageNames[nestedMsg.Name] = struct{}{}
		}
		for _, newMessage := range message.Messages {
			if _, exists := messageNames[newMessage.Name]; !exists {
				existingMessage.Messages = append(existingMessage.Messages, newMessage)
			}
		}

		// merge Enums
		enumNames := make(map[string]struct{})
		for _, enum := range existingMessage.Enums {
			enumNames[enum.Name] = struct{}{}
		}
		for _, newEnum := range message.Enums {
			if _, exists := enumNames[newEnum.Name]; !exists {
				existingMessage.Enums = append(existingMessage.Enums, newEnum)
			}
		}

		// merge Options
		optionNames := make(map[string]struct{})
		for _, option := range existingMessage.Options {
			optionNames[option.Name] = struct{}{}
		}
		for _, newOption := range message.Options {
			if _, exists := optionNames[newOption.Name]; !exists {
				existingMessage.Options = append(existingMessage.Options, newOption)
			}
		}
	} else {
		c.ProtoFile.Messages = append(c.ProtoFile.Messages, message)
	}

	return nil
}

// addEnumToProto adds an enum to the ProtoFile
func (c *ProtoConverter) addEnumToProto(enum *protobuf.ProtoEnum) {
	c.ProtoFile.Enums = append(c.ProtoFile.Enums, enum)
}

// AddProtoImport adds an import to the ProtoFile
func (c *ProtoConverter) AddProtoImport(importFile string) {
	if c.ProtoFile != nil {
		for _, existingImport := range c.ProtoFile.Imports {
			if existingImport == importFile {
				return
			}
		}
		c.ProtoFile.Imports = append(c.ProtoFile.Imports, importFile)
	}
}

// addFieldIfNotExists adds a field to Fields if it does not already exist
func (c *ProtoConverter) addFieldIfNotExists(fields *[]*protobuf.ProtoField, field *protobuf.ProtoField) {
	for _, existingField := range *fields {
		if existingField.Name == field.Name {
			return
		}
	}
	*fields = append(*fields, field)
}

// addMessageIfNotExists adds a message to Messages if it does not already exist
func (c *ProtoConverter) addMessageIfNotExists(messages *[]*protobuf.ProtoMessage, nestedMessage *protobuf.ProtoMessage) {
	for _, existingMessage := range *messages {
		if existingMessage.Name == nestedMessage.Name {
			return
		}
	}
	*messages = append(*messages, nestedMessage)
}

// methodExistsInService checks if a method exists in a service
func (c *ProtoConverter) methodExistsInService(service *protobuf.ProtoService, methodName string) bool {
	for _, method := range service.Methods {
		if method.Name == methodName {
			return true
		}
	}
	return false
}

// findOrCreateService finds an existing service by name or creates a new one if it doesn't exist
func (c *ProtoConverter) findOrCreateService(serviceName string) *protobuf.ProtoService {
	// Iterate over existing services to find a match
	for i := range c.ProtoFile.Services {
		if c.ProtoFile.Services[i].Name == serviceName {
			return c.ProtoFile.Services[i]
		}
	}

	// If no existing service is found, create a new one
	newService := &protobuf.ProtoService{Name: serviceName}
	c.ProtoFile.Services = append(c.ProtoFile.Services, newService)
	return newService
}
