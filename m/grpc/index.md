protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    hello/hello.proto

## gRPC Streaming 是基于 HTTP/2 的

## grpc stream的优势，适用场景
1. 大数据量
2. 实时性要求高