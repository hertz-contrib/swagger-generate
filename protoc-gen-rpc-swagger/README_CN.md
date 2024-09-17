# protoc-gen-rpc-swagger

[English](README.md) | 中文

适用于 [cloudwego/cwgo](https://github.com/cloudwego/cwgo) & [kitex](https://github.com/cloudwego/kitex) 的 rpc swagger 文档生成及 swagger-ui 访问调试插件。

## 安装

```sh
# 官方仓库安装
git clone https://github.com/hertz-contrib/swagger-generate
cd protoc-gen-rpc-swagger
go install
# 直接安装
go install github.com/hertz-contrib/swagger-generate/protoc-gen-rpc-swagger@latest
```

## 使用

### 生成 swagger 文档

```sh
protoc --rpc-swagger_out=swagger -I idl idl/hello.proto
```
### 在 Kitex Server 初始化中添加 option

```sh

svr := api.NewServer(new(HelloImpl), server.WithTransHandlerFactory(&swagger.MixTransHandlerFactory{}))

```

### 访问 swagger-ui (调试需启动 Kitex 服务)

```sh

http://127.0.0.1:8888/swagger/index.html
```

## 使用说明

### 生成说明
1. 插件会生成 swagger 文档，同时生成一个 http (Hertz) 服务, 用于提供 swagger 文档的访问及调试。 
2. 所有的 rpc 方法会转换成 http 的 `post` 方法，请求参数对应 Request body, content 类型为 `application/json` 格式，返回值同上。 
3. 可通过注解来补充 swagger 文档的信息，如 `openapi.operation`, `openapi.property`, `openapi.schema`, `api.base_domain`, `api.baseurl`。 
4. 如需使用`openapi.operation`, `openapi.property`, `openapi.schema`, `openpai.document` 注解，需引用 [annotations.proto](example/idl/openapi/annotations.proto)。

### 调试说明
1. 需保证 proto 文件与 `openapi.yaml`、 `swagger.go` 在同一目录下。
2. http 服务默认和 rpc 服务在一个端口, 通过嗅探协议实现。
3. swagger 文档的访问及 rpc 服务的调试需在 Kitex Server 初始化中加入 "server.WithTransHandlerFactory(&swagger.MixTransHandlerFactory{})"。

### 元信息传递
1. 支持元信息传递, 插件默认为每个方法生成一个`ttheader`的查询参数, 用于传递元信息, 格式需满足 json 格式, 如`{"p_k":"p_v","k":"v"}`。
2. 单跳透传元信息, 格式为 `"key":"value"`。
3. 持续透传元信息, 格式为 `"p_key":"value"`, 需添加前缀`p_`。
4. 支持反向透传元信息, 若设置则可在返回值中查看到元信息, 返回通过`"key":"value"`的格式附加在响应中。
5. 更多使用元信息可参考 [Metainfo](https://www.cloudwego.io/zh/docs/kitex/tutorials/advanced-feature/metainfo/)。

## 支持的注解

| 注解                  | 使用组件     | 说明                                                    |  
|---------------------|----------|-------------------------------------------------------|
| `openapi.operation` | Method   | 用于补充 `pathItem` 的 `operation`                         |
| `openapi.property`  | Field    | 用于补充 `schema` 的 `property`                            |
| `openapi.schema`    | Message  | 用于补充 `requestBody` 和 `response` 的 `schema`            |
| `openapi.document`  | Document | 用于补充 swagger 文档                                       |
| `api.base_domain`   | Service  | 对应 `server` 的 `url`, 用于指定 service 服务的 url             |
| `api.baseurl`       | Method   | 对应 `pathItem` 的 `server` 的 `url`, 用于指定单个 method 的 url |

## 更多信息

更多的使用方法请参考 [示例](example/idl/hello.proto)




