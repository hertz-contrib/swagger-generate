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
	"encoding/json"
	"errors"

	"github.com/cloudwego/thriftgo/extension/thrift_option"
	"github.com/cloudwego/thriftgo/thrift_reflection"
)

func ParseStructOption(descriptor *thrift_reflection.StructDescriptor, optionName string, obj interface{}) error {
	opt, err := thrift_option.ParseStructOption(descriptor, optionName)
	if errors.Is(err, thrift_option.ErrKeyNotMatch) ||
		errors.Is(err, thrift_option.ErrNotIncluded) ||
		errors.Is(err, thrift_option.ErrNotExistOption) {
		return nil
	}
	if err != nil {
		return err
	}
	mapVal := opt.GetValue()
	mapValMap := mapVal.(map[string]interface{})
	jsonData, err := json.Marshal(mapValMap)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, obj); err != nil {
		return err
	}
	return err
}

func ParseServiceOption(descriptor *thrift_reflection.ServiceDescriptor, optionName string, obj interface{}) error {
	opt, err := thrift_option.ParseServiceOption(descriptor, optionName)
	if errors.Is(err, thrift_option.ErrKeyNotMatch) ||
		errors.Is(err, thrift_option.ErrNotIncluded) ||
		errors.Is(err, thrift_option.ErrNotExistOption) {
		return nil
	}
	if err != nil {
		return err
	}
	mapVal := opt.GetValue()
	mapValMap := mapVal.(map[string]interface{})
	jsonData, err := json.Marshal(mapValMap)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, obj); err != nil {
		return err
	}
	return err
}

func ParseMethodOption(descriptor *thrift_reflection.MethodDescriptor, optionName string, obj interface{}) error {
	opt, err := thrift_option.ParseMethodOption(descriptor, optionName)
	if errors.Is(err, thrift_option.ErrKeyNotMatch) ||
		errors.Is(err, thrift_option.ErrNotIncluded) ||
		errors.Is(err, thrift_option.ErrNotExistOption) {
		return nil
	}
	if err != nil {
		return err
	}
	mapVal := opt.GetValue()
	mapValMap := mapVal.(map[string]interface{})
	jsonData, err := json.Marshal(mapValMap)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, obj); err != nil {
		return err
	}
	return err
}

func ParseFieldOption(descriptor *thrift_reflection.FieldDescriptor, optionName string, obj interface{}) error {
	opt, err := thrift_option.ParseFieldOption(descriptor, optionName)
	if errors.Is(err, thrift_option.ErrKeyNotMatch) ||
		errors.Is(err, thrift_option.ErrNotIncluded) ||
		errors.Is(err, thrift_option.ErrNotExistOption) {
		return nil
	}
	if err != nil {
		return err
	}
	mapVal := opt.GetValue()
	mapValMap := mapVal.(map[string]interface{})
	jsonData, err := json.Marshal(mapValMap)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jsonData, obj); err != nil {
		return err
	}
	return err
}
