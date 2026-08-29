-- TASK-085: Liên kết tracking_events với orders qua khoá ngoại order_id
-- (1 order có thể có nhiều tracking record; cột nullable để giữ tương thích
-- với tracking record cũ chưa có liên kết)
ALTER TABLE tracking_events ADD COLUMN order_id BIGINT DEFAULT NULL AFTER id;
ALTER TABLE tracking_events ADD INDEX idx_order_id (order_id);
ALTER TABLE tracking_events ADD CONSTRAINT fk_tracking_events_order
    FOREIGN KEY (order_id) REFERENCES orders(id);
