package service_test

import (
	"awesomeProject1/Internal/service"
	"awesomeProject1/internal/entity"
	InternalMock "awesomeProject1/internal/mock"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateOrder(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		repo.EXPECT().CreateOrder(mock.Anything, mock.Anything).Return(nil)
		repo.EXPECT().ProductExist(mock.Anything, mock.Anything).Return(true, nil)

		products := []string{"salam"}

		order, err := s.CreateOrder(context.Background(), &entity.CreateOrderRequest{
			UserID:       "123",
			Products:     products,
			Price:        100,
			DeliveryType: entity.Courier,
			AddressID:    "asd",
		})

		assert.NoError(t, err)
		assert.Equal(t, order.UserID, "123")
		assert.EqualValues(t, len(order.ProductIDs), len(products))
		for _, productID := range order.ProductIDs {
			assert.Contains(t, products, productID)
		}
		assert.Equal(t, order.Price, 100.00)
		assert.Equal(t, order.DeliveryType, entity.Courier)
		assert.Equal(t, order.Address, "asd")
		assert.Equal(t, order.OrderStatus, entity.Created)
	})

	t.Run("product does not exist", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		repo.EXPECT().ProductExist(mock.Anything, mock.Anything).Return(false, nil)

		products := []string{"salam"}

		order, err := s.CreateOrder(context.Background(), &entity.CreateOrderRequest{
			UserID:       "123",
			Products:     products,
			Price:        100,
			DeliveryType: entity.Courier,
			AddressID:    "asd",
		})

		assert.ErrorIs(t, err, entity.ProductDoesNotExistError)
		assert.Nil(t, order)
	})

	t.Run("impossible to check products", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		repo.EXPECT().ProductExist(mock.Anything, mock.Anything).Return(false, entity.ImpossibleToCheckProducts)

		products := []string{"salam"}

		order, err := s.CreateOrder(context.Background(), &entity.CreateOrderRequest{
			UserID:       "123",
			Products:     products,
			Price:        100,
			DeliveryType: entity.Courier,
			AddressID:    "asd",
		})

		assert.ErrorIs(t, err, entity.ImpossibleToCheckProducts)
		assert.Nil(t, order)
	})

	t.Run("delivery unavailable", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		repo.EXPECT().CreateOrder(mock.Anything, mock.Anything).Return(entity.DelTypeUnavailable)
		repo.EXPECT().ProductExist(mock.Anything, mock.Anything).Return(true, nil)

		products := []string{"salam"}

		order, err := s.CreateOrder(context.Background(), &entity.CreateOrderRequest{
			UserID:       "123",
			Products:     products,
			Price:        100,
			DeliveryType: entity.Drone,
			AddressID:    "asd",
		})

		assert.ErrorIs(t, err, entity.DelTypeUnavailable)
		assert.Nil(t, order)
	})

	t.Run("success get order", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		orderList := []entity.Order{
			{ID: "1", UserID: "1"},
		}
		getOrderList := &entity.GetOrders{
			UserID: "1",
		}

		repo.EXPECT().GetOrders(mock.Anything, getOrderList).Return(orderList, nil)
		orders, err := s.GetOrders(context.Background(), getOrderList)

		assert.NoError(t, err)
		assert.Equal(t, orderList, orders)

	})

	t.Run("error - order not found", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		getOrderList := &entity.GetOrders{
			UserID: "2",
		}

		repo.EXPECT().GetOrders(mock.Anything, getOrderList).Return(nil, entity.OrderNotFound)

		orders, err := s.GetOrders(context.Background(), getOrderList)

		assert.ErrorIs(t, err, entity.OrderNotFound)
		assert.Nil(t, orders)

	})

	t.Run("transition create to paid", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		order := &entity.Order{
			OrderStatus: entity.Created,
			ID:          "1",
		}

		nextOrder := &entity.Order{
			OrderStatus: entity.Paid,
			ID:          "1",
		}

		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Paid, "1")

		assert.NoError(t, err)
		assert.Equal(t, order, nextOrder)
	})

	t.Run("transition paid to collect", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		order := &entity.Order{
			OrderStatus: entity.Paid,
			ID:          "1",
		}

		nextOrder := &entity.Order{
			OrderStatus: entity.Collect,
			ID:          "1",
		}

		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Collect, "1")

		assert.NoError(t, err)
		assert.Equal(t, order, nextOrder)
	})

	t.Run("transition collect to collected", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		order := &entity.Order{
			OrderStatus: entity.Collect,
			ID:          "1",
		}

		nextOrder := &entity.Order{
			OrderStatus: entity.Collected,
			ID:          "1",
		}

		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Collected, "1")

		assert.NoError(t, err)
		assert.Equal(t, order, nextOrder)
	})

	t.Run("transition collected to delivery", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		order := &entity.Order{
			OrderStatus: entity.Collected,
			ID:          "1",
		}

		nextOrder := &entity.Order{
			OrderStatus: entity.Delivery,
			ID:          "1",
		}

		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Delivery, "1")

		assert.NoError(t, err)
		assert.Equal(t, order, nextOrder)
	})

	t.Run("transition delivery to done", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		order := &entity.Order{
			OrderStatus: entity.Delivery,
			ID:          "1",
		}

		nextOrder := &entity.Order{
			OrderStatus: entity.Done,
			ID:          "1",
		}

		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Done, "1")

		assert.NoError(t, err)
		assert.Equal(t, order, nextOrder)
	})

	t.Run("wrong transition delivery to cancelled", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		order := &entity.Order{
			OrderStatus: entity.Delivery,
			ID:          "1",
		}

		repo.EXPECT().UpdateOrder(mock.Anything, order).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Cancelled, "1")

		assert.Error(t, err)
	})

	t.Run("wrong transition done to cancelled", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		order := &entity.Order{
			OrderStatus: entity.Done,
			ID:          "1",
		}

		repo.EXPECT().UpdateOrder(mock.Anything, order).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Cancelled, "1")

		assert.Error(t, err)
	})

	t.Run("wrong transition collected to paid", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo)

		order := &entity.Order{
			OrderStatus: entity.Collected,
			ID:          "1",
		}

		repo.EXPECT().UpdateOrder(mock.Anything, order).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Paid, "1")

		assert.Error(t, err)
	})

}
