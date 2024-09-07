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

package generator

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

type ServerConfiguration struct {
	HertzAddr  *string
	KitexAddr  *string
	OutputMode *string
}

type ServerGenerator struct {
	IdlPath         []string
	HertzAddr       string
	KitexAddr       string
	SwaggerFileName string
	DirectUrl       string
}

func NewServerGenerator(conf ServerConfiguration, inputFiles []*protogen.File) (*ServerGenerator, error) {
	hertzAddr := conf.HertzAddr
	if hertzAddr == nil {
		*hertzAddr = DefaultHertzAddr
	}

	kitexAddr := conf.KitexAddr
	if kitexAddr == nil {
		*kitexAddr = DefaultKitexAddr
	}

	if conf.OutputMode == nil || (*conf.OutputMode != "merged" && *conf.OutputMode != "source_relative") {
		return nil, errors.New("invalid or missing output mode, must be 'merged' or 'source_relative'")
	}

	// Collect paths from input files that should be generated
	var idlPath []string
	for _, file := range inputFiles {
		if file.Generate {
			idlPath = append(idlPath, file.Desc.Path())
		}
	}

	// Check if there are any .proto files to process
	if len(idlPath) == 0 {
		return nil, errors.New("no .proto files marked for generation")
	}

	// Determine the swagger file name and direct URL based on the output mode
	swaggerFileName := MergedOpenapiFileName
	directUrl := MergedOpenapiFileName
	if *conf.OutputMode == "source_relative" {
		swaggerFileName = SourceRelativeOpenapiFileName
		directUrl = strings.TrimSuffix(idlPath[0], ProtoSuffix) + OpenapiSuffix
	}

	// Check if Hertz and Kitex addresses are valid (basic validation)
	if err := validateAddress(*hertzAddr); err != nil {
		return nil, fmt.Errorf("invalid Hertz address: %w", err)
	}
	if err := validateAddress(*kitexAddr); err != nil {
		return nil, fmt.Errorf("invalid Kitex address: %w", err)
	}

	return &ServerGenerator{
		IdlPath:         idlPath,
		HertzAddr:       *hertzAddr,
		KitexAddr:       *kitexAddr,
		SwaggerFileName: swaggerFileName,
		DirectUrl:       directUrl,
	}, nil
}

func validateAddress(addr string) error {
	if addr == "" {
		return errors.New("address cannot be empty")
	}
	if !strings.Contains(addr, ":") {
		return errors.New("address must include a port (e.g., '127.0.0.1:8080')")
	}
	return nil
}

func (g *ServerGenerator) Generate(outputFile *protogen.GeneratedFile) error {
	tmpl, err := template.New("server").Delims("{{", "}}").Parse(serverTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, g)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if _, err = outputFile.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}
	return nil
}

const (
	DefaultHertzAddr              = "127.0.0.1:8080"
	DefaultKitexAddr              = "127.0.0.1:8888"
	MergedOpenapiFileName         = "openapi.yaml"
	SourceRelativeOpenapiFileName = "*.openapi.yaml"
	OpenapiSuffix                 = ".openapi.yaml"
	ProtoSuffix                   = ".proto"
)

