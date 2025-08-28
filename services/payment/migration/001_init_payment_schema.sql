CREATE SCHEMA IF NOT EXISTS payment;

CREATE TABLE payment.payments (
     id BIGSERIAL PRIMARY KEY,
     payment_id UUID NOT NULL UNIQUE,
     external_order_id UUID NOT NULL,
     user_id BIGINT NOT NULL,
     idempotency_key UUID NOT NULL,
     amount NUMERIC(18,2) NOT NULL,
     currency CHAR(3) NOT NULL,
     status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
     created_at TIMESTAMPTZ NOT NULL,
     updated_at TIMESTAMPTZ NOT NULL
);

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

-- TODO: add index to payments and outbox.