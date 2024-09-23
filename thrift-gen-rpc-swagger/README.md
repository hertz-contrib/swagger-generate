# thrift-gen-rpc-swagger

English | [中文](README_CN.md)

This is a plugin for generating RPC Swagger documentation and providing Swagger-UI access and debugging for [cloudwego/cwgo](https://github.com/cloudwego/cwgo) & [kitex](https://github.com/cloudwego/kitex).

## Installation

```sh
# Install from the official repository

git clone https://github.com/hertz-contrib/swagger-generate
cd swagger-generate
go install

# Direct installation
go install github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger@latest

# Verify installation
thrift-gen-rpc-swagger --version
```

## Usage

### Generating Swagger Documentation

```sh
thriftgo -g go -p rpc-swagger hello.thrift
```

### Add the option during Kitex Server initialization

```sh
svr := api.NewServer(new(HelloImpl), server.WithTransHandlerFactory(&swagger.MixTransHandlerFactory{}))
```

### Access Swagger UI (Kitex service needs to be running for debugging)

```sh
http://127.0.0.1:8888/swagger/index.html
```

## Usage Instructions

### Debugging Notes
1. The plugin generates Swagger documentation and also sets up an HTTP (Hertz) service to provide access to the Swagger documentation and debugging.
2. The HTTP service defaults to the same port as the RPC service, implemented via protocol sniffing.
3. Accessing the Swagger documentation and debugging the RPC service requires adding `"server.WithTransHandlerFactory(&swagger.MixTransHandlerFactory{})"` to the Kitex Server initialization.

### Generation Notes
1. All RPC methods are converted into HTTP POST methods, with request parameters corresponding to the Request body in `application/json` format, and the same for the return value.
2. Swagger documentation can be supplemented with annotations such as `openapi.operation`, `openapi.property`, `openapi.schema`, `api.base_domain`, and `api.baseurl`.
3. To use annotations like `openapi.operation`, `openapi.property`, `openapi.schema`, and `openapi.document`, you need to import `openapi.thrift`.
4. Custom HTTP services are supported, and custom parts will not be overwritten during updates.
5. The RPC method request and response only support `struct` and empty types.

### Metadata Transmission
1. Metadata transmission is supported. By default, the plugin generates a `ttheader` query parameter for each method to transmit metadata, which should be in JSON format, e.g., `{"p_k":"p_v","k":"v"}`.
2. Single-hop metadata transmission format is `"key":"value"`.
3. Continuous metadata transmission format is `"p_key":"value"`, with a prefix `p_`.
4. Reverse metadata transmission is supported; if enabled, metadata can be viewed in the return value, appended to the response in `"key":"value"` format.
5. For more information on using metadata, refer to [Metainfo](https://www.cloudwego.io/zh/docs/kitex/tutorials/advanced-feature/metainfo/).

## Supported Annotations

| Annotation          | Component | Description                                                                              |
|---------------------|-----------|------------------------------------------------------------------------------------------|
| `openapi.operation` | Method    | Supplements the `operation` of `pathItem`                                                |
| `openapi.property`  | Field     | Supplements the `property` of `schema`                                                   |
| `openapi.schema`    | Struct    | Supplements the `schema` for `requestBody` and `response`                                |
| `openapi.document`  | Service   | Supplements Swagger documentation; add this annotation to any service                    |
| `api.base_domain`   | Service   | Corresponds to `server`'s `url`, specifies the URL for the service                       |
| `api.baseurl`       | Method    | Corresponds to `pathItem`'s `server`'s `url`, specifies the URL for an individual method |

## More Information

For more usage instructions, refer to [Example](example/hello.thrift).