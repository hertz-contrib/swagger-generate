// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

package openapi

import (
	"github.com/google/gnostic-models/compiler"
	"gopkg.in/yaml.v3"
)

// YAMLValue produces a serialized YAML representation of the document.
func (m *Document) YAMLValue(comment string) ([]byte, error) {
	rawInfo := m.ToRawInfo()
	rawInfo = &yaml.Node{
		Kind:        yaml.DocumentNode,
		Content:     []*yaml.Node{rawInfo},
		HeadComment: comment,
	}
	return yaml.Marshal(rawInfo)
}

// ToRawInfo returns a description of AdditionalPropertiesItem suitable for JSON or YAML export.
func (m *AdditionalPropertiesItem) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// AdditionalPropertiesItem
	// {Name:schemaOrReference Type:SchemaOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetSchemaOrReference()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:boolean Type:bool StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.Boolean
	if v1 != false {
		return compiler.NewScalarNodeForBool(v1)
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of Any suitable for JSON or YAML export.
func (m *Any) ToRawInfo() *yaml.Node {
	var err error
	var node yaml.Node
	err = yaml.Unmarshal([]byte(m.Yaml), &node)
	if err == nil {
		if node.Kind == yaml.DocumentNode {
			return node.Content[0]
		}
		return &node
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of AnyOrExpression suitable for JSON or YAML export.
func (m *AnyOrExpression) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// AnyOrExpression
	// {Name:any Type:Any StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetAny()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:expression Type:Expression StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.GetExpression()
	if v1 != nil {
		return v1.ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of Callback suitable for JSON or YAML export.
func (m *Callback) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Path != nil {
		for _, item := range m.Path {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of CallbackOrReference suitable for JSON or YAML export.
func (m *CallbackOrReference) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// CallbackOrReference
	// {Name:callback Type:Callback StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetCallback()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:reference Type:Reference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.GetReference()
	if v1 != nil {
		return v1.ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of CallbacksOrReferences suitable for JSON or YAML export.
func (m *CallbacksOrReferences) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Components suitable for JSON or YAML export.
func (m *Components) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Schemas != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("schemas"))
		info.Content = append(info.Content, m.Schemas.ToRawInfo())
	}
	if m.Responses != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("responses"))
		info.Content = append(info.Content, m.Responses.ToRawInfo())
	}
	if m.Parameters != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("parameters"))
		info.Content = append(info.Content, m.Parameters.ToRawInfo())
	}
	if m.Examples != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("examples"))
		info.Content = append(info.Content, m.Examples.ToRawInfo())
	}
	if m.RequestBodies != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("requestBodies"))
		info.Content = append(info.Content, m.RequestBodies.ToRawInfo())
	}
	if m.Headers != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("headers"))
		info.Content = append(info.Content, m.Headers.ToRawInfo())
	}
	if m.SecuritySchemes != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("securitySchemes"))
		info.Content = append(info.Content, m.SecuritySchemes.ToRawInfo())
	}
	if m.Links != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("links"))
		info.Content = append(info.Content, m.Links.ToRawInfo())
	}
	if m.Callbacks != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("callbacks"))
		info.Content = append(info.Content, m.Callbacks.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Contact suitable for JSON or YAML export.
func (m *Contact) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	if m.URL != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("url"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.URL))
	}
	if m.Email != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("email"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Email))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of DefaultType suitable for JSON or YAML export.
