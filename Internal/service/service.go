package service

import (
	"awesomeProject1/internal/entity"
	"awesomeProject1/internal/repository"
	"context"
	_ "errors"
	_ "fmt"
	"time"

	"github.com/google/uuid"
)

var _ OrderService = (*service)(nil)

type OrderService interface {
	CreateOrder(ctx context.Context, req *entity.CreateOrderRequest) (*entity.Order, error)
	UpdateOrderStatus(ctx context.Context, orderStatus entity.OrderStatus, orderID string) error
	GetOrders(ctx context.Context, req *entity.GetOrders) ([]entity.Order, error)
}

func NewOrderService(repo repository.DB) OrderService {
	return &service{repo: repo}
}

type service struct {
	repo repository.DB
}

func (s *service) CreateOrder(ctx context.Context, req *entity.CreateOrderRequest) (*entity.Order, error) {
	for _, p := range req.Products {
		ok, err := s.repo.ProductExist(ctx, p)
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, entity.ProductDoesNotExistError
		}
	}

	now := time.Now()

	order := entity.Order{
		ID:           uuid.New().String(),
		UserID:       req.UserID,
		ProductIDs:   req.Products,
		CreatedAt:    now,
		UpdatedAt:    now,
		Price:        req.Price,
		DeliveryType: req.DeliveryType,
		Address:      req.AddressID,
		OrderStatus:  entity.Created,
	}

	err := s.repo.CreateOrder(ctx, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *service) UpdateOrderStatus(ctx context.Context, orderStatus entity.OrderStatus, orderID string) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {

		return err
	}

	if order.OrderStatus == entity.Created {
		if orderStatus == entity.Paid {
			order.OrderStatus = entity.Paid
		} else {
			return entity.InvalidTransition
		}
	} else if order.OrderStatus == entity.Paid {
		if orderStatus == entity.Collect {
			order.OrderStatus = entity.Collect
		} else {
			return entity.InvalidTransition
		}
	} else if order.OrderStatus == entity.Collect {
		if orderStatus == entity.Collected {
			order.OrderStatus = entity.Collected
		} else {
			return entity.InvalidTransition
		}
	} else if order.OrderStatus == entity.Collected {
		if orderStatus == entity.Delivery {
			order.OrderStatus = entity.Delivery
		} else {
			return entity.InvalidTransition
		}
	} else if order.OrderStatus == entity.Delivery {
		if orderStatus == entity.Done {
			order.OrderStatus = entity.Done
		} else {
			return entity.InvalidTransition
		}
		if order.OrderStatus == entity.Delivery || order.OrderStatus == entity.Done {
			if orderStatus == entity.Cancelled {
				return entity.PozdnoNahui
			} else {
				order.OrderStatus = entity.Cancelled
			}
		}
	}

	return nil
}

func (s *service) GetOrders(ctx context.Context, req *entity.GetOrders) ([]entity.Order, error) {
	orders, err := s.repo.GetOrders(ctx, req)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
