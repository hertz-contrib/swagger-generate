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

package parser

import (
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

// LoadOpenAPISpec parses an OpenAPI spec from a file and returns it.
func LoadOpenAPISpec(filePath string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	var err error
	var spec *openapi3.T

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", filePath)
	}

	spec, err = loader.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %v", err)
	}

	if err = spec.Validate(loader.Context); err != nil {
		return nil, fmt.Errorf("failed to validate OpenAPI spec: %v", err)
	}

	return spec, nil
}
