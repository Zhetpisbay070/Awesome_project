package server

import (
	"awesomeProject1/internal/entity"
	"awesomeProject1/internal/repository"
	_ "awesomeProject1/internal/repository"
	"awesomeProject1/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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
	s.router.POST("/edit", s.EditOrder)

	return s.router
}

func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

func (s *Server) CreateOrder(ctx *gin.Context) {
	var req CreateOrderRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	order, err := s.service.CreateOrder(ctx, &entity.CreateOrderRequest{
		UserID:       req.UserID,
		Products:     req.Products,
		Price:        req.Price,
		DeliveryType: entity.DType(req.DeliveryType),
		AddressID:    req.AddressID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	orderDTO := Order{
		ID:               order.ID,
		UserID:           order.UserID,
		ProductIDs:       order.ProductIDs,
		CreatedAt:        order.CreatedAt,
		UpdatedAt:        order.UpdatedAt,
		DeliveryDeadLine: order.DeliveryDeadLine,
		Price:            order.Price,
		DeliveryType:     string(order.DeliveryType),
		Address:          order.Address,
		OrderStatus:      string(order.OrderStatus),
	}

	ctx.JSON(http.StatusOK, orderDTO)
}

func (s *Server) Update() {
	s.router.POST("/update-status", s.UpdateOrderStatus)
}

func (s *Server) UpdateOrderStatus(ctx *gin.Context) {

	var req UpdateOrderStatusRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.service.UpdateOrderStatus(ctx, req.OrderStatus, req.OrderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func (s *Server) GetOrders(ctx *gin.Context) {
	var req GetOrdersRequest

	err := ctx.BindJSON(&req)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders, err := s.service.GetOrders(ctx, &req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (s *Server) EditOrder(ctx *gin.Context) {
	var req EditOrderRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := s.service.EditOrder(ctx, &entity.EditOrderRequest{
		OrderID:  req.OrderID,
		Products: req.Products,
		Address:  req.Address,
	})

	if errors.Is(err, entity.ErrOrderNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, order)

}
