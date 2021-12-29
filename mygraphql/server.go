package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"pehks1980/be2_hw9/graphql1/graph"
	"pehks1980/be2_hw9/graphql1/graph/generated"
	"pehks1980/be2_hw9/graphql1/internal/pkg/repository"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var dbconn = flag.String("pg", "postgres://postuser:postpassword@192.168.1.204:5432/grpc?sslmode=disable", "postgres conn string: postgres://usr:pwd@localhost:5432/example?sslmode=disable")
	flag.Parse()

	ctx := context.Background()
	var err error
	repository.Reposit, err = repository.NewRepo(ctx, *dbconn)
	if err != nil {
		log.Fatalf("repository problem: %v", err)
	}

	err = repository.Reposit.InitSchema(ctx)
	if err != nil {
		log.Fatalf("repository DDL exec problem: %v", err)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
