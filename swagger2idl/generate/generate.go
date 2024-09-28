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
	"strings"

	"github.com/hertz-contrib/swagger-generate/swagger2idl/protobuf"
	"github.com/hertz-contrib/swagger-generate/swagger2idl/utils"
)

// Encoder is used to handle the encoding context
type Encoder struct {
	dst *strings.Builder // The target for output
}

// ConvertToProtoFile converts the ProtoFile structure into Proto file content
func ConvertToProtoFile(protoFile *protobuf.ProtoFile) string {
	var sb strings.Builder
	encoder := &Encoder{dst: &sb}

	encoder.dst.WriteString("syntax = \"proto3\";\n\n")
	encoder.dst.WriteString(fmt.Sprintf("package %s;\n\n", protoFile.PackageName))

	// Generate imports
	for _, importFile := range protoFile.Imports {
		encoder.dst.WriteString(fmt.Sprintf("import \"%s\";\n", importFile))
	}
	if len(protoFile.Imports) > 0 {
		encoder.dst.WriteString("\n")
	}

	// Generate file-level options
	for key, value := range protoFile.Options {
		encoder.dst.WriteString(fmt.Sprintf("option %s = %s;\n", key, utils.Stringify(value)))
	}
	if len(protoFile.Options) > 0 {
		encoder.dst.WriteString("\n")
	}

	// Sort messages by name
	sort.Slice(protoFile.Messages, func(i, j int) bool {
		return protoFile.Messages[i].Name < protoFile.Messages[j].Name
	})

	// Generate messages
	for _, message := range protoFile.Messages {
		encoder.encodeMessage(message, 0)
	}

	// Sort services by name
	sort.Slice(protoFile.Services, func(i, j int) bool {
		return protoFile.Services[i].Name < protoFile.Services[j].Name
	})

	// Generate services
	for _, service := range protoFile.Services {
		encoder.dst.WriteString(fmt.Sprintf("service %s {\n", service.Name))

		// Sort methods by name
		sort.Slice(service.Methods, func(i, j int) bool {
			return service.Methods[i].Name < service.Methods[j].Name
		})

		for _, method := range service.Methods {
			encoder.dst.WriteString(fmt.Sprintf("  rpc %s(%s) returns (%s)", method.Name, method.Input, method.Output))
			if len(method.Options) > 0 {
				encoder.dst.WriteString(" {\n")
				for _, option := range method.Options {
					encoder.dst.WriteString("     option ")
					encoder.encodeFieldOption(option)
					encoder.dst.WriteString(";\n")
				}
				encoder.dst.WriteString("  }\n")
			} else {
				encoder.dst.WriteString(";\n")
			}
		}
		encoder.dst.WriteString("}\n\n")
	}

	return encoder.dst.String()
}

// encodeMessage recursively encodes messages, including nested messages and enums
func (e *Encoder) encodeMessage(message *protobuf.ProtoMessage, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)
	e.dst.WriteString(fmt.Sprintf("%smessage %s {\n", indent, message.Name))

	// Generate message-level options
	if len(message.Options) > 0 {
		e.dst.WriteString(fmt.Sprintf("%s  option", indent))
		for _, option := range message.Options {
			e.encodeFieldOption(option)
			e.dst.WriteString(";\n")
		}
	}

	// Sort fields by name
	sort.Slice(message.Fields, func(i, j int) bool {
		return message.Fields[i].Name < message.Fields[j].Name
	})

	// Generate fields
	for i, field := range message.Fields {
		repeated := ""
		if field.Repeated {
			repeated = "repeated "
		}
		e.dst.WriteString(fmt.Sprintf("%s  %s%s %s = %d", indent, repeated, field.Type, field.Name, i+1))

		// Generate field-level options
		if len(field.Options) > 0 {
			e.dst.WriteString(" [\n    ")
			for j, option := range field.Options {
				e.encodeFieldOption(option)
				if j < len(field.Options)-1 {
					e.dst.WriteString(", ")
				}
			}
			e.dst.WriteString("]")
		}
		e.dst.WriteString(";\n")
	}

	// Recursively handle nested messages
	for _, nestedMessage := range message.Messages {
		e.encodeMessage(nestedMessage, indentLevel+1) // Increase indentation
	}

	e.dst.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

// encodeFieldOption encodes an option for a single field
func (e *Encoder) encodeFieldOption(opt *protobuf.Option) error {
	// Output the option name
	fmt.Fprintf(e.dst, "(%s) = ", opt.Name) // Add indentation for consistency

	// Check if the option value is a complex structure
	switch value := opt.Value.(type) {
	case map[string]interface{}:
		// If it's a map type, it needs to output as a nested structure
		fmt.Fprintf(e.dst, "{\n")        // Newline after {
		e.encodeFieldOptionMap(value, 6) // Output map content, passing the current indentation level
		fmt.Fprintf(e.dst, "    }")      // Indent and output the closing }, with the appropriate indentation level
	default:
		fmt.Fprintf(e.dst, "%s", value) // For simple types, output directly
	}

	return nil
}

// encodeFieldOptionMap encodes a complex map type option value
func (e *Encoder) encodeFieldOptionMap(optionMap map[string]interface{}, indent int) error {
	keys := make([]string, 0, len(optionMap))
	for k := range optionMap {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Sort keys to ensure consistent output order

	indentSpace := strings.Repeat(" ", indent) // Dynamically generate indent spaces

	for _, key := range keys {
		value := optionMap[key]
		// Output key-value pairs with appropriate indentation
		fmt.Fprintf(e.dst, "%s%s: %s", indentSpace, key, utils.Stringify(value)) // Add deeper indentation
		// Don't add a semicolon after the last item, maintain correct format
		fmt.Fprintf(e.dst, ";\n")
	}

	return nil
}
