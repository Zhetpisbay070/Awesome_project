package service_test

import (
	"awesomeProject1/internal/entity"
	InternalMock "awesomeProject1/internal/mock"
	"awesomeProject1/internal/service"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestService_CreateOrder(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

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
		assert.Equal(t, order.ID, "salam")
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

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

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

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

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

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

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

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

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

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

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

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		order := &entity.Order{
			OrderStatus: entity.Created,
			ID:          "1",
		}

		repo.EXPECT().GetOrderByID(mock.Anything, "1").Return(order, nil)
		fmt.Println("Before transition:", order.OrderStatus)

		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Paid, "1")
		fmt.Println("After transition:", order.OrderStatus)

		assert.NoError(t, err)

	})

	t.Run("transition paid to collect", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		order := &entity.Order{
			OrderStatus: entity.Paid,
			ID:          "1",
		}

		repo.EXPECT().GetOrderByID(mock.Anything, "1").Return(order, nil)
		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Collect, "1")

		assert.NoError(t, err)
	})

	t.Run("transition collect to collected", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		order := &entity.Order{
			OrderStatus: entity.Collect,
			ID:          "1",
		}

		repo.EXPECT().GetOrderByID(mock.Anything, "1").Return(order, nil)
		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Collected, "1")

		assert.NoError(t, err)
	})

	t.Run("transition collected to delivery", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		order := &entity.Order{
			OrderStatus: entity.Collected,
			ID:          "1",
		}

		repo.EXPECT().GetOrderByID(mock.Anything, "1").Return(order, nil)
		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Delivery, "1")

		assert.NoError(t, err)
	})

	t.Run("transition delivery to done", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		order := &entity.Order{
			OrderStatus: entity.Delivery,
			ID:          "1",
		}

		repo.EXPECT().GetOrderByID(mock.Anything, "1").Return(order, nil)

		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Done, "1")

		assert.NoError(t, err)
	})

	t.Run("wrong transition delivery to cancelled", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		order := &entity.Order{
			OrderStatus: entity.Delivery,
			ID:          "1",
		}

		repo.EXPECT().GetOrderByID(mock.Anything, "1").Return(order, nil)
		fmt.Println("Current order status2:", order.OrderStatus)

		err := s.UpdateOrderStatus(context.Background(), entity.Cancelled, "1")

		assert.ErrorIs(t, err, entity.Pozdno)
	})

	t.Run("wrong transition done to cancelled", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		order := &entity.Order{
			OrderStatus: entity.Done,
			ID:          "1",
		}

		repo.EXPECT().GetOrderByID(mock.Anything, "1").Return(order, nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Cancelled, "1")

		assert.ErrorIs(t, err, entity.Pozdno)
	})

	t.Run("wrong transition collected to paid", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		order := &entity.Order{
			OrderStatus: entity.Collected,
			ID:          "1",
		}

		repo.EXPECT().GetOrderByID(mock.Anything, "1").Return(order, nil)

		err := s.UpdateOrderStatus(context.Background(), entity.Paid, "1")

		assert.ErrorIs(t, err, entity.InvalidTransition)
	})

}

func TestEditOrder(t *testing.T) {
	//
	t.Run("new products", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		repo.EXPECT().GetOrderByID(mock.Anything, mock.Anything).Return(&entity.Order{
			ID:         "1",
			Address:    "old address",
			ProductIDs: []string{"prod1, prod2"},
		}, nil)

		repo.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)

		editReq := &entity.EditOrderRequest{
			OrderID:  "1",
			Address:  "new address",
			Products: []string{"prod3", "prod4"},
		}

		_, err := s.EditOrder(context.Background(), editReq)
		assert.NoError(t, err)

	})

	t.Run("order cannot be edited in this current status", func(t *testing.T) {
		repo := InternalMock.NewDB(t)

		s := service.NewOrderService(repo, func() string {
			return "salam"
		}, func() time.Time {
			return time.Now()
		}, logrus.New())

		repo.EXPECT().GetOrderByID(mock.Anything, mock.Anything).Return(&entity.Order{
			ID:               "1",
			UserID:           "001",
			ProductIDs:       []string{"prod1", "prod2"},
			CreatedAt:        time.Time{},
			UpdatedAt:        time.Time{},
			DeliveryDeadLine: time.Time{},
			Price:            0,
			DeliveryType:     "Drone",
			Address:          "Chkalovo",
			OrderStatus:      entity.Delivery,
		}, nil)

		editReq := &entity.EditOrderRequest{
			OrderID:  "1",
			Products: []string{"prod3", "prod4"},
			Address:  "new address",
		}

		_, err := s.EditOrder(context.Background(), editReq)

		assert.ErrorIs(t, err, entity.OrderCannotBeEdited)

	})

}
