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

package plugins

import (
	"errors"
	"io"
	"os"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/thriftgo/plugin"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger/args"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger/generator"
)

func Run() int {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		logs.Errorf("Failed to get input: %v", err.Error())
		os.Exit(1)
	}

	req, err := plugin.UnmarshalRequest(data)
	if err != nil {
		logs.Errorf("Failed to unmarshal request: %v", err.Error())
		os.Exit(1)
	}

	if err := handleRequest(req); err != nil {
		logs.Errorf("Failed to handle request: %v", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
	return 0
}

func handleRequest(req *plugin.Request) (err error) {
	if req == nil {
		return errors.New("request is nil")
	}

	args := new(args.Arguments)
	if err = args.Unpack(req.PluginParameters); err != nil {
		return err
	}

	ast := req.GetAST()

	og := generator.NewOpenAPIGenerator(ast)
	openapiContent := og.BuildDocument(args)

	sg, err := generator.NewServerGenerator(ast, args)
	if err != nil {
		return err
	}
	serverContent, err := sg.Generate()
	if err != nil {
		return err
	}

	res := &plugin.Response{
		Contents: append(openapiContent, serverContent...),
	}
	if err = handleResponse(res); err != nil {
		return err
	}

	return err
}

func handleResponse(res *plugin.Response) error {
	data, err := plugin.MarshalResponse(res)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		return err
	}
	return nil
}
