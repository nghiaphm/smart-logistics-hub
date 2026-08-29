ALTER TABLE tracking_events DROP FOREIGN KEY fk_tracking_events_order;
ALTER TABLE tracking_events DROP INDEX idx_order_id;
ALTER TABLE tracking_events DROP COLUMN order_id;
