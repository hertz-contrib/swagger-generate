namespace go example

include "openapi.thrift"

// Hello - request
struct FormReq {
    1: string FormValue (
        api.form = "form1",
        openapi.property = '{
            title: "this is an override field schema title",
            max_length: 255
        }'
    )

    2: InnerForm FormValue1 (
        api.form = "form3"
    )
}(
    openapi.schema = '{
          title: "Hello - request",
          description: "Hello - request",
          required: [
             "form1"
          ]
       }'
)

struct InnerForm {
        1: string InnerFormValue(
                api.form = "form2"
            )
    }

// QueryReq
struct QueryReq {
    1: string QueryValue (
        api.query = "query2",
        openapi.parameter = '{
            required: true
        }',
        openapi.property = '{
            title: "Name",
            description: "Name",
            type: "string",
            min_length: 1,
            max_length: 50
        }'
    )
    6: list<string> items (
        api.query = "items"
    )
    7: map<string, string> strings_map (
        api.query = "query1"
    )
}

// PathReq
struct PathReq {
    //field: path描述
    1: string PathValue (
        api.path = "path1"
    )
}

//BodyReq
struct BodyReq {
    //field: body描述
    1: string BodyValue (
        api.body = "body"
    )
    //field: query描述
    2: string QueryValue (
        api.query = "query2"
    )
    //field: body1描述
    3: string Body1Value (
        api.body = "body1"
    )
}

// HelloReq
struct HelloReq {
    1: string Name (
        api.query = "name",
        openapi.property = '{
            title: "Name",
            description: "Name",
            type: "string",
            min_length: 1,
            max_length: 50
        }'
    )
}

// HelloResp
struct HelloResp {
    1: string RespBody (
        api.body = "body",
        openapi.property = '{
            title: "response content",
            description: "response content",
            type: "string",
            min_length: 1,
            max_length: 80
        }'
    )
    2: string token (
        api.header = "token",
        openapi.property = '{
            title: "token",
            description: "token",
            type: "string"
        }'
    )
}(
    openapi.schema = '{
      title: "Hello - response",
      description: "Hello - response",
      required: [
         "body"
      ]
   }'
)

// HelloService1描述
service HelloService1 {
    HelloResp QueryMethod(1: QueryReq req) (
        api.get = "/hello1"
    )

    HelloResp FormMethod(1: FormReq req) (
        api.post = "/form"
    )

    HelloResp PathMethod(1: PathReq req) (
        api.get = "/path:path1"
    )

    HelloResp BodyMethod(1: BodyReq req) (
        api.post = "/body"
    )
}(
    api.base_domain = "127.0.0.1:8888",
    openapi.document = '{
       info: {
          title: "example swagger doc",
          version: "Version from annotation"
       }
    }'
)

// HelloService2描述
service HelloService2 {
    HelloResp QueryMethod(1: QueryReq req) (
        api.get = "/hello2"
        api.baseurl = "127.0.0.1:8889"
        openapi.operation = '{
            summary: "Hello - Get",
            description: "Hello - Get"
        }'
    )
}
