# swagger2idl

ENGLISH | [中文](README_CN.md)

`swagger2idl` is a tool designed to convert Swagger documentation into Thrift or Proto files. It supports relevant annotations from [swagger-generate](https://github.com/hertz-contrib/swagger-generate), [cloudwego/cwgo](https://github.com/cloudwego/cwgo), [hertz](https://github.com/cloudwego/hertz), and [kitex](https://github.com/cloudwego/kitex).

## Installation

```sh
# Install from the official repository

git clone https://github.com/hertz-contrib/swagger-generate
cd swagger2idl
go install

# Direct installation
go install github.com/hertz-contrib/swagger-generate/swagger2idl@latest
```

## Usage

### Parameter Description

| Parameter       | Abbreviation | Default Value                  | Description                                                                                                        |
|-----------------|--------------|--------------------------------|--------------------------------------------------------------------------------------------------------------------|
| `--type`        | `-t`         | Inferred from the output file extension | Specify the output type, either `'proto'` or `'thrift'`. If not provided, it is inferred from the output file extension. |
| `--output`      | `-o`         | `filename.proto` or `filename.thrift` | Specify the output file path. If not provided, it defaults to `output.proto` or `output.thrift`, depending on the output type. |
| `--openapi`     | `-oa`        | `false`                        | Includes OpenAPI-specific annotations and adds references. The related reference files can be found in [idl](https://github.com/hertz-contrib/swagger-generate/idl). |
| `--api`         | `-a`         | `false`                        | Adds annotations for compatibility with Cwgo/Hertz and adds references. The related reference files are in [idl](https://github.com/hertz-contrib/swagger-generate/idl). |
| `--naming`      | `-n`         | `true`                         | Use naming conventions in the output IDL file.                                                                     |

### Usage Examples

1. Convert to Protobuf format and specify the output path:
```bash
   swagger2idl --output my_output.proto --openapi --api --naming=false openapi.yaml
```
or
```bash
   swagger2idl -o my_output.proto -oa -a -n=false openapi.yaml
```

### Extensions
You can add extensions like `x-options` to parameters in the `openapi.yaml` file. More extensions will be supported in the future.

For Proto files:
```yaml
x-options:
  go_package: myawesomepackage
```
Generates:
```protobuf
option go_package = "myawesomepackage";
```

For Thrift files:
```yaml
x-options:
  go: myawesomepackage
```
Generates:
```thrift
namespace go myawesomepackage
```

### Naming Conventions

| **Category**                       | **Thrift/Proto Naming Rules**                                                  |
|------------------------------------|-------------------------------------------------------------------------------|
| **Struct/Message**                 | - Use **PascalCase**. <br> Example: `UserInfo`                                |
| **Field**                          | - Use **snake_case**. <br> Example: `user_id`. If a field name contains a number, the number should follow a letter, not an underscore. |
| **Enum**, **Service**, **Union**   | - Use **PascalCase**. <br> Example: `UserType`                                |
| **Enum Values**                    | - Use **UPPER_SNAKE_CASE**. <br> Example: `ADMIN_USER`                        |
| **RPC Methods**                    | - Use **PascalCase**. <br> Example: `GetUserInfo`                             |
| **Package/Namespace**              | - Use **snake_case**, typically based on the project structure. <br> Example: `com.project.service` |

#### Naming Conventions Explained:
- **PascalCase**: Capitalize the first letter of each word, such as `UserInfo`.
- **snake_case**: All lowercase with underscores separating words, such as `user_info`.
- **UPPER_SNAKE_CASE**: All uppercase letters with underscores separating words, such as `ADMIN_USER`.

## More Information

For more usage details, refer to the [Examples](example).