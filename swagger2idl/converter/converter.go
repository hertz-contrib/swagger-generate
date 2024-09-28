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
	"github.com/hertz-contrib/swagger-generate/swagger2idl/protobuf"
	"github.com/hertz-contrib/swagger-generate/swagger2idl/utils"
)

// ProtoConverter struct, used to convert OpenAPI specifications into Proto files
type ProtoConverter struct {
	ProtoFile       *protobuf.ProtoFile
	converterOption *ConvertOption
}

// ConvertOption adds a struct for conversion options
type ConvertOption struct {
	openapiOption bool
	apiOption     bool
}

const (
	apiProtoFile     = "api.proto"
	openapiProtoFile = "openapi.proto"
	EmptyProtoFile   = "google.protobuf.empty"
)

// NewProtoConverter creates and initializes a ProtoConverter
func NewProtoConverter(packageName string) *ProtoConverter {
	return &ProtoConverter{
		ProtoFile: &protobuf.ProtoFile{
			PackageName: packageName,
			Messages:    []*protobuf.ProtoMessage{},
			Services:    []*protobuf.ProtoService{},
		},
		converterOption: &ConvertOption{
			openapiOption: false,
			apiOption:     true,
		},
	}
}

// Convert converts the OpenAPI specification to a Proto file
func (c *ProtoConverter) Convert(spec *openapi3.T) error {
	// Convert components into Proto messages
	err := c.convertComponentsToProtoMessages(spec.Components)
	if err != nil {
		return fmt.Errorf("error converting components to proto messages: %w", err)
	}

	// Convert paths into Proto services
	err = c.convertPathsToProtoServices(spec.Paths)
	if err != nil {
		return fmt.Errorf("error converting paths to proto services: %w", err)
	}

	if c.converterOption.openapiOption {
		c.AddProtoImport(openapiProtoFile)
	}
	if c.converterOption.apiOption {
		c.AddProtoImport(apiProtoFile)
	}

	return nil
}

// convertComponentsToProtoMessages converts OpenAPI components into Proto messages and stores them in the ProtoFile
func (c *ProtoConverter) convertComponentsToProtoMessages(components *openapi3.Components) error {
	if components == nil {
		return nil
	}
	if components.Schemas == nil {
		return nil
	}
	for name, schemaRef := range components.Schemas {
		schema := schemaRef
		fieldOrMessage, err := c.ConvertSchemaToProtoFieldOrMessage(schema, name, nil)
		if err != nil {
			return fmt.Errorf("error converting schema %s: %w", name, err)
		}
		switch v := fieldOrMessage.(type) {
		case *protobuf.ProtoField:
			message := &protobuf.ProtoMessage{
				Name:   name,
				Fields: []*protobuf.ProtoField{v},
			}
			c.addMessageToProto(message)
		case *protobuf.ProtoMessage:
			c.addMessageToProto(v)
		}
	}
	return nil
}

// convertPathsToProtoServices converts OpenAPI path items into Proto services and stores them in the ProtoFile
func (c *ProtoConverter) convertPathsToProtoServices(paths *openapi3.Paths) error {
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

	methodToOption := map[string]string{
		"GET":    "api.get",
		"POST":   "api.post",
		"PUT":    "api.put",
		"PATCH":  "api.patch",
		"DELETE": "api.delete",
	}

	for path, pathItem := range paths.Map() {
		for method, operation := range pathItem.Operations() {
			serviceName := utils.GetServiceName(operation.Tags)

			methodName := utils.GenerateMethodName(operation.OperationID, method)

			inputMessage, err := c.generateRequestMessage(operation)
			if err != nil {
				return nil, fmt.Errorf("error generating request message for %s: %w", methodName, err)
			}

			outputMessage, err := c.generateResponseMessage(operation)
			if err != nil {
				return nil, fmt.Errorf("error generating response message for %s: %w", methodName, err)
			}

			service := findOrCreateService(&services, serviceName)

			if !methodExistsInService(service, methodName) {
				protoMethod := &protobuf.ProtoMethod{
					Name:   methodName,
					Input:  inputMessage,
					Output: outputMessage,
				}

				if c.converterOption.apiOption {
					if optionName, ok := methodToOption[method]; ok {
						option := &protobuf.Option{
							Name:  optionName,
							Value: fmt.Sprintf("%q", utils.ConvertPath(path)),
						}
						protoMethod.Options = append(protoMethod.Options, option)
					}
				}

				if c.converterOption.openapiOption {
					optionStr := utils.StructToProtobuf(operation, "  ")

					schemaOption := &protobuf.Option{
						Name:  "openapi.operation",
						Value: optionStr,
					}
					protoMethod.Options = append(protoMethod.Options, schemaOption)

				}
				service.Methods = append(service.Methods, protoMethod)
			}
		}
	}

	return services, nil
}

