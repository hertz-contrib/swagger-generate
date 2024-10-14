# Swagger Generate

English | [中文](README_CN.md)

**Swagger Generate** is a set of plugin tools designed for HTTP and RPC services, supporting the automatic generation of Swagger documentation and integration with Swagger-UI for debugging. Additionally, it provides the ability to convert Swagger documents into Protobuf or Thrift IDL files, simplifying the API development and maintenance process.

This project is compatible with the [CloudWeGo](https://www.cloudwego.io) ecosystem frameworks such as [Cwgo](https://github.com/cloudwego/cwgo), [Hertz](https://github.com/cloudwego/hertz), and [Kitex](https://github.com/cloudwego/kitex). It offers a convenient toolset for developers to automatically generate Swagger documentation, simplifying the API documentation and debugging process.

## Included Plugins

- **[protoc-gen-http-swagger](https://github.com/hertz-contrib/swagger-generate/tree/main/thrift-gen-rpc-swagger)**: Generates Swagger documentation and provides Swagger UI debugging for HTTP services based on Protobuf.
- **[thrift-gen-http-swagger](https://github.com/hertz-contrib/swagger-generate/tree/main/thrift-gen-http-swagger)**: Generates Swagger documentation and provides Swagger UI debugging for HTTP services based on Thrift.
- **[protoc-gen-rpc-swagger](https://github.com/hertz-contrib/swagger-generate/tree/main/protoc-gen-rpc-swagger)**: Generates Swagger documentation and provides Swagger UI debugging for RPC services based on Protobuf.
- **[thrift-gen-rpc-swagger](https://github.com/hertz-contrib/swagger-generate/tree/main/thrift-gen-rpc-swagger)**: Generates Swagger documentation and provides Swagger UI debugging for RPC services based on Thrift.
- **[swagger2idl](https://github.com/hertz-contrib/swagger-generate/tree/main/swagger2idl)**: Converts Swagger documents into Protobuf or Thrift IDL files.

## Key Advantages

- **Automated Generation**: Supports generating complete Swagger documentation from Protobuf and Thrift files, simplifying API documentation maintenance.
- **Integrated Debugging**: The generated Swagger UI can be used directly for service debugging, supporting both HTTP and RPC modes.
- **Hertz and Kitex Integration**: Provides seamless documentation generation and debugging support for [Hertz](https://github.com/cloudwego/hertz) and [Kitex](https://github.com/cloudwego/kitex).
- **Flexible Annotation Support**: Allows extending the generated Swagger documentation through annotations, supporting OpenAPI annotations such as `openapi.operation`, `openapi.schema`, etc.
- **IDL Conversion**: Supports converting Swagger documents into Protobuf or Thrift IDL files, making it easier for developers to switch between different frameworks.

## Installation

You can install the plugins using the following methods:

```sh
# Install from the official repository
git clone https://github.com/hertz-contrib/swagger-generate
cd <plugin-directory>
go install

# Direct installation
go install github.com/hertz-contrib/swagger-generate/<plugin-name>@latest
```

## Usage Examples

### Generating HTTP Swagger Documentation

For HTTP services based on Protobuf:

```sh
protoc --http-swagger_out=swagger -I idl hello.proto
```

For HTTP services based on Thrift:

```sh
thriftgo -g go -p http-swagger hello.thrift
```

### Generating RPC Swagger Documentation

For RPC services based on Protobuf:

```sh
protoc --rpc-swagger_out=swagger -I idl idl/hello.proto
```

For RPC services based on Thrift:

```sh
thriftgo -g go -p rpc-swagger hello.thrift
```

### Integrating Swagger-UI in Hertz or Kitex Services

In a Hertz service:

```go
func main() {
    h := server.Default()
    swagger.BindSwagger(h) // Add this line
    register(h)
    h.Spin()
}
```
Or:

```go
func register(r *server.Hertz) {
    swagger.BindSwagger(r) // Add this line
    
    router.GeneratedRegister(r)
    
    customizedRegister(r)
}
```

In a Kitex service:

```go
func main() {
	svr := example.NewServer(new(HelloService1Impl), server.WithTransHandlerFactory(&swagger.MixTransHandlerFactory{})) // Modify this line

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
```
For more examples, please refer to [kitex_swagger_gen](https://github.com/cloudwego/kitex-examples/tree/main/bizdemo/kitex_swagger_gen) and [hertz_swagger_gen](https://github.com/cloudwego/hertz-examples/tree/main/bizdemo/hertz_swagger_gen).

### Converting Swagger Documents to IDL Files

```sh
swagger2idl -o my_output.proto -oa -a openapi.yaml
```

## More Information

Refer to the README of each plugin for more detailed usage information.