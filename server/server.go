package server

import (
	"log"
	"crypto/tls"
	"net"
	"net/http"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "xiamei.guo/grpc-hello-world/proto"
	"xiamei.guo/grpc-hello-world/pkg/util"
)

var (
	ServerPort string
	CertName string
	CertPemPath string
	CertKeyPath string
	EndPoint string
)

func Serve() (err error) {
	// 1.启动监听
	EndPoint = ":" + ServerPort
	conn, err := net.Listen("tcp", EndPoint)  //用于监听本地的网络地址通知，返回一个监听器的结构体
	if err != nil {
		log.Printf("TCP Listen err: %v\n", err)
	}
    
	// 2. 获取TLS
	tlsConfig := util.GetTLSConfig(CertPemPath, CertKeyPath)  // 解析得到tls.Config，传达给http.Server服务的TLSConfig配置项使用
	// 3. 创建内部服务
	srv := createInternalServer(conn, tlsConfig)

	log.Printf("gRPC and https listen on: %s\n", ServerPort)
	
	// 创建tls.NewListener
	// 服务开始接受请求
	if err = srv.Serve(tls.NewListener(conn, tlsConfig)); err != nil {
		log.Printf("ListenAndServe: %v\n", err)
	}

	return err
}

func createInternalServer(conn net.Listener, tlsConfig *tls.Config) (*http.Server) {
	var opts []grpc.ServerOption

	// 1. grpc server：创建grpc的TLS认证凭证
	creds, err := credentials.NewServerTLSFromFile(CertPemPath, CertKeyPath)  //从输入证书文件和服务器的密钥文件构造TLS证书凭证
	if err != nil {
		log.Printf("Failed to create server TLS credentials %v", err)
	}

	// 2. 设置grpc ServerOption
	opts = append(opts, grpc.Creds(creds))
	// 3.创建grpc服务端
	grpcServer := grpc.NewServer(opts...)

	// 4.注册grpc服务
	pb.RegisterHelloWorldServer(grpcServer, NewHelloService())

	// 5. 创建grpc-gateway关联组件
	ctx := context.Background()  //返回一个非空的空上下文
	dcreds, err := credentials.NewClientTLSFromFile(CertPemPath, CertName)  //从客户机的输入证书文件构造TLS凭证
	if err != nil {
		log.Printf("Failed to create client TLS credentials %v", err)
	}
	// grpc.WithTransportCredentials：配置一个连接级别的安全凭据(例：TLS、SSL)，返回值为type DialOption
	// grpc.DialOption：配置我们如何设置连接（其内部具体由多个的DialOption组成，决定其设置连接的内容）
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
	// 6.创建HTTP NewServeMux及注册grpc-gateway逻辑
	gwmux := runtime.NewServeMux()  //返回一个新的ServeMux，它的内部映射是空的

	// 7.注册具体服务
	if err := pb.RegisterHelloWorldHandlerFromEndpoint(ctx, gwmux, EndPoint, dopts); err != nil {  //注册HelloWorld服务的HTTP Handle到grpc端点
		log.Printf("Failed to register gw server: %v\n", err)
	}

	//http服务
	mux := http.NewServeMux()  //分配并返回一个新的ServeMux
	mux.Handle("/", gwmux)  //为给定模式注册处理程序

	return &http.Server {
		Addr: EndPoint,
		Handler: util.GrpcHandlerFunc(grpcServer, mux),
		TLSConfig: tlsConfig,
	}
}
	








