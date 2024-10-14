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

package generate

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	common "github.com/hertz-contrib/swagger-generate/common/utils"
	"github.com/hertz-contrib/swagger-generate/swagger2idl/thrift"
)

// ThriftGenerate handles the encoding context for Thrift files
type ThriftGenerate struct {
	dst *strings.Builder // Output destination
}

// NewThriftGenerate creates a new instance of ThriftGenerate
func NewThriftGenerate() *ThriftGenerate {
	return &ThriftGenerate{dst: &strings.Builder{}}
}

// Generate converts a ThriftFile structure into Thrift file content
func (e *ThriftGenerate) Generate(fileContent interface{}) (string, error) {
	thriftFile, ok := fileContent.(*thrift.ThriftFile)
	if !ok {
		return "", fmt.Errorf("invalid type: expected *ThriftFile")
	}

	if len(thriftFile.Namespace) == 0 {
		e.dst.WriteString("namespace go example\n\n")
	} else {
		for language, ns := range thriftFile.Namespace {
			e.dst.WriteString(fmt.Sprintf("namespace %s %s\n", language, ns))
		}
		e.dst.WriteString("\n")
	}

	// Generate includes
	if len(thriftFile.Includes) > 0 {
		for _, include := range thriftFile.Includes {
			e.dst.WriteString(fmt.Sprintf("include \"%s\"\n", include))
		}

		e.dst.WriteString("\n")
	}

	// sort the enums
	sort.Slice(thriftFile.Enums, func(i, j int) bool {
		return thriftFile.Enums[i].Name < thriftFile.Enums[j].Name
	})

	// Generate enums
	for _, enum := range thriftFile.Enums {
		e.encodeEnum(enum, 0)
	}

	// sort the structs
	sort.Slice(thriftFile.Structs, func(i, j int) bool {
		return thriftFile.Structs[i].Name < thriftFile.Structs[j].Name
	})

	// Generate structs
	for _, message := range thriftFile.Structs {
		e.encodeMessage(message, 0)
	}

	// sort the unions
	sort.Slice(thriftFile.Unions, func(i, j int) bool {
		return thriftFile.Unions[i].Name < thriftFile.Unions[j].Name
	})

	// Generate unions
	for _, union := range thriftFile.Unions {
		e.encodeUnion(union, 0)
	}

	// sort the services
	sort.Slice(thriftFile.Services, func(i, j int) bool {
		return thriftFile.Services[i].Name < thriftFile.Services[j].Name
	})

	// Generate services
	for _, service := range thriftFile.Services {
		e.encodeService(service)
	}

	return e.dst.String(), nil
}

// encodeService encodes service definitions
func (e *ThriftGenerate) encodeService(service *thrift.ThriftService) {
	if service.Description != "" {
		e.dst.WriteString(fmt.Sprintf("// %s\n", service.Description))
	}

	e.dst.WriteString(fmt.Sprintf("service %s {\n", service.Name)) // Service declaration

	// Methods
	for _, method := range service.Methods {
		e.encodeMethod(method)
	}

	e.dst.WriteString("}")

	// Service options
	if len(service.Options) > 0 {
		e.dst.WriteString("(")
		for i, option := range service.Options {
			if i > 0 {
				e.dst.WriteString(", ")
			}
			e.encodeOption(option)
		}
		e.dst.WriteString(")\n")
	} else {
		e.dst.WriteString("\n")
	}

	e.dst.WriteString("\n")
}

// encodeMethod encodes methods within a service
func (e *ThriftGenerate) encodeMethod(method *thrift.ThriftMethod) {
	if method.Description != "" {
		e.dst.WriteString(fmt.Sprintf("  // %s\n", method.Description))
	}

	e.dst.WriteString(fmt.Sprintf("    %s %s (", method.Output, method.Name)) // Method signature

	// Input parameters
	for i, input := range method.Input {
		if i > 0 {
			e.dst.WriteString(", ")
		}
		e.dst.WriteString(fmt.Sprintf("%d: %s req", i+1, input))
	}

	e.dst.WriteString(")")

	// Method options
	if len(method.Options) > 0 {
		e.dst.WriteString(" (\n")
		for i, option := range method.Options {
			if i > 0 {
				e.dst.WriteString(",\n")
			}
			e.dst.WriteString("        ")
			e.encodeOption(option)
		}
		e.dst.WriteString("\n    )")
	}

	e.dst.WriteString("\n")
}

