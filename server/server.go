package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/graph"
	"github.com/minskylab/supersense/graph/generated"
	log "github.com/sirupsen/logrus"
)

const defaultPort = 8080

// LaunchServer launch the graphQL server
func LaunchServer(mux *supersense.Mux, port int64, withGraphQLPlayground bool) error {
	if port <= 0 {
		port = defaultPort
	}

	resolver := graph.NewResolver(mux)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolver }))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})

	http.Handle("/graphql", srv)
	if withGraphQLPlayground {
		http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
		log.Infof("connect to http://localhost:%d/ for GraphQL playground", port)
	}


	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
