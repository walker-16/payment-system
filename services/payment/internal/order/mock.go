package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type MockResponseType int

const (
	MockSuccess MockResponseType = iota
	MockErrorNotFound
	MockErrorUserMismatch
	MockErrorInternal
	MockErrorBadRequest
)

// MockOrderService is a mock implementation of the Order Service interface.
type MockOrderService struct {
	ResponseType MockResponseType
}

// NewMockOrderService creates a new mock order service with the specified
// response type.
func NewMockOrderService(responseType MockResponseType) *MockOrderService {
	return &MockOrderService{ResponseType: responseType}
}

// GetOrderByExternalIDForUser returns a mock order or an error depending
// on the configured ResponseType
func (m *MockOrderService) GetOrderByExternalIDForUser(ctx context.Context,
	externalOrderID uuid.UUID, userID uint32) (*Order, error) {
	switch m.ResponseType {
	case MockErrorNotFound:
		return nil, errors.New("order not found")
	case MockErrorUserMismatch:
		return nil, errors.New("order does not belong to the user")
	case MockErrorInternal:
		return nil, errors.New("internal error occurred")
	case MockErrorBadRequest:
		return nil, errors.New("bad request: invalid external order id")
	case MockSuccess:
		return &Order{
			ExternalID:  externalOrderID,
			UserID:      userID,
			ServiceName: "Service A",
			Amount:      99.99,
			Currency:    "USD",
			BankAccount: "123456789",
			BankCode:    "XY",
		}, nil
	default:
		return nil, errors.New("unknown mock error type")
	}
}
