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

type MockOrderService struct {
	ResponseType MockResponseType
}

func NewMockOrderService(responseType MockResponseType) *MockOrderService {
	return &MockOrderService{ResponseType: responseType}
}

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
