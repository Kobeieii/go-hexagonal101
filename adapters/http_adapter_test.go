package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"kobeieii/core"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type mockOrderService struct {
	mock.Mock
}

func (m *mockOrderService) CreateOrder(order *core.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func TestCreateOrderHandler(t *testing.T) {
	mockService := new(mockOrderService)
	handler := NewHttpOrderHandler(mockService)

	app := fiber.New()
	app.Post("/orders", handler.CreateOrder)

	tests := []struct {
		name             string
		requestBody      string
		mockReturn       error
		expectedStatus   int
		expectedError    string
	}{
		{
			name:           "successful order creation",
			requestBody:    `{"total": 100}`,
			mockReturn:     nil,
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name:             "failed order creation (total less than 0)",
			requestBody:      `{"total": -200}`,
			mockReturn:       errors.New("total must be positive"),
			expectedStatus:   http.StatusInternalServerError,
			expectedError:    "total must be positive",
		},
		{
			name:           "failed order creation (invalid JSON input)",
			requestBody:    `{"total": "invalid"}`,
			mockReturn:     nil, // service didn't return anything
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				mockService.ExpectedCalls = nil // Reset mock expectations
			})

			mockService.On("CreateOrder", mock.AnythingOfType("*core.Order")).Return(tt.mockReturn)

			req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)

			assert.NoError(t, err, "expected no error in making the request")
			defer res.Body.Close()

			assert.Equal(t, tt.expectedStatus, res.StatusCode, "expected status code to match")
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"), "expected Content-Type to be application/json")

			body, err := io.ReadAll(res.Body)
			assert.NoError(t, err, "expected no error in reading the response body")

			if tt.expectedError != "" {
				var result ErrorResponse
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "expected no error in unmarshaling the response body")
				assert.Equal(t, tt.expectedError, result.Error, "expected error message to match")
			}

			mockService.AssertExpectations(t)
		})
	}
}