// generateRequestMessage generates a request message for an operation
func (c *ProtoConverter) generateRequestMessage(operation *openapi3.Operation) (string, error) {
	messageName := operation.OperationID + "Request"
	message := &protobuf.ProtoMessage{Name: messageName}

	if operation.RequestBody == nil && len(operation.Parameters) == 0 {
		c.AddProtoImport(EmptyProtoFile)
		return EmptyProtoFile, nil
	}

	if operation.RequestBody != nil {
		if operation.RequestBody.Ref != "" {
			return utils.ExtractMessageNameFromRef(operation.RequestBody.Ref), nil
		}

		if operation.RequestBody.Value != nil && len(operation.RequestBody.Value.Content) > 0 {
			for mediaTypeStr, mediaType := range operation.RequestBody.Value.Content {
				schema := mediaType.Schema
				if schema != nil {
					fieldOrMessage, err := c.ConvertSchemaToProtoFieldOrMessage(schema, utils.SanitizeName(messageName+mediaTypeStr), message)
					if err != nil {
						return "", err
					}

					switch v := fieldOrMessage.(type) {
					case *protobuf.ProtoField:
						addFieldIfNotExists(&message.Fields, v)
					case *protobuf.ProtoMessage:
						addMessageIfNotExists(&message.Messages, v)
					}
				}
			}
		}
	}

	if len(operation.Parameters) > 0 {
		for _, param := range operation.Parameters {
			if param.Value.Schema != nil {
				fieldOrMessage, err := c.ConvertSchemaToProtoFieldOrMessage(param.Value.Schema, param.Value.Name, message)
				if err != nil {
					return "", err
				}

				switch v := fieldOrMessage.(type) {
				case *protobuf.ProtoField:
					addFieldIfNotExists(&message.Fields, v)
				case *protobuf.ProtoMessage:
					addMessageIfNotExists(&message.Messages, v)
				}
			}
		}
	}

	// if there are no fields or messages, return an empty message
	if len(message.Fields) > 0 || len(message.Messages) > 0 {
		c.addMessageToProto(message)
		return message.Name, nil
	}

	return "", nil
}

// generateResponseMessage generates a response message for an operation
func (c *ProtoConverter) generateResponseMessage(operation *openapi3.Operation) (string, error) {
	if operation.Responses == nil {
		return "", nil
	}

	responses := operation.Responses.Map()
	responseCount := len(responses)

	if responseCount == 1 {
		for statusCode, responseRef := range responses {
			if responseRef.Ref == "" && (responseRef.Value == nil || len(responseRef.Value.Content) == 0) {
				c.AddProtoImport(EmptyProtoFile)
				return EmptyProtoFile, nil
			}
			return c.processSingleResponse(statusCode, responseRef, operation)
		}
	}

	// create a wrapper message for multiple responses
	wrapperMessageName := operation.OperationID
	wrapperMessage := &protobuf.ProtoMessage{Name: wrapperMessageName}

	emptyFlag := true

	for statusCode, responseRef := range responses {
		if responseRef.Ref == "" && (responseRef.Value == nil || len(responseRef.Value.Content) == 0) {
			break
		}
		emptyFlag = false
		messageName, err := c.processSingleResponse(statusCode, responseRef, operation)
		if err != nil {
			return "", err
		}

		field := &protobuf.ProtoField{
			Name: "response_" + statusCode,
			Type: messageName,
		}
		wrapperMessage.Fields = append(wrapperMessage.Fields, field)
	}

	if emptyFlag {
		c.AddProtoImport(EmptyProtoFile)
		return EmptyProtoFile, nil
	}

	c.addMessageToProto(wrapperMessage)

	return wrapperMessageName, nil
}

