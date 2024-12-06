package main

import (
	"log"
	"net"
	"os"

	pb "github.com/araminian/grpcch4/proto/todo/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	_ "google.golang.org/grpc/encoding/gzip"
)

type server struct {
	d db
	pb.UnimplementedTodoServiceServer
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatalln("usage: server [IP_ADDR]")
	}

	addr := args[0]
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	defer func(lis net.Listener) {
		if err := lis.Close(); err != nil {
			log.Fatalf("unexpected error: %v", err)
		}
	}(lis)

	// Read Certificates and keys
	creds, err := credentials.NewServerTLSFromFile("../certs/server_cert.pem", "../certs/server_key.pem")
	if err != nil {
		log.Fatalf("failed to load certificates: %v", err)
	}

	// Add interceptors
	var opts []grpc.ServerOption = []grpc.ServerOption{
		grpc.UnaryInterceptor(unaryAuthInterceptor),
		grpc.StreamInterceptor(streamAuthInterceptor),
		grpc.Creds(creds), // Add TLS credentials
	}
	s := grpc.NewServer(opts...)

	//registration of endpoints
	pb.RegisterTodoServiceServer(s, &server{d: New()})

	log.Printf("listening at %s\n", addr)

	defer s.Stop()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
