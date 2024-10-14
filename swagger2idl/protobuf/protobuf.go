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
	PackageName string          // The package name of the Proto file
	Messages    []*ProtoMessage // List of Proto messages
	Services    []*ProtoService // List of Proto services
	Enums       []*ProtoEnum    // List of Proto enums
	Imports     []string        // List of imported Proto files
	Options     []*Option       // File-level options
}

// ProtoService represents a Proto service
type ProtoService struct {
	Name        string         // Name of the service
	Description string         // Description for the service
	Methods     []*ProtoMethod // List of methods in the service
	Options     []*Option      // Service-level options
}

// ProtoMethod represents a method in a Proto service
type ProtoMethod struct {
	Name        string    // Name of the method
	Description string    // Description for the method
	Input       string    // Input message type
	Output      string    // Output message type
	Options     []*Option // Options for the method
}

// ProtoMessage represents a Proto message
type ProtoMessage struct {
	Name        string
	Description string          // Description for the Proto message
	Fields      []*ProtoField   // List of fields in the Proto message
	Messages    []*ProtoMessage // Nested Proto messages
	Enums       []*ProtoEnum    // Enums within the Proto message
	OneOfs      []*ProtoOneOf   // OneOfs within the Proto message
	Options     []*Option       // Options specific to this Proto message
}

// ProtoField represents a field in a Proto message
type ProtoField struct {
	Name        string    // Name of the field
	Type        string    // Type of the field
	Description string    // Description for the field
	Repeated    bool      // Indicates if the field is repeated (array)
	Options     []*Option // Additional options for this field
}

// Option represents an option in a Proto field or message
type Option struct {
	Name  string      // Name of the option
	Value interface{} // Value of the option
}

// ProtoEnum represents a Proto enum
type ProtoEnum struct {
	Name        string            // Name of the enum
	Description string            // Description for the enum
	Values      []*ProtoEnumValue // Values within the enum
	Options     []*Option         // Enum-level options
}

// ProtoEnumValue represents a value in a Proto enum
type ProtoEnumValue struct {
	Index int // index of the enum value
	Value any // Corresponding integer value for the enum
}

// ProtoOneOf represents a oneof in a Proto message
type ProtoOneOf struct {
	Name    string        // Name of the oneof
	Fields  []*ProtoField // List of fields in the oneof
	Options []*Option     // Options specific to this oneof
}
