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
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cloudwego/thriftgo/extension/thrift_option"
	"github.com/cloudwego/thriftgo/parser"
	"github.com/cloudwego/thriftgo/thrift_reflection"
)

// Contains returns true if an array Contains a specified string.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ParseStructOption(descriptor *thrift_reflection.StructDescriptor, optionName string, obj interface{}) error {
	opt, err := thrift_option.ParseStructOption(descriptor, optionName)
	if err != nil {
		return err
	}
	mapVal := opt.GetValue()
	mapValMap := mapVal.(map[string]interface{})
	jsonData, err := json.Marshal(mapValMap)
	err = json.Unmarshal(jsonData, obj)
	return err
}

func ParseServiceOption(descriptor *thrift_reflection.ServiceDescriptor, optionName string, obj interface{}) error {
	opt, err := thrift_option.ParseServiceOption(descriptor, optionName)
	if err != nil {
		return err
	}
	mapVal := opt.GetValue()
	mapValMap := mapVal.(map[string]interface{})
	jsonData, err := json.Marshal(mapValMap)
	err = json.Unmarshal(jsonData, obj)
	return err
}

func ParseMethodOption(descriptor *thrift_reflection.MethodDescriptor, optionName string, obj interface{}) error {
	opt, err := thrift_option.ParseMethodOption(descriptor, optionName)
	if err != nil {
		return err
	}
	mapVal := opt.GetValue()
	mapValMap := mapVal.(map[string]interface{})
	jsonData, err := json.Marshal(mapValMap)
	err = json.Unmarshal(jsonData, obj)
	return err
}

func ParseFieldOption(descriptor *thrift_reflection.FieldDescriptor, optionName string, obj interface{}) error {
	opt, err := thrift_option.ParseFieldOption(descriptor, optionName)
	if err != nil {
		return err
	}
	mapVal := opt.GetValue()
	mapValMap := mapVal.(map[string]interface{})
	jsonData, err := json.Marshal(mapValMap)
	err = json.Unmarshal(jsonData, obj)
	return err
}

func UnpackArgs(args []string, c interface{}) error {
	m, err := MapForm(args)
	if err != nil {
		return fmt.Errorf("unmarshal args failed, err: %v", err.Error())
	}

	t := reflect.TypeOf(c).Elem()
	v := reflect.ValueOf(c).Elem()
	if t.Kind() != reflect.Struct {
		return errors.New("passed c must be struct or pointer of struct")
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		x := v.Field(i)
		n := f.Name
		values, ok := m[n]
		if !ok || len(values) == 0 || values[0] == "" {
			continue
		}
		switch x.Kind() {
		case reflect.Bool:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			x.SetBool(values[0] == "true")
		case reflect.String:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			x.SetString(values[0])
		case reflect.Slice:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			ss := strings.Split(values[0], ";")
			if x.Type().Elem().Kind() == reflect.Int {
				n := reflect.MakeSlice(x.Type(), len(ss), len(ss))
				for i, s := range ss {
					val, err := strconv.ParseInt(s, 10, 64)
					if err != nil {
						return err
					}
					n.Index(i).SetInt(val)
				}
				x.Set(n)
			} else {
				for _, s := range ss {
					val := reflect.Append(x, reflect.ValueOf(s))
					x.Set(val)
				}
			}
		case reflect.Map:
			if len(values) != 1 {
				return fmt.Errorf("field %s can't be assigned multi values: %v", n, values)
			}
			ss := strings.Split(values[0], ";")
			out := make(map[string]string, len(ss))
			for _, s := range ss {
				sk := strings.SplitN(s, "=", 2)
				if len(sk) != 2 {
					return fmt.Errorf("map filed %v invalid key-value pair '%v'", n, s)
				}
				out[sk[0]] = sk[1]
			}
			x.Set(reflect.ValueOf(out))
		default:
			return fmt.Errorf("field %s has unsupported type %+v", n, f.Type)
		}
	}
	return nil
}

func MapForm(input []string) (map[string][]string, error) {
	out := make(map[string][]string, len(input))

	for _, str := range input {
		parts := strings.SplitN(str, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid argument: '%s'", str)
		}
		key, val := parts[0], parts[1]
		out[key] = append(out[key], val)
	}

	return out, nil
}

// MergeStructs merges non-zero fields from src into dst.
func MergeStructs(dst, src interface{}) error {
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)

	// Ensure both dst and src are pointers to structs.
	if dstVal.Kind() != reflect.Ptr || srcVal.Kind() != reflect.Ptr {
		return errors.New("both dst and src must be pointers")
	}
	if dstVal.Elem().Kind() != reflect.Struct || srcVal.Elem().Kind() != reflect.Struct {
		return errors.New("both dst and src must be pointers to structs")
	}

	dstVal = dstVal.Elem()
	srcVal = srcVal.Elem()

	for i := 0; i < dstVal.NumField(); i++ {
		field := dstVal.Field(i)
		srcField := srcVal.Field(i)

		if !srcField.IsZero() {
			field.Set(srcField)
		}
	}

	return nil
}

func GetAnnotation(input parser.Annotations, target string) []string {
	if len(input) == 0 {
		return nil
	}
	for _, anno := range input {
		if strings.ToLower(anno.Key) == target {
			return anno.Values
		}
	}

	return []string{}
}

func GetAnnotations(input parser.Annotations, targets map[string]string) map[string][]string {
	if len(input) == 0 || len(targets) == 0 {
		return nil
	}
	out := map[string][]string{}
	for k, t := range targets {
		var ret *parser.Annotation
		for _, anno := range input {
			if strings.ToLower(anno.Key) == k {
				ret = anno
				break
			}
		}
		if ret == nil {
			continue
		}
		out[t] = ret.Values
	}
	return out
}

func AppendUnique(s []string, e string) []string {
	if !Contains(s, e) {
		return append(s, e)
	}
	return s
}
