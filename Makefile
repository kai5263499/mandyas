.DEFAULT_GOAL := build_service
.PHONY: all clean get_dependencies generate_proto build_service

get_dependencies:
	dep ensure

generate_proto:
	protoc protos/minecraft.proto --proto_path=protos --go_out=plugins=grpc:generated

build_service: get_dependencies generate_proto
	go build -o mandyas cmd/server/*.go

all: build_service

clean:
	rm mandyas