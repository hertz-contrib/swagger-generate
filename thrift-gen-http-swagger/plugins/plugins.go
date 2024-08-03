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
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cloudwego/thriftgo/plugin"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-http-swagger/args"
	"github.com/hertz-contrib/swagger-generate/thrift-gen-http-swagger/generator"
)

var flags flag.FlagSet

func Run() int {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		println("Failed to get input:", err.Error())
		os.Exit(1)
	}

	req, err := plugin.UnmarshalRequest(data)
	if err != nil {
		println("Failed to unmarshal request:", err.Error())
		os.Exit(1)
	}

	if err := handleRequest(req); err != nil {
		println("Failed to handle request:", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
	return 0
}

func handleRequest(req *plugin.Request) (err error) {
	args := new(args.Arguments)
	if err := args.Unpack(req.PluginParameters); err != nil {
		log.Printf("[Error]: unpack args failed: %s", err.Error())
		return err
	}

	if req == nil {
		fmt.Fprintf(os.Stderr, "unexpected nil request")
	}
	ast := req.GetAST()
	g := generator.NewOpenAPIGenerator(ast)
	contents := g.BuildDocument(args)
	res := &plugin.Response{
		Contents: contents,
	}
	if err := handleResponse(res); err != nil {
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
