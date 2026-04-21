package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"rate-limited-api/internal/adapters/http/handler"
	"rate-limited-api/internal/adapters/http/middleware"
	"rate-limited-api/internal/core/ports"
)

type Server struct {
	router      *gin.Engine
	port        int
	service     ports.RequestService
	rateLimiter ports.RateLimiter
}

func NewServer(port int, service ports.RequestService, rateLimiter ports.RateLimiter) *Server {
	router := gin.Default()
	
	s := &Server{
		router:      router,
		port:        port,
		service:     service,
		rateLimiter: rateLimiter,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	reqHandler := handler.NewRequestHandler(s.service)
	statsHandler := handler.NewStatsHandler(s.service)

	api := s.router.Group("/api")
	{
		api.GET("/stats", statsHandler.HandleGetStats)
		api.POST("/request", middleware.RateLimitMiddleware(s.rateLimiter), reqHandler.HandlePostRequest)
	}
}

func (s *Server) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.port))
}
