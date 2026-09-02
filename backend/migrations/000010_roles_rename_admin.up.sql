-- Đổi tên role seed "admin" → "system_admin".
-- LÝ DO CHỌN ĐỔI TÊN (không INSERT thêm): role 'admin' trong seed migration
-- 000006 được tạo để đại diện "System Administrator with full access to
-- logistics operations" — chính là system admin; tên "admin" chỉ là đặt khác
-- đi trước khi chuẩn hoá. Mọi code đều dùng tên chuẩn "system_admin"
-- (backend RequireRole("system_admin"), frontend isSystemAdmin(), realm role
-- Keycloak) — xem FIX-004 / WN-006. user_roles tham chiếu qua FK id (không
-- phải name) nên rename không phá vỡ dữ liệu. Role "user" giữ nguyên.
UPDATE roles SET name = 'system_admin' WHERE name = 'admin';
