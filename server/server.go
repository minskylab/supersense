package server

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/minskylab/supersense"
	"github.com/minskylab/supersense/graph"
	"github.com/minskylab/supersense/graph/generated"
	"github.com/minskylab/supersense/persistence"
	log "github.com/sirupsen/logrus"
)

const defaultPort = "8080"

type Server struct {
	mux *supersense.Mux
	db *persistence.Persistence
	router *gin.Engine
}

// Defining the Graphql handler
func (s *Server)  graphqlHandler() gin.HandlerFunc {
	resolver := graph.NewResolver(s.mux, s.db)

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

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func (s *Server) playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/graphql")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}


type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
// LaunchServer launch the graphQL server
func (s *Server) LaunchServer(port string ) {
	if port == "" {
		port = defaultPort
	}

	identityKey := "id"

	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "supersense",
		Key:         []byte("secret key"),
		Timeout:     1*time.Hour,
		MaxRefresh:  1*time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*persistence.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &persistence.User{
				ID: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginPayload login
			if err := c.ShouldBind(&loginPayload); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			user, err := s.db.LoginWithUserPassword(loginPayload.Username, loginPayload.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication // errors.WithStack(err)
			}
			return user, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*persistence.User); ok && v.Username == "admin" {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Panic(err)
	}

	// return s.db.GetSecret(), nil
	s.router.POST("/graphql", s.graphqlHandler())
	s.router.GET("/", s.playgroundHandler())

	s.router.POST("/login", authMiddleware.LoginHandler)
	s.router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := s.router.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	// auth.GET("/hello", )


	log.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	if err := s.router.Run(); err != nil {
		log.Panic(err)
	}
}
