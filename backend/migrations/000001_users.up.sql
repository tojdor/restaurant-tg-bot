CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    nickname TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'waiter', 'kitchen'))
);