func (m *DefaultType) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// DefaultType
	if m.Number != 0 {
		return compiler.NewScalarNodeForFloat(m.Number)
	}
	if m.Boolean != false {
		return compiler.NewScalarNodeForBool(m.Boolean)
	}
	if m.String_ != "" {
		return compiler.NewScalarNodeForString(m.String_)
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of Discriminator suitable for JSON or YAML export.
func (m *Discriminator) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("propertyName"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.PropertyName))
	if m.Mapping != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("mapping"))
		info.Content = append(info.Content, m.Mapping.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Document suitable for JSON or YAML export.
func (m *Document) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("openapi"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Openapi))
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("info"))
	info.Content = append(info.Content, m.Info.ToRawInfo())
	if len(m.Servers) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.Servers {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("servers"))
		info.Content = append(info.Content, items)
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("paths"))
	info.Content = append(info.Content, m.Paths.ToRawInfo())
	if m.Components != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("components"))
		info.Content = append(info.Content, m.Components.ToRawInfo())
	}
	if len(m.Security) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.Security {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("security"))
		info.Content = append(info.Content, items)
	}
	if len(m.Tags) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.Tags {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("tags"))
		info.Content = append(info.Content, items)
	}
	if m.ExternalDocs != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("externalDocs"))
		info.Content = append(info.Content, m.ExternalDocs.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Encoding suitable for JSON or YAML export.
func (m *Encoding) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.ContentType != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("contentType"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.ContentType))
	}
	if m.Headers != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("headers"))
		info.Content = append(info.Content, m.Headers.ToRawInfo())
	}
	if m.Style != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("style"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Style))
	}
	if m.Explode != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("explode"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Explode))
	}
	if m.AllowReserved != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("allowReserved"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.AllowReserved))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Encodings suitable for JSON or YAML export.
func (m *Encodings) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Example suitable for JSON or YAML export.
func (m *Example) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Summary != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("summary"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Summary))
	}
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.Value != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("value"))
		info.Content = append(info.Content, m.Value.ToRawInfo())
	}
	if m.ExternalValue != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("externalValue"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.ExternalValue))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of ExampleOrReference suitable for JSON or YAML export.
func (m *ExampleOrReference) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// ExampleOrReference
	// {Name:example Type:Example StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetExample()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:reference Type:Reference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.GetReference()
	if v1 != nil {
		return v1.ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of ExamplesOrReferences suitable for JSON or YAML export.
func (m *ExamplesOrReferences) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Expression suitable for JSON or YAML export.
func (m *Expression) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of ExternalDocs suitable for JSON or YAML export.
func (m *ExternalDocs) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("url"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.URL))
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Header suitable for JSON or YAML export.
func (m *Header) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.Required != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("required"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Required))
	}
	if m.Deprecated != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("deprecated"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Deprecated))
	}
	if m.AllowEmptyValue != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("allowEmptyValue"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.AllowEmptyValue))
	}
	if m.Style != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("style"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Style))
	}
	if m.Explode != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("explode"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Explode))
	}
	if m.AllowReserved != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("allowReserved"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.AllowReserved))
	}
	if m.Schema != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("schema"))
		info.Content = append(info.Content, m.Schema.ToRawInfo())
	}
	if m.Example != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("example"))
		info.Content = append(info.Content, m.Example.ToRawInfo())
	}
	if m.Examples != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("examples"))
		info.Content = append(info.Content, m.Examples.ToRawInfo())
	}
	if m.Content != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("content"))
		info.Content = append(info.Content, m.Content.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of HeaderOrReference suitable for JSON or YAML export.
func (m *HeaderOrReference) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// HeaderOrReference
	// {Name:header Type:Header StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetHeader()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:reference Type:Reference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.GetReference()
	if v1 != nil {
		return v1.ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of HeadersOrReferences suitable for JSON or YAML export.
func (m *HeadersOrReferences) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Info suitable for JSON or YAML export.
func (m *Info) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("title"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Title))
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.TermsOfService != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("termsOfService"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.TermsOfService))
	}
	if m.Contact != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("contact"))
		info.Content = append(info.Content, m.Contact.ToRawInfo())
	}
	if m.License != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("license"))
		info.Content = append(info.Content, m.License.ToRawInfo())
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("version"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Version))
	if m.Summary != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("summary"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Summary))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of ItemsItem suitable for JSON or YAML export.
func (m *ItemsItem) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if len(m.SchemaOrReference) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.SchemaOrReference {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("schemaOrReference"))
		info.Content = append(info.Content, items)
	}
	return info
}

