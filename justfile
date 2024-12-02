generate:
  protoc --go_out=. \
  --go_opt=module=github.com/araminian/grpcch4 \
  --go-grpc_out=. \
  --go-grpc_opt=module=github.com/araminian/grpcch4 \
  proto/dummy/v1/dummy.proto