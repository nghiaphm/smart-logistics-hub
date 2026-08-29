-- TASK-087: Domain vehicle — bảng phương tiện (Fleet), CRUD cơ bản
-- Ràng buộc task: KHÔNG có quan hệ gán tài xế/chuyến (tính năng Dispatch để sau)
CREATE TABLE vehicles (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    license_plate VARCHAR(50) NOT NULL,
    type VARCHAR(50) DEFAULT '',
    capacity DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'ACTIVE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT '',
    UNIQUE KEY uq_license_plate (license_plate),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
