CREATE SCHEMA IF NOT EXISTS payment;

CREATE TABLE payment.payments (
     id BIGSERIAL PRIMARY KEY,
     payment_id UUID NOT NULL UNIQUE,
     external_order_id UUID NOT NULL,
     user_id BIGINT NOT NULL,
     idempotency_key UUID NOT NULL UNIQUE,
     amount NUMERIC(18,2) NOT NULL,
     currency CHAR(3) NOT NULL,
     status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
     created_at TIMESTAMPTZ NOT NULL,
     updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX idx_payments_payment_id
ON payment.payments (payment_id);

CREATE UNIQUE INDEX idx_payments_idempotency_key
ON payment.payments (idempotency_key);

CREATE INDEX idx_payments_created_at
ON payment.payments (created_at);

CREATE TABLE payment.outbox (
    id BIGSERIAL PRIMARY KEY,            
    aggregate_id UUID NOT NULL,          
    aggregate_type VARCHAR(50) NOT NULL, 
    event_type VARCHAR(50) NOT NULL,     
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',        
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_outbox_status_created_at 
ON payment.outbox (status, created_at);

-- // TODO: review again table for pending index.