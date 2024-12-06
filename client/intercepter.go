package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const authTokenKey = "auth_token"
const authTokenValue = "secret"

func unaryAuthInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	ctx = metadata.AppendToOutgoingContext(ctx, authTokenKey, authTokenValue)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func streamAuthInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, authTokenKey, authTokenValue)
	return streamer(ctx, desc, cc, method, opts...)
}
