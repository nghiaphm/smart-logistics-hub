DROP TABLE IF EXISTS workspace_users;
ALTER TABLE workspaces
    DROP COLUMN keycloak_group_id,
    DROP COLUMN is_deleted;
-- Lưu ý: KHÔNG drop is_active — cột này có từ 000002, không phải của 000012.
