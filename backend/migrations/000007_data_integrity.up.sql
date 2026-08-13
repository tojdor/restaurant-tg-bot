ALTER TABLE menu_items
    ALTER COLUMN price TYPE INTEGER USING price::integer,
    ADD CONSTRAINT menu_items_price_positive CHECK (price > 0);

ALTER TABLE order_items
    ADD CONSTRAINT order_items_primary_key PRIMARY KEY (order_id, menu_item_id),
    ADD CONSTRAINT order_items_count_positive CHECK (count > 0);

CREATE INDEX orders_waiter_id_index ON orders (waiter_id);
CREATE INDEX orders_table_number_index ON orders (table_number);
CREATE INDEX menu_items_category_index ON menu_items (category);
