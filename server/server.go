package server

import (

	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/asdine/storm/v3"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/graph"
	"github.com/minskylab/supersense/graph/generated"
	log "github.com/sirupsen/logrus"
)

const defaultPort = "8080"

type Server struct {
	mux *supersense.Mux
	secret []byte
	db *storm.DB
}

// LaunchServer launch the graphQL server
func (s *Server) LaunchServer(port string) {
	if port == "" {
		port = defaultPort
	}

	resolver := graph.NewResolver(s.mux, s)

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", srv)

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return s.secret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	app := jwtMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user")
		log.Info("This is an authenticated request")
		log.Info("Claim content:\n")
		for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
			log.Info("%s :\t%#v\n", k, v)
		}
	}))

	log.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, app))
}
