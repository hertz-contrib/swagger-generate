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

package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/hertz-contrib/cors"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

//go:embed openapi.yaml
var openapiYAML []byte

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	h.Use(cors.Default())

	cli := initializeGenericClient()
	setupSwaggerRoutes(h)
	setupProxyRoutes(h, cli)

	hlog.Info("Swagger UI is available at: http://127.0.0.1:8080/swagger/index.html")

	h.Spin()
}

func findThriftFile(fileName string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	foundPath := ""
	err = filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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
	for parentDir != "/" && parentDir != "." {
		filePath := filepath.Join(parentDir, fileName)
		if _, err := os.Stat(filePath); err == nil {
			return filePath, nil
		}
		parentDir = filepath.Dir(parentDir)
	}

	return "", errors.New("thrift file not found: " + fileName)
}

func initializeGenericClient() genericclient.Client {
	thriftFile, err := findThriftFile("hello.thrift")
	if err != nil {
		hlog.Fatal("Failed to locate Thrift file:", err)
	}

	p, err := generic.NewThriftFileProvider(thriftFile)
	if err != nil {
		hlog.Fatal("Failed to create ThriftFileProvider:", err)
	}

	g, err := generic.HTTPThriftGeneric(p)
	if err != nil {
		hlog.Fatal("Failed to create HTTPThriftGeneric:", err)
	}

	cli, err := genericclient.NewClient("swagger", g, client.WithHostPorts("127.0.0.1:8888"))
	if err != nil {
		hlog.Fatal("Failed to create generic client:", err)
	}

	return cli
}

func setupSwaggerRoutes(h *server.Hertz) {
	h.GET("swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, swagger.URL("/openapi.yaml")))

	h.GET("/openapi.yaml", func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Content-Type", "application/x-yaml")
		ctx.Write(openapiYAML)
	})
}

func setupProxyRoutes(h *server.Hertz, cli genericclient.Client) {
	h.Any("/*ServiceMethod", func(c context.Context, ctx *app.RequestContext) {
		serviceMethod := ctx.Param("ServiceMethod")
		if serviceMethod == "" {
			handleError(ctx, "ServiceMethod not provided", http.StatusBadRequest)
			return
		}

		queryString := formatQueryParams(ctx)
		bodyBytes := ctx.Request.Body()
		contentType := string(ctx.Request.Header.ContentType())

		url := "http://127.0.0.1:8080/" + serviceMethod
		if len(queryString) > 0 {
			url += "?" + queryString
		}

		req, err := http.NewRequest(string(ctx.Request.Method()), url, bytes.NewBuffer(bodyBytes))
		if err != nil {
			handleError(ctx, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx.Request.Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})

		req.Header.Set("Content-Type", contentType)

		handleProxyRequest(ctx, cli, req)
	})
}

func formatQueryParams(ctx *app.RequestContext) string {
	var newQueryParams []string
	ctx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		newQueryParams = append(newQueryParams, string(key)+"="+string(value))
	})
	return strings.Join(newQueryParams, "&")
}

func handleProxyRequest(ctx *app.RequestContext, cli genericclient.Client, req *http.Request) {
	customReq, err := generic.FromHTTPRequest(req)
	if err != nil {
		handleError(ctx, "Failed to create generic request", http.StatusInternalServerError)
		return
	}

	resp, err := cli.GenericCall(context.Background(), "", customReq)
	if err != nil {
		handleError(ctx, "GenericCall error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if resp == nil {
		handleError(ctx, "Received nil response from the service", http.StatusInternalServerError)
		return
	}

	realResp, ok := resp.(*generic.HTTPResponse)
	if !ok {
		handleError(ctx, "Invalid response format", http.StatusInternalServerError)
		return
	}

	sendResponse(ctx, realResp)
}

func sendResponse(ctx *app.RequestContext, realResp *generic.HTTPResponse) {
	if realResp.StatusCode == 0 {
		realResp.StatusCode = http.StatusOK
	}

	for key, values := range realResp.Header {
		for _, value := range values {
			ctx.Response.Header.Set(key, value)
		}
	}

	respBody, err := json.Marshal(realResp.Body)
	if err != nil {
		handleError(ctx, "Failed to marshal response body", http.StatusInternalServerError)
		return
	}

	ctx.Data(int(realResp.StatusCode), string(realResp.ContentType), respBody)
}

func handleError(ctx *app.RequestContext, errMsg string, statusCode int) {
	hlog.Errorf("Error: %s", errMsg)
	ctx.JSON(statusCode, map[string]interface{}{
		"error": errMsg,
	})
}
