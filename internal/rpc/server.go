package rpc

import (
	"context"
	"fmt"
	"github.com/fullstorydev/grpcui/standalone"
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

func srvInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := srvAuth(ctx)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func StartSrv(grpcAddr, gatewayAddr, grpcUiAddr string) error {
	return startSrv(grpcAddr, gatewayAddr, grpcUiAddr)
}

func startSrv(grpcAddr, gatewayAddr, grpcUiAddr string) error {
	grpcListener, grpcListenerErr := func(addr string) (net.Listener, error) {
		lis, err := net.Listen("tcp", addr)
		return lis, err
	}(grpcAddr)
	if grpcListenerErr != nil {
		return grpcListenerErr
	}

	grpcServer, grpcServerErr := func() (*grpc.Server, error) {
		var opts []grpc.ServerOption
		opts = append(opts, grpc.UnaryInterceptor(srvInterceptor))
		s := grpc.NewServer(opts...)
		proto.RegisterGreeterServer(s, &server{})
		reflection.Register(s)
		return s, nil
	}()
	if grpcServerErr != nil {
		return grpcServerErr
	}

	mux := http.NewServeMux()
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
	opts := func() []grpc.DialOption {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		return opts
	}()

	prettier := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := r.URL.Query()["pretty"]; ok {
				r.Header.Set("Accept", "application/json+pretty")
			}
			h.ServeHTTP(w, r)
		})
	}

	startGrpcServer := func(gs *grpc.Server, lis net.Listener) error {
		fmt.Printf("start grpc server in:%v\n", lis.Addr().String())
		return gs.Serve(lis)
	}
	startGrpcUi := func(grpcAddr, grpcUiAddr string) error {
		conn, err := grpc.DialContext(
			context.Background(),
			grpcAddr,
			opts...,
		)
		if err != nil {
			return fmt.Errorf("failed to dial server:%w", err)
		}
		h, err := standalone.HandlerViaReflection(context.Background(), conn, grpcAddr)
		fmt.Printf("start grpc ui server on:%v\n", grpcUiAddr)

		serveMux := http.NewServeMux()

		serveMux.Handle("/grpcui/", http.StripPrefix("/grpcui", h))

		serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "html")
			_, _ = fmt.Fprintf(w, `
			<html>
			<head><title>Example server</title></head>
			<body>
			<h1>Hello, world!</h1>
			<p>Check out the gRPC UI <a href="/grpcui/">here</a>.</p>
			</body>
			</html>
		`)
		})

		l, err := net.Listen("tcp", grpcUiAddr)
		if err != nil {
			return err
		}

		err = http.Serve(l, serveMux)
		if err != nil {
			return err
		}
		return nil
	}
	startGateway := func(grpcAddr, gatewayAddr string, gwMux *runtime.ServeMux) error {
		conn, err := grpc.DialContext(
			context.Background(),
			grpcAddr,
			opts...,
		)
		if err != nil {
			return fmt.Errorf("failed to dial server:%w", err)
		}
		err = proto.RegisterGreeterHandler(context.Background(), gwMux, conn)
		if err != nil {
			return fmt.Errorf("failed to register gateway:%w", err)
		}
		mux.Handle("/", gwMux)
		fmt.Printf("start grpc gateway server on:%v\n", gatewayAddr)
		return http.ListenAndServe(gatewayAddr, prettier(gwMux))
	}
	tasks := func() []func(*sync.WaitGroup) {
		var tasks []func(*sync.WaitGroup)
		tasks = append(
			tasks,
			func(wg *sync.WaitGroup) {
				defer wg.Done()
				grpcErr := startGrpcServer(grpcServer, grpcListener)
				if grpcErr != nil {
					panic(grpcErr)
				}
			},
		)
		tasks = append(
			tasks,
			func(wg *sync.WaitGroup) {
				defer wg.Done()
				gatewayErr := startGateway(grpcAddr, gatewayAddr, gwMux)
				panic(gatewayErr)
			},
		)
		tasks = append(
			tasks,
			func(wg *sync.WaitGroup) {
				defer wg.Done()
				grpcUiErr := startGrpcUi(grpcAddr, grpcUiAddr)
				panic(grpcUiErr)
			},
		)
		return tasks
	}()

	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go task(&wg)
	}
	wg.Wait()

	return nil
}
