package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

// GraphQLHandler handles GraphQL requests
func GraphQLHandler(resolver *Resolver) gin.HandlerFunc {
	h := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))

	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			h := playground.Handler("GraphQL", "/query")
			h.ServeHTTP(c.Writer, c.Request)
			return
		}

		h.ServeHTTP(c.Writer, c.Request)
	}
}