// ToRawInfo returns a description of License suitable for JSON or YAML export.
func (m *License) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	if m.URL != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("url"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.URL))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Link suitable for JSON or YAML export.
func (m *Link) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.OperationRef != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("operationRef"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.OperationRef))
	}
	if m.OperationID != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("operationId"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.OperationID))
	}
	if m.Parameters != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("parameters"))
		info.Content = append(info.Content, m.Parameters.ToRawInfo())
	}
	if m.RequestBody != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("requestBody"))
		info.Content = append(info.Content, m.RequestBody.ToRawInfo())
	}
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.Server != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("server"))
		info.Content = append(info.Content, m.Server.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of LinkOrReference suitable for JSON or YAML export.
func (m *LinkOrReference) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// LinkOrReference
	// {Name:link Type:Link StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetLink()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:reference Type:Reference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.GetReference()
	if v1 != nil {
		return v1.ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of LinksOrReferences suitable for JSON or YAML export.
func (m *LinksOrReferences) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of MediaType suitable for JSON or YAML export.
func (m *MediaType) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Schema != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("schema"))
		info.Content = append(info.Content, m.Schema.ToRawInfo())
	}
	if m.Example != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("example"))
		info.Content = append(info.Content, m.Example.ToRawInfo())
	}
	if m.Examples != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("examples"))
		info.Content = append(info.Content, m.Examples.ToRawInfo())
	}
	if m.Encoding != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("encoding"))
		info.Content = append(info.Content, m.Encoding.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of MediaTypes suitable for JSON or YAML export.
func (m *MediaTypes) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of NamedAny suitable for JSON or YAML export.
func (m *NamedAny) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	if m.Value != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("value"))
		info.Content = append(info.Content, m.Value.ToRawInfo())
	}
	return info
}

