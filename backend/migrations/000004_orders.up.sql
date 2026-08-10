CREATE TABLE orders(
    id SERIAL PRIMARY KEY,
    table_number INTEGER NOT NULL REFERENCES tables(number),
    waiter_id INTEGER NOT NULL REFERENCES users(id),
    is_served BOOLEAN DEFAULT FALSE,
    is_payed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMP,
);