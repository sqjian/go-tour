package rpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sqjian/go-tour/internal/rpc/idl/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"sync"
)

type server struct {
	proto.UnimplementedGreeterServer
}

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
		return status.Errorf(codes.Unauthenticated, "No Token authentication information")
	}
	var (
		appid  string
		appkey string
	)

	fmt.Printf("got md:%+v\n", md)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}
	if appid != "101010" || appkey != "i am key" {
		return status.Errorf(codes.Unauthenticated, "Token authentication information is invalid: appid =%s, appkey =%s", appid, appkey)
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

func StartSrv(grpcAddr, gatewayAddr string) error {
	return startSrv(grpcAddr, gatewayAddr)
}
func startSrv(grpcAddr, gatewayAddr string) error {
	startGrpcServer := func() error {
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			return err
		}
		var opts []grpc.ServerOption
		opts = append(opts, grpc.UnaryInterceptor(srvInterceptor))
		s := grpc.NewServer(opts...)
		proto.RegisterGreeterServer(s, &server{})
		reflection.Register(s)
		fmt.Printf("start grpc server in:%v\n", grpcAddr)
		if err := s.Serve(lis); err != nil {
			return err
		}
		return nil
	}
	startGrpcGateway := func() error {
		conn, err := grpc.DialContext(
			context.Background(),
			grpcAddr,
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return fmt.Errorf("failed to dial server:%w", err)
		}

		gwMux := runtime.NewServeMux(
			runtime.WithMarshalerOption("application/json+pretty", &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					Indent:    "  ",
					Multiline: true, // Optional, implied by presence of "Indent".
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			}),
			runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
				md := metadata.MD{}
				metadata.Pairs()
				{
					appidHeader := request.Header.Get("appid")
					md.Set("appid", appidHeader)
				}
				{
					appkeyHeader := request.Header.Get("appkey")
					md.Set("appkey", appkeyHeader)
				}
				return md
			}),
		)
		err = proto.RegisterGreeterHandler(context.Background(), gwMux, conn)
		if err != nil {
			return fmt.Errorf("failed to register gateway:%w", err)
		}

		prettier := func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, ok := r.URL.Query()["pretty"]; ok {
					r.Header.Set("Accept", "application/json+pretty")
				}
				h.ServeHTTP(w, r)
			})
		}

		gwServer := &http.Server{
			Addr:    gatewayAddr,
			Handler: prettier(gwMux),
		}
		fmt.Printf("start grpc gateway server in:%v\n", gatewayAddr)
		return gwServer.ListenAndServe()
	}

	{
		var wg sync.WaitGroup
		wg.Add(2)

		go func(m *sync.WaitGroup) {
			defer wg.Done()
			grpcErr := startGrpcServer()
			if grpcErr != nil {
				panic(grpcErr)
			}
		}(&wg)
		go func(m *sync.WaitGroup) {
			defer wg.Done()
			gatewayErr := startGrpcGateway()
			panic(gatewayErr)
		}(&wg)

		wg.Wait()
	}

	return nil
}
