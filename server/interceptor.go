package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const authTokenKey = "auth_token"
const authTokenValue = "secret"

func validateAuthToken(ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("Server: metadata: %v", md)

	if t, ok := md[authTokenKey]; ok {
		switch {
		case len(t) != 1:
			return status.Error(codes.InvalidArgument, "auth token must be provided exactly once")
		case t[0] != authTokenValue:
			return status.Error(codes.Unauthenticated, "invalid auth token")
		}
	} else {
		return status.Error(codes.Unauthenticated, "auth token is required")
	}

	return nil
}

func unaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := validateAuthToken(ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func streamAuthInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := validateAuthToken(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}