// ToRawInfo returns a description of NamedCallbackOrReference suitable for JSON or YAML export.
func (m *NamedCallbackOrReference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:CallbackOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedEncoding suitable for JSON or YAML export.
func (m *NamedEncoding) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:Encoding StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedExampleOrReference suitable for JSON or YAML export.
func (m *NamedExampleOrReference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:ExampleOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedHeaderOrReference suitable for JSON or YAML export.
func (m *NamedHeaderOrReference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:HeaderOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedLinkOrReference suitable for JSON or YAML export.
func (m *NamedLinkOrReference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:LinkOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedMediaType suitable for JSON or YAML export.
func (m *NamedMediaType) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:MediaType StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedParameterOrReference suitable for JSON or YAML export.
func (m *NamedParameterOrReference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:ParameterOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedPathItem suitable for JSON or YAML export.
func (m *NamedPathItem) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:PathItem StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedRequestBodyOrReference suitable for JSON or YAML export.
func (m *NamedRequestBodyOrReference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:RequestBodyOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedResponseOrReference suitable for JSON or YAML export.
func (m *NamedResponseOrReference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:ResponseOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedSchemaOrReference suitable for JSON or YAML export.
func (m *NamedSchemaOrReference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:SchemaOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedSecuritySchemeOrReference suitable for JSON or YAML export.
func (m *NamedSecuritySchemeOrReference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:SecuritySchemeOrReference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedServerVariable suitable for JSON or YAML export.
func (m *NamedServerVariable) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:ServerVariable StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of NamedString suitable for JSON or YAML export.
func (m *NamedString) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	if m.Value != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("value"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Value))
	}
	return info
}

// ToRawInfo returns a description of NamedStringArray suitable for JSON or YAML export.
func (m *NamedStringArray) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	// &{Name:value Type:StringArray StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:Mapped value}
	return info
}

// ToRawInfo returns a description of OauthFlow suitable for JSON or YAML export.
func (m *OauthFlow) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AuthorizationURL != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("authorizationUrl"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.AuthorizationURL))
	}
	if m.TokenURL != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("tokenUrl"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.TokenURL))
	}
	if m.RefreshURL != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("refreshUrl"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.RefreshURL))
	}
	if m.Scopes != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("scopes"))
		info.Content = append(info.Content, m.Scopes.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of OauthFlows suitable for JSON or YAML export.
func (m *OauthFlows) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Implicit != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("implicit"))
		info.Content = append(info.Content, m.Implicit.ToRawInfo())
	}
	if m.Password != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("password"))
		info.Content = append(info.Content, m.Password.ToRawInfo())
	}
	if m.ClientCredentials != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("clientCredentials"))
		info.Content = append(info.Content, m.ClientCredentials.ToRawInfo())
	}
	if m.AuthorizationCode != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("authorizationCode"))
		info.Content = append(info.Content, m.AuthorizationCode.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Object suitable for JSON or YAML export.
func (m *Object) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Operation suitable for JSON or YAML export.
func (m *Operation) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if len(m.Tags) != 0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("tags"))
		info.Content = append(info.Content, compiler.NewSequenceNodeForStringArray(m.Tags))
	}
	if m.Summary != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("summary"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Summary))
	}
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.ExternalDocs != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("externalDocs"))
		info.Content = append(info.Content, m.ExternalDocs.ToRawInfo())
	}
	if m.OperationID != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("operationId"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.OperationID))
	}
	if len(m.Parameters) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.Parameters {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("parameters"))
		info.Content = append(info.Content, items)
	}
	if m.RequestBody != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("requestBody"))
		info.Content = append(info.Content, m.RequestBody.ToRawInfo())
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("responses"))
	info.Content = append(info.Content, m.Responses.ToRawInfo())
	if m.Callbacks != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("callbacks"))
		info.Content = append(info.Content, m.Callbacks.ToRawInfo())
	}
	if m.Deprecated != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("deprecated"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Deprecated))
	}
	if len(m.Security) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.Security {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("security"))
		info.Content = append(info.Content, items)
	}
	if len(m.Servers) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.Servers {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("servers"))
		info.Content = append(info.Content, items)
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Parameter suitable for JSON or YAML export.
func (m *Parameter) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("in"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.In))
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.Required != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("required"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Required))
	}
	if m.Deprecated != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("deprecated"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Deprecated))
	}
	if m.AllowEmptyValue != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("allowEmptyValue"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.AllowEmptyValue))
	}
	if m.Style != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("style"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Style))
	}
	if m.Explode != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("explode"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Explode))
	}
	if m.AllowReserved != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("allowReserved"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.AllowReserved))
	}
	if m.Schema != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("schema"))
		info.Content = append(info.Content, m.Schema.ToRawInfo())
	}
	if m.Example != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("example"))
		info.Content = append(info.Content, m.Example.ToRawInfo())
	}
	if m.Examples != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("examples"))
		info.Content = append(info.Content, m.Examples.ToRawInfo())
	}
	if m.Content != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("content"))
		info.Content = append(info.Content, m.Content.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of ParameterOrReference suitable for JSON or YAML export.
func (m *ParameterOrReference) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// ParameterOrReference
	// {Name:parameter Type:Parameter StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetParameter()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:reference Type:Reference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.GetReference()
	if v1 != nil {
		return v1.ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of ParametersOrReferences suitable for JSON or YAML export.
func (m *ParametersOrReferences) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of PathItem suitable for JSON or YAML export.
func (m *PathItem) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Xref != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("$ref"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Xref))
	}
	if m.Summary != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("summary"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Summary))
	}
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.Get != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("get"))
		info.Content = append(info.Content, m.Get.ToRawInfo())
	}
	if m.Put != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("put"))
		info.Content = append(info.Content, m.Put.ToRawInfo())
	}
	if m.Post != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("post"))
		info.Content = append(info.Content, m.Post.ToRawInfo())
	}
	if m.Delete != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("delete"))
		info.Content = append(info.Content, m.Delete.ToRawInfo())
	}
	if m.Options != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("options"))
		info.Content = append(info.Content, m.Options.ToRawInfo())
	}
	if m.Head != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("head"))
		info.Content = append(info.Content, m.Head.ToRawInfo())
	}
	if m.Patch != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("patch"))
		info.Content = append(info.Content, m.Patch.ToRawInfo())
	}
	if m.Trace != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("trace"))
		info.Content = append(info.Content, m.Trace.ToRawInfo())
	}
	if len(m.Servers) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.Servers {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("servers"))
		info.Content = append(info.Content, items)
	}
	if len(m.Parameters) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.Parameters {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("parameters"))
		info.Content = append(info.Content, items)
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Paths suitable for JSON or YAML export.
func (m *Paths) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Path != nil {
		for _, item := range m.Path {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Properties suitable for JSON or YAML export.
func (m *Properties) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Reference suitable for JSON or YAML export.
func (m *Reference) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("$ref"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Xref))
	if m.Summary != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("summary"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Summary))
	}
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	return info
}

