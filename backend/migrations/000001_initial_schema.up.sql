CREATE TABLE warehouses (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    warehouse_code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    lat DOUBLE DEFAULT 0,
    lng DOUBLE DEFAULT 0,
    contact_phone VARCHAR(20) DEFAULT '',
    manager_name VARCHAR(255) DEFAULT '',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_warehouse_code (warehouse_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) DEFAULT '',
    price DECIMAL(15,2) DEFAULT 0,
    weight_gram INT DEFAULT 0,
    length_cm DECIMAL(8,2) DEFAULT 0,
    width_cm DECIMAL(8,2) DEFAULT 0,
    height_cm DECIMAL(8,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT '',
    UNIQUE KEY uq_sku (sku)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE drivers (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    driver_code VARCHAR(50) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) DEFAULT '',
    vehicle_type VARCHAR(50) DEFAULT '',
    license_plate VARCHAR(20) DEFAULT '',
    status VARCHAR(20) DEFAULT 'AVAILABLE',
    current_lat DOUBLE DEFAULT 0,
    current_lng DOUBLE DEFAULT 0,
    warehouse_id BIGINT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT '',
    UNIQUE KEY uq_driver_code (driver_code),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE orders (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_code VARCHAR(50) NOT NULL,
    sender_name VARCHAR(255) DEFAULT '',
    sender_phone VARCHAR(20) DEFAULT '',
    sender_address TEXT,
    sender_province VARCHAR(100) DEFAULT '',
    sender_district VARCHAR(100) DEFAULT '',
    sender_ward VARCHAR(100) DEFAULT '',
    sender_postal_code VARCHAR(20) DEFAULT '',
    receiver_name VARCHAR(255) DEFAULT '',
    receiver_phone VARCHAR(20) DEFAULT '',
    receiver_address TEXT,
    receiver_province VARCHAR(100) DEFAULT '',
    receiver_district VARCHAR(100) DEFAULT '',
    receiver_ward VARCHAR(100) DEFAULT '',
    receiver_postal_code VARCHAR(20) DEFAULT '',
    status VARCHAR(30) DEFAULT 'PENDING',
    assigned_driver_id BIGINT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT '',
    UNIQUE KEY uq_order_code (order_code),
    INDEX idx_status (status),
    INDEX idx_driver (assigned_driver_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE order_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT DEFAULT NULL,
    product_name VARCHAR(255) DEFAULT '',
    quantity INT DEFAULT 1,
    weight_gram INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE inventory (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    available_qty INT DEFAULT 0,
    reserved_qty INT DEFAULT 0,
    damaged_qty INT DEFAULT 0,
    hold_qty INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(255) DEFAULT '',
    UNIQUE KEY uq_product_warehouse (product_id, warehouse_id),
    INDEX idx_product (product_id),
    INDEX idx_warehouse (warehouse_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE trips (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    trip_code VARCHAR(50) NOT NULL,
    driver_id BIGINT DEFAULT NULL,
    vehicle_license_plate VARCHAR(20) DEFAULT '',
    status VARCHAR(30) DEFAULT 'PLANNED',
    total_distance_km DECIMAL(8,2) DEFAULT 0,
    estimated_duration_min INT DEFAULT 0,
    actual_start_at TIMESTAMP NULL,
    actual_end_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by VARCHAR(255) DEFAULT '',
    UNIQUE KEY uq_trip_code (trip_code),
    INDEX idx_driver (driver_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE trip_stops (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    trip_id BIGINT NOT NULL,
    order_code VARCHAR(50) DEFAULT '',
    stop_type VARCHAR(30) DEFAULT 'PICKUP',
    address TEXT,
    lat DOUBLE DEFAULT 0,
    lng DOUBLE DEFAULT 0,
    status VARCHAR(30) DEFAULT 'PENDING',
    planned_at TIMESTAMP NULL,
    arrived_at TIMESTAMP NULL,
    departure_at TIMESTAMP NULL,
    FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE tracking_events (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    order_code VARCHAR(50) NOT NULL,
    driver_code VARCHAR(50) DEFAULT '',
    status_update VARCHAR(100) NOT NULL,
    lat DOUBLE DEFAULT 0,
    lng DOUBLE DEFAULT 0,
    note TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_order_code (order_code),
    INDEX idx_driver_code (driver_code),
    INDEX idx_timestamp (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE billing (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    billing_code VARCHAR(50) NOT NULL,
    order_code VARCHAR(50) NOT NULL,
    amount_total DECIMAL(15,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'VND',
    payment_method VARCHAR(30) DEFAULT 'COD',
    payment_status VARCHAR(30) DEFAULT 'PENDING',
    transaction_id VARCHAR(255) DEFAULT '',
    payer_name VARCHAR(255) DEFAULT '',
    payer_phone VARCHAR(20) DEFAULT '',
    payer_email VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    paid_at TIMESTAMP NULL,
    created_by VARCHAR(255) DEFAULT '',
    UNIQUE KEY uq_billing_code (billing_code),
    INDEX idx_order_code (order_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE inbounds (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    receipt_code VARCHAR(50) NOT NULL,
    supplier_name VARCHAR(255) DEFAULT '',
    status VARCHAR(30) DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    created_by VARCHAR(255) DEFAULT '',
    UNIQUE KEY uq_receipt_code (receipt_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE inbound_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    inbound_id BIGINT NOT NULL,
    product_id BIGINT DEFAULT NULL,
    expected_qty INT DEFAULT 0,
    received_qty INT DEFAULT 0,
    rejected_qty INT DEFAULT 0,
    qc_passed INT DEFAULT 0,
    FOREIGN KEY (inbound_id) REFERENCES inbounds(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE ai_events (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    event_code VARCHAR(50) NOT NULL,
    license_plate VARCHAR(20) NOT NULL,
    confidence_score DECIMAL(5,4) DEFAULT 0,
    event_type VARCHAR(20) DEFAULT 'INBOUND',
    gate_id VARCHAR(50) DEFAULT '',
    timestamp TIMESTAMP NULL,
    matched_driver_id BIGINT DEFAULT NULL,
    matched_trip_id BIGINT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_event_code (event_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