// processSingleResponse deals with a single response in an operation
func (c *ProtoConverter) processSingleResponse(statusCode string, responseRef *openapi3.ResponseRef, operation *openapi3.Operation) (string, error) {
	if responseRef.Ref != "" {
		return utils.ExtractMessageNameFromRef(responseRef.Ref), nil
	}

	response := responseRef.Value
	messageName := operation.OperationID + "Response_" + statusCode
	newMessage := &protobuf.ProtoMessage{Name: messageName}

	if len(response.Headers) > 0 {
		for headerName, headerRef := range response.Headers {
			if headerRef != nil {

				fieldOrMessage, err := c.ConvertSchemaToProtoFieldOrMessage(headerRef.Value.Schema, headerName, newMessage)
				if err != nil {
					return "", err
				}

				switch v := fieldOrMessage.(type) {
				case *protobuf.ProtoField:
					addFieldIfNotExists(&newMessage.Fields, v)
				case *protobuf.ProtoMessage:
					addMessageIfNotExists(&newMessage.Messages, v)
				}
			}
		}
	}

	for mediaTypeStr, mediaType := range response.Content {
		schema := mediaType.Schema
		if schema != nil {

			fieldOrMessage, err := c.ConvertSchemaToProtoFieldOrMessage(schema, mediaTypeStr, newMessage)
			if err != nil {
				return "", err
			}

			switch v := fieldOrMessage.(type) {
			case *protobuf.ProtoField:
				addFieldIfNotExists(&newMessage.Fields, v)
			case *protobuf.ProtoMessage:
				addMessageIfNotExists(&newMessage.Messages, v)
			}
		}
	}

	if len(newMessage.Fields) > 0 || len(newMessage.Messages) > 0 {
		c.addMessageToProto(newMessage)
		return newMessage.Name, nil
	}
	return "", nil
}

