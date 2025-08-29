package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payment-system/pkg/logger"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/test-go/testify/require"
	"github.com/walker-16/payment-system/services/payment/internal/domain"
	"github.com/walker-16/payment-system/services/payment/internal/order"
)

type MockLogger struct{}

func (l *MockLogger) Debug(msg string, args ...any)  {}
func (l *MockLogger) Info(msg string, args ...any)   {}
func (l *MockLogger) Warn(msg string, args ...any)   {}
func (l *MockLogger) Error(msg string, args ...any)  {}
func (l *MockLogger) Fatal(msg string, args ...any)  {}
func (l *MockLogger) With(args ...any) logger.Logger { return l }

type MockRepo struct {
	InsertFunc func(ctx context.Context, p *domain.Payment) error
}

func (m *MockRepo) InsertPayment(ctx context.Context, p *domain.Payment) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, p)
	}
	return nil
}

// TestCreatePayment_Success verifies that a valid request with all required headers
// and a correct external order ID creates a payment successfully and returns StatusAccepted.
func TestCreatePayment_Success(t *testing.T) {
	app := fiber.New()

	mockRepo := &MockRepo{}
	mockOrderService := order.NewMockOrderService(order.MockSuccess)
	mockLogger := &MockLogger{}

	h := NewPaymentHandler(mockOrderService, mockRepo, mockLogger)
	app.Post("/payments", h.CreatePayment)

	externalOrderID := uuid.New()
	reqBody, _ := json.Marshal(map[string]string{
		"external_order_id": externalOrderID.String(),
	})

	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(reqBody))
	req.Header.Set("idempotency-key", uuid.New().String())
	req.Header.Set("x-user-id", "1")
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, fiber.StatusAccepted, resp.StatusCode)
}

// TestCreatePayment_MissingHeaders checks that if required headers are missing,
// the handler returns StatusBadRequest instead of proceeding.
func TestCreatePayment_MissingHeaders(t *testing.T) {
	app := fiber.New()

	mockRepo := &MockRepo{}
	mockOrderService := order.NewMockOrderService(order.MockSuccess)
	h := NewPaymentHandler(mockOrderService, mockRepo, &MockLogger{})

	app.Post("/payments", h.CreatePayment)

	req := httptest.NewRequest(http.MethodPost, "/payments", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}

// TestCreatePayment_OrderServiceError simulates an error returned by the order service.
// The handler should catch the error and return StatusBadRequest.
func TestCreatePayment_OrderServiceError(t *testing.T) {
	app := fiber.New()

	mockRepo := &MockRepo{}
	mockOrder := order.NewMockOrderService(order.MockErrorInternal)

	h := NewPaymentHandler(mockOrder, mockRepo, &MockLogger{})

	app.Post("/payments", h.CreatePayment)

	reqBody, _ := json.Marshal(map[string]string{
		"external_order_id": uuid.New().String(),
	})
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(reqBody))
	req.Header.Set("idempotency-key", uuid.New().String())
	req.Header.Set("x-user-id", "1")
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}
