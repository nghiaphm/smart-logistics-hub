-- TASK-096: mở rộng bảng workspaces + tạo bảng workspace_users
--
-- 1) workspaces: chỉ thêm cột CÒN THIẾU. `is_active BOOLEAN DEFAULT TRUE` đã có
--    sẵn từ migration 000002 (cả DB dev lẫn fresh migration đều có trước khi
--    000012 chạy) — không thêm lại. KHÔNG thêm created_by/updated_by ở bảng
--    này (chưa có luồng tạo workspace thật, xem mô tả task).
ALTER TABLE workspaces
    ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE,
    ADD COLUMN keycloak_group_id VARCHAR(255) NULL;

-- 2) workspace_users: membership workspace ↔ user (1 user có thể thuộc nhiều
--    workspace). user_id/created_by/updated_by khớp khoá chính bảng users
--    (BIGINT — đã xác nhận bảng users thật). created_by/updated_by là FK users
--    vì Super Admin gán quyền cần ghi người thao tác. KHÔNG có is_deleted
--    (giữ đúng thiết kế bảng gốc).
CREATE TABLE workspace_users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    workspace_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_by BIGINT DEFAULT NULL,
    updated_by BIGINT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_workspace_user (workspace_id, user_id),
    INDEX idx_workspace_users_user_id (user_id),
    CONSTRAINT fk_workspace_users_workspace FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE,
    CONSTRAINT fk_workspace_users_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_workspace_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
    CONSTRAINT fk_workspace_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
