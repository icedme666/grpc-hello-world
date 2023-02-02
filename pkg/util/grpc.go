package util

import (
	"net/http"
	"strings"
	"google.golang.org/grpc"
)

// 判断请求是来源于Rpc客户端还是Restful Api的请求，根据不同的请求注册不同的ServeHTTP服务
func GrpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	if otherHandler == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			grpcServer.ServeHTTP(w, r)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {  //表示请求必须基于HTTP/2
			grpcServer.ServeHTTP(w, r)
		}else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}