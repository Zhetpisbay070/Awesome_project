package server_test

import (
	"awesomeProject1/internal/entity"
	_ "awesomeProject1/internal/entity"
	internalMock "awesomeProject1/internal/mock"
	"awesomeProject1/internal/server"
	"bytes"

	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_CreateOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		serviceMock := internalMock.NewOrderService(t)

		reqDTO := &server.CreateOrderRequest{
			UserID:       "USR78910",
			Price:        99.99,
			DeliveryType: "Standard",
		}

		reqJSON, err := json.Marshal(reqDTO)
		assert.NoError(t, err)

		timeNow := time.Now()

		respDTO := &server.Order{
			ID:               "ORD123456",
			UserID:           "USR78910",
			ProductIDs:       []string{"PROD001", "PROD002", "PROD003"},
			CreatedAt:        server.AwesomeTime(timeNow),
			UpdatedAt:        server.AwesomeTime(timeNow),
			DeliveryDeadLine: server.AwesomeTime(timeNow.Add(48 * time.Hour)),
			Price:            99.99,
			DeliveryType:     "Standard",
			Address:          "123 Main Street, City, Country",
			OrderStatus:      "Processing",
		}

		respJSON, err := json.Marshal(respDTO)
		assert.NoError(t, err)

		serviceMock.EXPECT().CreateOrder(mock.Anything, mock.Anything).Return(&entity.Order{
			ID:               "ORD123456",
			UserID:           "USR78910",
			ProductIDs:       []string{"PROD001", "PROD002", "PROD003"},
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			DeliveryDeadLine: time.Now().Add(48 * time.Hour),
			Price:            99.99,
			DeliveryType:     "Standard",
			Address:          "123 Main Street, City, Country",
			OrderStatus:      "Processing",
		}, nil)

		s := server.NewServer(serviceMock, logrus.New())

		r := s.SetupRouter()

		w := httptest.NewRecorder()

		body := bytes.NewBuffer(reqJSON)

		req, err := http.NewRequest("POST", "/create", body)
		assert.NoError(t, err)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(respJSON), w.Body.String())
	})
}
func TestServer_UpdateOrderStatus(t *testing.T) {
	t.Run("order update success", func(t *testing.T) {
		serviceMock := internalMock.NewOrderService(t)

		reqDTO := &server.UpdateOrderStatusRequest{
			OrderID:     "1",
			OrderStatus: "Delivered",
		}

		reqJSON, err := json.Marshal(reqDTO)
		assert.NoError(t, err)

		serviceMock.EXPECT().UpdateOrderStatus(mock.Anything, mock.Anything, mock.Anything).Return(nil)

		s := server.NewServer(serviceMock, logrus.New())

		r := s.SetupRouter()

		w := httptest.NewRecorder()

		body := bytes.NewBuffer(reqJSON)

		req, err := http.NewRequest("POST", "/update", body)
		assert.NoError(t, err)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})
}

func TestServer_GetOrders(t *testing.T) {
	t.Run("get orders success", func(t *testing.T) {
		serviceMock := internalMock.NewOrderService(t)

		reqDTO := &entity.GetOrders{
			UserID: "1",
			Limit:  10,
			Page:   1,
			Asc:    true,
		}

		reqJSON, err := json.Marshal(reqDTO)
		assert.NoError(t, err)

		serviceMock.EXPECT().GetOrders(mock.Anything, mock.Anything).Return([]entity.Order{}, nil)

		s := server.NewServer(serviceMock, logrus.New())

		r := s.SetupRouter()

		w := httptest.NewRecorder()

		body := bytes.NewBuffer(reqJSON)

		req, err := http.NewRequest("POST", "/getOrders", body)
		assert.NoError(t, err)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})

}

func TestServer_EditOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		serviceMock := internalMock.NewOrderService(t)

		reqDTO := &entity.EditOrderRequest{
			OrderID:  "1",
			Products: []string{"prod1", "prod2"},
			Address:  "New Address",
		}

		reqJSON, err := json.Marshal(reqDTO)
		assert.NoError(t, err)

		serviceMock.EXPECT().EditOrder(mock.Anything, mock.Anything).Return(&entity.Order{
			ID:         "1",
			ProductIDs: []string{"prod1", "prod2"},
			Address:    "New Address",
		}, nil)

		s := server.NewServer(serviceMock, logrus.New())

		r := s.SetupRouter()

		w := httptest.NewRecorder()

		body := bytes.NewBuffer(reqJSON)

		req, err := http.NewRequest("POST", "/edit", body)
		assert.NoError(t, err)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("order not found", func(t *testing.T) {
		serviceMock := internalMock.NewOrderService(t)

		reqDTO := &entity.EditOrderRequest{
			OrderID:  "0000",
			Products: []string{"prod1", "prod2"},
			Address:  "New Address",
		}

		reqJSON, err := json.Marshal(reqDTO)
		assert.NoError(t, err)

		serviceMock.EXPECT().EditOrder(mock.Anything, mock.Anything).Return(nil, entity.ErrOrderNotFound)

		s := server.NewServer(serviceMock, logrus.New())

		r := s.SetupRouter()

		w := httptest.NewRecorder()

		body := bytes.NewBuffer(reqJSON)

		req, err := http.NewRequest("POST", "/edit", body)
		assert.NoError(t, err)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, "{\"error\":\"order not found\"}", w.Body.String())

	})

	t.Run("server error", func(t *testing.T) {
		serviceMock := internalMock.NewOrderService(t)

		reqDTO := &entity.EditOrderRequest{
			OrderID:  "1",
			Products: []string{"prod1", "prod2"},
			Address:  "New Address",
		}

		reqJSON, err := json.Marshal(reqDTO)
		assert.NoError(t, err)

		serviceMock.EXPECT().EditOrder(mock.Anything, mock.Anything).Return(nil, errors.New("server error"))

		s := server.NewServer(serviceMock, logrus.New())

		r := s.SetupRouter()

		w := httptest.NewRecorder()

		body := bytes.NewBuffer(reqJSON)

		req, err := http.NewRequest("POST", "/edit", body)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, `{"error":"server error"}`, w.Body.String())

	})
}
