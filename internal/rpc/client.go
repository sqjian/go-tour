package rpc

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/sqjian/go-tour/internal/rpc/idl/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	var copts []grpc_retry.CallOption
	copts = append(copts, grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100*time.Millisecond)))
	copts = append(copts, grpc_retry.WithCodes(codes.NotFound, codes.Aborted))

	var dopts []grpc.DialOption
	dopts = append(dopts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	dopts = append(dopts, grpc.WithPerRPCCredentials(new(cliCustomCredential)))
	dopts = append(dopts, grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		cliInterceptor,
		grpc_retry.UnaryClientInterceptor(copts...),
		grpc_validator.UnaryClientInterceptor(),
	)))
	dopts = append(dopts, grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
		grpc_retry.StreamClientInterceptor(copts...),
	)))

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), dopts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	req := &proto.HelloRequest{Name: "world"}
	res, err := c.SayHello(
		context.Background(),
		req,
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(1*time.Second),
	)
	if err != nil {
		return err
	}
	log.Printf(res.Message)
	return nil
}
