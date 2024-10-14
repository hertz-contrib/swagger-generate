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

package thrift

// ThriftFile represents a complete Thrift file
type ThriftFile struct {
	Namespace map[string]string // Namespace for the Thrift file
	Includes  []string          // List of included Thrift files
	Structs   []*ThriftStruct   // List of Thrift structs
	Unions    []*ThriftUnion    // List of Thrift unions
	Enums     []*ThriftEnum     // List of Thrift enums
	Services  []*ThriftService  // List of Thrift services
}

// ThriftService represents a Thrift service
type ThriftService struct {
	Name        string          // Name of the service
	Description string          // Description of the service
	Methods     []*ThriftMethod // List of methods in the service
	Options     []*Option       // Service-level options
}

// ThriftMethod represents a method in a Thrift service
type ThriftMethod struct {
	Name        string    // Name of the method
	Description string    // Description of the method
	Input       []string  // List of input fields for the method
	Output      string    // Output field for the method
	Options     []*Option // Options for the method
}

// ThriftStruct represents a Thrift struct
type ThriftStruct struct {
	Name        string         // Name of the struct
	Description string         // Description of the struct
	Fields      []*ThriftField // List of fields in the struct
	Options     []*Option      // Options specific to this struct
}

// ThriftField represents a field in a Thrift struct or union
type ThriftField struct {
	ID          int       // Field ID for Thrift
	Name        string    // Name of the field
	Description string    // Description of the field
	Type        string    // Type of the field (Thrift types)
	Optional    bool      // Indicates if the field is optional
	Repeated    bool      // Indicates if the field is repeated (list)
	Options     []*Option // Additional options for this field
}

// ThriftUnion represents a Thrift union (similar to a struct but only one field can be set at a time)
type ThriftUnion struct {
	Name    string         // Name of the union
	Fields  []*ThriftField // List of fields in the union
	Options []*Option      // Options specific to this union
}

// ThriftEnum represents a Thrift enum
type ThriftEnum struct {
	Name        string             // Name of the enum
	Description string             // Description of the enum
	Values      []*ThriftEnumValue // Values within the enum
	Options     []*Option          // Enum-level options
}

// ThriftEnumValue represents a value in a Thrift enum
type ThriftEnumValue struct {
	Index int // Index of the enum value
	Value any // Enum values are integers in Thrift
}

// Option represents an option in a Thrift field or struct
type Option struct {
	Name  string      // Name of the option
	Value interface{} // Value of the option
}
