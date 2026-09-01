ALTER TABLE orders
    DROP FOREIGN KEY fk_orders_sender_workspace,
    DROP KEY idx_orders_sender_workspace,
    DROP COLUMN sender_workspace_id;
