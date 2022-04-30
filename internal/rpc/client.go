package rpc

import (
	"fmt"
	"github.com/sqjian/go-tour/internal/rpc/idl/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type cliCustomCredential struct{}

func (c cliCustomCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return credential, nil
}

// RequireTransportSecurity 是否开启TLS
func (c cliCustomCredential) RequireTransportSecurity() bool {
	return false
}

func cliInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("method=%s req=%v rep=%v duration=%s error=%v	\n", method, req, reply, time.Since(start), err)
	return err
}

func StartCli(addr, port string) error {
	return startCli(addr, port)
}
func startCli(addr, port string) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithPerRPCCredentials(new(cliCustomCredential)))
	opts = append(opts, grpc.WithUnaryInterceptor(cliInterceptor))
	// 连接
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), opts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	// 初始化客户端
	c := proto.NewGreeterClient(conn)
	// 调用方法
	req := &proto.HelloRequest{Name: "world"}
	res, err := c.SayHello(context.Background(), req)
	if err != nil {
		return err
	}
	log.Printf(res.Message)
	return nil
}