// ToRawInfo returns a description of RequestBodiesOrReferences suitable for JSON or YAML export.
func (m *RequestBodiesOrReferences) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of RequestBody suitable for JSON or YAML export.
func (m *RequestBody) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("content"))
	info.Content = append(info.Content, m.Content.ToRawInfo())
	if m.Required != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("required"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Required))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of RequestBodyOrReference suitable for JSON or YAML export.
func (m *RequestBodyOrReference) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// RequestBodyOrReference
	// {Name:requestBody Type:RequestBody StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetRequestBody()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:reference Type:Reference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.GetReference()
	if v1 != nil {
		return v1.ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of Response suitable for JSON or YAML export.
func (m *Response) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	if m.Headers != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("headers"))
		info.Content = append(info.Content, m.Headers.ToRawInfo())
	}
	if m.Content != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("content"))
		info.Content = append(info.Content, m.Content.ToRawInfo())
	}
	if m.Links != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("links"))
		info.Content = append(info.Content, m.Links.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of ResponseOrReference suitable for JSON or YAML export.
func (m *ResponseOrReference) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// ResponseOrReference
	// {Name:response Type:Response StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetResponse()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:reference Type:Reference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.GetReference()
	if v1 != nil {
		return v1.ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of Responses suitable for JSON or YAML export.
func (m *Responses) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Default != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("default"))
		info.Content = append(info.Content, m.Default.ToRawInfo())
	}
	if m.ResponseOrReference != nil {
		for _, item := range m.ResponseOrReference {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of ResponsesOrReferences suitable for JSON or YAML export.
func (m *ResponsesOrReferences) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Schema suitable for JSON or YAML export.
func (m *Schema) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Nullable != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("nullable"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Nullable))
	}
	if m.Discriminator != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("discriminator"))
		info.Content = append(info.Content, m.Discriminator.ToRawInfo())
	}
	if m.ReadOnly != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("readOnly"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.ReadOnly))
	}
	if m.WriteOnly != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("writeOnly"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.WriteOnly))
	}
	if m.XML != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("xml"))
		info.Content = append(info.Content, m.XML.ToRawInfo())
	}
	if m.ExternalDocs != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("externalDocs"))
		info.Content = append(info.Content, m.ExternalDocs.ToRawInfo())
	}
	if m.Example != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("example"))
		info.Content = append(info.Content, m.Example.ToRawInfo())
	}
	if m.Deprecated != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("deprecated"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Deprecated))
	}
	if m.Title != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("title"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Title))
	}
	if m.MultipleOf != 0.0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("multipleOf"))
		info.Content = append(info.Content, compiler.NewScalarNodeForFloat(m.MultipleOf))
	}
	if m.Maximum != 0.0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("maximum"))
		info.Content = append(info.Content, compiler.NewScalarNodeForFloat(m.Maximum))
	}
	if m.ExclusiveMaximum != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("exclusiveMaximum"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.ExclusiveMaximum))
	}
	if m.Minimum != 0.0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("minimum"))
		info.Content = append(info.Content, compiler.NewScalarNodeForFloat(m.Minimum))
	}
	if m.ExclusiveMinimum != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("exclusiveMinimum"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.ExclusiveMinimum))
	}
	if m.MaxLength != 0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("maxLength"))
		info.Content = append(info.Content, compiler.NewScalarNodeForInt(m.MaxLength))
	}
	if m.MinLength != 0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("minLength"))
		info.Content = append(info.Content, compiler.NewScalarNodeForInt(m.MinLength))
	}
	if m.Pattern != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("pattern"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Pattern))
	}
	if m.MaxItems != 0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("maxItems"))
		info.Content = append(info.Content, compiler.NewScalarNodeForInt(m.MaxItems))
	}
	if m.MinItems != 0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("minItems"))
		info.Content = append(info.Content, compiler.NewScalarNodeForInt(m.MinItems))
	}
	if m.UniqueItems != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("uniqueItems"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.UniqueItems))
	}
	if m.MaxProperties != 0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("maxProperties"))
		info.Content = append(info.Content, compiler.NewScalarNodeForInt(m.MaxProperties))
	}
	if m.MinProperties != 0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("minProperties"))
		info.Content = append(info.Content, compiler.NewScalarNodeForInt(m.MinProperties))
	}
	if len(m.Required) != 0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("required"))
		info.Content = append(info.Content, compiler.NewSequenceNodeForStringArray(m.Required))
	}
	if len(m.Enum) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.Enum {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("enum"))
		info.Content = append(info.Content, items)
	}
	if m.Type != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("type"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Type))
	}
	if len(m.AllOf) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.AllOf {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("allOf"))
		info.Content = append(info.Content, items)
	}
	if len(m.OneOf) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.OneOf {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("oneOf"))
		info.Content = append(info.Content, items)
	}
	if len(m.AnyOf) != 0 {
		items := compiler.NewSequenceNode()
		for _, item := range m.AnyOf {
			items.Content = append(items.Content, item.ToRawInfo())
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("anyOf"))
		info.Content = append(info.Content, items)
	}
	if m.Not != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("not"))
		info.Content = append(info.Content, m.Not.ToRawInfo())
	}
	if m.Items != nil {

		items := compiler.NewSequenceNode()
		for _, item := range m.Items.SchemaOrReference {
			//if item == nil {
			//	continue
			//}
			items.Content = append(items.Content, item.ToRawInfo())
		}
		if len(items.Content) == 1 {
			items = items.Content[0]
		}
		info.Content = append(info.Content, compiler.NewScalarNodeForString("items"))
		info.Content = append(info.Content, items)
	}
	if m.Properties != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("properties"))
		info.Content = append(info.Content, m.Properties.ToRawInfo())
	}
	if m.AdditionalProperties != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("additionalProperties"))
		info.Content = append(info.Content, m.AdditionalProperties.ToRawInfo())
	}
	if m.Default != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("default"))
		info.Content = append(info.Content, m.Default.ToRawInfo())
	}
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.Format != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("format"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Format))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of SchemaOrReference suitable for JSON or YAML export.
func (m *SchemaOrReference) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// SchemaOrReference
	// {Name:schema Type:Schema StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}

	if m.IsSetSchema() == true {
		return m.GetSchema().ToRawInfo()
	}
	// {Name:reference Type:Reference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	if m.IsSetReference() == true {
		return m.GetReference().ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of SchemasOrReferences suitable for JSON or YAML export.
