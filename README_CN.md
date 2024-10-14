# Swagger Generate

[English](README.md) | 中文

**Swagger Generate** 是一组插件工具，专为 HTTP 和 RPC 服务设计，支持自动生成 Swagger 文档，并集成 Swagger-UI 进行调试。此外，它还提供将 Swagger 文档转换为 Protobuf 或 Thrift IDL 文件的功能，简化了 API 开发与维护的流程。

该项目适用于 [CloudWeGo](https://www.cloudwego.io) 生态下的 [Cwgo](https://github.com/cloudwego/cwgo)、 [Hertz](https://github.com/cloudwego/hertz) 和 [Kitex](https://github.com/cloudwego/kitex) 框架。它提供了一套便捷的工具来帮助开发者自动生成 Swagger 文档，从而简化 API 文档编写及调试过程。

## 包含的插件

- **[protoc-gen-http-swagger](https://github.com/hertz-contrib/swagger-generate/tree/main/thrift-gen-rpc-swagger)**：为基于 Protobuf 的 HTTP 服务生成 Swagger 文档和 Swagger UI 进行调试。
- **[thrift-gen-http-swagger](https://github.com/hertz-contrib/swagger-generate/tree/main/thrift-gen-http-swagger)**：为基于 Thrift 的 HTTP 服务生成 Swagger 文档和 Swagger UI 进行调试。
- **[protoc-gen-rpc-swagger](https://github.com/hertz-contrib/swagger-generate/tree/main/protoc-gen-rpc-swagger)**：为基于 Protobuf 的 RPC 服务生成 Swagger 文档和 Swagger UI 进行调试。
- **[thrift-gen-rpc-swagger](https://github.com/hertz-contrib/swagger-generate/tree/main/thrift-gen-rpc-swagger)**：为基于 Thrift 的 RPC 服务生成 Swagger 文档和 Swagger UI 进行调试。
- **[swagger2idl](https://github.com/hertz-contrib/swagger-generate/tree/main/swagger2idl)**：将 Swagger 文档转换为 Protobuf 或 Thrift IDL 文件。

## 项目优势

- **自动化生成**：支持通过 Protobuf 和 Thrift 文件生成完整的 Swagger 文档，简化了 API 文档的维护。
- **集成调试**：生成的 Swagger UI 能直接用于调试服务，支持 HTTP 和 RPC 两种模式。
- **Hertz 和 Kitex 集成**：为 [Hertz](https://github.com/cloudwego/hertz) 和 [Kitex](https://github.com/cloudwego/kitex) 提供了无缝的文档生成和调试支持。
- **灵活的注解支持**：允许通过注解扩展生成的 Swagger 文档内容，支持 `openapi.operation`、`openapi.schema` 等 OpenAPI 注解。
- **IDL 转换**：支持将 Swagger 文档转换为 Protobuf 或 Thrift IDL 文件，方便开发者在不同框架间切换。

## 安装

可以通过以下方式安装各个插件：

```sh
# 官方仓库安装
git clone https://github.com/hertz-contrib/swagger-generate
cd <plugin-directory>
go install

# 直接安装
go install github.com/hertz-contrib/swagger-generate/<plugin-name>@latest
```

## 使用示例

### 生成 HTTP Swagger 文档

对于基于 Protobuf 的 HTTP 服务：

```sh
protoc --http-swagger_out=swagger -I idl hello.proto
```

对于基于 Thrift 的 HTTP 服务：

```sh
thriftgo -g go -p http-swagger hello.thrift
```

### 生成 RPC Swagger 文档

对于基于 Protobuf 的 RPC 服务：

```sh
protoc --rpc-swagger_out=swagger -I idl idl/hello.proto
```

对于基于 Thrift 的 RPC 服务：

```sh
thriftgo -g go -p rpc-swagger hello.thrift
```

### 在 Hertz 或 Kitex 服务中集成 Swagger-UI

在 Hertz 服务中：

```go
func main() {
    h := server.Default()
    swagger.BindSwagger(h) //增加改行
    register(h)
    h.Spin()
}
```
或者

```go
func register(r *server.Hertz) {
    swagger.BindSwagger(r) //增加改行
    
    router.GeneratedRegister(r)
    
    customizedRegister(r)
}
```
在 Kitex 服务中：

```go
func main() {
	svr := example.NewServer(new(HelloService1Impl), server.WithTransHandlerFactory(&swagger.MixTransHandlerFactory{})) //改动改行

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
```
请参考 [kitex_swagger_gen](https://github.com/cloudwego/kitex-examples/tree/main/bizdemo/kitex_swagger_gen) 和 [hertz_swagger_gen](https://github.com/cloudwego/hertz-examples/tree/main/bizdemo/hertz_swagger_gen) 获取更多使用场景示例。

### 将 Swagger 文档转换为 IDL 文件

```sh
swagger2idl -o my_output.proto -oa -a openapi.yaml
```

## 更多信息

请参考各个插件的 README 文档获取更多使用细节。