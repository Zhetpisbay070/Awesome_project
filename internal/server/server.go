package server

import (
	"awesomeProject1/internal/repository"
	_ "awesomeProject1/internal/repository"
	"awesomeProject1/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	service service.OrderService
	router  *gin.Engine
	logger  *logrus.Logger
	repo    repository.DB
}

func NewServer(service service.OrderService, logger *logrus.Logger) *Server {
	return &Server{router: gin.New(), service: service, logger: logger}
}

func (s *Server) Run(port string) error {

	return s.router.Run(":" + port)
}

func (s *Server) SetupRouter() *gin.Engine {
	s.router = gin.New()

	s.router.POST("/create", s.CreateOrder)
	s.router.POST("/update", s.UpdateOrderStatus)
	s.router.POST("/edit", s.EditOrder)
	s.router.POST("/getOrders", s.GetOrders)

	return s.router
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