func (m *SchemasOrReferences) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of SecurityRequirement suitable for JSON or YAML export.
func (m *SecurityRequirement) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of SecurityScheme suitable for JSON or YAML export.
func (m *SecurityScheme) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("type"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m._Type))
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	if m._In != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("in"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m._In))
	}
	if m.Scheme != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("scheme"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Scheme))
	}
	if m.BearerFormat != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("bearerFormat"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.BearerFormat))
	}
	if m.Flows != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("flows"))
		info.Content = append(info.Content, m.Flows.ToRawInfo())
	}
	if m.OpenIDConnectURL != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("openIdConnectUrl"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.OpenIDConnectURL))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of SecuritySchemeOrReference suitable for JSON or YAML export.
func (m *SecuritySchemeOrReference) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// SecuritySchemeOrReference
	// {Name:securityScheme Type:SecurityScheme StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v0 := m.GetSecurityScheme()
	if v0 != nil {
		return v0.ToRawInfo()
	}
	// {Name:reference Type:Reference StringEnumValues:[] MapType: Repeated:false Pattern: Implicit:false Description:}
	v1 := m.GetReference()
	if v1 != nil {
		return v1.ToRawInfo()
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of SecuritySchemesOrReferences suitable for JSON or YAML export.
func (m *SecuritySchemesOrReferences) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Server suitable for JSON or YAML export.
func (m *Server) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("url"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.URL))
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.Variables != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("variables"))
		info.Content = append(info.Content, m.Variables.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of ServerVariable suitable for JSON or YAML export.
func (m *ServerVariable) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if len(m.Enum) != 0 {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("enum"))
		info.Content = append(info.Content, compiler.NewSequenceNodeForStringArray(m.Enum))
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("default"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m._Default))
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of ServerVariables suitable for JSON or YAML export.
func (m *ServerVariables) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.AdditionalProperties != nil {
		for _, item := range m.AdditionalProperties {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of SpecificationExtension suitable for JSON or YAML export.
func (m *SpecificationExtension) ToRawInfo() *yaml.Node {
	// ONE OF WRAPPER
	// SpecificationExtension
	if m.Number != 0 {
		return compiler.NewScalarNodeForFloat(m.Number)
	}

	if m.Boolean != false {
		return compiler.NewScalarNodeForBool(m.Boolean)
	}

	if m.String_ != "" {
		return compiler.NewScalarNodeForString(m.String_)
	}
	return compiler.NewNullNode()
}

// ToRawInfo returns a description of StringArray suitable for JSON or YAML export.
func (m *StringArray) ToRawInfo() *yaml.Node {
	return compiler.NewSequenceNodeForStringArray(m.Values)
}

// ToRawInfo returns a description of Strings suitable for JSON or YAML export.
func (m *Strings) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// &{Name:additionalProperties Type:NamedString StringEnumValues:[] MapType:string Repeated:true Pattern: Implicit:true Description:}
	return info
}

// ToRawInfo returns a description of Tag suitable for JSON or YAML export.
func (m *Tag) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	// always include this required field.
	info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
	info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	if m.Description != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("description"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Description))
	}
	if m.ExternalDocs != nil {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("externalDocs"))
		info.Content = append(info.Content, m.ExternalDocs.ToRawInfo())
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}

// ToRawInfo returns a description of Xml suitable for JSON or YAML export.
func (m *Xml) ToRawInfo() *yaml.Node {
	info := compiler.NewMappingNode()
	if m == nil {
		return info
	}
	if m.Name != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("name"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Name))
	}
	if m.Namespace != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("namespace"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Namespace))
	}
	if m.Prefix != "" {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("prefix"))
		info.Content = append(info.Content, compiler.NewScalarNodeForString(m.Prefix))
	}
	if m.Attribute != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("attribute"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Attribute))
	}
	if m.Wrapped != false {
		info.Content = append(info.Content, compiler.NewScalarNodeForString("wrapped"))
		info.Content = append(info.Content, compiler.NewScalarNodeForBool(m.Wrapped))
	}
	if m.SpecificationExtension != nil {
		for _, item := range m.SpecificationExtension {
			info.Content = append(info.Content, compiler.NewScalarNodeForString(item.Name))
			info.Content = append(info.Content, item.Value.ToRawInfo())
		}
	}
	return info
}
