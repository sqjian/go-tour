package rpc

import (
	"context"
	"fmt"
	"github.com/sqjian/go-tour/internal/rpc/idl/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	return &proto.HelloReply{Message: "Hello " + in.Name}, nil
}

func srvAuth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	var (
		appid  string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}
	if appid != "101010" || appkey != "i am key" {
		return status.Errorf(codes.Unauthenticated, "Token认证信息无效:appid =%s, appkey =%s", appid, appkey)
	}
	return nil
}

// interceptor 拦截器
func srvInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := srvAuth(ctx)
	if err != nil {
		return nil, err
	}
	// 继续处理请求
	return handler(ctx, req)
}

func StartSrv(addr, port string) error {
	return startSrv(addr, port)
}
func startSrv(addr, port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", addr, port))
	if err != nil {
		return err
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(srvInterceptor))
	s := grpc.NewServer(opts...)
	proto.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
