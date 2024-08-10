# thrift-gen-rpc-swagger

English | [中文](README_CN.md)

This is a plugin for generating RPC Swagger documentation and providing Swagger-UI access and debugging for [cloudwego/cwgo](https://github.com/cloudwego/cwgo) & [kitex](https://github.com/cloudwego/kitex).

## Supported Annotations

### Request Specifications

1. Interface request fields need to be associated with certain HTTP parameters and parameter names using annotations. Fields without annotations will not be processed.
2. The `method` request `message` is used to generate the `parameters` and `requestBody` for `operation` in Swagger.
3. If the HTTP request uses `GET`, `HEAD`, or `DELETE` methods, the `api.body` annotation in the `request` definition will be invalid, and only `api.query`, `api.path`, `api.cookie`, and `api.header` will be valid.

#### Annotation Descriptions

| Annotation     | Description                                                                                                                                                                |  
|----------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `api.query`    | `api.query` corresponds to the `in: query` parameter in `parameter`, supports basic types and `list` (but not `object` or `map`)                                           |  
| `api.path`     | `api.path` corresponds to the `in: path` parameter in `parameter`, `required` is `true`, supports basic types                                                              |
| `api.header`   | `api.header` corresponds to the `in: header` parameter in `parameter`, supports basic types and `list`                                                                     |       
| `api.cookie`   | `api.cookie` corresponds to the `in: cookie` parameter in `parameter`, supports basic types                                                                                |
| `api.body`     | `api.body` corresponds to the `content` in `requestBody` as `application/json`                                                                                             |
| `api.form`     | `api.form` corresponds to the `content` in `requestBody` as `multipart/form-data` or `application/x-www-form-urlencoded`, reserved for future use, Kitex not yet supported | 
| `api.raw_body` | `api.raw_body` corresponds to the `content` in `requestBody` as `text/plain`                                                                                               |

### Response Specifications

1. Interface response fields need to be associated with certain HTTP parameters and parameter names using annotations. Fields without annotations will not be processed.
2. The `method` response `message` is used to generate the `responses` for `operation` in Swagger.

#### Annotation Descriptions

| Annotation   | Description                                                                                             |  
|--------------|---------------------------------------------------------------------------------------------------------|
| `api.header` | `api.header` corresponds to `header` in `response`, supports only basic types and comma-separated lists |
| `api.body`   | `api.body` corresponds to the `content` in `response` as `application/json`                             |

### Method Specifications

1. Each `method` is associated with a `pathItem` using annotations.

#### Annotation Descriptions

| Annotation    | Description                                                                              |  
|---------------|------------------------------------------------------------------------------------------|
| `api.get`     | `api.get` corresponds to a `GET` request, only `parameter` is used                       |
| `api.put`     | `api.put` corresponds to a `PUT` request                                                 |
| `api.post`    | `api.post` corresponds to a `POST` request                                               |
| `api.patch`   | `api.patch` corresponds to a `PATCH` request                                             |
| `api.delete`  | `api.delete` corresponds to a `DELETE` request, only `parameter` is used                 |
| `api.baseurl` | `api.baseurl` corresponds to the `url` of `server` in `pathItem`, not a Kitex annotation |

### Service Specifications

#### Annotation Descriptions

| Annotation        | Description                                                                    |  
|-------------------|--------------------------------------------------------------------------------|
| `api.base_domain` | `api.base_domain` corresponds to the `url` of `server`, not a Kitex annotation |

## openapi Annotations

| Annotation          | Used For | Description                                                                      |  
|---------------------|----------|----------------------------------------------------------------------------------|
| `openapi.operation` | Method   | Used to supplement the `operation` in `pathItem`                                 |
| `openapi.property`  | Field    | Used to supplement the `property` in `schema`                                    |
| `openapi.schema`    | Struct   | Used to supplement the `schema` in `requestBody` and `response`                  |
| `openapi.document`  | Service  | Used to supplement the Swagger documentation, add this annotation to any service |
| `openapi.parameter` | Field    | Used to supplement `parameter`                                                   |

For more usage examples, please refer to the [example](example/hello.thrift).

## Installation

```sh

# Install from official repository

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

## Additional Information

1. The plugin generates Swagger documentation and an HTTP (Hertz) service for accessing and debugging the Swagger documentation.
2. To access the Swagger documentation, you need to start the `swagger.go` HTTP service. The address of the service can be specified using the `HertzAddr` parameter, defaulting to 127.0.0.1:8080. The `server` in the Swagger documentation must match the `HertzAddr` for debugging to work, visit /swagger/index.html after startup.
3. To debug the Swagger documentation, you also need to start the Kitex service. The `KitexAddr` parameter specifies the address of the Kitex service, defaulting to 127.0.0.1:8888. This address must match the actual Kitex service address.
4. Debugging RPC services is based on Kitex's HTTP generic calls. For more information, please refer to [Kitex Generic Calls](https://www.cloudwego.io/en/docs/kitex/tutorials/advanced-feature/generic-call/thrift_idl_annotation_standards/).

## More Information

See the [example](example/hello.thrift).