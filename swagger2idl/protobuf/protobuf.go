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

package protobuf

// ProtoFile represents a complete Proto file
type ProtoFile struct {
	PackageName string            // The package name of the Proto file
	Messages    []*ProtoMessage   // List of Proto messages
	Services    []*ProtoService   // List of Proto services
	Enums       []*ProtoEnum      // List of Proto enums
	Imports     []string          // List of imported Proto files
	Options     map[string]string // File-level options
}

// ProtoMessage represents a Proto message
type ProtoMessage struct {
	Name     string
	Fields   []*ProtoField   // List of fields in the Proto message
	Messages []*ProtoMessage // Nested Proto messages
	Enums    []*ProtoEnum    // Enums within the Proto message
	Options  []*Option       // Options specific to this Proto message
}

// ProtoField represents a field in a Proto message
type ProtoField struct {
	Name     string
	Type     string
	Repeated bool            // Indicates if the field is repeated (array)
	OneOf    []*ProtoField   // Collection of oneof fields (for representing oneOf types)
	Fields   []*ProtoField   // Nested fields (for object types)
	Messages []*ProtoMessage // Nested Proto messages within the field
	Enums    []*ProtoEnum    // Nested enums within the field
	Options  []*Option       // Additional options for this field
}

// Option represents an option in a Proto field or message
type Option struct {
	Name  string
	Value interface{}
}

// ProtoMethod represents a method in a Proto service
type ProtoMethod struct {
	Name    string
	Input   string    // Input message type
	Output  string    // Output message type
	Options []*Option // Options for the method
}

// ProtoService represents a Proto service
type ProtoService struct {
	Name    string         // Name of the service
	Methods []*ProtoMethod // List of methods in the service
	Options []*Option      // Service-level options
}

// ProtoEnum represents a Proto enum
type ProtoEnum struct {
	Name    string            // Name of the enum
	Values  []*ProtoEnumValue // Values within the enum
	Options []*Option         // Enum-level options
}

// ProtoEnumValue represents a value in a Proto enum
type ProtoEnumValue struct {
	Name  string // Name of the enum value
	Value int32  // Corresponding integer value for the enum
}
