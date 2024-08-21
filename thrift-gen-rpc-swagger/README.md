# thrift-gen-rpc-swagger

English | [中文](README_CN.md)

This is a plugin for generating RPC Swagger documentation and providing Swagger-UI access and debugging for [cloudwego/cwgo](https://github.com/cloudwego/cwgo) & [kitex](https://github.com/cloudwego/kitex).

## Installation

```sh
# Install from the official repository

git clone https://github.com/hertz-contrib/swagger-generate
cd thrift-gen-rpc-swagger
go install

# Direct installation
go install github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger@latest

# Verify installation
thrift-gen-rpc-swagger --version
```

## Usage

### Generate Swagger Documentation

```sh
thriftgo -g go -p rpc-swagger:OutputDir=./output,HertzAddr=127.0.0.1:8080,KitexAddr=127.0.0.1:8888 hello.thrift
```

### Start the Swagger-UI Service

```sh
go run ./output/swagger.go
```

### Access Swagger-UI (Kitex service needs to be running for debugging)

```sh
http://127.0.0.1:8080/swagger/index.html
```

## Usage Instructions

### Debugging Instructions
1. The plugin generates Swagger documentation and also creates an HTTP (Hertz) service to provide access and debugging of the Swagger documentation.
2. To access the Swagger documentation, start `swagger.go`. The HTTP service address can be specified via the `HertzAddr` parameter, which defaults to 127.0.0.1:8080. Ensure that the `server` URL in the Swagger documentation matches `HertzAddr` for debugging purposes. After starting, access `/swagger/index.html`.
3. To debug the Swagger documentation, the Kitex service must also be running. The `KitexAddr` parameter specifies the Kitex service address, which defaults to 127.0.0.1:8888. Ensure that this matches the actual Kitex service address.

### Generation Instructions
1. All RPC methods will be converted to HTTP POST methods. The request parameters correspond to the request body, with the content type being `application/json`. The response is handled similarly.
2. Annotations can be used to supplement information in the Swagger documentation, such as `openapi.operation`, `openapi.property`, `openapi.schema`, `api.base_domain`, and `api.baseurl`.
3. To use the `openapi.operation`, `openapi.property`, `openapi.schema`, and `openapi.document` annotations, you need to import `openapi.thrift`.

### Metadata Passing
1. Metadata passing is supported. The plugin generates a `theader` query parameter for each method by default, used for metadata transmission. The format should be in JSON, e.g., `{"p_k":"p_v","k":"v"}`.
2. For single-hop metadata passthrough, the format is `"key":"value"`.
3. For continuous metadata passthrough, the format is `"p_key":"value"`, with the `p_` prefix added.
4. Reverse metadata passthrough is supported. If set, metadata can be viewed in the response, attached in the format `"key":"value"`.
5. For more information on using metadata, please refer to [Metainfo](https://www.cloudwego.io/en/docs/kitex/tutorials/advanced-feature/metainfo/).

## Supported Annotations

| Annotation          | Component | Description                                                                                            |  
|---------------------|-----------|--------------------------------------------------------------------------------------------------------|
| `openapi.operation` | Method    | Used to supplement the `operation` of a `pathItem`                                                     |
| `openapi.property`  | Field     | Used to supplement the `property` of a `schema`                                                        |
| `openapi.schema`    | Struct    | Used to supplement the `schema` of a `requestBody` and `response`                                      |
| `openapi.document`  | Service   | Used to supplement the Swagger document; can be added in any service                                   |
| `api.base_domain`   | Service   | Corresponds to the `url` of the `server`, used to specify the service's URL                            |
| `api.baseurl`       | Method    | Corresponds to the `url` of the `server` for a `pathItem`, used to specify the URL for a single method |

## Additional Information

For more usage details, please refer to the [example](example/hello.thrift).