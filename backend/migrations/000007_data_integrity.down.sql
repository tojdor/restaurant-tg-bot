DROP INDEX menu_items_category_index;
DROP INDEX orders_table_number_index;
DROP INDEX orders_waiter_id_index;

ALTER TABLE order_items
    DROP CONSTRAINT order_items_count_positive,
    DROP CONSTRAINT order_items_primary_key;

ALTER TABLE menu_items
    DROP CONSTRAINT menu_items_price_positive,
    ALTER COLUMN price TYPE NUMERIC(10, 2) USING price::numeric;
