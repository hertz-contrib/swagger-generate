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
protoc --rpc-swagger_out=swagger -I idl idl/hello.proto
```
Here’s the translation of the given section into English:

### Add Option During Kitex Server Initialization

```sh
svr := api.NewServer(new(HelloImpl), server.WithTransHandlerFactory(&swagger.MixTransHandlerFactory{}))
```

### Access Swagger-UI (Kitex service must be running for debugging)

```sh
http://127.0.0.1:8888/swagger/index.html
```

## Instructions

### Generation Instructions
1. The plugin will generate Swagger documentation and simultaneously generate an HTTP (Hertz) service to provide access to and debugging of the Swagger documentation.
2. All RPC methods will be converted into HTTP `POST` methods. The request parameters correspond to the Request body, and the content type is in `application/json` format. The response follows the same format.
3. Annotations can be used to supplement the Swagger documentation with information, such as `openapi.operation`, `openapi.property`, `openapi.schema`, `api.base_domain`, `api.baseurl`.
4. To use annotations like `openapi.operation`, `openapi.property`, `openapi.schema`, and `openapi.document`, you need to reference [annotations.proto](example/idl/openapi/annotations.proto).

### Debugging Instructions
1. Ensure that the proto files, `openapi.yaml`, and `swagger.go` are in the same directory.
2. By default, the HTTP service runs on the same port as the RPC service, with protocol sniffing implemented.
3. To access the Swagger documentation and debug the RPC service, you must add "server.WithTransHandlerFactory(&swagger.MixTransHandlerFactory{})" during Kitex Server initialization.

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

For more usage examples, please refer to the [examples](example/idl/hello.proto).