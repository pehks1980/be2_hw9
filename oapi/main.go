package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"pehks1980/be2_hw9/oapi/api"

	middleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
)

// taken from codegen example of petstore with minimal change
func main() {
	var port = flag.Int("port", 9000, "Port for test HTTP server")
	var dbconn = flag.String("pg", "postgres://postuser:postpassword@192.168.1.204:5432/grpc?sslmode=disable", "postgres conn string: postgres://usr:pwd@localhost:5432/example?sslmode=disable")
	flag.Parse()

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create an instance of our handler which satisfies the generated interface
	ctx := context.Background()
	myapi, _ := api.NewApi(ctx, *dbconn)

	// This is how you set up a basic chi router
	r := chi.NewRouter()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))

	// We now register our myapi above as the handler for the interface
	api.HandlerFromMux(myapi, r)
	log.Printf("started http server at port :%d", *port)
	s := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
	}

	// And we serve HTTP until the world ends.
	log.Fatal(s.ListenAndServe())
}
