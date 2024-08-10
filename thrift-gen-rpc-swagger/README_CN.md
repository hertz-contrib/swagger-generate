# thrift-gen-rpc-swagger

[English](README.md) | 中文

适用于 [cloudwego/cwgo](https://github.com/cloudwego/cwgo) & [kitex](https://github.com/cloudwego/kitex) 的 rpc swagger 文档生成及 swagger-ui 访问调试插件。

## 支持的注解

### Request 规范

1. 接口请求字段需要使用注解关联到 HTTP 的某类参数和参数名称, 没有注解的字段不做处理。
2. 根据 `method` 中的请求 `message` 生成 swagger 中 `operation` 的 `parameters` 和 `requestBody`。
3. 如果 HTTP 请求是采用 `GET`、`HEAD`、`DELETE` 方式的，那么 `request` 定义中出现的 `api.body` 注解无效，只有`api.query`, `api.path`, `api.cookie`, `api.header` 有效。

#### 注解说明

| 注解             | 说明                                                                                                                   |  
|----------------|----------------------------------------------------------------------------------------------------------------------|
| `api.query`    | `api.query` 对应 `parameter` 中 `in: query` 参数, 支持基本类型和`list`(`object`, `map`暂不支持）                                      |  
| `api.path`     | `api.path` 对应 `parameter` 中 `in: path` 参数, `required` 为 `true`, 支持基本类型                                               |
| `api.header`   | `api.header` 对应 `parameter` 中 `in: header` 参数, 支持基本类型和`list`                                                         |       
| `api.cookie`   | `api.cookie` 对应 `parameter` 中 `in: cookie` 参数, 支持基本类型                                                                |
| `api.body`     | `api.body` 对应 `requestBody` 中 `content` 为 `application/json`                                                         |
| `api.form`     | `api.form` 对应 `requestBody` 中 `content` 为 `multipart/form-data` 或 `application/x-www-form-urlencoded`, 预留, Kitex暂不支持 | 
| `api.raw_body` | `api.body` 对应 `requestBody` 中 `content` 为 `text/plain`                                                               |

### Response 规范

1. 接口响应字段需要使用注解关联到 HTTP 的某类参数和参数名称, 没有注解的字段不做处理。
2. 根据 `method` 中的响应 `message` 生成 swagger 中 `operation` 的 `responses`。

#### 注解说明

| 注解             | 说明                                                        |  
|----------------|-----------------------------------------------------------|
| `api.header`   | `api.header` 对应 `response` 中 `header`, 只支持基本类型和逗号分隔的list  |
| `api.body`     | `api.body` 对应 `response` 中 `content` 为 `application/json` |

### Method 规范

1. 每个 `method` 通过注解来关联 `pathItem`

#### 注解说明

| 注解            | 说明                                                         |  
|---------------|------------------------------------------------------------|
| `api.get`     | `api.get` 对应 `GET` 请求，只有 `parameter`                       |
| `api.put`     | `api.put` 对应 `PUT` 请求                                      |
| `api.post`    | `api.post` 对应 `POST` 请求                                    |
| `api.patch`   | `api.patch` 对应 `PATCH` 请求                                  |
| `api.delete`  | `api.delete` 对应 `DELETE` 请求，只有 `parameter`                 |
| `api.baseurl` | `api.baseurl` 对应 `pathItem` 的 `server` 的 `url`, 非 Kitex 注解 |

### Service 规范

#### 注解说明

| 注解                | 说明                                                |  
|-------------------|---------------------------------------------------|
| `api.base_domain` | `api.base_domain` 对应 `server` 的 `url`, 非 Kitex 注解 |

## openapi 注解

| 注解                  | 使用组件    | 说明                                         |  
|---------------------|---------|--------------------------------------------|
| `openapi.operation` | Method  | 用于补充 `pathItem` 的 `operation`              |
| `openapi.property`  | Field   | 用于补充 `schema` 的 `property`                 |
| `openapi.schema`    | Struct  | 用于补充 `requestBody` 和 `response` 的 `schema` |
| `openapi.document`  | Service | 用于补充 swagger 文档，任意service中添加该注解即可          |
| `openapi.parameter` | Field   | 用于补充 `parameter`                           |

更多的使用方法请参考 [示例](example/hello.thrift)

## 安装

```sh

# 官方仓库安装

git clone https://github.com/hertz-contrib/swagger-generate
cd thrift-gen-rpc-swagger
go install

# 直接安装
go install github.com/hertz-contrib/swagger-generate/thrift-gen-rpc-swagger@latest

# 验证安装
thrift-gen-rpc-swagger --version
```

## 使用

### 生成 swagger 文档

```sh

thriftgo -g go -p rpc-swagger:OutputDir=./output,HertzAddr=127.0.0.1:8080,KitexAddr=127.0.0.1:8888 hello.thrift

```
### 启动 swagger-ui 服务

```sh

go run ./output/swagger.go

```

## 补充说明

1. 插件会生成 swagger 文档，并且会生成一个 http (Hertz) 服务, 用于提供 swagger 文档的访问及调试。
2. swagger 文档的访问需启动 swagger.go, http 服务的地址可以通过参数 `HertzAddr` 参数指定, 默认为127.0.0.1:8080, 需要保持 swagger 文档中的 `server` 与 `HertzAddr` 一致才可以调试, 启动后访问访问/swagger/index.html。
3. swagger 文档的调试还需启动 Kitex 服务, `KitexAddr`用于指定 Kitex 服务的地址, 默认为127.0.0.1:8888, 需要保持与实际的 Kitex 服务地址一致。
4. 对 rpc 服务的调试基于 Kitex 的 http 泛化调用, 更多的信息请参考 [Kitex泛化调用](https://www.cloudwego.io/zh/docs/kitex/tutorials/advanced-feature/generic-call/thrift_idl_annotation_standards/)。

## 更多信息

查看 [示例](example/hello.thrift)




