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
protoc --rpc-swagger_out=. --rpc-swagger_opt=output_mode=merged,HertzAddr=127.0.0.1:8080,KitexAddr=127.0.0.1:8888 -I . hello.proto
```
### 启动 swagger-ui 服务

```sh
go run swagger.go
```

### 访问 swagger-ui (调试需启动Kitex服务)

```sh
http://127.0.0.1:8080/swagger/index.html
```

## 使用说明

### 生成说明
1. 插件会生成 swagger 文档，同时生成一个 http (Hertz) 服务, 用于提供 swagger 文档的访问及调试。
2. 当 `output_mode` 为 `merged` 时，会聚合生成一个 `openapi.yaml`；当 `output_mode` 为 `source_relative` 时，会为每个 proto 文件生成一个对应的 `{proto文件名称}.openapi.yaml`。
3. 所有的 rpc 方法会转换成 http 的 `post` 方法，请求参数对应 Request body, content 类型为 `application/json` 格式，返回值同上。
4. 可通过注解来补充 swagger 文档的信息，如 `openapi.operation`, `openapi.property`, `openapi.schema`, `api.base_domain`, `api.baseurl`。
5. 如需使用`openapi.operation`, `openapi.property`, `openapi.schema`, `openpai.document` 注解，需引用 [annotations.proto](example/openapi/annotations.proto)。

### 调试说明
1. 需保证 proto 文件与 `*.openapi.yaml`、 `swagger.go` 在同一目录下。
2. swagger 文档的访问需启动 swagger.go, http 服务的地址可以通过参数 `HertzAddr` 参数指定, 默认为 127.0.0.1:8080, 需要保持 swagger 文档中 `server` 的 url 与 `HertzAddr` 一致才可以调试, 启动后访问访问/swagger/index.html。
3. swagger 文档的调试还需启动 Kitex 服务, `KitexAddr`用于指定 Kitex 服务的地址, 默认为 127.0.0.1:8888, 需要保持与实际的 Kitex 服务地址一致。
4. 当 `output_mode` 为 `source_relative` 时, 启动调试默认会定位到第一个 proto 文件的 swagger 文档, 您也可以通过{proto文件名称}.openapi.yaml来访问其他 proto 文件的 swagger 文档。

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

更多的使用方法请参考 [示例](example/hello.proto)




