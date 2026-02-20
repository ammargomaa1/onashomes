ALTER TABLE customers ADD COLUMN store_front_id BIGINT;
CREATE INDEX idx_customers_store_front_id ON customers(store_front_id);
-- Optional: Add Foreign Key constraint if store_fronts table exists and you want strict integrity
-- ALTER TABLE customers ADD CONSTRAINT fk_customers_store_front FOREIGN KEY (store_front_id) REFERENCES store_fronts(id);
