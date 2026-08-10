CREATE TABLE order_items(
    order_id INTEGER NOT NULL REFERENCES orders(id),
    menu_item_id INTEGER NOT NULL REFERENCES menu_items(id),
    count INTEGER NOT NULL,
    is_ready BOOLEAN DEFAULT FALSE
);