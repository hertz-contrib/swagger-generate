# swagger2idl

[English](README.md) | 中文

swagger2idl 是一个用于将 Swagger 文档转换为 Thrift 或 Proto 文件的工具。
适配了[swagger-generate](https://github.com/hertz-contrib/swagger-generate)、[cloudwego/cwgo](https://github.com/cloudwego/cwgo)、[hertz](https://github.com/cloudwego/hertz)及[kitex](https://github.com/cloudwego/kitex)中的相关注解。

## 安装

```sh
# 官方仓库安装

git clone https://github.com/hertz-contrib/swagger-generate
cd swagger2idl
go install

# 直接安装
go install github.com/hertz-contrib/swagger-generate/swagger2idl@latest
```

## 使用
### 参数说明

| 参数名称        | 缩写    | 默认值                        | 说明                                                                                                    |
|-------------|-------|----------------------------|-------------------------------------------------------------------------------------------------------|
| `--type`    | `-t`  | 自动根据输出文件扩展名推断              | 指定输出类型，可选值为 `'proto'` 或 `'thrift'`。如果未提供，则从输出文件扩展名推断。                                                 |
| `--output`  | `-o`  | `文件名.proto` 或 `文件名.thrift` | 指定输出文件的路径。如果未提供，默认为 `output.proto` 或 `output.thrift`，具体取决于输出类型。                                       |
| `--openapi` | `-oa` | `false`                    | 会生成相应的openapi注解，并添加引用，相关引用文件可以在[idl](https://github.com/hertz-contrib/swagger-generate/idl)中找到。       |
| `--api`     | `-a`  | `false`                    | 会生成相应的适配Cwgo/Hertz的注解，并添加引用，相关引用文件可以在[idl](https://github.com/hertz-contrib/swagger-generate/idl)中找到。 |
| `--naming`  | `-n`  | `true`                     | 在输出的 IDL 文件中使用命名约定。                                                                                   |

### 使用示例

1. 指定输出为 Protobuf 格式，并输出到指定路径：
```bash
   swagger2idl --output my_output.proto --openapi --api --naming=false openapi.yaml
```
or
```bash
   swagger2idl -o my_output.proto -oa -a -n=false openapi.yaml
```

### 扩展
支持向openapi.yaml中的参数添加扩展，如`x-options`，后面会增加更多扩展。

如果是proto文件
```yaml
x-options:
  go_package: myawesomepackage
```
会生成
```protobuf
option go_package = "myawesomepackage";
```
如果是thrift文件
```yaml
x-options:
  go: myawesomepackage
```
会生成
```thrift
namespace go myawesomepackage
```
### 命名约定

| **类别**                           | **Thrift/Proto 命名规范**                                                         |
|----------------------------------|-------------------------------------------------------------------------------|
| **Struct/Message**               | - 使用 **PascalCase** 命名。<br> - 例：`UserInfo`                                    |
| **Field**                        | - 使用 **snake_case** 命名。<br> - 例：`user_id`, 如果你的字段名包含一个数字，数字应该出现在字母后面，而不是下划线后面 |
| **Enum**, **Service**, **Union** | - 使用 **PascalCase**。<br> - 例：`UserType`                                       |
| **Enum 值**                       | - 使用 **UPPER_SNAKE_CASE** 命名。<br> - 例：`ADMIN_USER`                            |
| **RPC 方法**                       | - 使用 **PascalCase** 命名。<br> - 例：`GetUserInfo`                                 |
| **Package/Namespace**            | - 使用 **snake_case**，通常基于项目结构命名。<br> 例：`com.project.service`                   |

#### 详细说明：
- **PascalCase**: 首字母大写，每个单词的首字母都大写，例如 `UserInfo`。
- **snake_case**: 全部小写，单词之间使用下划线分隔，例如 `user_info`。
- **UPPER_SNAKE_CASE**: 全部字母大写，单词之间用下划线分隔，例如 `ADMIN_USER`。

## 更多信息

更多的使用方法请参考 [示例](example)