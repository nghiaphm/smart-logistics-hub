ALTER TABLE orders
    ADD COLUMN sender_workspace_id BIGINT DEFAULT NULL,
    ADD KEY idx_orders_sender_workspace (sender_workspace_id),
    ADD CONSTRAINT fk_orders_sender_workspace FOREIGN KEY (sender_workspace_id) REFERENCES workspaces(id);
