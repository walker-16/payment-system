package repository

import (
	"context"
	"encoding/json"
	"payment-system/pkg/db"
	"time"

	"github.com/walker-16/payment-system/services/payment/internal/domain"
)

type PaymentRepo interface {
	InsertPayment(ctx context.Context, p *domain.Payment) error
}

type PaymentRepository struct {
	db db.DB
}

func NewPaymentRepository(db db.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) InsertPayment(ctx context.Context, p *domain.Payment) error {
	// start transaction
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// insert payment
	now := time.Now()
	paymentInsert := `
		INSERT INTO payment.payments
		(payment_id, external_order_id, user_id, idempotency_key, amount, currency, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`

	if _, err := tx.Exec(ctx, paymentInsert,
		p.PaymentID,
		p.ExternalOrderID,
		p.UserID,
		p.IdempotencyKey,
		p.Amount,
		p.Currency,
		p.Status,
		now,
		now,
	); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// marshal payment as event payload
	payload, err := json.Marshal(p)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// insert outbox event
	outboxInsert := `
		INSERT INTO payment.outbox
		(aggregate_id, aggregate_type, event_type, payload, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	if _, err := tx.Exec(ctx, outboxInsert,
		p.PaymentID,
		"payment",
		"payment_created",
		payload,
		"PENDING",
		now,
		now,
	); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return nil
}
