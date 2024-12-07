package main

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const authTokenKey = "auth_token"
const authTokenValue = "secret"

func validateAuthToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	log.Printf("Server: metadata: %v", md)

	if !ok {
		return nil, status.Error(codes.Unauthenticated, "auth token is required")
	}

	if t, ok := md[authTokenKey]; ok {
		switch {
		case len(t) != 1:
			return nil, status.Error(codes.InvalidArgument, "auth token must be provided exactly once")
		case t[0] != authTokenValue:
			return nil, status.Error(codes.Unauthenticated, "invalid auth token")
		}
	} else {
		return nil, status.Error(codes.Unauthenticated, "auth token is required")
	}

	return ctx, nil
}
