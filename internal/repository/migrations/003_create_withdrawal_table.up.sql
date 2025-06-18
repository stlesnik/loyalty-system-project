CREATE TABLE withdrawals (
     user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
     order_number VARCHAR PRIMARY KEY,
     amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
     processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);