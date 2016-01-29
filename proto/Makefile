PROTO_PATH := ${GOPATH}/src:.
PROTO_PATH := ${PROTO_PATH}:${GOPATH}/src/github.com/gogo/protobuf/protobuf
PROTO_PATH := ${PROTO_PATH}:${GOPATH}/src/github.com/gogo/protobuf/gogoproto

.PHONY: all

all:
	protoc --proto_path=${PROTO_PATH} --gogo_out=. ./service/*.proto
	protoc --proto_path=${PROTO_PATH} --gogo_out=. ./order/*.proto
	protoc --proto_path=${PROTO_PATH} --gogo_out=. ./pricing/*.proto