// encodeMessage recursively encodes structs, including nested structs and enums
func (e *ThriftGenerate) encodeMessage(message *thrift.ThriftStruct, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	if message.Description != "" {
		e.dst.WriteString(fmt.Sprintf("%s// %s\n", indent, message.Description))
	}

	e.dst.WriteString(fmt.Sprintf("%sstruct %s {\n", indent, message.Name))

	// Fields: Traverse the fields and assign indexes
	for i, field := range message.Fields {
		e.encodeField(field, i+1, indentLevel+1) // Use 1-based indexing with `i+1`
	}

	e.dst.WriteString(fmt.Sprintf("%s}", indent))

	// Struct options
	if len(message.Options) > 0 {
		e.dst.WriteString(indent + "(\n")
		for i, option := range message.Options {
			if i > 0 {
				e.dst.WriteString(",\n")
			}
			e.dst.WriteString(indent + "    ") // Increase indentation
			e.encodeOption(option)
		}
		e.dst.WriteString("\n" + indent + ")\n")
	} else {
		e.dst.WriteString("\n")
	}

	e.dst.WriteString("\n")
}

// encodeUnion encodes a Thrift union
func (e *ThriftGenerate) encodeUnion(union *thrift.ThriftUnion, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	e.dst.WriteString(fmt.Sprintf("%sunion %s {\n", indent, union.Name))

	// Traverse union fields
	for i, field := range union.Fields {
		e.encodeField(field, i+1, indentLevel+1) // Use 1-based indexing with `i+1`
	}

	e.dst.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

// encodeEnum encodes enum types
func (e *ThriftGenerate) encodeEnum(enum *thrift.ThriftEnum, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	e.dst.WriteString(fmt.Sprintf("%senum %s {\n", indent, enum.Name))

	for _, value := range enum.Values {
		valueStr := fmt.Sprintf("%v", value.Value) // Convert the value to a string

		// Check if the value is numeric and generate a name if needed
		enumValueName := valueStr
		if _, err := strconv.Atoi(valueStr); err == nil {
			enumValueName = fmt.Sprintf("%s%s", enum.Name, valueStr)
		}

		enumValueName = strings.ToUpper(common.FormatStr(enumValueName))

		e.dst.WriteString(fmt.Sprintf("%s  %s = %d;\n", indent, enumValueName, value.Index))
	}

	e.dst.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

// encodeField encodes a single field within a struct
func (e *ThriftGenerate) encodeField(field *thrift.ThriftField, index, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	if field.Description != "" {
		e.dst.WriteString(fmt.Sprintf("%s// %s\n", indent, field.Description))
	}

	// Field index and type
	fieldType := field.Type
	if field.Repeated {
		fieldType = fmt.Sprintf("list<%s>", field.Type)
	}

	// Handle optional fields
	optionalFlag := ""
	if field.Optional {
		optionalFlag = "optional "
	}

	// Assign the provided index to the field
	e.dst.WriteString(fmt.Sprintf("%s%d: %s%s %s", indent, index, optionalFlag, fieldType, common.FormatStr(field.Name)))

	// Field options
	if len(field.Options) > 0 {
		e.dst.WriteString(" (")
		for i, option := range field.Options {
			if i > 0 {
				e.dst.WriteString(",\n" + indent)
			}
			e.encodeOption(option)
		}
		e.dst.WriteString(")\n") // Close parentheses aligned with the field
	} else {
		e.dst.WriteString("\n")
	}
}

// encodeOption handles the encoding of options for methods, structs, and fields
func (e *ThriftGenerate) encodeOption(option *thrift.Option) {
	// If the option key starts with "api", use double quotes for the value
	if strings.HasPrefix(option.Name, "api.") {
		e.dst.WriteString(fmt.Sprintf("%s = %s", option.Name, option.Value))
	} else if strings.HasPrefix(option.Name, "openapi.") {
		e.dst.WriteString(fmt.Sprintf("%s = '%s'", option.Name, option.Value))
	} else {
		e.dst.WriteString(fmt.Sprintf("%s = %s", option.Name, option.Value))
	}
}
