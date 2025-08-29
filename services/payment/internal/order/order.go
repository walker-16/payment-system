package order

import (
	"context"

	"github.com/google/uuid"
)

// Order represents a user's order with relevant payment and service details.
type Order struct {
	ExternalID  uuid.UUID // External order ID
	ServiceName string    // Name of the service
	Amount      float64   // Amount
	Currency    string    // Currency code
	BankAccount string    // Bank account of the service
	BankCode    string    // Bank code
	UserID      uint32    // User ID
}

// Service defines the interface for retrieving order details.
type Service interface {
	GetOrderByExternalIDForUser(ctx context.Context, externalOrderID uuid.UUID, userID uint32) (*Order, error)
}
