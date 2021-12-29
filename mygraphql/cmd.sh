#install all needed
#!/bin/bash
go get github.com/99designs/gqlgen/graphql/handler/transport@v0.14.0
go get github.com/99designs/gqlgen/graphql/handler/lru@v0.14.0
go get github.com/99designs/gqlgen/graphql/handler/extension@v0.14.0
#go get github.com/99designs/gqlgen
#go run github.com/99designs/gqlgen init
# make schema.graphqls !!! при изменение схемы обработчики в schema.resolvers.go не затираются!!!
go get github.com/99designs/gqlgen/cmd@v0.14.0
go get github.com/99designs/gqlgen/internal/imports@v0.14.0
go get github.com/99designs/gqlgen/internal/code@v0.14.0
go run github.com/99designs/gqlgen generate
go run ./server.go
