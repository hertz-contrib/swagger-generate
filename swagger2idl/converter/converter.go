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

import "github.com/hertz-contrib/swagger-generate/common/consts"

// Converter is an interface for converting files
type Converter interface {
	Convert() error
	GetIdl() interface{}
}

// ConvertOption adds a struct for conversion options
type ConvertOption struct {
	OpenapiOption bool
	ApiOption     bool
	NamingOption  bool
}

var MethodToOption = map[string]string{
	consts.HttpMethodGet:     consts.ApiGet,
	consts.HttpMethodPost:    consts.ApiPost,
	consts.HttpMethodPut:     consts.ApiPut,
	consts.HttpMethodPatch:   consts.ApiPatch,
	consts.HttpMethodDelete:  consts.ApiDelete,
	consts.HttpMethodOptions: consts.ApiOptions,
	consts.HttpMethodHead:    consts.ApiHEAD,
}
