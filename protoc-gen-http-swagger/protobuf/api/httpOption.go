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

package api

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

var HttpMethodOptions = map[*protoimpl.ExtensionInfo]string{
	E_Get:     "GET",
	E_Post:    "POST",
	E_Put:     "PUT",
	E_Patch:   "PATCH",
	E_Delete:  "DELETE",
	E_Options: "OPTIONS",
	E_Head:    "HEAD",
}

func GetAllOptions(extensions map[*protoimpl.ExtensionInfo]string, opts ...protoreflect.ProtoMessage) map[string]interface{} {
	out := map[string]interface{}{}
	for _, opt := range opts {
		for e, t := range extensions {
			if proto.HasExtension(opt, e) {
				v := proto.GetExtension(opt, e)
				out[t] = v
			}
		}
	}
	return out
}
