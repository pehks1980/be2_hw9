package serve

import (
	"lesson09/internal/api/grpc/apisrv"
	"lesson09/internal/api/grpc/myservice"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RunMyServer() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	// register server
	server := &apisrv.MyServer{}

	grpcServer := grpc.NewServer()
	myservice.RegisterMessageServiceServer(grpcServer, server)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(grpcServer)

	// start to serve
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
