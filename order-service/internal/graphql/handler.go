package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// Handler sets up the GraphQL server and playground
func Handler(resolver ResolverRoot) http.Handler {
	cfg := Config{Resolvers: resolver}
	srv := handler.NewDefaultServer(NewExecutableSchema(cfg))

	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	mux.Handle("/graphql", srv)

	return mux
}
