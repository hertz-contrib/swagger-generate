namespace go example

include "openapi.thrift"

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
    2: list<string> items (
        api.query = "items"
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

    HelloResp PathMethod(1: PathReq req) (
        api.get = "/path:path1"
    )

    HelloResp BodyMethod(1: BodyReq req) (
        api.post = "/body"
    )
}(
    api.base_domain = "127.0.0.1:8080",
    openapi.document = '{
       info: {
          title: "example swagger doc",
          version: "Version from annotation"
       }
    }'
)