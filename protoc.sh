#!/bin/sh

protoc \
  --proto_path ./server/.protos/v1 \
  --go_opt paths=source_relative \
  --go_out ./server/grpc \
  --go-grpc_opt paths=source_relative \
  --go-grpc_out ./server/grpc \
  --plugin protoc-gen-ts_proto=./client/node_modules/.bin/protoc-gen-ts_proto \
  --ts_proto_opt outputServices=grpc-js,env=node,esModuleInterop=true \
  --ts_proto_out ./client/src/grpc \
  trading.proto \
  board.proto
