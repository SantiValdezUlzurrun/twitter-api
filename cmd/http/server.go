package http

import (
	"github.com/gin-gonic/gin"
	"twitter-api/internal/tweets"
	"twitter-api/internal/users"
)

type Server struct {
	engine        *gin.Engine
	usersHandler  *users.Handler
	tweetsHandler *tweets.Handler
}

func NewServer(engine *gin.Engine, usersHandler *users.Handler, tweetsHandler *tweets.Handler) *Server {
	server := &Server{
		engine:        engine,
		usersHandler:  usersHandler,
		tweetsHandler: tweetsHandler,
	}
	server.registerRoutes()
	return server
}

func (s *Server) Run(port string) {
	if err := s.engine.Run(port); err != nil {
		panic(err)
	}
}

func (s *Server) registerRoutes() {
	usersRouter := s.engine.Group("/users")
	usersRouter.POST("/", s.usersHandler.Create)
	usersRouter.POST("/:id/follow/:followerId", s.usersHandler.Follow)
	usersRouter.POST("/:id/unfollow/:followerId", s.usersHandler.Unfollow)

	tweetsRouter := s.engine.Group("/tweets")
	tweetsRouter.POST("/", s.tweetsHandler.Create)
	tweetsRouter.DELETE("/:id", s.tweetsHandler.Delete)
	tweetsRouter.GET("/feed/:id", s.tweetsHandler.GetFeed)
}
