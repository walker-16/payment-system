package domain

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID              int64     `db:"id"`
	PaymentID       uuid.UUID `db:"payment_id"`
	ExternalOrderID uuid.UUID `db:"external_order_id"`
	UserID          uint32    `db:"user_id"`
	IdempotencyKey  uuid.UUID `db:"idempotency_key"`
	Amount          float64   `db:"amount"`
	Currency        string    `db:"currency"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
