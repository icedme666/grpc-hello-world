# 制作证书
* 私钥
  ```bash
  openssl genrsa -out certs/server.key 2048
  openssl ecparam -genkey -name secp384r1 -out certs/server.key
  ```
* 自签名公钥
  ```bash
  openssl req -new -x509 -sha256 -key certs/server.key -out certs/server.pem -days 3650 -extensions req_ext -config certs/server.conf
  ```
# proto
  1. google.api
     + proto/google/api/annotations.proto
     + proto/google/api/http.proto
  2. proto/hello.proto
     + 定义了一个serviceRPC服务HelloWorld，在其内部定义了一个HTTP Option的POST方法，HTTP响应路径为/hello_world
     + 定义message类型HelloWorldRequest、HelloWorldResponse，用于响应请求和返回结果
  3. 编译
     ```bash
     cd proto
     # 编译google.api
     protoc -I . --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:google/api google/api/*.proto
     
     #编译hello_http.proto为hello_http.pb.proto
     protoc -I . --go_out=plugins=grpc,Mgoogle/api/annotations.proto=grpc-hello-world/proto/google/api:. ./hello.proto
     
     #编译hello_http.proto为hello_http.pb.gw.proto
     protoc --grpc-gateway_out=logtostderr=true:. ./hello.proto
     ```

# 命令行模块cmd
* 实现server：server/server.gp
* 实现cmd
  - cmd/root.go
  - cmd/server.go
* 实现main.go
* 验证
  ```bash
  go run main.go server
  go run main.go server --port=8000 --cert-pem=test-pem --cert-key=test-key --cert-name=test-name
  ```

# 服务端模块server
1. server/hello.go
2. 实现服务端程序：
   - pkg/util/grpc.go
   - pkg/util/tls.go
3. 修改：server/server.go
# 验证
* 客户端实现：client/client.go
* 访问
  ```bash
  go run main.go server
  go run client/client.go
  curl -X POST -k https://localhost:50052/hello_world -d '{"referer": "restful_api"}'
  ```
