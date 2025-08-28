package domain

import (
	"time"

	"github.com/google/uuid" // TODO:
)

type Outbox struct {
	ID            int64     `db:"id"`
	AggregateID   uuid.UUID `db:"aggregate_id"`
	AggregateType string    `db:"aggregate_type"`
	EventType     string    `db:"event_type"`
	Payload       []byte    `db:"payload"` // JSON serializado
	Status        string    `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
