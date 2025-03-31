package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func Handler(resolver *Resolver) gin.HandlerFunc {
	h := handler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers: resolver,
		Schema:    parsedSchema,
	}))

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			playground.Handler("GraphQL", "/query").ServeHTTP(c.Writer, c.Request)
			return
		}

		h.ServeHTTP(c.Writer, c.Request)
	}
}
