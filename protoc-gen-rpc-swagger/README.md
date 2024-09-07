# protoc-gen-rpc-swagger

English | [中文](README_CN.md)

This is a plugin for generating RPC Swagger documentation and providing Swagger-UI access and debugging for [cloudwego/cwgo](https://github.com/cloudwego/cwgo) & [kitex](https://github.com/cloudwego/kitex).

## Installation

```sh
# Install from the official repository
git clone https://github.com/hertz-contrib/swagger-generate
cd protoc-gen-rpc-swagger
go install
# Direct installation
go install github.com/hertz-contrib/swagger-generate/protoc-gen-rpc-swagger@latest
```

## Usage

### Generate Swagger Documentation

```sh
protoc --rpc-swagger_out=. --rpc-swagger_opt=output_mode=merged,HertzAddr=127.0.0.1:8080,KitexAddr=127.0.0.1:8888 -I . hello.proto
```

### Start Swagger-UI Server

```sh
go run swagger.go
```

### Access Swagger-UI (Kitex service must be running for debugging)

```sh
http://127.0.0.1:8080/swagger/index.html
```

## Instructions

### Generation Guide
1. The plugin will generate Swagger documentation and a HTTP (Hertz) service for accessing and debugging the Swagger docs.
2. When `output_mode` is set to `merged`, a single `openapi.yaml` is generated; when set to `source_relative`, a corresponding `{proto_filename}.openapi.yaml` is generated for each proto file.
3. All RPC methods are converted to HTTP `POST` methods. The request parameters map to the Request body, with the content type as `application/json`, and the response follows the same format.
4. You can use annotations to supplement the Swagger documentation, such as `openapi.operation`, `openapi.property`, `openapi.schema`, `api.base_domain`, `api.baseurl`.
5. To use annotations like `openapi.operation`, `openapi.property`, `openapi.schema`, and `openpai.document`, include the [annotations.proto](example/openapi/annotations.proto).

### Debugging Guide
1. Ensure that proto files, `*.openapi.yaml`, and `swagger.go` are in the same directory.
2. To access the Swagger documentation, start the `swagger.go` HTTP service. The service address can be specified using the `HertzAddr` parameter (default: 127.0.0.1:8080). Ensure that the `server` URL in the Swagger documentation matches the `HertzAddr` for debugging. Access the Swagger-UI at `/swagger/index.html` after starting the service.
3. Kitex service must be running for Swagger documentation debugging. The `KitexAddr` parameter specifies the Kitex service address (default: 127.0.0.1:8888). Ensure it matches the actual Kitex service address.
4. When `output_mode` is set to `source_relative`, the debugging process defaults to the Swagger documentation for the first proto file. You can also access other proto files’ Swagger documentation by using `{proto_filename}.openapi.yaml`.

### Metadata Transmission
1. Metadata transmission is supported. The plugin generates a `ttheader` query parameter for each method by default, used for passing metadata. The format should comply with JSON, like `{"p_k":"p_v","k":"v"}`.
2. Single-hop metadata transmission uses the format `"key":"value"`.
3. Persistent metadata transmission uses the format `"p_key":"value"` and requires the prefix `p_`.
4. Reverse metadata transmission is supported. If set, metadata will be included in the response and returned in the `"key":"value"` format.
5. For more details on using metadata, refer to [Metainfo](https://www.cloudwego.io/docs/kitex/tutorials/advanced-feature/metainfo/).

## Supported Annotations

| Annotation          | Component | Description                                                          |  
|---------------------|-----------|----------------------------------------------------------------------|
| `openapi.operation` | Method    | Supplements `operation` in `pathItem`                                |
| `openapi.property`  | Field     | Supplements `property` in `schema`                                   |
| `openapi.schema`    | Message   | Supplements `schema` in `requestBody` and `response`                 |
| `openapi.document`  | Document  | Supplements the Swagger documentation                                |
| `api.base_domain`   | Service   | Specifies the service `url` corresponding to the `server`            |
| `api.baseurl`       | Method    | Specifies the method’s `url` corresponding to `server` in `pathItem` |

## More Information

For more usage examples, please refer to the [examples](example/hello.proto).