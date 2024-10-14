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

package utils

import (
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hertz-contrib/swagger-generate/common/utils"
)

// GetMethodName generates a method name from the OpenAPI spec
func GetMethodName(operation *openapi3.Operation, path, method string) string {
	if operation.OperationID != "" {
		return operation.OperationID
	}
	if operation.Tags != nil {
		return operation.Tags[0]
	}
	if path != "" {
		// Convert path to PascalCase, replacing placeholders with a suitable format
		convertedPath := ConvertPathToPascalCase(path)
		return convertedPath + strings.Title(strings.ToLower(method))
	}
	// If no OperationID, generate using HTTP method
	return strings.Title(strings.ToLower(method)) + "Method"
}

// GetServiceName generates a service name from the OpenAPI spec
func GetServiceName(operation *openapi3.Operation) string {
	if len(operation.Tags) > 0 {
		return operation.Tags[0]
	}
	return "DefaultService"
}

// GetMessageName generates a message name from the OpenAPI spec
func GetMessageName(operation *openapi3.Operation, methodName, suffix string) string {
	if operation.OperationID != "" {
		return operation.OperationID + suffix
	}
	return methodName + suffix
}

// GetPackageName generates a package name from the OpenAPI spec
func GetPackageName(spec *openapi3.T) string {
	if spec.Info.Title != "" {
		return utils.FormatStr(utils.ToSnakeCase(spec.Info.Title))
	}
	if spec.Info.Description != "" {
		return utils.FormatStr(utils.ToSnakeCase(spec.Info.Description))
	}
	return "DefaultPackage"
}

// ConvertPathToPascalCase converts a path with placeholders to PascalCase
func ConvertPathToPascalCase(path string) string {
	// Replace placeholders like {orderId} with OrderId
	re := regexp.MustCompile(`\{(\w+)\}`)
	path = re.ReplaceAllStringFunc(path, func(s string) string {
		return utils.ToPascaleCase(strings.Trim(s, "{}"))
	})

	// Split the path by '/' and convert each segment to PascalCase
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		segments[i] = utils.ToPascaleCase(segment)
	}

	// Join the segments back together
	return strings.Join(segments, "")
}

// ExtractMessageNameFromRef extracts the name of a message from a reference
func ExtractMessageNameFromRef(ref string) string {
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1] // Return the last part, usually the name of the reference
}

// ConvertPath converts a path with placeholders to a format that can be used in a URL
func ConvertPath(path string) string {
	// Regular expression to match content inside {}
	re := regexp.MustCompile(`\{(\w+)\}`)
	// Replace {param} with :param
	result := re.ReplaceAllString(path, ":$1")
	return result
}
