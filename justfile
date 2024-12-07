generate:
  protoc --go_out=. \
  --go_opt=module=github.com/araminian/grpcch4 \
  --go-grpc_out=. \
  --go-grpc_opt=module=github.com/araminian/grpcch4 \
  proto/dummy/v1/dummy.proto

generate-validate:
  protoc -Iproto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative --validate_out="lang=go,paths=source_relative:proto" proto/todo/v2/*.proto