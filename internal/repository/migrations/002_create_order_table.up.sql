CREATE TYPE order_status AS ENUM (
    'NEW',
    'PROCESSING',
    'INVALID',
    'PROCESSED'
);

CREATE TABLE orders (
    number VARCHAR PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status order_status NOT NULL DEFAULT 'NEW',
    accrual DECIMAL(10,2),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
