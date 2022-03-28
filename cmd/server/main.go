package main

import (
	"log"
	"net"

	. "challenge/pkg/proto"
	. "challenge/pkg/server"

	grpc "google.golang.org/grpc"
)

func RegisterServer(s grpc.ServiceRegistrar, srv TestRPCServer) {
	s.RegisterService(&ChallengeService_ServiceDesc, srv)
}

func main() {
	listener, err := net.Listen("tcp", ":8088")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	RegisterServer(grpcServer, TestRPCServer{})
	grpcServer.Serve(listener)
	print("qwerty")
}