const serverTemplate = `// Code generated by thrift-gen-rpc-swagger.
package main

import (
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bytedance/gopkg/cloud/metainfo"
	dproto "github.com/cloudwego/dynamicgo/proto"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	"github.com/emicklei/proto"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

//go:embed {{.SwaggerFileName}}
var files embed.FS

type ClientPool struct {
	serviceMap map[string]genericclient.Client
	mutex      sync.RWMutex
}

func NewClientPool(protoFiles []string) *ClientPool {
	clientPool := &ClientPool{
		serviceMap: make(map[string]genericclient.Client),
	}

	for _, protoFile := range protoFiles {
		filePath, err := findPbFile(protoFile)
		if err != nil {
			hlog.Fatalf("Error finding proto file: %v", err)
		}

		err = clientPool.GetServicesFromIDL(filePath)
		if err != nil {
			hlog.Fatalf("Error loading protobuf files from directory: %v", err)
		}
	}

	return clientPool
}

func newClient(pbFilePath, svcName string) genericclient.Client {
	dOpts := dproto.Options{}
	p, err := generic.NewPbFileProviderWithDynamicGo(pbFilePath, context.Background(), dOpts)
	if err != nil {
		hlog.Fatalf("Failed to create protobufFileProvider for %s: %v", svcName, err)
	}

	g, err := generic.JSONPbGeneric(p)
	if err != nil {
		hlog.Fatalf("Failed to create JSONPbGeneric for %s: %v", svcName, err)
	}

	cli, err := genericclient.NewClient(svcName, g,
		client.WithTransportProtocol(transport.TTHeader),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		client.WithHostPorts("{{.KitexAddr}}"),
	)
	if err != nil {
		hlog.Fatalf("Failed to create generic client for %s: %v", svcName, err)
	}

	return cli
}

func (cp *ClientPool) getClient(svcName string) (genericclient.Client, error) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	client, ok := cp.serviceMap[svcName]
	if !ok {
		return nil, errors.New("service not found: " + svcName)
	}
	return client, nil
}

func (cp *ClientPool) GetServicesFromIDL(idlPath string) error {
	reader, err := os.Open(idlPath)
	if err != nil {
		return fmt.Errorf("failed to open proto file: %w", err)
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse proto file: %w", err)
	}

	proto.Walk(definition,
		proto.WithService(func(s *proto.Service) {
			cp.serviceMap[s.Name] = newClient(idlPath, s.Name)
		}),
	)

	return nil
}

func main() {
	h := server.Default(server.WithHostPorts("{{.HertzAddr}}"))
	h.Use(cors.Default())

	protoFiles := []string{
		{{- range .IdlPath }}
		"{{ . }}",
		{{- end }}
	}

	clientPool := NewClientPool(protoFiles)

	setupSwaggerRoutes(h)
	setupProxyRoutes(h, clientPool)

	hlog.Info("Swagger UI is available at: http://{{.HertzAddr}}/swagger/index.html")
	h.Spin()
}

func setupSwaggerRoutes(h *server.Hertz) {
	h.GET("swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, swagger.URL("/{{.DirectUrl}}")))

	h.GET("/:filename", func(c context.Context, ctx *app.RequestContext) {
		filename := ctx.Param("filename")

		if !strings.HasSuffix(filename, ".openapi.yaml") && filename != "openapi.yaml" {
			handleError(ctx, "Invalid file name", http.StatusBadRequest)
			return
		}

		data, err := files.ReadFile(filename)
		if err != nil {
			handleError(ctx, "File not found: "+filename, http.StatusNotFound)
			return
		}

		ctx.Header("Content-Type", "application/x-yaml")
		ctx.Write(data)
	})
}


func setupProxyRoutes(h *server.Hertz, cp *ClientPool) {
	h.Any("/:serviceName/:methodName", func(c context.Context, ctx *app.RequestContext) {
		serviceName := ctx.Param("serviceName")
		methodName := ctx.Param("methodName")

		if serviceName == "" || methodName == "" {
			handleError(ctx, "ServiceName or MethodName not provided", http.StatusBadRequest)
			return
		}
		
		cli, err := cp.getClient(serviceName)
		if err != nil {
			handleError(ctx, err.Error(), http.StatusNotFound)
			return
		}

		bodyBytes := ctx.Request.Body()

		queryMap := formatQueryParams(ctx)
		
		for k, v := range queryMap {
			if strings.HasPrefix(k, "p_") {
				c = metainfo.WithPersistentValue(c, k, v)
			} else {
				c = metainfo.WithValue(c, k, v)
			}
		}

		c = metainfo.WithBackwardValues(c)

		jReq := string(bodyBytes)
		
		jRsp, err := cli.GenericCall(c, methodName, jReq)
		if err != nil {
			hlog.Errorf("GenericCall error: %v", err)
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		result := make(map[string]interface{})
		if err := json.Unmarshal([]byte(jRsp.(string)), &result); err != nil {
			hlog.Errorf("Failed to unmarshal response body: %v", err)
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to unmarshal response body",
			})
			return
		}

		m := metainfo.RecvAllBackwardValues(c)

		for key, value := range m {
			result[key] = value
		}

		respBody, err := json.Marshal(result)
		if err != nil {
			hlog.Errorf("Failed to marshal response body: %v", err)
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to marshal response body",
			})
			return
		}

		ctx.Data(http.StatusOK, "application/json", respBody)
	})
}

func formatQueryParams(ctx *app.RequestContext) map[string]string {
	var QueryParams = make(map[string]string)
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		QueryParams[string(key)] = string(value)
	})
	return QueryParams
}

func handleError(ctx *app.RequestContext, errMsg string, statusCode int) {
	hlog.Errorf("Error: %s", errMsg)
	ctx.JSON(statusCode, map[string]interface{}{
		"error": errMsg,
	})
}

func findPbFile(fileName string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	foundPath := ""
	err = filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking through files: %w", err)
		}
		if !info.IsDir() && info.Name() == fileName {
			foundPath = path
			return filepath.SkipDir
		}
		return nil
	})

	if err == nil && foundPath != "" {
		return foundPath, nil
	}

	parentDir := filepath.Dir(workingDir)
	for parentDir != "/" && parentDir != "." && parentDir != workingDir {
		filePath := filepath.Join(parentDir, fileName)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, fmt.Errorf("file found at: %s", filePath)
		}
		workingDir = parentDir
		parentDir = filepath.Dir(parentDir)
	}

	return "", errors.New("pb file not found: " + fileName)
}
`
