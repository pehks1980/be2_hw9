package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	pricelists "pehks1980/be2_hw9/grpc"
	"pehks1980/be2_hw9/grpc/api"
	"pehks1980/be2_hw9/internal/pkg/repository"
)


func main() {
	var addrport = flag.String("addrport", "127.0.0.1:9090", "Addr:Port for test gRPC server")
	var dbconn = flag.String("pg", "postgres://postuser:postpassword@192.168.1.204:5432/grpc?sslmode=disable", "postgres conn string: postgres://usr:pwd@localhost:5432/example?sslmode=disable")
	flag.Parse()

	lis, err := net.Listen("tcp", *addrport)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	ctx := context.Background()
	repo, err1 := repository.NewRepository(ctx, *dbconn)
	if err1 != nil {
		log.Fatalf("repository problem: %v", err1)
	}

	err = repo.InitSchema(ctx)
	if err != nil {
		log.Fatalf("repository DDL exec problem: %v", err)
	}
	// register server
	server := &api.Server{
		Repo: repo,
		Ctx:  ctx,
	}

	grpcServer := grpc.NewServer()
	pricelists.RegisterPricelistsServer(grpcServer, server)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(grpcServer)

	log.Printf("started grpc server at: %s", *addrport)
	// start to serve
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