// ConvertSchemaToProtoFieldOrMessage converts an OpenAPI schema to a Proto field or message
func (c *ProtoConverter) ConvertSchemaToProtoFieldOrMessage(schemaRef *openapi3.SchemaRef, protoName string, parentMessage *protobuf.ProtoMessage) (interface{}, error) {
	if schemaRef.Ref != "" {
		return &protobuf.ProtoField{Name: utils.ExtractMessageNameFromRef(schemaRef.Ref), Type: utils.ExtractMessageNameFromRef(schemaRef.Ref)}, nil
	}

	if schemaRef.Value != nil {
		schema := schemaRef.Value
		if schema.Type == nil {
			return nil, errors.New("schema type is required")
		}
		if schema.Type != nil {
			var protoType string
			if schema.Type.Includes("string") {
				if schema.Format == "date" || schema.Format == "date-time" {
					protoType = "google.protobuf.Timestamp"
					c.AddProtoImport("google/protobuf/timestamp.proto")
				} else {
					protoType = "string"
				}
			} else if schema.Type.Includes("integer") {
				if schema.Format == "int32" {
					protoType = "int32"
				} else {
					protoType = "int64"
				}
			} else if schema.Type.Includes("number") {
				if schema.Format == "float" {
					protoType = "float"
				} else {
					protoType = "double"
				}
			} else if schema.Type.Includes("boolean") {
				protoType = "bool"
			} else if schema.Type.Includes("array") {
				if schema.Items != nil {
					itemSchema := schema.Items
					// recursive call to handle array items
					fieldOrMessage, err := c.ConvertSchemaToProtoFieldOrMessage(itemSchema, protoName+"Item", parentMessage)
					if err != nil {
						return nil, err
					}

					if field, ok := fieldOrMessage.(*protobuf.ProtoField); ok {
						return &protobuf.ProtoField{
							Name:     protoName,
							Type:     field.Type, // 这里直接生成 repeated 类型
							Repeated: true,
						}, nil
					} else if nestedMessage, ok := fieldOrMessage.(*protobuf.ProtoMessage); ok {
						repeatedField := &protobuf.ProtoField{
							Name:     protoName,
							Type:     nestedMessage.Name,
							Repeated: true,
						}

						c.addNestedMessageToParent(parentMessage, nestedMessage)

						return repeatedField, nil
					}
				}
			} else if schema.Type.Includes("object") {
				message := &protobuf.ProtoMessage{Name: protoName}
				for propName, propSchema := range schema.Properties {
					// recursive call to handle object properties
					fieldOrMessage, err := c.ConvertSchemaToProtoFieldOrMessage(propSchema, propName, message)
					if err != nil {
						return nil, err
					}

					if field, ok := fieldOrMessage.(*protobuf.ProtoField); ok {
						message.Fields = append(message.Fields, field)
					} else if nestedMessage, ok := fieldOrMessage.(*protobuf.ProtoMessage); ok {
						c.addNestedMessageToParent(message, nestedMessage)
						message.Fields = append(message.Fields, &protobuf.ProtoField{
							Name: propName + "Field",
							Type: nestedMessage.Name,
						})
					}
				}

				if schema.AdditionalProperties.Schema != nil {
					mapValueType := "string"
					additionalPropMessage, err := c.ConvertSchemaToProtoFieldOrMessage(schema.AdditionalProperties.Schema, protoName+"AdditionalProperties", parentMessage)
					if err != nil {
						return nil, err
					}
					if message, ok := additionalPropMessage.(*protobuf.ProtoMessage); ok {
						mapValueType = message.Name
					}
					message.Fields = append(message.Fields, &protobuf.ProtoField{
						Name: "additionalProperties",
						Type: "map<string, " + mapValueType + ">",
					})
				}

				return message, nil
			}

			return &protobuf.ProtoField{
				Name: protoName,
				Type: protoType,
			}, nil
		}
	}
	return nil, nil
}

// addNestedMessageToParent adds a nested message to a parent message
func (c *ProtoConverter) addNestedMessageToParent(parentMessage, nestedMessage *protobuf.ProtoMessage) {
	if parentMessage != nil && nestedMessage != nil {
		parentMessage.Messages = append(parentMessage.Messages, nestedMessage)
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
			fieldNames[field.Name] = struct{}{} // 记录已有字段名称
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
func addFieldIfNotExists(fields *[]*protobuf.ProtoField, field *protobuf.ProtoField) {
	for _, existingField := range *fields {
		if existingField.Name == field.Name {
			return
		}
	}
	*fields = append(*fields, field)
}

// addMessageIfNotExists adds a message to Messages if it does not already exist
func addMessageIfNotExists(messages *[]*protobuf.ProtoMessage, nestedMessage *protobuf.ProtoMessage) {
	for _, existingMessage := range *messages {
		if existingMessage.Name == nestedMessage.Name {
			return
		}
	}
	*messages = append(*messages, nestedMessage)
}

// methodExistsInService checks if a method exists in a service
func methodExistsInService(service *protobuf.ProtoService, methodName string) bool {
	for _, method := range service.Methods {
		if method.Name == methodName {
			return true
		}
	}
	return false
}

// findOrCreateService finds or creates a service
func findOrCreateService(services *[]*protobuf.ProtoService, serviceName string) *protobuf.ProtoService {
	for i := range *services {
		if (*services)[i].Name == serviceName {
			return (*services)[i]
		}
	}

	newService := &protobuf.ProtoService{Name: serviceName}
	*services = append(*services, newService)
	return (*services)[len(*services)-1]
}
