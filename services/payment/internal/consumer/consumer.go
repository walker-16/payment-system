package consumer

import (
	"context"
	"fmt"
	"payment-system/pkg/db"
	"payment-system/pkg/domain"
	"payment-system/pkg/kafka"
	"payment-system/pkg/logger"
	"time"
)

// TODO: add env var to modify default batch size.
const batchSize = 10

// OutboxConsumer is responsible for polling the outbox table, processing events,
// and publishing them to Kafka.
type OutboxConsumer struct {
	db       db.DB
	producer *kafka.Producer
	logger   logger.Logger
	interval time.Duration
}

// NewOutboxConsumer creates a new OutboxConsumer with the given database, Kafka producer,
// logger, and processing interval.
func NewOutboxConsumer(db db.DB, producer *kafka.Producer,
	logger logger.Logger, interval time.Duration) *OutboxConsumer {
	return &OutboxConsumer{
		db:       db,
		producer: producer,
		logger:   logger,
		interval: interval,
	}
}

// Start launches the consumer loop which continuously polls and processes outbox events
// until the provided context is canceled.
func (c *OutboxConsumer) Start(ctx context.Context) {
	c.logger.Info("starting outbox consumer")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("outbox consumer stopped due to context cancellation")
			return
		default:
			if err := c.processBatch(ctx); err != nil {
				c.logger.Error("failed to process outbox batch", logger.Error(err))
			}
			time.Sleep(c.interval)
		}
	}
}

// processBatch retrieves a batch of pending outbox events within a transaction
// and processes each event.
func (c *OutboxConsumer) processBatch(ctx context.Context) error {
	tx, err := c.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	var outboxes []domain.Outbox
	query := `
		SELECT id, aggregate_id, aggregate_type, event_type, payload, status, created_at, updated_at
		FROM payment.outbox
		WHERE status = 'PENDING'
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	if err := tx.Select(ctx, &outboxes, query, batchSize); err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("select outbox: %w", err)
	}

	for _, o := range outboxes {
		if err := c.processOutbox(ctx, tx, &o); err != nil {
			c.logger.Error("failed to process outbox event", logger.Error(err))
		}
	}

	return tx.Commit(ctx)
}

func (c *OutboxConsumer) processOutbox(ctx context.Context,
	tx db.Tx, o *domain.Outbox) error {

	//TODO: continue implementation.
	return nil
}
