.PHONY: all go

all:go

out_dir="./proto"

go:
	-mkdir ${out_dir}
	protoc \
      -I . \
      --go_out=${out_dir} --go_opt=paths=source_relative \
      --go-grpc_out=${out_dir} --go-grpc_opt=paths=source_relative \
      --validate_out="lang=go:." \
      --grpc-gateway_out ${out_dir} --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
      *.proto

