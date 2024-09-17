namespace go example

include "openapi.thrift"

// QueryReq
struct QueryReq {
    1: string QueryValue (
        openapi.property = '{
            title: "Name",
            description: "Name",
            type: "string",
            min_length: 1,
            max_length: 50
        }'
    )
    2: list<string> Items ()
}

// PathReq
struct PathReq {
    //field: path描述
    1: string PathValue ()
}

//BodyReq
struct BodyReq {
    //field: body描述
    1: string BodyValue ()

    //field: query描述
    2: string QueryValue ()
}

// HelloResp
struct HelloResp {
    1: string RespBody (
        openapi.property = '{
            title: "response content",
            description: "response content",
            type: "string",
            min_length: 1,
            max_length: 80
        }'
    )
    2: string token (
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
         "RespBody"
      ]
   }'
)

// HelloService1描述
service HelloService1 {
    HelloResp QueryMethod(1: QueryReq req) ()

    HelloResp PathMethod(1: PathReq req) ()

    HelloResp BodyMethod(1: BodyReq req) ()
}(
    api.base_domain = "127.0.0.1:8888",
    openapi.document = '{
       info: {
          title: "example swagger doc",
          version: "Version from annotation"
       }
    }'
)