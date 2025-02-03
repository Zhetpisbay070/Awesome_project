package server_test

import (
	"awesomeProject1/internal/entity"
	_ "awesomeProject1/internal/entity"
	internalMock "awesomeProject1/internal/mock"
	"awesomeProject1/internal/server"
	"bytes"
	_ "bytes"
	"encoding/json"
	_ "encoding/json"
	_ "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "github.com/stretchr/testify/mock"
	"net/http"
	_ "net/http"
	"net/http/httptest"
	_ "net/http/httptest"
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

		respDTO := &server.Order{
			ID:               "ORD123456",
			UserID:           "USR78910",
			ProductIDs:       []string{"PROD001", "PROD002", "PROD003"},
			CreatedAt:        time.Time(server.AwesomeTime(time.Now())),
			UpdatedAt:        time.Time(server.AwesomeTime(time.Now())),
			DeliveryDeadLine: time.Time(server.AwesomeTime(time.Now().Add(48 * time.Hour))),
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
			CreatedAt:        time.Time(server.AwesomeTime(time.Now())),
			UpdatedAt:        time.Time(server.AwesomeTime(time.Now())),
			DeliveryDeadLine: time.Time(server.AwesomeTime(time.Now().Add(48 * time.Hour))), // Delivery deadline set to 2 days from now
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

func TestServer_EditOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repoMock := internalMock.NewDB(t)
		serviceMock := internalMock.NewOrderService(t)

		reqDTO := &server.EditOrderRequest{
			OrderID:  "1",
			Products: []string{"prod1", "prod2"},
			Address:  "Old address",
		}

		reqJSON, err := json.Marshal(reqDTO)
		assert.NoError(t, err)

		repoMock.EXPECT().GetOrderByID(mock.Anything, mock.Anything).Return(&entity.Order{
			ID: "1", ProductIDs: []string{"prod1", "prod2"}, Address: "Old address", OrderStatus: entity.Created,
		}, nil)

		repoMock.EXPECT().UpdateOrder(mock.Anything, mock.Anything).Return(nil)
		//	ID:          "1",
		//	Address:     "New address",
		//	ProductIDs:  []string{"prod1", "prod2"},
		//	OrderStatus: entity.Created,
		//}).Return(nil)

		s := server.NewServer(serviceMock, logrus.New())

		r := s.EditRouter()

		w := httptest.NewRecorder()

		body := bytes.NewBuffer(reqJSON)

		req, err := http.NewRequest("POST", "/edit-router", body)
		assert.NoError(t, err)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, string(reqJSON), w.Body.String())

	})
}
