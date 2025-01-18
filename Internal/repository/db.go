//go:generate mockery --name=DB --with-expecter --output=../mock --outpkg=mock --case=underscore

package repository

import (
	"awesomeProject1/internal/entity"
	"context"
)

type DB interface {
	CreateOrder(ctx context.Context, order *entity.Order) error
	GetOrderByID(ctx context.Context, id string) (*entity.Order, error)
	ProductExist(ctx context.Context, productID string) (bool, error)
	UpdateOrder(ctx context.Context, order *entity.Order) error
	GetOrders(ctx context.Context, req *entity.GetOrders) ([]entity.Order, error)
}
