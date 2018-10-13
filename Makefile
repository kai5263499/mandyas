.PHONY: all get_dependencies generate_proto build_service

get_dependencies:
	dep ensure

generate_proto:
	protoc protos/minecraft.proto --proto_path=protos --go_out=plugins=grpc:generated

build_service: get_dependencies generate_proto
    go build -o mandyas .

all: build_service
