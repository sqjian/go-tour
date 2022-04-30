export GO111MODULE=on
export GOPROXY=https://goproxy.cn

.PHONY:idl
idl:
	cd idl && make

.PHONY:env
env:
	go install github.com/golang/protobuf/protoc-gen-go@latest
	go install github.com/golang/protobuf/proto@latest
	go install google.golang.org/grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest