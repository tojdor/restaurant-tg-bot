CREATE TABLE menu_items(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    price NUMERIC(10, 2) NOT NULL,
    category TEXT NOT NULL
)