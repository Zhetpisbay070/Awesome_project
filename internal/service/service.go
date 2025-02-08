package service

import (
	"awesomeProject1/internal/entity"
	"awesomeProject1/internal/repository"
	"context"
	_ "errors"
	_ "fmt"
	"github.com/sirupsen/logrus"
	"time"
)

var _ OrderService = (*service)(nil)

//go:generate mockery --name=OrderService --with-expecter --output=../mock --outpkg=mock --case=underscore

type OrderService interface {
	CreateOrder(ctx context.Context, req *entity.CreateOrderRequest) (*entity.Order, error)
	UpdateOrderStatus(ctx context.Context, orderStatus entity.OrderStatus, orderID string) error
	GetOrders(ctx context.Context, req *entity.GetOrders) ([]entity.Order, error)
	EditOrder(ctx context.Context, req *entity.EditOrderRequest) (*entity.Order, error)
}

func NewOrderService(repo repository.DB, uuidFunc func() string, logger *logrus.Logger) OrderService {
	return &service{repo: repo, uuidFunc: uuidFunc, logger: logger}
}

type service struct {
	logger   *logrus.Logger
	repo     repository.DB
	uuidFunc func() string
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
		ID:           s.uuidFunc(),
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
		} else if orderStatus == entity.Cancelled {
			order.OrderStatus = entity.Cancelled
		} else {
			return entity.InvalidTransition
		}
	}

	if order.OrderStatus == entity.Paid {
		if orderStatus == entity.Collect {
			order.OrderStatus = entity.Collect
		} else if orderStatus == entity.Cancelled {
			order.OrderStatus = entity.Cancelled
		} else {
			return entity.InvalidTransition
		}

	} else if order.OrderStatus == entity.Collect {
		if orderStatus == entity.Collected {
			order.OrderStatus = entity.Collected
		} else if orderStatus == entity.Cancelled {
			order.OrderStatus = entity.Cancelled
		}
	} else {
		return entity.InvalidTransition
	}

	if order.OrderStatus == entity.Collected {
		if orderStatus == entity.Delivery {
			order.OrderStatus = entity.Delivery
		} else if orderStatus == entity.Cancelled {
			order.OrderStatus = entity.Cancelled
		}
	} else {
		return entity.InvalidTransition
	}

	if order.OrderStatus == entity.Delivery {
		if orderStatus == entity.Done {
		} else if orderStatus == entity.Cancelled {
			order.OrderStatus = entity.Cancelled
		}
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

	err = s.repo.UpdateOrder(ctx, order)
	if err != nil {
		return err
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

func (s *service) EditOrder(ctx context.Context, req *entity.EditOrderRequest) (*entity.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return nil, entity.ErrOrderNotFound
	}

	if order.OrderStatus == entity.Delivery || order.OrderStatus == entity.Done {
		return nil, entity.OrderCannotBeEdited
	}

	if order.OrderStatus == entity.Done {
		return nil, entity.AddressCannotBeEdited
	}

	if req.Address != "" {
		order.Address = req.Address
	}

	if len(req.Products) > 0 {
		order.ProductIDs = req.Products
	}

	err = s.repo.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, err
}
